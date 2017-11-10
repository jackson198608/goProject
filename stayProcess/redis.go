package stayProcess

import (
	// "bufio"
	"errors"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/eventLog/task"
	// recomPosition "github.com/jackson198608/goProject/recomPosition/task"
	// recommendTask "github.com/jackson198608/goProject/recommend/task"
	redis "gopkg.in/redis.v4"
	// "os"
	"database/sql"
	mgo "gopkg.in/mgo.v2"
	// "reflect"
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
	}
	defer db.Close()
	// mongoConn := t.taskNewArgs[3]
	// session, err := mgo.Dial(mongoConn)
	// if err != nil {
	// 	logger.Error("[error] connect mongodb err")
	// 	return
	// }
	// defer session.Close()
	mongoConnArr := strings.Split(t.taskNewArgs[3], ",")
	if len(mongoConnArr) < 3 {
		logger.Error("[error] mongo config error")
		return
	}
	Host := []string{
		mongoConnArr[0],
		mongoConnArr[1],
		mongoConnArr[2],
	}
	const (
		Database       = "FansData"
		ReplicaSetName = "goumin"
	)

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          Host,
		Database:       Database,
		ReplicaSetName: ReplicaSetName,
	})

	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.SecondaryPreferred, false)

	for {
		logger.Info("pop ", t.queueName)
		redisStr := (*t.client).LPop(t.queueName).Val()
		if redisStr == "" {
			logger.Info("got nothing", t.queueName)
			c <- 1
			return
		}
		logger.Info("got redisStr ", redisStr)
		task := task.NewTask(t.logLevel, redisStr, db, session)
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

func (t *RedisEngine) PushData() bool {
	redisStart, _ := strconv.Atoi(t.taskNewArgs[4])
	redisEnd, _ := strconv.Atoi(t.taskNewArgs[5])
	logger.Info("RPush queueName string")
	for i := redisStart; i <= redisEnd; i++ {
		queueName := t.queueName //t.getTaskQueueName(realTasks[i])
		err := (*t.client).RPush(queueName, i).Err()
		if err != nil {
			logger.Error("insert redis error", err)
		}
	}
	return true
}

func (t *RedisEngine) PushFansData() bool {
	redisStart, _ := strconv.Atoi(t.taskNewArgs[4])
	redisEnd, _ := strconv.Atoi(t.taskNewArgs[5])
	for i := redisStart; i <= redisEnd; i++ {
		queueName := t.queueName //t.getTaskQueueName(realTasks[i])
		err := (*t.client).RPush(queueName, strconv.Itoa(i)+"|2").Err()
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
			queueName := "followData" //t.getTaskQueueName(realTasks[i])
			err := (*t.client).RPush(queueName, realTasks[i]).Err()
			if err != nil {
				logger.Error("insert redis error", err)
			}
		}

	case []int64:
		logger.Info("this is int task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			// redisStr := strconv.Itoa(int(realTasks[i]))
			queueName := "followData" //t.getTaskQueueName(redisStr)
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

func (t *RedisEngine) LoopPush() {
	for {
		t.getPushTaskNum()
		if t.taskNum == 0 {
			logger.Info("got nothing followData queue")
			time.Sleep(5 * time.Second)
			continue
		}
		t.doOneLoopPush()
	}
}

func (t *RedisEngine) LoopPushRecommend() {
	t.getPushRecommendTaskNum()
	t.doOneLoopPushRecommend()
}

func (t *RedisEngine) RemoveFansData() {
	// for {
	t.getPushTaskNum()
	if t.taskNum == 0 {
		logger.Info("got nothing followData queue")
		return
		// time.Sleep(5 * time.Second)
		// continue
	}
	t.removeDataLoop()
	// }
}

func (t *RedisEngine) removeDataLoop() {
	logger.Info("do in oneloop taskNum", t.taskNum)
	c := make(chan int, t.taskNum)
	for i := 0; i < t.taskNum; i++ {
		go t.croutinePopJobRemoveFansData(c, i)
	}

	for i := 0; i < t.taskNum; i++ {
		<-c
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

//it's for doing job at one time using tasknum's croutine
func (t *RedisEngine) doOneLoopPushRecommend() {
	logger.Info("do in oneloop taskNum", t.taskNum)
	c := make(chan int, t.taskNum)
	for i := 0; i < t.taskNum; i++ {
		go t.croutinePopJobRecommendActiveUserData(c, i)
	}

	for i := 0; i < t.taskNum; i++ {
		<-c
	}
}

func (t *RedisEngine) getPushRecommendTaskNum() {
	len := (*t.client).LLen(t.queueName).Val()
	if int(len) > t.numForOneLoop {
		t.taskNum = t.numForOneLoop
	} else {
		t.taskNum = int(len)
	}
}

func (t *RedisEngine) getPushTaskNum() {
	len := (*t.client).LLen("followData").Val()
	if int(len) > t.numForOneLoop {
		t.taskNum = t.numForOneLoop
	} else {
		t.taskNum = int(len)
	}
}

func (t *RedisEngine) croutinePopJobRecommendActiveUserData(x chan int, i int) {
	// fmt.Println(t.taskNewArgs)
	dbAuth := t.taskNewArgs[0]
	dbDsn := t.taskNewArgs[1]
	dbName := t.taskNewArgs[2]

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	defer db.Close()

	mongoConn := t.taskNewArgs[3]
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		logger.Error("[error] connect mongodb err")
		return
	}
	defer session.Close()

	for {
		//doing until got nothing]
		activeUserQueue := t.queueName
		redisStr := (*t.client).LPop(activeUserQueue).Val()
		// fmt.Println("***************redisStr: *************"+redisStr)
		if redisStr == "" {
			logger.Info("got nothing ", activeUserQueue)
			x <- 1
			return
		}

		// recommendTask := recommendTask.NewTask(t.logLevel, redisStr, db, session)
		// if recommendTask != nil {
		// 	recommendTask.Dopush()
		// }

	}
}

func (t *RedisEngine) croutinePopJobFollowData(x chan int, i int) {
	// fmt.Println(t.taskNewArgs)
	dbAuth := t.taskNewArgs[0]
	dbDsn := t.taskNewArgs[1]
	dbName := t.taskNewArgs[2]
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	defer db.Close()
	// mongoConn := t.taskNewArgs[3]
	// session, err := mgo.Dial(mongoConn)
	// if err != nil {
	// 	logger.Error("[error] connect mongodb err")
	// 	return
	// }
	// defer session.Close()

	mongoConnArr := strings.Split(t.taskNewArgs[3], ",")
	if len(mongoConnArr) < 3 {
		logger.Error("[error] mongo config error")
		return
	}
	Host := []string{
		mongoConnArr[0],
		mongoConnArr[1],
		mongoConnArr[2],
	}
	const (
		Database       = "FansData"
		ReplicaSetName = "goumin"
	)

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          Host,
		Database:       Database,
		ReplicaSetName: ReplicaSetName,
	})

	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.SecondaryPreferred, false)

	for {
		//doing until got nothing]
		followQueue := "followData"
		redisStr := (*t.client).LPop(followQueue).Val()
		if redisStr == "" {
			logger.Info("got nothing", followQueue)
			x <- 1
			return
		}

		task := task.NewTask(t.logLevel, redisStr, db, session)
		if task != nil {
			task.Dopush(t.taskNewArgs[4], t.numForOneLoop, t.taskNewArgs[7], t.taskNewArgs[8], t.taskNewArgs[9])
		}

	}
}

func (t *RedisEngine) croutinePopJobRemoveFansData(x chan int, i int) {
	dbAuth := t.taskNewArgs[0]
	dbDsn := t.taskNewArgs[1]
	dbName := t.taskNewArgs[2]
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	defer db.Close()
	/*mongoConn := t.taskNewArgs[3]
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		logger.Error("[error] connect mongodb err")
		return
	}
	defer session.Close()*/
	mongoConnArr := strings.Split(t.taskNewArgs[3], ",")
	if len(mongoConnArr) < 3 {
		logger.Error("[error] mongo config error")
		return
	}
	Host := []string{
		mongoConnArr[0],
		mongoConnArr[1],
		mongoConnArr[2],
	}
	const (
		Database       = "FansData"
		ReplicaSetName = "goumin"
	)

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          Host,
		Database:       Database,
		ReplicaSetName: ReplicaSetName,
	})

	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.SecondaryPreferred, false)

	for {
		//doing until got nothing]
		followQueue := "followData"
		redisStr := (*t.client).LPop(followQueue).Val()
		if redisStr == "" {
			logger.Info("got nothing", followQueue)
			x <- 1
			return
		}

		task := task.NewTask(t.logLevel, redisStr, db, session)
		if task != nil {
			//t.taskNewArgs[8]: eventLimit动态限制数量 sleeptime
			task.Doremove(t.numForOneLoop, t.taskNewArgs[8], t.taskNewArgs[10])
		}

	}
}

