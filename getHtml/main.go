package main

import (
	"fmt"
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
	"getHtml",
	"/data/thread/",
	"/tmp/create_thread.log",
	1,
	"300", "500", "http://m.goumin.com/bbs/"}

var offset = 1000

func Init(args []string) {

	loadConfig(args)
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
	logger.Error(logger.DEBUG)

}

func saveHtmlUrl(jobType string, cat string) {
	r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", c.numloops, c.dbAuth, c.dbDsn, c.dbName)
	// page := 1
	intIdStart, _ := strconv.Atoi(c.tidStart)
	startId := intIdStart
	endId := startId + 1000
	maxId := getMaxId(jobType)
	fmt.Println(maxId)
	for {
		var ids []string
		if jobType == "threadsave" {
			ids = getThreadTask(startId, endId, cat)
		}
		if jobType == "asksave" {
			ids = getAskList(startId, endId, cat)
		}
		// setMaxid, _ := strconv.Atoi(c.tidEnd)
		// if startId > setMaxid {
		// 	break
		// }
		// if ids == nil {
		// 	break
		// }
		r.PushTaskData(ids)
		if cat == "update" {
			break
		}
		if startId > maxId {
			break
		}
		startId += offset
		endId += offset
	}
}

func createHtmlByUrl(jobType string) {
	r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, jobType, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.saveDir, c.tidStart, c.tidEnd, c.domain)
	r.Loop()
}

func main() {
	Init(os.Args)
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
		logger.Info("in the thread html get ")
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
