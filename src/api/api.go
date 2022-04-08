package api

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yiuked/go-novel/src/storege"
	"strconv"
)

// BookDetail 书本详情
type BookDetail struct {
	BookID      string  `json:"book_id"`      // 书籍ID
	BookTags    string  `json:"book_tags"`    // 标签，以","号分割
	CatID       string  `json:"cat_id"`       // 分类ID
	ChannelID   string  `json:"channel_id"`   // 频道ID
	BookProcess string  `json:"book_process"` // 连载状态
	BookName    string  `json:"book_name"`    // 书名
	AuthorName  string  `json:"author_name"`  // 作者
	BookCover   string  `json:"book_cover"`   // 封面
	Score       float64 `json:"score"`        // 评分
	VisitCount  int     `json:"visit_count"`  // 浏览次数
	WordsCount  string  `json:"words_count"`  // 字数
	BookDesc    string  `json:"book_desc"`    // 书籍描述
}

// BookChapter 书本章节信息
type BookChapter struct {
	BookID          string `json:"book_id"`           // 书籍ID
	CatID           string `json:"cat_id"`            // 分类ID
	BookChapterID   string `json:"book_chapter_id"`   // 书籍章节ID
	BookChapterName string `json:"book_chapter_name"` // 章节名称
	BookChapterURL  string `json:"book_chapter_url"`  // 本地存储路径
}

func Routes() {
	r := gin.Default()

	r.GET("/api/books", func(c *gin.Context) {
		var books []storege.BookDetail
		storege.DB().Limit(20).Find(&books)
		var details []BookDetail

		for _, book := range books {
			details = append(details, BookDetail{
				BookID:      book.BookID,
				BookTags:    book.BookTags,
				CatID:       book.CatID,
				ChannelID:   book.ChannelID,
				BookProcess: book.ChannelID,
				BookName:    book.BookName,
				AuthorName:  book.AuthorName,
				BookCover:   oss(book.BookCover),
				Score:       book.Score,
				VisitCount:  book.VisitCount,
				WordsCount:  book.WordsCount,
				BookDesc:    book.BookDesc,
			})
		}

		c.JSON(200, details)
	})
	r.GET("/api/chapters", func(c *gin.Context) {
		var chapters []storege.BookChapter
		storege.DB().Where("book_id=?", c.Query("book_id")).Order("book_chapter_id ASC").Limit(20).Find(&chapters)

		var bookChapters []BookChapter
		for _, chapter := range chapters {
			bookChapters = append(bookChapters, BookChapter{
				BookID:          chapter.BookID,
				CatID:           chapter.CatID,
				BookChapterID:   chapter.BookChapterID,
				BookChapterName: chapter.BookChapterName,
				BookChapterURL:  oss(chapter.BookChapterLocal),
			})
		}
		c.JSON(200, bookChapters)
	})

	ginPort := flag.Int("port", 7040,
		fmt.Sprintf("get ginServerPort from cmd,default %d as port", 7040))
	flag.Parse()
	err := r.Run(":" + strconv.Itoa(*ginPort))
	if err != nil {
		panic(err)
	}
}

func oss(path string) string {
	return "http://192.168.3.135:7080/" + path
}
