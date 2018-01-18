package main

import (
	"errors"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/Recommend_1.0/task"
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/jackson198608/goProject/common/tools"
	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v4"
	"os"
	"strings"
)

var c Config = Config{
	"192.168.86.193:3307", //mysql dsn
	"new_dog123",          //mysql dbName
	"dog123:dog123",       //mysqldbAuth
	"127.0.0.1:6379",      //redis info
	1,                     //thread num
	"pushContentCenter",   //queuename
	"192.168.86.192:27017",
	"192.168.86.5:9200/"} // mongo

func init() {
	loadConfig()
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

//全部活跃用户
func getAllActiveUsers(mongoConn string) []int {
	var user []int
	var session *mgo.Session
	var err error
	mgoInfos := strings.Split(mongoConn, ",")
	if len(mgoInfos) == 1 {
		session, err = tools.GetStandAloneConnecting(mongoConn)
	} else {
		session, err = tools.GetReplicaConnecting(mgoInfos)
	}
	if err != nil {
		return user
	}

	defer session.Close()

	c := session.DB("ActiveUser").C("active_user")
	err = c.Find(nil).Distinct("uid", &user)
	if err != nil {
		panic(err)
	}
	return user
}

func pushAllActiveUserToRedis(queueName string) bool {
	rc := createClient()
	realTasks := getAllActiveUsers(c.mongoConn)
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

func main() {
	jobType := os.Args[1]
	switch jobType {
	case "recommend": //push content conter
		var mongoConnInfo []string
		mongoConnInfo = append(mongoConnInfo, c.mongoConn)
		var mysqlInfo []string
		mysqlInfo = append(mysqlInfo, c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")

		redisInfo := redis.Options{
			Addr: c.redisConn,
		}
		//生产任务
		pushAllActiveUserToRedis(c.queueName)

		logger.Info("start work")
		r, err := redisEngine.NewRedisEngine(c.queueName, &redisInfo, mongoConnInfo, mysqlInfo, c.coroutinNum, jobFuc, c.elkDsn)
		if err != nil {
			logger.Error("[NewRedisEngine] ", err)
		}

		err = r.Do()
		if err != nil {
			logger.Error("[redisEngine Do] ", err)
		}
	case "--help":
		help()
	default:
		fmt.Println("unsupported params")
	}
}

func jobFuc(job string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {
	if (mysqlConns == nil) || (mgoConns == nil) {
		return errors.New("mysql or mongo conn error")
	}
	t, err := task.NewTask(job, mysqlConns, mgoConns, taskarg)
	if err != nil {
		return err
	}
	err = t.Do()
	if err != nil {
		return err
	}
	return err
}

func help() {
	fmt.Println("usage: pushApp [options]")
	fmt.Println("Options:")
	// fmt.Println("  allindex\t\t\t\tThe index of all collections is deleted, and then a new index is created")
	// fmt.Println("  singleindex_userId\t\t\tSpecify a collection to create an index")
	fmt.Println("  recommend\t\t\t\t\trecommend data")
	fmt.Println("  --help\t\t\t\tshow this usage information")
}
