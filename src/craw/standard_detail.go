package craw

import (
	"github.com/yiuked/go-novel/src/storege"
	"github.com/yiuked/go-novel/src/utils"
	"log"
	"strings"
	"time"
)

func (c *StandardCrawAction) GetBooksSummary(rule *BookCrawRule, processCheck bool) {
	c.init(rule)
	tryAgingCnt := 0
TRY:
	pageSize, page := 100, 0
	for {
		// 采用单协程分配，多协程处理，避免资源分配不重复
		var bookList []storege.BookList
		if processCheck {
			storege.DB().Where("book_state=1 AND book_process_state=0 AND book_platform=?", c.rule.PlatformName).Offset(page * pageSize).Limit(pageSize).Find(&bookList)
		} else {
			storege.DB().Where("book_state=0 AND book_platform=?", c.rule.PlatformName).Offset(page * pageSize).Limit(pageSize).Find(&bookList)
		}
		if len(bookList) <= 0 {
			if !processCheck {
				log.Println("get book summary task done,wait 3 seconds try aging ...")
				time.Sleep(3 * time.Second)
				tryAgingCnt++
				if tryAgingCnt >= 5 {
					log.Println("try aging finished")
					break
				}
				goto TRY
			} else {
				log.Println("update book process task done")
				break
			}
		}
		tryAgingCnt = 0
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
				updates := make(map[string]interface{})
				book.BookProcess, isBool = GetSinglePatten(bytes, c.rule.BookProcessPatten)
				if !isBool {
					continue
				}
				if !processCheck {
					book.BookPlatform = c.rule.PlatformName
					book.BookID = s.BookID
					book.CatID = s.CatID
					book.BookName, isBool = GetSinglePatten(bytes, c.rule.BookNamePatten)
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
					if err := storege.DB().Create(&book).Error; err != nil {
						log.Println(err)
						continue
					}

					updates["book_state"] = 1
				}
				if book.BookProcess == c.rule.BookProcessUpdate.EndName {
					updates["book_process_state"] = 1
				}
				if len(updates) > 0 {
					storege.DB().Model(&storege.BookList{}).Where("id=?", s.ID).Update("book_state", 1)
				}
				log.Println("saved done ", book.BookName)
			}
		}(bookList)
	}
}
