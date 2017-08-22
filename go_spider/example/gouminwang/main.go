package main

import (
	"fmt"
	"github.com/jackson198608/gotest/go_spider/core/pipeline"
	"github.com/jackson198608/gotest/go_spider/core/spider"
	"log"
	"os"
	"strconv"
)

// var dbAuth string = "dog123:dog123"

// var dbDsn string = "192.168.5.86:3306"

var dbAuth string = "root:goumin123"

var dbDsn string = "127.0.0.1:3306"

var dbName string = "big_data_mall"

var redisConn string = "192.168.86.68:6379"

var redisQueueName string = "threadInfo"

var saveDir string = "/data/gouminwang/"

var imgUrl string = "http://dev.img.goumintest.com"

var threaNum int = 1000

var logger *log.Logger
var logPath string = "/tmp/gouminwang_spider.log"

func load() {
	if checkFileIsExist(logPath) {
		fmt.Println("file exist")
		file, err := os.OpenFile(logPath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			os.Exit(1)
		}
		logger = log.New(file, "", log.LstdFlags)

	} else {
		fmt.Println("new file")
		file, err := os.Create(logPath)
		if err != nil {
			os.Exit(1)
		}
		logger = log.New(file, "", log.LstdFlags)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func main() {

	if len(os.Args) != 6 {
		fmt.Println("useage: datadir startUrl startUrlType logfile threadnum")
		os.Exit(1)
	}
	saveDir = os.Args[1]
	startUrl := os.Args[2]
	startUrlTag := os.Args[3]
	logPath = os.Args[4]
	threaNum, _ = strconv.Atoi(os.Args[5])

	load()
	logger.Println("[info]start ", startUrl)
	req := newRequest(startUrlTag, startUrl)
	spider.NewSpider(NewMyPageProcesser(), "TaskName").
		AddRequest(req).
		AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		SetThreadnum(uint(threaNum)).               // Crawl request by three Coroutines
		Run()
}
