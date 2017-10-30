package main

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"os"
)

var c Config = Config{
	"192.168.86.193:3307",
	"new_dog123",
	"dog123:dog123",
	"192.168.86.193:3307",
	"process",
	"dog123:dog123",
	10,
	"127.0.0.1:6379",
	"getHtml",
	"/data/thread/",
	"/tmp/create_thread.log",
	1,
	"300", "500", "http://m.goumin.com/bbs/", "zhidao.goumin.com"}

var offset = 1000

func Init(args []string) {

	loadConfig(args)
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
	logger.Error(logger.DEBUG)

}

func main() {
	Init(os.Args)
	fmt.Println(len(os.Args))
	jobType := os.Args[1]
	cat := "all"
	if len(os.Args) < 3 {
		logger.Error("args error")
		return
	}
	if len(os.Args) == 4 {
		cat = os.Args[3]
	}
	switch jobType {
	case "thread": //thread
		logger.Info("in the ask html get ")
		createHtmlByUrl(jobType)
	case "ask": //thread
		logger.Info("in the ask html get ")
		createHtmlByUrl(jobType)
	case "asksave": //create html url
		logger.Info("in the html url save ")
		saveHtmlUrl(jobType, cat)
	case "threadsave": //create html url
		logger.Info("in the html url save ")
		saveHtmlUrl(jobType, cat)
	default:

	}
}
