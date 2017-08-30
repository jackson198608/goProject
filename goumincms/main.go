package main

import (
	"github.com/donnie4w/go-logger/logger"
	"os"
)

var logLevel int = 1

var dbAuth string = "dog123:dog123"

var dbDsn string = "192.168.86.193:3307"

var dbName string = "new_dog123"

var saveDir string = "/data/cms/"

var numloops int = 10

var queueName string = "createHtml"

var redisConn string = "127.0.0.1:6379"

// var logger *log.Logger
var logPath string = "/tmp/create_thread.log"

var templatefile = "/data/thread/template.html"

func Init() {

	// loadConfig()
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
	logger.Error(logger.DEBUG)

}

func createThreadHtml() {
	r := NewRedisEngine(logLevel, queueName, redisConn, "", 0, numloops, dbAuth, dbDsn, dbName)
	r.Loop()
}

func main() {
	Init()
	jobType := os.Args[1]
	switch jobType {
	case "thread": //thread
		logger.Info("in the thread html create ")
		createThreadHtml()
	case "info": //资讯
		logger.Info("in the info html create ")
		// createInfoHtml()
	default:

	}

}
