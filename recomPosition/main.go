package main

import (
	"github.com/donnie4w/go-logger/logger"
	// "github.com/jackson198608/goProject/eventLog/task"
	// "fmt"
	"github.com/jackson198608/goProject/recomPosition/Pushdata"
	"github.com/jackson198608/goProject/stayProcess"
	redis "gopkg.in/redis.v4"
	// "os"
	// "time"
)

var c Config = Config{
	"192.168.86.193:3307",
	"new_dog123",
	"dog123:dog123",
	100,
	"192.168.86.56:6379",
	"recommendActiveUser",
	"/tmp/recommend.log",
	1,
	"192.168.86.122:27017",
	"RecommendData"}

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

func pushAllActiveUserToRedis(queueName string) bool {
	rc := createClient()
	realTasks := Pushdata.GetAllActiveUsers(c.mongoConn)
	if len(realTasks) == 0 {
		logger.Error("active user data is empty")
		return false
	}
	if realTasks == nil {
		logger.Error("not get active user data")
		return false
	}

	logger.Info("this is int task", realTasks)
	for i := 0; i < len(realTasks); i++ {
		err := rc.RPush(queueName, realTasks[i]).Err()
		if err != nil {
			logger.Error("insert redis error", err)
			return false
		}
	}
	return true
}

func push() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.mongoConn)

	//生产任务
	pushAllActiveUserToRedis(c.queueName)

	//处理任务
	r.LoopPushRecommendPosition()
}

func Init() {
	loadConfig()
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
}
func main() {
	Init()
	push()
}
