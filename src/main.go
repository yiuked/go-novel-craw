package main

import (
	"flag"
	"github.com/yiuked/go-novel/src/craw"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	bytes, err := ioutil.ReadFile("./src/test.yaml")
	if err != nil {
		panic(err)
	}
	rule := craw.BookCrawRule{}
	if err := yaml.Unmarshal(bytes, &rule); err != nil {
		panic(err)
	}
	g := flag.Int64("g", 20, "groutine work limited number")

	bookCraw := craw.NewBookCraw(&rule, &craw.StandardCrawAction{RequestGLimit: *g})
	bookCraw.StartBookIDCraw()
	//bookCraw.StartBookSummaryCraw()

	signal := make(chan os.Signal)
	read := <-signal
	log.Println(read)
}
