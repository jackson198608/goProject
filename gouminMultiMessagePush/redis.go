package main

import (
	"fmt"
	redis "gopkg.in/redis.v4"
	"log"
	"strconv"
	"strings"
)

func putFailOneBack(i int) {

	if tasks[i].times == 5 {
		log.Println("[Error] fail exceed,drop: ", tasks[i].pushStr, tasks[i].insertStr)
		return
	}
	client := connect(c.redisConn)
	var pushStr string = ""
	if (jobType == "multi") || (jobType == "single") {
		pushStr = tasks[i].pushStr
	} else {
		pushStr = tasks[i].insertStr
	}

	//add times into it
	newTimes := tasks[i].times + 1
	pushStr = pushStr + "^" + strconv.Itoa(newTimes)

	err := (*client).RPush(redisQueueName, pushStr).Err()
	if err != nil {
		log.Println("[Error] push str into redis error:  ", pushStr)
	}

	client.Close()

}

func connect(conn string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Println("[Error] redis connect error")
	}
	return client
}

func testLlen(client *redis.Client) {
	len := (*client).LLen(redisQueueName).Val()
	if int(len) > numForOneLoop {
		taskNum = numForOneLoop
	} else {
		taskNum = int(len)
	}
}

func croutinePopRedisMultiData(c chan int, client *redis.Client, i int) {
	log.Println("[notice] pop mcMulti")
	redisStr := (*client).LPop("mcMulti").Val()
	if redisStr == "" {
		log.Println("[notice] got nothing")
		c <- 1
		return
	}
	redisArr := strings.Split(redisStr, "^")
	tasks[i].pushStr = redisArr[0]
	tasks[i].insertStr = ""
	if len(redisArr) == 2 {
		tasks[i].times, _ = strconv.Atoi(redisArr[1])
	} else {
		tasks[i].times = 1
	}

	c <- 1
}

func lopMulti(client *redis.Client) {
	c := make(chan int, taskNum)
	for i := 0; i < taskNum; i++ {
		go croutinePopRedisMultiData(c, client, i)
	}

	for i := 0; i < taskNum; i++ {
		<-c
	}
}

func croutinePopRedisSingleData(c chan int, client *redis.Client, i int) {
	redisStr := (*client).LPop("mcSingle").Val()
	if redisStr == "" {
		log.Println("[notice] got nothing")
		c <- 1
		return
	}
	redisArr := strings.Split(redisStr, "^")
	tasks[i].pushStr = redisArr[0]
	tasks[i].insertStr = ""
	if len(redisArr) == 2 {
		tasks[i].times, _ = strconv.Atoi(redisArr[1])
	} else {
		tasks[i].times = 1
	}

	c <- 1
}

func lopSingle(client *redis.Client) {
	c := make(chan int, taskNum)
	for i := 0; i < taskNum; i++ {
		go croutinePopRedisSingleData(c, client, i)
	}

	for i := 0; i < taskNum; i++ {
		<-c
	}
}

func croutinePopRedisInsertData(c chan int, client *redis.Client, i int) {
	redisStr := (*client).LPop("mcInsert").Val()
	if redisStr == "" {
		log.Println("[notice] got nothing")
		c <- 1
		return
	}
	redisArr := strings.Split(redisStr, "^")
	tasks[i].insertStr = redisArr[0]
	tasks[i].pushStr = ""
	if len(redisArr) == 2 {
		tasks[i].times, _ = strconv.Atoi(redisArr[1])
	} else {
		tasks[i].times = 1
	}

	c <- 1
}

func lopInsert(client *redis.Client) {
	c := make(chan int, taskNum)
	for i := 0; i < taskNum; i++ {
		go croutinePopRedisInsertData(c, client, i)
	}

	for i := 0; i < taskNum; i++ {
		<-c
	}
}

func loadDataFromRedis() {
	client := connect(c.redisConn)
	testLlen(client)
	fmt.Println(taskNum)
	switch jobType {
	case "multi":
		lopMulti(client)
	case "single":
		lopSingle(client)
	case "insert":
		lopInsert(client)

	default:
		fmt.Println("[notice] no use to get data from redis")
	}

	client.Close()
}
