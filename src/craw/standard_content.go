package craw

import (
	"fmt"
	"github.com/yiuked/go-novel/src/storege"
	"github.com/yiuked/go-novel/src/utils"
	"log"
	"os"
	"path"
	"strings"
)

// GetBooksContent 书本内容采集
func (c *StandardCrawAction) GetBooksContent(rule *BookCrawRule) {
	c.init(rule)
	pageSize, page := 100, 0
	for {
		// 采用单协程分配，多协程处理，避免资源分配不重复
		var chapters []storege.BookChapter
		storege.DB().Where("book_chapter_state=0 AND book_platform=?", c.rule.PlatformName).Offset(page * pageSize).Limit(pageSize).Find(&chapters)
		if len(chapters) <= 0 {
			log.Println("books cover has download done")
			break
		}
		page++
		c.Wait()
		go func(chapters []storege.BookChapter) {
			defer c.Done()
			for _, chapter := range chapters {
				abs, rel := c.chapterFileName(chapter)
				_, err := os.Lstat(abs)
				if err == nil {
					chapter.BookChapterState = 1
					chapter.BookChapterLocal = rel
					storege.DB().Save(&chapter)
					log.Println(chapter.BookChapterID, "download")
					continue
				}
				url := strings.ReplaceAll(c.rule.BookChapter.BookChapterURL, CrawCatID, chapter.CatID)
				url = strings.ReplaceAll(url, CrawBookID, chapter.BookID)
				url = strings.ReplaceAll(url, CrawChapterID, chapter.BookChapterID)
				bytes, err := utils.Get(url, nil)
				if err != nil {
					log.Println(err)
					continue
				}
				var content string
				var isBool bool
				// 如果存在二级页采集
				if len(c.rule.BookContentPatten.NewPage.NewPageURLPatten.Patten) > 0 {
					url, isBool = GetSinglePatten(bytes, c.rule.BookContentPatten.NewPage.NewPageURLPatten)
					if !isBool {
						continue
					}
					newPageBytes, err := utils.Get(url, nil)
					if err != nil {
						log.Println(err)
						continue
					}
					content, isBool = GetBetweenPatten(newPageBytes, c.rule.BookContentPatten.NewPage.ContentPatten)
					if !isBool {
						continue
					}
				} else {
					content, isBool = GetBetweenPatten(bytes, c.rule.BookContentPatten.CurrentPage)
					if !isBool {
						continue
					}
				}

				err = os.WriteFile(abs, []byte(content), 0644)
				if err != nil {
					log.Println(err)
					continue
				}
				chapter.BookChapterState = 1
				chapter.BookChapterLocal = rel
				storege.DB().Save(&chapter)
				log.Println(chapter.BookChapterID, "download")
			}
		}(chapters)
	}
}

func (c *StandardCrawAction) chapterFileName(chapter storege.BookChapter) (absolute, relative string) {
	// 获取相对路径
	if len(c.rule.BookChapter.BookChapterLocalPath) > 0 {
		relative = strings.ReplaceAll(c.rule.BookChapter.BookChapterLocalPath, CrawCatID, chapter.CatID)
		relative = strings.ReplaceAll(relative, CrawBookID, chapter.BookID)
		relative = strings.ReplaceAll(relative, CrawChapterID, chapter.BookChapterID)
	} else {
		relative = fmt.Sprintf("/chapter/%s/%s/%s.html", chapter.CatID, chapter.BookID, chapter.BookChapterID)
	}

	// 获取绝对路径
	absolute = utils.GetDataDir() + relative
	_, err := os.Stat(path.Dir(absolute))
	if err != nil {
		err := os.MkdirAll(path.Dir(absolute), 0777)
		if err != nil {
			return "", ""
		}
	}
	return absolute, relative
}
