package main

import (
	"fmt"
	"github.com/jackson198608/gotest/go_spider/core/pipeline"
	"github.com/jackson198608/gotest/go_spider/core/spider"
	"log"
	"os"
	"strconv"
)

var saveDir string = "/data/aigouwang"

var threadNum int = 1000

var logPath string = "/var/spider.log"
var result *os.File

var logger *log.Logger

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

func loop(startUrl string, startUrlTag string, threadNum int) {
	logger.Println("do in oneloop taskNum", threadNum)
	x := make(chan int, threadNum)
	for i := 0; i < threadNum; i++ {
		go getRequestUrlTime(startUrl, startUrlTag, threadNum)
	}

	for i := 0; i < threadNum; i++ {
		<-x
	}
}

func getRequestUrlTime(startUrl string, startUrlTag string, threadNum int) {
	for {
		logger.Println("[info]start ", startUrl)
		// req := newRequest(startUrlTag, startUrl)
		// spider.NewSpider(NewMyPageProcesser(), "StressTest").
		// 	AddRequest(req).
		// 	AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		// 	SetThreadnum(uint(threadNum)).              // Crawl request by three Coroutines
		// 	Run()
	}
}

func main() {

	if len(os.Args) != 4 {
		fmt.Println("useage: startUrl startUrlType logfile threadnum ")
		os.Exit(1)
	}
	startUrl := os.Args[1]
	startUrlTag := os.Args[2]
	logPath = os.Args[3]
	threadNum, _ = strconv.Atoi(os.Args[4])

	load()
	// loop(startUrl, startUrlTag, threadNum)
	logger.Println("[info]start ", startUrl)
	req := newRequest(startUrlTag, startUrl)
	spider.NewSpider(NewMyPageProcesser(), "StressTest").
		AddRequest(req).
		AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		SetThreadnum(5).                            // Crawl request by three Coroutines
		Run()
}
