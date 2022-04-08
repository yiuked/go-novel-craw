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

func (c *StandardCrawAction) GetBooksCover(rule *BookCrawRule) {
	c.init(rule)
	tryAgingCnt := 0
TRY:
	pageSize, page := 100, 0
	for {
		// 采用单协程分配，多协程处理，避免资源分配不重复
		var books []storege.BookDetail
		storege.DB().Where("book_cover_download=0 AND book_platform=?", c.rule.PlatformName).Offset(page * pageSize).Limit(pageSize).Find(&books)
		if len(books) <= 0 {
			log.Println("get book cover task done,wait 3 seconds try aging ...")
			tryAgingCnt++
			if tryAgingCnt >= 5 {
				log.Println("try aging finished")
				break
			}
			goto TRY
		}
		page++
		c.Wait()
		go func(books []storege.BookDetail) {
			defer c.Done()
			for _, book := range books {
				abs, rel := c.coverFileName(book)
				_, err := os.Lstat(abs)
				if err == nil {
					book.BookCoverDownload = 1
					book.BookCoverLocal = rel
					storege.DB().Save(&book)
					log.Println(book.BookCover, "download")
					continue
				}

				bytes, err := utils.Get(book.BookCover, nil)
				if err != nil {
					log.Println(err)
					continue
				}
				if !isImg(bytes) {
					continue
				}
				err = os.WriteFile(abs, bytes, 0644)
				if err != nil {
					log.Println(err)
					continue
				}
				book.BookCoverDownload = 1
				book.BookCoverLocal = rel
				storege.DB().Save(&book)
				log.Println(book.BookCover, "download")
			}
		}(books)
	}
}

func isImg(img []byte) bool {
	return true
}

func (c *StandardCrawAction) coverFileName(book storege.BookDetail) (absolute, relative string) {
	// 获取相对路径
	if len(c.rule.BookCoverLocalPath) > 0 {
		relative = strings.ReplaceAll(c.rule.BookCoverLocalPath, CrawCatID, book.CatID)
		relative = strings.ReplaceAll(relative, CrawBookID, book.BookID)
	} else {
		relative = fmt.Sprintf("/img/%s", book.CatID)
	}
	relative = fmt.Sprintf("%s/%s%s", relative, book.BookDetailHash, path.Ext(book.BookCover))

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
