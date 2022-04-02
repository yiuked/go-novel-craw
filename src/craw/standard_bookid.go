package craw

import (
	"github.com/yiuked/go-novel/src/storege"
	"github.com/yiuked/go-novel/src/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func (c *StandardCrawAction) GetBooksID(rule *BookCrawRule, maxPage int) {
	c.init(rule)
	if len(rule.BookList.Channels) <= 0 {
		return
	}
	// 此处需要注意，全局协程限制只会限制请求，这里如果分类过多，可能会导致协程过多
	for _, channel := range rule.BookList.Channels {
		if len(channel.CatID) <= 0 {
			continue
		}
		for _, catID := range channel.CatID {
			go func(cat string) {
				c.crawCat(channel.ChannelID, cat, maxPage)
			}(catID)
		}
	}
}

// crawCat 采集分类下所有书籍ID
func (c *StandardCrawAction) crawCat(channelID, catID string, maxPage int) {
	baseURL := strings.Replace(c.rule.BookList.BookListURL, CrawChannelId, channelID, -1)
	baseURL = strings.Replace(baseURL, CrawCatID, catID, -1)
	url := strings.Replace(baseURL, CrawPage, "1", 1)
	log.Println("start craw cat:", url)
	listHtml, err := c.getCatPage(channelID, catID, url)
	if err != nil {
		log.Println("get total page err:", err)
		return
	}

	// 总页
	totalPageReg := regexp.MustCompile(c.rule.BookList.TotalPagePatten)
	totalPageResult := totalPageReg.FindSubmatch(listHtml)
	if len(totalPageResult) < 2 {
		log.Println("get total page err")
		return
	}
	totalPage := utils.StringToInt(string(totalPageResult[1]))
	if maxPage > 0 {
		totalPage = maxPage
	}

	log.Printf("cat[%s] total page %d\n", catID, totalPage)

	for i := 2; i <= totalPage; i++ {
		c.Wait()
		urlTemp := strings.Replace(baseURL, CrawPage, strconv.Itoa(i), 1)
		go func() {
			defer func() {
				c.Done()
				log.Println("catID--", catID)
			}()
			_, err := c.getCatPage(channelID, catID, urlTemp)
			if err != nil {
				log.Println("request err:", err)
			}
		}()
	}
	log.Println("craw book id finished ", catID)
	// 存储到数据库
}

func (c *StandardCrawAction) getCatPage(channelID, catID, url string) ([]byte, error) {
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
				ChannelID:    channelID,
			}
			storege.DB().Create(&book)
		}
	}
	return listHtml, nil
}
