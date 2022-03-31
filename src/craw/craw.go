package craw

const (
	CRAW_CATID   = "{catId}"
	CRAW_PAGE    = "{page}"
	CRAW_BOOK_ID = "{bookId}"
)

type BookCraw struct {
	Rule   *BookCrawRule
	action BookCrawAction
}

// Book 清洗后的数据格式
type Book struct {
	BookID      string            // 原平台书籍ID
	CatID       string            // 原平台ID
	BookProcess string            // 原平台状态
	BookName    string            // 原平台书本名称
	AuthorName  string            // 原平台作者
	BookCover   string            // 书籍封面地址
	Score       string            // 书籍评分
	VisitCount  string            // 阅读量
	BookDesc    string            // 书籍简介
	BookChapter map[string]string // key 章节ID，value章节名称
}

// BookCrawRule 采集规则
type BookCrawRule struct {
	PlatformName string `yaml:"platform_name"` // 平台名称
	BookList     struct {
		BookListURL     string            `yaml:"book_list_url"`     // 书本标签列表
		CatIDRelation   map[string]string `yaml:"cat_id_relation"`   // 分类与平台的关联：`key` 为平台分类ID，`value`为所采集平台的分类ID
		PagePatten      string            `yaml:"page_patten"`       // 获取当前页规则
		TotalPagePatten string            `yaml:"total_page_patten"` // 获取总页数规则
		BookIDPatten    string            `yaml:"book_id_patten"`    // 获取采集书本的ID正则
	} `yaml:"book_list"`
	BookChapter struct {
		BookChapterURL       string              `yaml:"book_chapter_url"`    // 书本章节地址
		BookChapterPatten    string              `yaml:"book_chapter_patten"` // 获取章节ID\名称规则
		Extend               []map[string]string `yaml:"extend"`
		ProcessCheckKeywords []string            `yaml:"process_check_keywords"`
	} `yaml:"book_chapter"`
	BookProcessRelation struct {
		Finished struct {
			Name  string
			Value string
		}
		Unfinished struct {
			Name  string
			Value string
		}
	} `yaml:"book_process_relation"`
	BookDetailURL     string `yaml:"book_detail_url"`     // 书本地址
	BookNamePatten    string `yaml:"book_name_patten"`    // 获取书本名称规则
	AuthorNamePatten  string `yaml:"author_name_patten"`  // 获取作者规则
	BookCoverHost     string `yaml:"book_cover_host"`     // 封面图片主机地址（如果是相对路径的情况需要填写）
	BookCoverPatten   string `yaml:"book_cover_patten"`   // 获取封面图片规则
	BookProcessPatten string `yaml:"book_process_patten"` // 获取小说连载进度规则
	ScorePatten       string `yaml:"score_patten"`        // 获取小说评价规则
	VisitCountPatten  string `yaml:"visit_count_patten"`  // 获取小说阅读量规则
	BookDescPatten    struct {
		Start  string              `yaml:"start"` // 开始标签
		End    string              `yaml:"end"`   // 结束标签
		Extend []map[string]string `yaml:"extend"`
	} `yaml:"book_desc_patten"` // 获取小说描述规则
	BookContentURL    string `yaml:"book_content_url"` // 书籍内容页地址
	BookContentPatten struct {
		Start string `yaml:"start"` // 开始标签
		End   string `yaml:"end"`   // 结束标签
	} `yaml:"book_content_patten"` // 获取小说描述规则
}

type BookCrawAction interface {
	// GetBooksID 获取指定分类下的所有书籍ID
	GetBooksID(rule *BookCrawRule)
	// GetBooksSummary 获取书本详情（基本信息+章节）
	GetBooksSummary(rule *BookCrawRule)
}

func NewBookCraw(rule *BookCrawRule, action BookCrawAction) *BookCraw {
	return &BookCraw{Rule: rule, action: action}
}

// StartBookIDCraw 开始采集书籍ID
func (c *BookCraw) StartBookIDCraw() {
	c.action.GetBooksID(c.Rule)
}

// StartBookSummaryCraw 开始采集书籍介绍+章节信息
func (c *BookCraw) StartBookSummaryCraw() {
	c.action.GetBooksSummary(c.Rule)
}
