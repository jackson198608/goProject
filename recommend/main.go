package main

import (
	"github.com/donnie4w/go-logger/logger"
	// "github.com/jackson198608/goProject/eventLog/task"
	// "fmt"
	"github.com/jackson198608/goProject/recommend/Pushdata"
	"github.com/jackson198608/goProject/stayProcess"
	"os"
)

var c Config = Config{
	"210.14.154.198:33068",
	"new_dog123",
	"dog123:dog123",
	100,
	"192.168.86.56:6379",
	"recommendActiveUser",
	"/tmp/recommend.log",
	1,
	"192.168.86.192:27017,192.168.86.192:27017,192.168.86.192:27017",
	"BidData",
	"1000"}

func pushAllActiveUserToRedis() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.mongoConn)
	ids := Pushdata.GetAllActiveUsers()
	if len(ids) == 0 {
		return 
	}
	if ids == nil {
		return
	}
	r.PushActiveUserTaskData(ids)
}

func push() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.mongoConn)
	r.LoopPushRecommend()
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
	case "activeuser":
		logger.Info("in the create active user", 10)
		pushAllActiveUserToRedis()
	case "push":
		logger.Info("in the do")
		push()
	default:
	}
}
