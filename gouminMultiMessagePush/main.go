package main

import (
	"fmt"
	"github.com/jackson198608/gotest/appPush"
	"io/ioutil"
	"log"
	"os"
	"time"
)

//define the config var
var c Config = Config{5, 100, "127.0.0.1:6379", "127.0.0.1:27017"}
var numForOneLoop int = c.currentNum
var p12Bytes []byte
var timeout time.Duration = c.httpTimeOut

//define the tasks array for each loop
var tasks []redisData
var taskNum int = 0

func Init() {
	cBytes, err := ioutil.ReadFile("/etc/pro-lingdang.pem")
	if err != nil {
		log.Fatal("[Error] read cert file error")
		return
	}

	p12Bytes = cBytes

	loadConfig()

	appPush.Init(timeout)
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
			log.Println("[Notice] sleep for", 5, " second")
			time.Sleep(5 * time.Second)
		}

		//insert data into mongo
		insertMongo()

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
			log.Println("[Notice] sleep for", 5, " second")
			time.Sleep(5 * time.Second)
		}

		//push data to app
		push()
	}

}

func main() {

	//init the system process
	Init()

	switch os.Args[1] {
	case "test":
		fmt.Println("[test]")
		testCreateTestData()

	case "multi":
		fmt.Println("[multi-push]")
		pushCenter()
	case "single":
		fmt.Println("[single-push]")
		singlePush()

	default:
		fmt.Println("do nothing")
	}

}
