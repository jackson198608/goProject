package main

import (
	// "bufio"
	"errors"
	"github.com/donnie4w/go-logger/logger"
	// "github.com/jackson198608/goProject/eventLog/task"
	"fmt"
	redis "gopkg.in/redis.v4"
	// "os"
	"database/sql"
	// "reflect"
	mgo "gopkg.in/mgo.v2"
	"strconv"
	"strings"
	"time"
)

type RedisEngine struct {
	logLevel      int
	queueName     string
	connstr       string
	password      string
	db            int
	client        *redis.Client
	taskNum       int
	numForOneLoop int
	taskNewArgs   []string
}

func NewRedisEngine(
	logLevel int,
	queueName string,
	connstr string,
	password string,
	db int,
	numForOneLoop int, taskarg ...string) *RedisEngine {

	logger.SetLevel(logger.LEVEL(logLevel))

	t := new(RedisEngine)

	if queueName == "" || connstr == "" || numForOneLoop <= 0 {
		return nil
	}

	t.logLevel = logLevel
	t.queueName = queueName
	t.connstr = connstr
	t.password = password
	t.db = db
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
		Password: t.password, // no password set
		DB:       t.db,       // use default DB
	})
	_, err := t.client.Ping().Result()
	if err != nil {
		return errors.New("[Error] redis connect error")
	}
	return nil
}

func (t *RedisEngine) getTaskQueueName(redisStr string) string {
	redisInt, err := strconv.Atoi(redisStr)
	if err != nil {
		logger.Error("not normal redisStr", redisStr)
		return ""
	} else {
		tableNameInt := redisInt % 100
		if tableNameInt == 0 {
			tableNameInt = 100
		}
		// talbeNameStr := strconv.Itoa(tableNameInt)
		return t.queueName //+ "_" + talbeNameStr

	}
}
func (t *RedisEngine) PushTaskData(tasks interface{}) bool {
	switch realTasks := tasks.(type) {
	case []string:
		logger.Info("this is string task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			queueName := t.queueName //t.getTaskQueueName(realTasks[i])
			err := (*t.client).RPush(queueName, realTasks[i]).Err()
			if err != nil {
				logger.Error("insert redis error", err)
			}
		}

	case []int64:
		logger.Info("this is int task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			// redisStr := strconv.Itoa(int(realTasks[i]))
			queueName := t.queueName //t.getTaskQueueName(redisStr)
			err := (*t.client).RPush(queueName, realTasks[i]).Err()
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

func (t *RedisEngine) PushData() bool {
	for i := c.redisStart; i < c.redisEnd; i++ {
		queueName := t.queueName //t.getTaskQueueName(realTasks[i])
		err := (*t.client).RPush(queueName, i).Err()
		if err != nil {
			logger.Error("insert redis error", err)
		}
	}
	return true

}

func (t *RedisEngine) PushFollowTaskData(tasks interface{}) bool {
	switch realTasks := tasks.(type) {
	case []string:
		logger.Info("this is string task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			queueName := followQueue //t.getTaskQueueName(realTasks[i])
			err := (*t.client).RPush(queueName, realTasks[i]).Err()
			if err != nil {
				logger.Error("insert redis error", err)
			}
		}

	case []int64:
		logger.Info("this is int task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			// redisStr := strconv.Itoa(int(realTasks[i]))
			queueName := followQueue //t.getTaskQueueName(redisStr)
			err := (*t.client).RPush(queueName, realTasks[i]).Err()
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

func (t *RedisEngine) getTaskNum() {
	len := (*t.client).LLen(t.queueName).Val()
	if int(len) > t.numForOneLoop {
		t.taskNum = t.numForOneLoop
	} else {
		t.taskNum = int(len)
	}
}

func (t *RedisEngine) croutinePopJobData(x chan int, i int) {
	// tableNumInt := i + 1
	// tableNumStr := strconv.Itoa(tableNumInt)
	queueName := t.queueName //+ "_" + tableNumStr
	logger.Info("pop ", queueName)
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	defer db.Close()
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		return
	}
	defer session.Close()
	for {
		//doing until got nothing
		redisStr := (*t.client).LPop(queueName).Val()
		if redisStr == "" {
			logger.Info("got nothing", queueName)
			x <- 1
			return
		}
		redisArr := strings.Split(redisStr, "|")
		if len(redisArr) == 2 {
			if redisArr[1] == "3" {
				uids := strings.Split(redisArr[0], "&")
				RemoveFansEventLog(uids[0], uids[1], session) //fuid,uid
			} else {
				u := LoadMongoById(redisArr[0], session)
				fans := GetFansData(u.Uid, db)
				status := redisArr[1] //要执行的操作:0:删除,-1隐藏,1显示,2动态推送给粉丝
				UpdateMongoEventLogStatus(u, fans, status, session)
			}
		}
		//doing job
		if len(redisArr) == 1 {
			id, _ := strconv.Atoi(redisStr)
			u := LoadById(id, db)
			// fans := GetFansData(u.uid, db)
			var fans []*Follow
			SaveMongoEventLog(u, fans, session)
		}
	}
}
func (t *RedisEngine) Loop() {
	logger.Info("do in the loop")
	// t.taskNum = t.numForOneLoop
	// t.doOneLoop()

	// fmt.Println(reflect.TypeOf(db))
	for {
		t.getTaskNum()
		fmt.Println(t.taskNum)
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

func (t *RedisEngine) LoopPush() {
	for {
		t.getPushTaskNum()
		if t.taskNum == 0 {
			logger.Info("no push data")
			time.Sleep(5 * time.Second)
			continue
		}
		t.doOneLoopPush()
	}
}

//it's for doing job at one time using tasknum's croutine
func (t *RedisEngine) doOneLoopPush() {
	logger.Info("do in oneloop taskNum", t.taskNum)
	c := make(chan int, t.taskNum)
	for i := 0; i < t.taskNum; i++ {
		go t.croutinePopJobFollowData(c, i)
	}

	for i := 0; i < t.taskNum; i++ {
		<-c
	}
}

func (t *RedisEngine) getPushTaskNum() {
	len := (*t.client).LLen(followQueue).Val()
	if int(len) > t.numForOneLoop {
		t.taskNum = t.numForOneLoop
	} else {
		t.taskNum = int(len)
	}
}

func (t *RedisEngine) croutinePopJobFollowData(x chan int, i int) {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	defer db.Close()
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		return
	}
	defer session.Close()
	for {
		//doing until got nothing
		redisStr := (*t.client).LPop(followQueue).Val()
		if redisStr == "" {
			logger.Info("got nothing", followQueue)
			x <- 1
			return
		}
		redisArr := strings.Split(redisStr, "|")
		if len(redisArr) == 3 {
			uid := redisArr[0]   //用户uid
			fans := redisArr[1]  //粉丝数
			count := redisArr[2] //粉丝数
			pushEventToFansTask(fans, uid, count, session, db)
		}

	}
}