//baidu rank
func (t *RedisEngine) GetKeywordData(queueName string, start int64, stop int64) []string {
	arr := (*t.client).LRange(queueName, start, stop).Val()
	return arr
}

//baidu rank
func (t *RedisEngine) GetKeywordDataNew(queueName string) string {
	arr := (*t.client).LPop(queueName).Val()
	return arr
}

func (t *RedisEngine) SaveKeywordRedis(queueName string, value string) int64 {
	size, _ := (*t.client).LPush(queueName, value).Result()
	return size
}

func (t *RedisEngine) KeywordExist(queueName string) int64 {
	len := (*t.client).LLen(queueName).Val()
	return len
}

func (t *RedisEngine) CompensateKeyword(queueName string, keyword string, url string, num int) {
	//关键词补偿获取次数
	date := time.Now().Format("0102")
	key := keyword + "_times" + date
	timesStr := (*t.client).Get(key).Val()
	fmt.Println("key name:", key, len(timesStr))
	value := 1
	if len(timesStr) > 0 {
		times, _ := strconv.Atoi(timesStr)
		value = times + 1
	}
	if value <= 5 {
		(*t.client).LPush(queueName, keyword).Result()
		fmt.Println("times value is :", value, keyword)
		_, err := (*t.client).Set(key, value, 86400*time.Second).Result()
		if err != nil {
			fmt.Println("set key error", err)
		}
		//关键词补偿获取的开始url地址
		numstr := strconv.Itoa(num)
		urlkey := keyword + "_url" + date
		urlvalue := url + "|" + numstr
		_, err1 := (*t.client).Set(urlkey, urlvalue, 86400*time.Second).Result()
		if err1 != nil {
			fmt.Println("set urlkey error:", err1)
		}
	}
}

func (t *RedisEngine) GetTimes(keyword string) int {
	date := time.Now().Format("0102")
	key := keyword + "_time" + date
	timesStr := (*t.client).Get(key).Val()
	if len(timesStr) == 0 {
		return 0
	}
	times, _ := strconv.Atoi(timesStr)
	return times
}

func (t *RedisEngine) GetUrl(keyword string) (string, int) {
	date := time.Now().Format("0102")
	urlkey := keyword + "_url" + date
	urlvalue := (*t.client).Get(urlkey).Val()
	if len(urlvalue) == 0 {
		return "", 0
	}
	urlArr := strings.Split(urlvalue, "|")
	num := 0
	if len(urlArr) == 2 {
		num, _ = strconv.Atoi(urlArr[1])
	}
	url := urlArr[0]
	return url, num
}
