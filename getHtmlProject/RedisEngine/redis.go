package RedisEngine

import (
	"errors"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"github.com/jackson198608/goProject/getHtmlProject/Task"
	redis "gopkg.in/redis.v4"
	"time"
)

type Engine struct {
	logLevel      int
	queueName     string
	connstr       string
	jobType       string
	client        *redis.Client
	taskNum       int
	numForOneLoop int
	taskNewArgs   []string
}

func NewEngine(
	logLevel int,
	queueName string,
	connstr string,
	jobType string,
	numForOneLoop int, taskarg ...string) *Engine {

	logger.SetLevel(logger.LEVEL(logLevel))

	t := new(Engine)

	if queueName == "" || connstr == "" || numForOneLoop <= 0 {
		return nil
	}

	t.logLevel = logLevel
	t.queueName = queueName
	t.connstr = connstr
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

func (t *Engine) connect() error {
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
func (t *Engine) PushTaskData(tasks interface{}) bool {
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

func (t *Engine) getTaskNum() {
	len := (*t.client).LLen(t.queueName).Val()
	if int(len) > t.numForOneLoop {
		t.taskNum = t.numForOneLoop
	} else {
		t.taskNum = int(len)
	}
}
func (t *Engine) Loop() {
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
func (t *Engine) doOneLoop() {
	logger.Info("do in oneloop taskNum", t.taskNum)
	c := make(chan int, t.taskNum)
	for i := 0; i < t.taskNum; i++ {
		go t.croutinePopJobData(c, i)
	}

	for i := 0; i < t.taskNum; i++ {
		<-c
	}
}

func (t *Engine) croutinePopJobData(c chan int, i int) {
	abuyun := setAbuyun()
	for {
		logger.Info("pop ", t.queueName)
		redisStr := (*t.client).LPop(t.queueName).Val()
		if redisStr == "" {
			logger.Info("got nothing ", t.queueName)
			c <- 1
			return
		}
		logger.Info("got redisStr ", redisStr)
		task := Task.NewTask(t.logLevel, t.queueName, redisStr, t.taskNewArgs, t.client, abuyun)
		if task != nil {
			task.Do()
		}
	}
}

const proxyServer = "http-pro.abuyun.com:9010"
const proxyUser = "HK71T41EZ21304GP"
const proxyPasswd = "75FE0C4E23EEA0E7"

// const proxyServer = ""
// const proxyUser = ""
// const proxyPasswd = ""

func setAbuyun() *abuyunHttpClient.AbuyunProxy {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy(proxyServer, proxyUser, proxyPasswd)
	return abuyun
}
