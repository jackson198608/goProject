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
	"recommendActiveUserByDog",
	"/tmp/recommend.log",
	1,
	"192.168.86.192:27017", //BidData
	"BidData",
	"192.168.86.104:27017", //mongo ActiveUser RecommendData
	"RecommendData",
	"50",
	"3"}

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
	realTasks := Pushdata.GetAllActiveUsers(c.mongoConn1)
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

func pushUser() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.mongoConn, c.pushLimit, c.mongoConn1, c.pushDogLimit)
	rc := createClient()
	queueName := "user_recomment_data_status"
	num := 0
	for {
		v, _ := rc.Get(queueName).Result()

		if v == "" {
			logger.Info("got nothing user recommend queue")
			time.Sleep(3600 * time.Second)
			num++
			if num>23 {
				break
			}
			continue
		}

		//生产任务
		pushAllActiveUserToRedis(c.queueName)

		//处理任务
		r.LoopPushRecommend()

		rc.Del(queueName)

		logger.Info("*********  start dog push *******")
		pushDog()
		break;
	}
}

func pushDog() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName1, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName, c.mongoConn1, c.pushLimit, c.mongoConn1, c.pushDogLimit)
	
	//生产任务
	pushAllActiveUserToRedis(c.queueName1)

	//处理任务
	r.LoopPushRecommend()
}


func Init() {
	loadConfig()
	logger.SetConsole(true)
	logger.SetLevel(logger.DEBUG)
}
func main() {
	Init()
	jobType := os.Args[1]
	switch jobType {
	// case "pushdog":
	// 	logger.Info("in the do dog")
	// 	pushDog()
	case "pushuser":
		logger.Info("in the do")
		pushUser()
	default:
	}
}
