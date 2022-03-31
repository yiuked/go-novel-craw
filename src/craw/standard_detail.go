package craw

import (
	"github.com/yiuked/go-novel/src/storege"
	"github.com/yiuked/go-novel/src/utils"
	"log"
	"regexp"
	"strings"
)

func (c *StandardCrawAction) GetBooksSummary(rule *BookCrawRule) {
	c.init(rule)
	if len(rule.BookList.CatIDRelation) <= 0 {
		return
	}
	for catID, _ := range rule.BookList.CatIDRelation {
		catID := catID
		go func() {
			var bookList []storege.BookList
			storege.DB.Where("book_state=0").Limit(100).Find(&bookList)
			if bookList != nil {
				for _, s := range bookList {
					url := strings.Replace(c.rule.BookDetailURL, CRAW_BOOK_ID, s.BookID, 1)
					bytes, err := utils.Get(url, nil)
					if err != nil {
						continue
					}

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
					if err := storege.DB.Create(&book).Error; err != nil {
						log.Println(err)
					}
				}
			}
		}()
	}
}

func (c *StandardCrawAction) getValueByPatten(listHtml []byte, patten string) string {
	if len(patten) <= 0 {
		return ""
	}
	reg := regexp.MustCompile(patten)
	regResult := reg.FindSubmatch(listHtml)
	if len(regResult) > 1 {
		return string(regResult[1])
	}
	return ""
}

func (c *StandardCrawAction) getBookDesc(listHtml []byte) string {
	startHtml := strings.SplitAfterN(string(listHtml), c.rule.BookDescPatten.Start, 2)
	if len(startHtml) < 1 {
		return ""
	}
	endHtml := strings.SplitN(startHtml[1], c.rule.BookDescPatten.End, 2)
	if len(endHtml) <= 0 {
		return ""
	}

	html := strings.TrimSpace(endHtml[0])
	if len(c.rule.BookDescPatten.Extend) > 0 {
		for _, ext := range c.rule.BookDescPatten.Extend {
			reg := regexp.MustCompile(ext["patten"])
			html = reg.ReplaceAllString(html, ext["replace"])
		}
	}
	return html
}
