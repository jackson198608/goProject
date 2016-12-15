package main

import (
	"fmt"
	"github.com/hu17889/go_spider/core/pipeline"
	"github.com/hu17889/go_spider/core/spider"
	"log"
	"os"
	"strconv"
)

var Type int = 1

var City string = "北京"

var saveDir string = "/data/aigouwang"

var threaNum int = 1000

var logger *log.Logger
var logPath string = "/var/spider.log"

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

	// Spider input:
	//  PageProcesser ;
	//  Task name used in Pipeline for record;

	if len(os.Args) != 8 {
		fmt.Println("useage: datadir startUrl startUrlType logfile threadnum type city")
		os.Exit(1)
	}
	saveDir = os.Args[1]
	startUrl := os.Args[2]
	startUrlTag := os.Args[3]
	logPath = os.Args[4]
	threaNum, _ = strconv.Atoi(os.Args[5])
	Type, _ = strconv.Atoi(os.Args[6])
	City = os.Args[7]

	load()
	logger.Println("[info]start ", startUrl)
	req := newRequest(startUrlTag, startUrl)
	spider.NewSpider(NewMyPageProcesser(), "TaskName").
		AddRequest(req).
		AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		SetThreadnum(uint(threaNum)).               // Crawl request by three Coroutines
		Run()
}
