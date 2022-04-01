package main

import (
	"flag"
	"github.com/yiuked/go-novel/src/craw"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

// init 初始化日志存储
func init() {
	if os.Getenv("DEBUG") == "false" {
		_, err := os.Stat("./logs")
		if err != nil {
			err := os.MkdirAll("./logs", 0777)
			if err != nil {
				panic(err)
				return
			}
		}
		logFie, err := os.OpenFile("./logs/info.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}

		log.SetOutput(logFie)
	}
}

func main() {
	bytes, err := ioutil.ReadFile("./yaml/tudu.yaml")
	if err != nil {
		panic(err)
	}
	rule := craw.BookCrawRule{}
	if err := yaml.Unmarshal(bytes, &rule); err != nil {
		panic(err)
	}
	g := flag.Int64("g", 1, "groutine work limited number")
	maxPage := flag.Int("p", 2, "max page,default 0,craw all page")

	bookCraw := craw.NewBookCraw(&rule, &craw.StandardCrawAction{RequestGLimit: *g, MaxPage: *maxPage})
	//bookCraw.StartBookIDCraw()
	//bookCraw.StartBookSummaryCraw()
	//bookCraw.StartBookChapterCraw()
	//bookCraw.StartBookCoverDownload()
	bookCraw.StartBookContentCraw()

	signal := make(chan os.Signal)
	read := <-signal
	log.Println("exit", read.String())
}
