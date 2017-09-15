package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v4"
	"strconv"
	"time"
)

type RedisEngine struct {
	logLevel      int
	queueName     string
	connstr       string
	mongoConn     string
	jobType       string
	client        *redis.Client
	taskNum       int
	numForOneLoop int
	taskNewArgs   []string
}

func NewRedisEngine(
	logLevel int,
	queueName string,
	connstr string,
	mongoConn string,
	jobType string,
	numForOneLoop int, taskarg ...string) *RedisEngine {

	logger.SetLevel(logger.LEVEL(logLevel))

	t := new(RedisEngine)

	if queueName == "" || connstr == "" || numForOneLoop <= 0 {
		return nil
	}

	t.logLevel = logLevel
	t.queueName = queueName
	t.connstr = connstr
	t.mongoConn = mongoConn
	t.jobType = jobType
	t.numForOneLoop = numForOneLoop
	t.taskNewArgs = taskarg
	err := t.connect()
	if err != nil {
		logger.Error("redis connect error", err)
		return nil
	}

	return t
}

func (t *RedisEngine) connect() error {
	t.client = redis.NewClient(&redis.Options{
		Addr:     t.connstr,
		Password: "", // no password set
		DB:       0,  //e.db use default DB
	})
	_, err := t.client.Ping().Result()
	if err != nil {
		return errors.New("[Error] redis connect error")
	}
	return nil
}

func (t *RedisEngine) getTaskNum() {
	len := (*t.client).LLen(t.queueName).Val()
	if int(len) > t.numForOneLoop {
		t.taskNum = t.numForOneLoop
	} else {
		t.taskNum = int(len)
	}
}

func (t *RedisEngine) croutinePopJobData(c chan int, i int) {
	dbAuth := t.taskNewArgs[0]
	dbDsn := t.taskNewArgs[1]
	dbName := t.taskNewArgs[2]
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("[error] connect db err")
		return
	}
	defer db.Close()

	session, err := mgo.Dial(t.mongoConn)
	if err != nil {
		logger.Error("[error] connect mongodb err")
		return
	}
	defer session.Close()
	relateAsk := ""
	if t.jobType == "thread" {
		fmt.Println("&&&&")
		relateAsk = LoadDefaultRelateAsk(db)
	}
	doctors := ""
	if t.jobType == "ask" {
		fmt.Println("defalt doctor data")
		doctors = DefaultDoctors(db)
	}
	for {
		logger.Info("pop ", t.queueName)
		redisStr := (*t.client).LPop(t.queueName).Val()
		fmt.Println(redisStr)
		if redisStr == "" {
			logger.Info("got nothing ", t.queueName)
			c <- 1
			return
		}
		logger.Info("got redisStr ", redisStr)
		task := NewTask(t.logLevel, redisStr, db, session, t.taskNewArgs, relateAsk, t.jobType)
		if task != nil {
			task.Do()
		}
	}
}

func (t *RedisEngine) Loop() {
	logger.Info("do in the loop")
	for {
		t.getTaskNum()
		logger.Info("got nothing", t.queueName)
		if t.taskNum == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		t.doOneLoop()
	}
}

//it's for doing job at one time using tasknum's croutine
func (t *RedisEngine) doOneLoop() {
	logger.Info("do in oneloop taskNum", t.taskNum)
	c := make(chan int, t.taskNum)
	for i := 0; i < t.taskNum; i++ {
		go t.croutinePopJobData(c, i)
	}

	for i := 0; i < t.taskNum; i++ {
		<-c
	}
}

func (t *RedisEngine) PushThreadTaskData(tasks interface{}) bool {
	switch realTasks := tasks.(type) {
	case []string:
		logger.Info("this is string task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			err := (*t.client).RPush(t.queueName, realTasks[i]).Err()
			if err != nil {
				logger.Error("insert redis error", err)
			}
		}

	case []int:
		logger.Info("this is int task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			err := (*t.client).RPush(t.queueName, realTasks[i]).Err()
			if err != nil {
				logger.Error("insert redis error", err)
			}
		}

	default:
		logger.Error("this is not normal format", realTasks)
		return false
	}

	return true
}

func (t *RedisEngine) PushIdData() bool {
	redisStart, _ := strconv.Atoi(t.taskNewArgs[3])
	redisEnd, _ := strconv.Atoi(t.taskNewArgs[4])
	logger.Info("RPush queueName string", t.queueName)
	for i := redisStart; i <= redisEnd; i++ {
		err := (*t.client).RPush(t.queueName, i).Err()
		if err != nil {
			logger.Error("insert redis error", err)
		}
	}
	return true
}
