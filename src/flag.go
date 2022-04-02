package main

import (
	"github.com/urfave/cli/v2"
	"github.com/yiuked/go-novel/src/craw"
	"github.com/yiuked/go-novel/src/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func commandParse() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "goroutine",
				Aliases:     []string{"g"},
				Value:       2,
				Usage:       "set goroutine craw limited",
				Destination: &goroutine,
			},
			&cli.StringFlag{
				Name:        "data_dir",
				Aliases:     []string{"c"},
				Usage:       "set（sqlite、cover、HTML）storage path,default current",
				Destination: &utils.DataDir,
			},
			&cli.StringFlag{
				Name:        "source",
				Aliases:     []string{"s"},
				Value:       "./yaml/tudu.yaml",
				Usage:       "set need craw source yaml file",
				Destination: &source,
			},
		},
		Commands: []*cli.Command{
			{
				Name: "bookid",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "page",
						Value: 1,
						Usage: "set max craw page",
					},
				},
				Usage:  "craw books ID to local",
				Action: StartBookIDCraw,
			},
			{
				Name:   "detail",
				Usage:  "craw books detail to local",
				Action: StartBookDetailCraw,
			},
			{
				Name:   "cover",
				Usage:  "craw books cover image to local",
				Action: StartBookCoverCraw,
			},
			{
				Name: "chapter",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "start",
						Value: 20,
						Usage: "开始时间",
					},
					&cli.StringFlag{
						Name:  "end",
						Usage: "结束时间",
					},
				},
				Usage:  "craw books chapter to local",
				Action: StartBookChapterCraw,
			},
			{
				Name:   "content",
				Usage:  "craw books content to local,save to HTML file",
				Action: StartBookContentCraw,
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	app.EnableBashCompletion = true
	if err != nil {
		log.Fatal(err)
	}
}

func CrawClient(c *cli.Context) (*craw.BookCraw, error) {
	bytes, err := ioutil.ReadFile(source)
	if err != nil {
		panic(err)
	}
	rule := craw.BookCrawRule{}
	if err := yaml.Unmarshal(bytes, &rule); err != nil {
		panic(err)
	}
	bookCraw := craw.NewBookCraw(&rule, &craw.StandardCrawAction{RequestGLimit: goroutine})
	return bookCraw, nil
}

func StartBookIDCraw(c *cli.Context) error {
	client, err := CrawClient(c)
	if err != nil {
		return err
	}
	client.StartBookIDCraw(c.Int("page"))
	waitSignal()
	return nil
}

func StartBookDetailCraw(c *cli.Context) error {
	client, err := CrawClient(c)
	if err != nil {
		return err
	}
	client.StartBookSummaryCraw()
	waitSignal()
	return nil
}

func StartBookCoverCraw(c *cli.Context) error {
	client, err := CrawClient(c)
	if err != nil {
		return err
	}
	client.StartBookCoverDownload()
	waitSignal()
	return nil
}

func StartBookChapterCraw(c *cli.Context) error {
	client, err := CrawClient(c)
	if err != nil {
		return err
	}
	client.StartBookChapterCraw()
	waitSignal()
	return nil
}

func StartBookContentCraw(c *cli.Context) error {
	client, err := CrawClient(c)
	if err != nil {
		return err
	}
	client.StartBookContentCraw()
	waitSignal()
	return nil
}

func waitSignal() {
	signal := make(chan os.Signal)
	read := <-signal
	log.Println("exit", read.String())
}
