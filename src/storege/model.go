package storege

import "gorm.io/gorm"

// BookList 书籍列表
type BookList struct {
	gorm.Model
	BookPlatform     string `gorm:"type:varchar(32) not null default '';"`                            // 平台名称
	CatID            string `gorm:"type:varchar(32) not null default ''"`                             // 分类ID
	BookID           string `gorm:"type:varchar(32) not null default '';index:idx_list"`              // 书籍ID
	ChannelID        string `gorm:"type:varchar(32) not null default ''"`                             // 频道ID
	BookListHash     string `gorm:"type:varchar(32) not null default '';uniqueIndex:idx_list_unique"` // 唯一标识
	BookState        int    `gorm:"type:tinyint(1) not null default 0"`                               // 书籍采集状态(0未采集详情,1已采集详情)
	BookProcessState int    `gorm:"type:tinyint(1) not null default 0"`                               // 书籍连载状态(0连载中,1已完结书籍)
	Version          int    `gorm:"type:int(11) not null default 0"`                                  // 更新版本
}

// BookDetail 书本详情
type BookDetail struct {
	gorm.Model
	BookPlatform      string  `gorm:"type:varchar(32) not null default ''"`                               // 平台名称
	BookID            string  `gorm:"type:varchar(32) not null default '';index:idx_detail"`              // 书籍ID
	BookTags          string  `gorm:"type:varchar(128) not null default ''"`                              // 标签，以","号分割
	BookDetailHash    string  `gorm:"type:varchar(32) not null default '';uniqueIndex:idx_detail_unique"` // 唯一标识
	CatID             string  `gorm:"type:varchar(32) not null default ''"`                               // 分类ID
	ChannelID         string  `gorm:"type:varchar(32) not null default ''"`                               // 频道ID
	BookProcess       string  `gorm:"type:varchar(32) not null default ''"`                               // 连载状态
	BookName          string  `gorm:"type:varchar(128) not null default ''"`                              // 书名
	AuthorName        string  `gorm:"type:varchar(64) not null default ''"`                               // 作者
	BookCover         string  `gorm:"type:varchar(256) not null default ''"`                              // 封面
	BookCoverDownload int     `gorm:"type:tinyint(1) not null default 0"`                                 // 封面是否已下载，0未下载，1已下载
	BookCoverLocal    string  `gorm:"type:varchar(256) not null default ''"`                              // 封面本地路径
	Score             float64 `gorm:"type:decimal(5,2) not null"`                                         // 评分
	VisitCount        int     `gorm:"type:int(10) default 0 not null"`                                    // 浏览次数
	WordsCount        string  `gorm:"type:varchar(64) default 0 not null"`                                // 字数
	BookDesc          string  `gorm:"type:text not null default ''"`                                      // 书籍描述
}

// BookChapter
//书本详情
type BookChapter struct {
	gorm.Model
	BookPlatform     string `gorm:"type:varchar(32) not null default ''"`                                //平台名称
	BookID           string `gorm:"type:varchar(32) not null default '';index:idx_chapter_book_id"`      //书籍ID
	CatID            string `gorm:"type:varchar(32) not null default ''"`                                // 分类ID
	BookChapterID    string `gorm:"type:varchar(32) not null default '';uniqueIndex:idx_chapter_id"`     //书籍章节ID
	BookChapterHash  string `gorm:"type:varchar(32) not null default '';uniqueIndex:idx_chapter_unique"` //唯一标识
	BookChapterName  string `gorm:"type:varchar(128) not null default ''"`                               //章节名称
	BookChapterState int    `gorm:"type:tinyint(1) not null default 2"`                                  //是否已处理:1=已处理,2=未处理
	BookChapterLocal string `gorm:"type:varchar(256) not null default ''"`                               // 本地存储路径
}
