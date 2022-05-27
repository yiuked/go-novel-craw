package main

import (
	"github.com/urfave/cli/v2"
	"github.com/yiuked/go-novel/src/api"
	"github.com/yiuked/go-novel/src/craw"
	"github.com/yiuked/go-novel/src/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

func commandParse() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "goroutine",
				Aliases:     []string{"g"},
				Value:       20,
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
				Value:       "./yaml/default.yaml",
				Usage:       "set need craw source yaml file",
				Destination: &source,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "api",
				Usage:  "start web api service",
				Action: StartWebApi,
			},
			{
				Name: "bookid",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "page",
						Value: 5,
						Usage: "set max craw page",
					},
				},
				Usage:  "craw books ID to local",
				Action: StartBookIDCraw,
			},
			{
				Name:      "detail",
				UsageText: "craw books detail to local",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "update",
						Usage: "if true, check the existence record and update the process status",
					},
				},
				Action: StartBookDetailCraw,
			},
			{
				Name:   "cover",
				Usage:  "craw books cover image to local",
				Action: StartBookCoverCraw,
			},
			{
				Name:      "chapter",
				UsageText: "If the start or end parameter is set, the chapter data is updated.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "start",
						Usage: "Seconds, data less than the current time minus this time will re-read the chapter.",
					},
					&cli.StringFlag{
						Name:  "end",
						Usage: "Seconds, data greater than the current time minus this time will re-read the chapter.",
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

func StartWebApi(c *cli.Context) error {
	api.Routes()
	return nil
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
	processCheck := c.Bool("update")
	client.StartBookSummaryCraw(processCheck)
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
	var start, end time.Duration
	if len(c.String("start")) > 0 {
		start, err = utils.String2TimeDuration(c.String("start"))
		if err != nil {
			return err
		}
	}
	if len(c.String("end")) > 0 {
		end, err = utils.String2TimeDuration(c.String("end"))
		if err != nil {
			return err
		}
	}

	client.StartBookChapterCraw(start, end)
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
