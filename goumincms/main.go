package main

import (
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	"os"
)

var logLevel int = 1

var dbAuth string = "dog123:dog123"

var dbDsn string = "192.168.86.193:3307"

var dbName string = "new_dog123"

var saveDir string = "/data/thread/"

var staticH5Url string = "http://m.goumin.com/bbs/"

var numloops int = 10

var queueName string = "createHtml"

var redisConn string = "127.0.0.1:6379"

var mongoConn string = "192.168.86.192:27017"

// var logger *log.Logger
var logPath string = "/tmp/create_thread.log"

var h5templatefile = "/data/thread/h5template.html"
var miptemplatefile = "/data/thread/miptemplate.html"

var diaryDomain = "http://c1.cdn.goumin.com/diary/"

var offset = 10
var maxThreadid = 100

func Init() {

	// loadConfig()
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
	logger.Error(logger.DEBUG)

}

func createThreadHtml(templateType string) {
	r := NewRedisEngine(logLevel, queueName, redisConn, "", 0, numloops, dbAuth, dbDsn, dbName, templateType)
	r.Loop()
}

func createRedis() {
	r := NewRedisEngine(logLevel, queueName, redisConn, "", 0, numloops, dbAuth, dbDsn, dbName)
	page := 1
	for {
		tids := getThreadTask(page)
		offset := page * offset
		if offset > maxThreadid {
			break
		}
		if len(tids) == 0 {
			break
		}
		if tids == nil {
			break
		}
		r.PushThreadTaskData(tids)
		page++
	}
}

func main() {
	Init()
	jobType := os.Args[1]
	templateType := ""
	if len(os.Args) == 3 {
		templateType = os.Args[2] //string 默认0不传参数或是参数是0为mip模板 ,1:h5
	}
	switch jobType {
	case "thread": //thread
		logger.Info("in the thread html create ")
		createThreadHtml(templateType)
	case "info": //资讯
		logger.Info("in the info html create ")
	case "create": //资讯
		logger.Info("in the info html create ")
		createRedis()
	default:

	}

}
