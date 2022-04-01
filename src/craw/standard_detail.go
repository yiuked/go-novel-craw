package craw

import (
	"github.com/yiuked/go-novel/src/storege"
	"github.com/yiuked/go-novel/src/utils"
	"log"
	"strings"
)

func (c *StandardCrawAction) GetBooksSummary(rule *BookCrawRule) {
	c.init(rule)
	pageSize, page := 100, 0
	for {
		// 采用单协程分配，多协程处理，避免资源分配不重复
		var bookList []storege.BookList
		storege.DB.Where("book_state=0 AND book_platform=?", c.rule.PlatformName).Offset(page * pageSize).Limit(pageSize).Find(&bookList)
		if len(bookList) <= 0 {
			log.Println("get book summary task done")
			break
		}
		page++
		c.Wait()
		go func(bookList []storege.BookList) {
			defer c.Done()
			for _, s := range bookList {
				url := strings.Replace(c.rule.BookDetailURL, CrawBookID, s.BookID, 1)
				bytes, err := utils.Get(url, nil)
				if err != nil {
					log.Println(err)
					continue
				}
				var isBool bool
				var book storege.BookDetail
				book.BookPlatform = c.rule.PlatformName
				book.BookID = s.BookID
				book.CatID = s.CatID
				book.BookName, isBool = GetSinglePatten(bytes, c.rule.BookNamePatten)
				if !isBool {
					continue
				}
				book.BookProcess, isBool = GetSinglePatten(bytes, c.rule.BookProcessPatten)
				if !isBool {
					continue
				}
				book.AuthorName, isBool = GetSinglePatten(bytes, c.rule.AuthorNamePatten)
				if !isBool {
					continue
				}
				book.BookCover, isBool = GetSinglePatten(bytes, c.rule.BookCoverPatten)
				if !isBool {
					continue
				}
				visitCount, isBool := GetSinglePatten(bytes, c.rule.VisitCountPatten)
				if !isBool {
					continue
				}
				book.VisitCount = utils.StringToInt(visitCount)
				score, isBool := GetSinglePatten(bytes, c.rule.ScorePatten)
				if !isBool {
					continue
				}
				book.Score = utils.StringToFloat64(score)
				tags, isBool := GetSinglePattenAll(bytes, c.rule.BookTagsPatten)
				if !isBool {
					continue
				}
				book.WordsCount, isBool = GetSinglePatten(bytes, c.rule.WordsCountPatten)
				if !isBool {
					continue
				}
				book.BookTags = strings.Join(tags, ",")
				book.Score = utils.StringToFloat64(score)
				book.BookDesc, isBool = GetBetweenPatten(bytes, c.rule.BookDescPatten)
				book.BookDetailHash = utils.Md5(c.rule.PlatformName + s.BookID)
				if err := storege.DB.Create(&book).Error; err != nil {
					log.Println(err)
					continue
				}
				storege.DB.Model(&storege.BookList{}).Where("id=?", s.ID).Update("book_state", 1)
				log.Println("saved done ", book.BookName)
			}
		}(bookList)
	}
}
