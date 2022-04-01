package craw

import (
	"regexp"
	"strings"
)

const (
	CrawCatID     = "{catId}"
	CrawPage      = "{page}"
	CrawBookID    = "{bookId}"
	CrawChapterID = "{chapterId}"
)

type BookCraw struct {
	Rule   *BookCrawRule
	action BookCrawAction
}

// SinglePatten 单项匹配
type SinglePatten struct {
	Require bool
	Patten  string
	Replace []map[string]string `yaml:"replace"`
}

type NewPagePatten struct {
	NewPageURLPatten SinglePatten  `yaml:"new_page_url_patten"`
	ContentPatten    BetweenPatten `yaml:"content_patten"`
}

// BetweenPatten 区间匹配
type BetweenPatten struct {
	Require bool
	Start   string              `yaml:"start"` // 开始标签
	End     string              `yaml:"end"`   // 结束标签
	Replace []map[string]string `yaml:"replace"`
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
		BookChapterURL       string              `yaml:"book_chapter_url"`        // 书本章节地址
		BookChapterLocalPath string              `yaml:"book_chapter_local_path"` // 章节存储路径，支持`{catId}`、`{bookId}`、`{chapterId}`
		BookChapterPatten    string              `yaml:"book_chapter_patten"`     // 获取章节ID\名称规则
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
	BookDetailURL      string        `yaml:"book_detail_url"`       // 书本地址
	BookNamePatten     SinglePatten  `yaml:"book_name_patten"`      // 获取书本名称规则
	AuthorNamePatten   SinglePatten  `yaml:"author_name_patten"`    // 获取作者规则
	BookCoverHost      string        `yaml:"book_cover_host"`       // 封面图片主机地址（如果是相对路径的情况需要填写）
	BookCoverLocalPath string        `yaml:"book_cover_local_path"` // 封面本地存储路径，支持`{catId}`、`{bookId}`
	BookCoverPatten    SinglePatten  `yaml:"book_cover_patten"`     // 获取封面图片规则
	BookProcessPatten  SinglePatten  `yaml:"book_process_patten"`   // 获取小说连载进度规则
	ScorePatten        SinglePatten  `yaml:"score_patten"`          // 获取小说评价规则
	BookTagsPatten     SinglePatten  `yaml:"book_tags_patten"`      // 标签匹配
	WordsCountPatten   SinglePatten  `yaml:"words_count_patten"`    // 字数匹配
	VisitCountPatten   SinglePatten  `yaml:"visit_count_patten"`    // 获取小说阅读量规则
	BookDescPatten     BetweenPatten `yaml:"book_desc_patten"`      // 获取小说描述规则
	BookContentURL     string        `yaml:"book_content_url"`      // 书籍内容页地址
	BookContentPatten  struct {
		CurrentPage BetweenPatten `yaml:"current_page"` // 在当前页采集
		NewPage     NewPagePatten `yaml:"new_page"`     // 如果需要二级采集页，设置了new_page则优先使用new_page，忽略current_page
	} `yaml:"book_content_patten"`
}

type BookCrawAction interface {
	// GetBooksID 获取指定分类下的所有书籍ID
	GetBooksID(rule *BookCrawRule)
	// GetBooksSummary 获取书本详情（基本信息）
	GetBooksSummary(rule *BookCrawRule)
	// GetBooksChapter 获取书本章节
	GetBooksChapter(rule *BookCrawRule)
	// GetBooksCover 下载封面图片
	GetBooksCover(rule *BookCrawRule)
	// GetBooksContent 获取书本内容
	GetBooksContent(rule *BookCrawRule)
}

func NewBookCraw(rule *BookCrawRule, action BookCrawAction) *BookCraw {
	return &BookCraw{Rule: rule, action: action}
}

// StartBookIDCraw 开始采集书籍ID
func (c *BookCraw) StartBookIDCraw() {
	c.action.GetBooksID(c.Rule)
}

// StartBookSummaryCraw 开始采集书籍介绍
func (c *BookCraw) StartBookSummaryCraw() {
	c.action.GetBooksSummary(c.Rule)
}

// StartBookChapterCraw 开始采集书籍章节信息
func (c *BookCraw) StartBookChapterCraw() {
	c.action.GetBooksChapter(c.Rule)
}

// StartBookCoverDownload 封面图片下载
func (c *BookCraw) StartBookCoverDownload() {
	c.action.GetBooksCover(c.Rule)
}

// StartBookContentCraw 开始采集书籍内容
func (c *BookCraw) StartBookContentCraw() {
	c.action.GetBooksContent(c.Rule)
}

func GetSinglePatten(listHtml []byte, patten SinglePatten) (string, bool) {
	if len(patten.Patten) <= 0 {
		return "", patten.Require == false
	}
	reg := regexp.MustCompile(patten.Patten)
	regResult := reg.FindSubmatch(listHtml)
	if len(regResult) > 1 {
		html := strings.TrimSpace(string(regResult[1]))
		if len(patten.Replace) > 0 {
			for _, ext := range patten.Replace {
				reg := regexp.MustCompile(ext["patten"])
				html = reg.ReplaceAllString(html, ext["replace"])
			}
		}
		return html, true
	}
	return "", patten.Require == false
}

func GetSinglePattenAll(listHtml []byte, patten SinglePatten) ([]string, bool) {
	var find []string
	if len(patten.Patten) <= 0 {
		return find, patten.Require == false
	}
	reg := regexp.MustCompile(patten.Patten)
	results := reg.FindAllSubmatch(listHtml, -1)
	if len(results) > 0 {
		for _, regResult := range results {
			html := strings.TrimSpace(string(regResult[1]))
			if len(patten.Replace) > 0 {
				for _, ext := range patten.Replace {
					reg := regexp.MustCompile(ext["patten"])
					html = reg.ReplaceAllString(html, ext["replace"])
				}
			}
			find = append(find, html)
		}

		return find, true
	}
	return find, patten.Require == false
}

func GetBetweenPatten(listHtml []byte, patten BetweenPatten) (string, bool) {
	startHtml := strings.SplitAfterN(string(listHtml), patten.Start, 2)
	if len(startHtml) < 2 {
		return "", patten.Require == false
	}
	endHtml := strings.SplitN(startHtml[1], patten.End, 2)
	if len(endHtml) < 1 {
		return "", patten.Require == false
	}

	html := strings.TrimSpace(endHtml[0])
	if len(patten.Replace) > 0 {
		for _, ext := range patten.Replace {
			reg := regexp.MustCompile(ext["patten"])
			html = reg.ReplaceAllString(html, ext["replace"])
		}
	}
	if len(html) <= 0 {
		return "", patten.Require == false
	}
	return html, true
}
