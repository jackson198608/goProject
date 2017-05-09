package main

import (
	"github.com/donnie4w/go-logger/logger"
	// "github.com/jackson198608/goProject/eventLog/task"
	"fmt"
	"os"
	// "log"
)

var c Config = Config{
	"192.168.86.72:3309",
	"test_dz2",
	"root:goumintech",
	1,
	10, //2545,
	1,
	"127.0.0.1:6379",
	"moveEvent",
	"/tmp/moveEdddvent.log", 0, "3", "2014-01-01", "1", "192.168.86.68:27017", "Event"}

func pushALLEventIdFromStartToEnd() {
	r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.fansLimit, c.dateLimit)
	page := 0
	for {
		ids := getTask(page)
		if page*c.numloops > c.lastId {
			break
		}
		if len(ids) == 0 {
			page++
			continue
		}
		if ids == nil {
			break
		}
		r.PushTaskData(ids)
		page++
	}
}

func do() {
	r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.fansLimit, c.dateLimit, c.logFile)
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
	// data := getEventLogData(1,10,10,0)
	// fmt.Println(data)
	// NewTask(1)
	jobType := os.Args[1]
	fmt.Println(jobType)
	switch jobType {
	case "create":
		logger.Info("in the create", 10)
		pushALLEventIdFromStartToEnd()

	case "do":
		logger.Info("in the do")
		do()
	default:

	}
}
