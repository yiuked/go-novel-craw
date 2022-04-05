package craw

import (
	"github.com/yiuked/go-novel/src/storege"
	"github.com/yiuked/go-novel/src/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetBooksChapter 通过 BookList 下载章节
func (c *StandardCrawAction) GetBooksChapter(rule *BookCrawRule, start, end time.Duration) {
	c.init(rule)
	pageSize, page := 100, 0
	version := time.Now().Unix()
	if start > 0 && end > 0 {
		storege.DB().Model(&storege.BookList{}).Where("book_state=1 AND book_platform=? AND updated_at BETWEEN ? AND ?",
			c.rule.PlatformName, time.Now().Add(-start), time.Now().Add(-end)).
			Update("version", version)
	}
	for {
		// 采用单协程分配，多协程处理，避免资源分配不重复
		var bookList []storege.BookList
		if start > 0 && end > 0 {
			storege.DB().Where("book_state=1 AND book_platform=? AND version=?", c.rule.PlatformName, version).Offset(page * pageSize).Limit(pageSize).Find(&bookList)
		} else {
			storege.DB().Where("book_state=1 AND book_platform=?", c.rule.PlatformName).Offset(page * pageSize).Limit(pageSize).Find(&bookList)
		}
		if len(bookList) <= 0 {
			log.Println("get book summary task done,start check update")
			break
		}
		page++
		c.Wait()
		go func(bookList []storege.BookList) {
			defer c.Done()
			for _, s := range bookList {
				url := strings.Replace(c.rule.BookChapter.BookChapterURL, CrawBookID, s.BookID, 1)
				bytes, err := utils.Get(url, nil)
				if err != nil {
					continue
				}

				chapters, process := c.getBookChapter(bytes)
				if len(chapters) <= 0 {
					continue
				}
				// 从章节关键词中更新进度
				if len(process) > 0 {
					storege.DB().Model(&storege.BookDetail{}).Where("book_id=? AND book_platform=? AND book_process=''", s.ID).
						Update("book_process", process)
					if process == c.rule.BookProcessUpdate.EndName {
						storege.DB().Model(&storege.BookList{}).Where("book_id=? AND book_platform=? AND book_process=''", s.ID).
							Update("book_process_state", 1)
					}
				}

				for id, name := range chapters {
					chapter := storege.BookChapter{
						BookChapterName: name,
						BookChapterID:   id,
						BookID:          s.BookID,
						CatID:           s.CatID,
						BookPlatform:    c.rule.PlatformName,
						BookChapterHash: utils.Md5(c.rule.PlatformName + s.BookID + id),
					}
					if err := storege.DB().Create(&chapter).Error; err != nil {
						log.Printf("save chapter to DB error,err[%v]\n", err)
					}
				}
				s.Version++
				storege.DB().Select("version").Save(&s)
				log.Println("saved done ", s.BookID)
			}
		}(bookList)
	}
}

func (c *StandardCrawAction) getBookChapter(detailHtml []byte) (map[string]string, string) {
	reg := regexp.MustCompile(c.rule.BookChapter.BookChapterPatten)
	chapters := reg.FindAllSubmatch(detailHtml, -1)

	chaptersMap := make(map[string]string)
	var isLastChapter string
	for _, chapter := range chapters {
		if len(chapter) < 3 {
			continue
		}
		isLastChapter = chapterCompare(isLastChapter, string(chapter[1]))
		if len(c.rule.BookChapter.Extend) > 0 {
			for _, ext := range c.rule.BookChapter.Extend {
				reg := regexp.MustCompile(ext["patten"])
				chapter[1] = reg.ReplaceAll(chapter[1], []byte(ext["replace"]))
				chapter[2] = []byte(reg.ReplaceAllString(string(chapter[2]), ext["replace"]))
			}
		}
		chaptersMap[string(chapter[1])] = string(chapter[2])
	}
	// 匹配连载状态
	var process string
	if len(c.rule.BookChapter.ProcessCheckKeywords) > 0 {
		process = c.rule.BookProcessUpdate.UpdatingName
		for _, keyword := range c.rule.BookChapter.ProcessCheckKeywords {
			if strings.IndexAny(chaptersMap[isLastChapter], keyword) > 0 {
				process = c.rule.BookProcessUpdate.EndName
				break
			}
		}
	}
	return chaptersMap, process
}

func chapterCompare(last, new string) string {
	lastInt, err := strconv.Atoi(last)
	if err != nil {
		return new
	}
	newInt, err := strconv.Atoi(new)
	if err != nil {
		return last
	}
	if lastInt > newInt {
		return last
	}
	return new
}
