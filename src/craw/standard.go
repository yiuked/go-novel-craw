package craw

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type StandardCrawAction struct {
	rule           *BookCrawRule
	RequestGLimit  int64
	MaxPage        int
	UpdateInterval time.Duration
	currentG       *int64
	wgs            map[string]*sync.WaitGroup
}

func (c *StandardCrawAction) init(rule *BookCrawRule) {
	var needWgs bool
	if c.wgs == nil {
		c.wgs = make(map[string]*sync.WaitGroup)
		needWgs = true
	}
	if c.currentG == nil {
		c.currentG = new(int64)
	}

	for catID, _ := range rule.BookList.CatIDRelation {
		if needWgs {
			c.wgs[catID] = &sync.WaitGroup{}
		}
	}
	if c.rule == nil {
		c.rule = rule
	}
}

func (c *StandardCrawAction) Wait() {
	if *c.currentG >= c.RequestGLimit {
		log.Printf("request limit %d,current %d\n", c.RequestGLimit, *c.currentG)
		time.Sleep(3 * time.Second)
		c.Wait()
	}
	atomic.AddInt64(c.currentG, 1)
}

func (c *StandardCrawAction) Done() {
	atomic.AddInt64(c.currentG, -1)
}
