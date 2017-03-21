package main

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/tableSplit/pre_forum_post/redisLoopTask"
	"os"
)

var c Config = Config{
	"210.14.154.198:3306",
	"new_dog123",
	"dog123:dog123",
	100,
	4392691,
	0,
	"127.0.0.1:6379",
	"movePost",
	"/tmp/move.log", 0}

func pushALLTidFromStartToEnd() {
	r := redisLoopTask.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName)

	page := 0
	for {
		tids := getTask(page)
		if tids == nil {
			break
		}
		r.PushTaskData(tids)
		page++
	}
}

func do() {
	r := redisLoopTask.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, 100, c.dbAuth, c.dbDsn, c.dbName)
	r.Loop()
}

func Init() {

	loadConfig()
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
	logger.Error(logger.DEBUG)

}

func main() {
	Init()
	jobType := os.Args[1]

	switch jobType {
	case "create":
		logger.Info("in the create", 10)
		pushALLTidFromStartToEnd()

	case "do":
		logger.Info("in the do")
		do()
	default:

	}
}
