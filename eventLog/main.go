package main

import (
	"github.com/donnie4w/go-logger/logger"
	// "github.com/jackson198608/goProject/eventLog/task"
	// "fmt"
	"github.com/jackson198608/goProject/stayProcess"
	"os"
	// "time"
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
	"/tmp/moveEdddvent.log", 0, "3", "2014-01-01", "1", "192.168.86.68:27017,192.168.86.68:27017,192.168.86.68:27017", "Event", 10, 1, "200", "1", "100", "1000"}

var followQueue = "followData"
var limit = 1000

// var pushLimit = 30

// func pushALLEventIdFromStartToEnd() {
// 	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.dateLimit)
// 	page := 0
// 	for {
// 		for {
// 			lens := (*r.client).LLen(c.queueName).Val()
// 			fmt.Println(lens)
// 			if int(lens) > 10000 {
// 				time.Sleep(2 * time.Second)
// 				continue
// 			} else {
// 				break
// 			}
// 		}
// 		ids := getTask(page)
// 		offset := page * c.numloops
// 		if offset > c.lastId {
// 			break
// 		}
// 		if len(ids) == 0 {
// 			page++
// 			continue
// 		}
// 		if ids == nil {
// 			break
// 		}
// 		// r.PushTaskData(ids)
// 		page++
// 	}
// }

func createRedis() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.dateLimit, c.redisStart, c.redisEnd)
	r.PushData()
}

//动态推送给用户的mongo动态id
func createPushRedis() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.dateLimit, c.redisStart, c.redisEnd)
	r.PushFansData()
}

func pushAllFollowUserToRedis() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.dateLimit, c.redisStart, c.redisEnd)
	page := 0
	for {
		ids := getFollowTask(page, limit)
		offset := page * limit
		if offset > c.followLastId {
			break
		}
		if len(ids) == 0 {
			break
		}
		if ids == nil {
			break
		}
		r.PushFollowTaskData(ids)
		page++
	}
}

func getFansUserToRedis() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.dateLimit, c.redisStart, c.redisEnd)
	page := 0
	for {
		ids := getFansTask(page, limit)
		offset := page * limit
		if offset > c.followLastId {
			break
		}
		if len(ids) == 0 {
			break
		}
		if ids == nil {
			break
		}
		r.PushFollowTaskData(ids)
		page++
	}
}

func do() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.mongoConn, c.dateLimit, c.redisStart, c.redisEnd)
	r.Loop()
}

func push() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.mongoConn, c.dateLimit, c.redisStart, c.redisEnd, c.fansLimit, c.eventLimit, c.pushLimit)
	r.LoopPush()
}

func removeData() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.mongoConn, c.dateLimit, c.redisStart, c.redisEnd, c.fansLimit, c.eventLimit, c.pushLimit)
	r.RemoveFansData()
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
	// case "create":
	// 	logger.Info("in the create", 10)
	// 	pushALLEventIdFromStartToEnd()
	case "newcreate":
		logger.Info("in the create", 10)
		createRedis()
	case "addcreate": //动态粉丝增量
		logger.Info("in the create", 10)
		createPushRedis()
	case "do":
		logger.Info("in the do")
		do()
	case "follow":
		logger.Info("in the follow create")
		pushAllFollowUserToRedis()
	case "push":
		logger.Info("in the push fans data")
		push()
	case "fans":
		logger.Info("in the fans create")
		getFansUserToRedis()
	case "remove":
		logger.Info("in the fans create")
		removeData()
	default:

	}
}
