package craw

import (
	"github.com/yiuked/go-novel/src/storege"
	"github.com/yiuked/go-novel/src/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func (c *StandardCrawAction) GetBooksID(rule *BookCrawRule) {
	c.init(rule)
	if len(rule.BookList.CatIDRelation) <= 0 {
		return
	}
	for catID, _ := range rule.BookList.CatIDRelation {
		c.Wait()
		go func(cat string) {
			defer func() {
				c.Done()
			}()
			c.crawCat(cat)
		}(catID)
	}
}

// crawCat 采集分类下所有书籍ID
func (c *StandardCrawAction) crawCat(catID string) {
	urlCat := strings.Replace(c.rule.BookList.BookListURL, CRAW_CATID, catID, -1)
	url := strings.Replace(urlCat, CRAW_PAGE, "1", 1)
	log.Println("start craw cat:", url)
	listHtml, err := c.getCatPage(catID, url)
	if err != nil {
		log.Println("get total page err:", err)
		return
	}

	// 总页
	totalPageReg := regexp.MustCompile(c.rule.BookList.TotalPagePatten)
	totalPageResult := totalPageReg.FindSubmatch(listHtml)
	totalPage := utils.StringToInt(string(totalPageResult[1]))
	if c.MaxPage > 0 {
		totalPage = c.MaxPage
	}

	log.Printf("cat[%s] total page %d\n", catID, totalPage)
	c.wgs[catID].Add(totalPage - 1)
	for i := 2; i <= totalPage; i++ {
		c.Wait()

		urlTemp := strings.Replace(urlCat, CRAW_PAGE, strconv.Itoa(i), 1)
		go func() {
			defer func() {
				c.Done()
				log.Println("catID--", catID)
				c.wgs[catID].Done()
			}()

			_, err := c.getCatPage(catID, urlTemp)
			if err != nil {
				log.Println("request err:", err)
			}
		}()
	}
	c.wgs[catID].Wait()
	log.Println("craw book id finished ", catID)
	// 存储到数据库
}

func (c *StandardCrawAction) getCatPage(catID, url string) ([]byte, error) {
	log.Println("Get request:", url)
	listHtml, err := utils.Get(url, nil)
	if err != nil {
		return nil, err
	}

	bookIDReg := regexp.MustCompile(c.rule.BookList.BookIDPatten)
	bookIDResult := bookIDReg.FindAllSubmatch(listHtml, -1)
	if len(bookIDResult) > 0 {
		for _, bookID := range bookIDResult {
			book := storege.BookList{
				BookPlatform: c.rule.PlatformName,
				CatID:        catID,
				BookID:       string(bookID[1]),
				BookListHash: utils.Md5(c.rule.PlatformName + string(bookID[1])),
			}
			storege.DB.Create(&book)
		}
	}
	return listHtml, nil
}
