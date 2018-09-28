package main

import (
	"fmt"
	"github.com/jackson198608/goProject/go_spider/core/pipeline"
	"github.com/jackson198608/goProject/go_spider/core/spider"
	"log"
	"os"
)

// var dbAuth string = "dog123:dog123"

// var dbDsn string = "192.168.5.86:3306"

// var dbAuth string = "dog123:dog123"

// var dbDsn string = "192.168.86.194:3307"

// var dbName string = "big_data_mall"

// var redisConn string = "192.168.86.80:6380,192.168.86.80:6381,192.168.86.81:6380,192.168.86.81:6381,192.168.86.82:6380,192.168.86.82:6381"

// var redisQueueName string = "threadInfo"

// var saveDir string = "/data/gouminwang/"

// var imgUrl string = "http://dev.img.goumintest.com"

// var threadNum int = 1000
// var logPath string = "/tmp/gouminwang_spider.log"

var logger *log.Logger

var c Config = Config{
	"192.168.86.193:3307", //mysql dsn
	"new_dog123",          //mysql dbName
	"dog123:dog123",       //mysqldbAuth
	"192.168.86.193:6380", //redis add
	"", //redisQueueName
	"", //saveDir
	"", //imgUrl
	0,  //threaNum
	"", //logPath
	"", //startUrl
	"", //startUrlTag
}

func load() {
	if checkFileIsExist(c.logPath) {
		fmt.Println("file exist")
		file, err := os.OpenFile(c.logPath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			os.Exit(1)
		}
		logger = log.New(file, "", log.LstdFlags)

	} else {
		fmt.Println("new file path", c.logPath)
		file, err := os.Create(c.logPath)
		if err != nil {
			fmt.Println("create file error", err)
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

func Init(args []string) {
	loadConfig(args)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("useage: datadir startUrl startUrlType logfile threadnum")
		os.Exit(1)
	}
	Init(os.Args)
	load()
	logger.Println("[info]start ", c.startUrl)
	req := newRequest(c.startUrlTag, c.startUrl)
	spider.NewSpider(NewMyPageProcesser(), "TaskName").
		AddRequest(req).
		AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		SetThreadnum(uint(c.threadNum)).            // Crawl request by three Coroutines
		Run()
}
