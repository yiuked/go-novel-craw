package craw

import (
	"github.com/yiuked/go-novel/src/storege"
	"github.com/yiuked/go-novel/src/utils"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (c *StandardCrawAction) GetBooksChapter(rule *BookCrawRule) {
	c.init(rule)
	if len(rule.BookList.CatIDRelation) <= 0 {
		return
	}
	for catID, _ := range rule.BookList.CatIDRelation {
		catID := catID
		go func() {
			var bookList []storege.BookList

			storege.DB.Where("book_state=0 AND updated_at<?", time.Now().Add(-c.UpdateInterval)).Limit(100).Find(&bookList)
			if bookList != nil {
				for _, s := range bookList {
					url := strings.Replace(c.rule.BookDetailURL, CRAW_BOOK_ID, s.BookID, 1)
					bytes, err := utils.Get(url, nil)
					if err != nil {
						continue
					}
					chapters, process := c.getBookChapter(bytes)

					var book storege.BookDetail
					book.BookPlatform = c.rule.PlatformName
					book.BookID = s.BookID
					book.CatID = catID
					book.BookName = c.getValueByPatten(bytes, c.rule.BookNamePatten)
					book.BookProcess = c.getValueByPatten(bytes, c.rule.BookProcessPatten)
					book.AuthorName = c.getValueByPatten(bytes, c.rule.AuthorNamePatten)
					book.BookCover = c.getValueByPatten(bytes, c.rule.BookCoverPatten)
					book.VisitCount = utils.StringToInt(c.getValueByPatten(bytes, c.rule.VisitCountPatten))
					book.Score = utils.StringToFloat64(c.getValueByPatten(bytes, c.rule.ScorePatten))
					book.BookDesc = c.getBookDesc(bytes)
					book.BookDetailHash = utils.Md5(c.rule.PlatformName + s.BookID)
					if len(book.BookProcess) <= 0 {
						book.BookProcess = process
					}
					var bookChapters []storege.BookChapter
					for id, name := range chapters {
						bookChapter := storege.BookChapter{
							BookPlatform:    book.BookPlatform,
							BookID:          book.BookID,
							BookChapterName: name,
							BookChapterID:   id,
							BookChapterHash: utils.Md5(book.BookPlatform + book.BookID + id),
						}
						bookChapters = append(bookChapters, bookChapter)
					}
					err = storege.DB.Transaction(func(tx *gorm.DB) error {
						if err := tx.Create(&book).Error; err != nil {
							return err
						}
						if err := tx.CreateInBatches(bookChapters, 100).Error; err != nil {
							return err
						}
						updates := make(map[string]interface{})
						if process == c.rule.BookProcessRelation.Finished.Name {
							updates["book_state"] = 1
						}
						updates["version"] = gorm.Expr("version+1")
						if err := tx.Model(&storege.BookList{}).Updates(updates).Error; err != nil {
							return err
						}
						return nil
					})
				}
			}
		}()
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
			chaptersMap[string(chapter[1])] = string(chapter[2])
		}
	}
	// 匹配连载状态
	var process string
	if len(c.rule.BookChapter.ProcessCheckKeywords) > 0 {
		process = c.rule.BookProcessRelation.Unfinished.Name
		for _, keyword := range c.rule.BookChapter.ProcessCheckKeywords {
			if strings.IndexAny(chaptersMap[isLastChapter], keyword) > 0 {
				process = c.rule.BookProcessRelation.Finished.Name
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
