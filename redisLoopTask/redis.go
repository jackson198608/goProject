package redisLoopTask

import (
	"errors"
	//"github.com/jackson198608/goProject/tableSplit/pre_forum_post/task"
	redis "gopkg.in/redis.v4"
	"strconv"
	"strings"
)

type RedisEngine struct {
	queueName     string
	connstr       string
	password      string
	db            int
	client        *redis.Client
	taskNum       int
	numForOneLoop int
}

func NewRedisEngine(queueName string, connstr string, password string, db int, numForOneLoop int) *RedisEngine {
	t := new(RedisEngine)

	if queueName == "" || connstr == "" || numForOneLoop <= 0 {
		return nil
	}

	t.queueName = queueName
	t.connstr = connstr
	t.password = password
	t.db = db
	t.numForOneLoop = numForOneLoop

	err := t.connect()
	if err != nil {
		fmt.Println("[Error] redis connect error", err)
		return nil
	}

	return t
}

func (t *RedisEngine) putFailOneBack(backstr string) {
}

func (t *RedisEngine) connect(conn string) error {
	t.client = redis.NewClient(&redis.Options{
		Addr:     t.conn,
		Password: t.password, // no password set
		DB:       t.db,       // use default DB
	})
	_, err := t.client.Ping().Result()
	if err != nil {
		return errors.New("[Error] redis connect error")
	}
	return nil
}

func (t *RedisEngine) getTaskNum() {
	len := (*t.client).LLen(redisQueueName).Val()
	if int(len) > numForOneLoop {
		t.taskNum = numForOneLoop
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

	c <- 1
}

func (t *RedisEngine) Loop() {
}

//it's for doing job at one time using tasknum's croutine
func (t *RedisEngine) doOneLoop() {
	c := make(chan int, taskNum)
	for i := 0; i < taskNum; i++ {
		go croutinePopJobData(c, i)
	}

	for i := 0; i < taskNum; i++ {
		<-c
	}
}
