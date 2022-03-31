package storege

import "gorm.io/gorm"

// BookList 书籍列表
type BookList struct {
	gorm.Model
	BookPlatform string `gorm:"type:varchar(32) not null;"`                            // 平台名称
	CatID        string `gorm:"type:varchar(32) not null"`                             // 分类ID
	BookID       string `gorm:"type:varchar(32) not null;index:idx_list"`              // 书籍ID
	BookListHash string `gorm:"type:varchar(32) not null;uniqueIndex:idx_list_unique"` // 唯一标识
	BookState    int    `gorm:"type:tinyint(1) not null default 0"`                    // 书籍状态
	Version      int    `gorm:"type:int(1) not null default 0"`                        // 更新版本
}

// BookDetail 书本详情
type BookDetail struct {
	gorm.Model
	BookPlatform   string  `gorm:"type:varchar(32) not null"`                               // 平台名称
	BookID         string  `gorm:"type:varchar(32) not null;index:idx_detail"`              // 书籍ID
	BookDetailHash string  `gorm:"type:varchar(32) not null;uniqueIndex:idx_detail_unique"` // 唯一标识
	CatID          string  `gorm:"type:varchar(32) not null"`                               // 分类ID
	BookProcess    string  `gorm:"type:varchar(32) not null"`                               // 连载状态
	BookName       string  `gorm:"type:varchar(128) not null"`                              // 书名
	AuthorName     string  `gorm:"type:varchar(64) not null"`                               // 作者
	BookCover      string  `gorm:"type:varchar(256) not null"`                              // 封面
	Score          float64 `gorm:"type:decimal(5,2) not null"`                              // 评分
	VisitCount     int     `gorm:"type:int(10) default 0 not null"`                         // 浏览次数
	BookDesc       string  `gorm:"type:text not null"`                                      // 书籍描述
}

// BookChapter
//书本详情
type BookChapter struct {
	gorm.Model
	BookPlatform     string `gorm:"type:varchar(32) not null"`                                //平台名称
	BookID           string `gorm:"type:varchar(32) not null;index:idx_chapter_book_id"`      //书籍ID
	BookChapterID    string `gorm:"type:varchar(32) not null;uniqueIndex:idx_chapter_id"`     //书籍章节ID
	BookChapterHash  string `gorm:"type:varchar(32) not null;uniqueIndex:idx_chapter_unique"` //唯一标识
	BookChapterName  string `gorm:"type:varchar(128) not null"`                               //章节名称
	BookChapterState int    `gorm:"type:tinyint(1) not null default 2"`                       //是否已处理:1=已处理,2=未处理
}
