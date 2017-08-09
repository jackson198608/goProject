package main

import (
	"github.com/donnie4w/go-logger/logger"
	// "github.com/jackson198608/goProject/eventLog/task"
	// "fmt"
	"github.com/jackson198608/goProject/recommend/Pushdata"
	"github.com/jackson198608/goProject/stayProcess"
	"time"
	"os"
	redis "gopkg.in/redis.v4"
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

func createClient() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     c.redisConn,
        Password: "",
        DB:       0,
    })

    // 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
    _, err := client.Ping().Result()
    if err != nil {
		logger.Error("redis connect error", err)
	}
    return client
}

func push() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.mongoConn)
	c := createClient()
	queueName := "user_recomment_data_status"

	for {
		v, _ := c.Get(queueName).Result()

		if v == "" {
			logger.Info("got nothing user recommend queue")
			time.Sleep(3 * time.Second)
			continue
		}

		//生产任务
		ids := Pushdata.GetAllActiveUsers()
		if len(ids) == 0 {
			return
		}
		if ids == nil {
			return
		}
		pushAllActiveUserToRedis()

		//处理任务
		r.LoopPushRecommend()

		c.Del(queueName)

		break;
	}
}

func Init() {
	loadConfig()
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
	// logger.Error(logger.DEBUG)

}
func main() {
	Init()
	jobType := os.Args[1]
	switch jobType {
	// case "activeuser":
	// 	logger.Info("in the create active user", 10)
	// 	pushAllActiveUserToRedis()
	case "push":
		logger.Info("in the do")
		push()
	default:
	}
}
