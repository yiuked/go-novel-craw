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
	UpdateInterval time.Duration
	currentG       *int64
	mu             sync.Mutex
}

func (c *StandardCrawAction) init(rule *BookCrawRule) {
	if c.currentG == nil {
		c.currentG = new(int64)
	}
	if c.rule == nil {
		c.rule = rule
	}
}

func (c *StandardCrawAction) Wait() {
	c.mu.Lock()
	if *c.currentG < c.RequestGLimit {
		atomic.AddInt64(c.currentG, 1)
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()
	log.Printf("request limit %d,current %d\n", c.RequestGLimit, *c.currentG)
	time.Sleep(3 * time.Second)
	c.Wait()

}

func (c *StandardCrawAction) Done() {
	atomic.AddInt64(c.currentG, -1)
}
