package main

import (
	"github.com/jackson198608/goProject/appPush"
	"io/ioutil"
	"os"
	"time"
)

//define the config var
var c Config = Config{5, 100, "192.168.86.68:6380,192.168.86.68:6381,192.168.86.68:6382,192.168.86.68:6383,192.168.86.68:6384,192.168.86.68:6385", "127.0.0.1:27017","http://192.168.86.230:9200,http://192.168.86.231:9200"}
var numForOneLoop int = c.currentNum
var p12Bytes []byte
var timeout time.Duration = c.httpTimeOut

/*
type of job
	multi: multi app push
	single: single app push
	insert: insert data into mongo
*/
var jobType string = "multi"
var redisQueueName string = "mcMulti"

//define the tasks array for each loop
var tasks []redisData
var taskNum int = 0

func Init() {
	getRedisQueueName()
	cBytes, err := ioutil.ReadFile("/etc/pro-lingdang.pem")
	if err != nil {
		return
	}
	p12Bytes = cBytes

	loadConfig()

	appPush.Init(timeout)
}

func getRedisQueueName() {
	switch os.Args[1] {
	case "multi":
		redisQueueName = "mcMulti"
	case "single":
		redisQueueName = "mcSingle"
	case "insert":
		redisQueueName = "mcInsert"

	default:
		redisQueueName = "mcMulti"
	}
}

func pushCenter() {
	//main loop
	for {
		//init data
		tasks = make([]redisData, numForOneLoop, numForOneLoop)
		taskNum = 0

		//load data from redis
		loadDataFromRedis()

		//if there is no data
		if taskNum == 0 {
			time.Sleep(5 * time.Second)
		}

		//push data to app
		push()
	}

}

func singlePush() {
	//main loop
	for {
		//init data
		tasks = make([]redisData, numForOneLoop, numForOneLoop)
		taskNum = 0

		//load data from redis
		loadDataFromRedis()

		//if there is no data
		if taskNum == 0 {
			time.Sleep(5 * time.Second)
		}

		//push data to app
		push()
	}

}

func onlyInsertMongo() {
	//main loop
	for {
		//init data
		tasks = make([]redisData, numForOneLoop, numForOneLoop)
		taskNum = 0

		//load data from redis
		loadDataFromRedis()

		//if there is no data
		if taskNum == 0 {
			time.Sleep(5 * time.Second)
		}

		//insert data into mongo
		insertMongo()
	}

}

func main() {

	//init the system process
	Init()
	jobType = os.Args[1]

	switch jobType {
	case "test":
		testCreateTestData()

	case "multi":
		pushCenter()

	case "single":
		singlePush()

	case "insert":
		onlyInsertMongo()

	default:
	}

}
