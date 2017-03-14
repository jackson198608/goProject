package redisLoopTask

import (
	"errors"
	"fmt"
	"github.com/jackson198608/goProject/tableSplit/pre_forum_post/task"
	redis "gopkg.in/redis.v4"
)

type RedisEngine struct {
	queueName     string
	connstr       string
	password      string
	db            int
	client        *redis.Client
	taskNum       int
	numForOneLoop int
	taskNewArgs   []interface{}
}

func NewRedisEngine(
	queueName string,
	connstr string,
	password string,
	db int,
	numForOneLoop int, taskarg ...interface{}) *RedisEngine {

	t := new(RedisEngine)

	if queueName == "" || connstr == "" || numForOneLoop <= 0 {
		return nil
	}

	t.queueName = queueName
	t.connstr = connstr
	t.password = password
	t.db = db
	t.numForOneLoop = numForOneLoop
	t.taskNewArgs = taskarg

	err := t.connect()
	if err != nil {
		fmt.Println("[Error] redis connect error", err)
		return nil
	}

	return t
}

func (t *RedisEngine) putFailOneBack(backstr string) {
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

func (t *RedisEngine) PushTaskData(tasks interface{}) bool {

	switch realTasks := tasks.(type) {
	case []string:
		fmt.Println("this is string task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			err := (*t.client).RPush(t.queueName, realTasks[i]).Err()
			if err != nil {
				fmt.Println("[error]insert redis error", err)
			}
		}

	case []int64:
		fmt.Println("this is int task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			err := (*t.client).RPush(t.queueName, realTasks[i]).Err()
			if err != nil {
				fmt.Println("[error]insert redis error", err)
			}
		}

	default:
		fmt.Println("task is no normal format", realTasks)
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
	fmt.Println("[notice] pop " + t.queueName)
	redisStr := (*t.client).LPop(t.queueName).Val()
	if redisStr == "" {
		fmt.Println("[notice] got nothing")
		c <- 1
		return
	}

	//doing job
	//task.NewTask(redisStr, "dog123:dog123", "210.14.154.198:3306", "new_dog123")
	task.NewTask(redisStr, t.taskNewArgs)
	c <- 1
}

func (t *RedisEngine) Loop() {
	fmt.Println("[notice]do in loop")
	t.doOneLoop()
}

//it's for doing job at one time using tasknum's croutine
func (t *RedisEngine) doOneLoop() {
	t.getTaskNum()
	fmt.Println("[notice]do in oneloop tasknum", t.taskNum)

	c := make(chan int, t.taskNum)
	for i := 0; i < t.taskNum; i++ {
		go t.croutinePopJobData(c, i)
	}

	for i := 0; i < t.taskNum; i++ {
		<-c
	}
}
