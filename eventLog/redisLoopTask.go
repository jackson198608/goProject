package main

import (
	// "bufio"
	"errors"
	"github.com/donnie4w/go-logger/logger"
	// "github.com/jackson198608/goProject/eventLog/task"
	"fmt"
	redis "gopkg.in/redis.v4"
	"os"
	// "reflect"
	"strconv"
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

// func (t *RedisEngine) putFailOneBack(backstr string) {
// }

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

func (t *RedisEngine) getTaskNum() {
	len := (*t.client).LLen(t.queueName).Val()
	if int(len) > t.numForOneLoop {
		t.taskNum = t.numForOneLoop
	} else {
		t.taskNum = int(len)
	}
}

func (t *RedisEngine) croutinePopJobData(c chan int, i int) {
	// tableNumInt := i + 1
	// tableNumStr := strconv.Itoa(tableNumInt)
	queueName := t.queueName //+ "_" + tableNumStr
	logger.Info("pop ", queueName)

	for {
		//doing until got nothing
		redisStr := (*t.client).LPop(queueName).Val()
		if redisStr == "" {
			logger.Info("got nothing", queueName)
			c <- 1
			return
		}

		//doing job
		id, _ := strconv.Atoi(redisStr)
		u := LoadById(id)
		fans := GetFansData(u.uid)
		// fmt.Println("type:", reflect.TypeOf(fans))
		var err1 error
		var f *os.File
		if checkFileIsExist(filename) { //如果文件存在
			f, err1 = os.OpenFile(filename, os.O_WRONLY, 0666) //打开文件
		} else {
			f, err1 = os.Create(filename) //创建文件
		}

		// f, err := os.OpenFile(filename, os.O_WRONLY, 0644)
		if err1 != nil {
			fmt.Println("cacheFileList.yml file create failed. err: " + err1.Error())
		}
		// fmt.Println("type:", reflect.TypeOf(f))
		SaveMongoEventLog(u, fans, f)
		defer f.Close()
		// w.Flush()
	}
}

func (t *RedisEngine) Loop() {
	logger.Info("do in the loop")
	t.taskNum = t.numForOneLoop
	t.doOneLoop()
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
