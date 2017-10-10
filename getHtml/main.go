package main

import (
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

// func getHtmlUrls(pm string, jobType string, page int) []string {
// 	var urls []string
// 	if jobType == "asksave" {
// 		ids := getAskList(page)
// 		domain := "http://ask.goumin.com/ask/"
// 		if pm == "h5" {
// 			domain = "http://ask.m.goumin.com/ask/"
// 		}
// 		for _, v := range ids {
// 			url := domain + strconv.Itoa(v) + ".html"
// 			urls = append(urls, url)
// 		}
// 	}
// 	if jobType == "threadsave" {
// 		ids := getThreadTask(page)
// 		domain := "http://bbs.goumin.com/"
// 		if pm == "h5" {
// 			domain = "http://m.goumin.com/bbs/"
// 		}
// 		for _, v := range ids {
// 			url := domain + "thread-" + strconv.Itoa(v) + "-1-1.html"
// 			urls = append(urls, url)
// 		}
// 	}
// 	return urls
// }

func saveHtmlUrl(jobType string) {
	r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", c.numloops, c.dbAuth, c.dbDsn, c.dbName)
	page := 1
	for {
		var ids []string
		// urls := getHtmlUrls(pm, jobType, page)
		if jobType == "threadsave" {
			ids = getThreadTask(page)
		}
		if jobType == "asksave" {
			ids = getAskList(page)
		}
		offset := page * offset
		maxThreadid, _ := strconv.Atoi(c.tidEnd)
		if offset > maxThreadid {
			break
		}
		if len(ids) == 0 {
			break
		}
		if ids == nil {
			break
		}
		r.PushTaskData(ids)
		page++
	}
}

func createHtmlByUrl(jobType string) {
	r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, jobType, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.saveDir, c.tidStart, c.tidEnd, c.domain)
	r.Loop()
}

func main() {
	Init(os.Args)
	jobType := os.Args[1]
	if len(os.Args) < 3 {
		logger.Error("args error")
		return
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
		saveHtmlUrl(jobType)
	case "threadsave": //create html url
		logger.Info("in the html url save ")
		saveHtmlUrl(jobType)
	default:

	}

}
