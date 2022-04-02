package main

import (
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

// 全局参数
var (
	source    string
	goroutine int64
)

func main() {
	commandParse()
}
