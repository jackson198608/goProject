package main

import (
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	"os"
	"strconv"
)

var c Config = Config{
	"192.168.86.193:3307",
	"new_dog123",
	"dog123:dog123",
	10,
	"127.0.0.1:6379",
	"createHtml",
	"/data/thread/",
	"/tmp/create_thread.log",
	1,
	"192.168.86.192:27017",
	"/data/thread/h5template.html",
	"/data/thread/miptemplate.html",
	"300", "500"}

// var logLevel int = 1

// var dbAuth string = "dog123:dog123"

// var dbDsn string = "192.168.86.193:3307"

// var dbName string = "new_dog123"

// var saveDir string = "/data/thread/"

// var numloops int = 10

// var redisConn string = "127.0.0.1:6379"

// var mongoConn string = "192.168.86.192:27017"

// // var logger *log.Logger
// var logPath string = "/tmp/create_thread.log"

// var h5templatefile = "/data/thread/h5template.html"
// var miptemplatefile = "/data/thread/miptemplate.html"

// var maxThreadid = 100

var diaryDomain = "http://c1.cdn.goumin.com/diary/"
var bbsDomain = "http://f1.cdn.goumin.com/attachments/"
var staticH5Url string = "http://m.goumin.com/bbs/"

// var queueName string = "createHtml"
var offset = 10

func Init(args []string) {

	loadConfig(args)
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
	logger.Error(logger.DEBUG)

}

func createThreadHtml(templateType string) {
	templatefile := c.miptemplatefile
	if templateType == "1" {
		templatefile = c.h5templatefile
	}

	r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, c.mongoConn, 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, templateType, templatefile, c.saveDir, c.tidStart, c.tidEnd)
	r.Loop()
}

func createRedis() {
	r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, c.mongoConn, 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName)
	page := 1
	for {
		tids := getThreadTask(page)
		offset := page * offset
		maxThreadid, _ := strconv.Atoi(c.tidEnd)
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

func createNumRedis() {
	r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, c.mongoConn, 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.tidStart, c.tidEnd)
	r.PushTidData()
}

func main() {
	Init(os.Args)
	jobType := os.Args[1]
	templateType := ""
	if len(os.Args) >= 4 {
		templateType = os.Args[3] //string 默认0不传参数或是参数是0为mip模板 ,1:h5
	}
	switch jobType {
	case "thread": //thread
		logger.Info("in the thread html create ")
		createThreadHtml(templateType)
	case "info": //资讯
		logger.Info("in the info html create ")
	case "create": //创建帖子idredis
		logger.Info("in the info html create ")
		createRedis()
	case "createnum": //创建帖子idredis
		logger.Info("in the info html create ")
		createNumRedis()
	default:

	}

}
