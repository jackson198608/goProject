package main

import (
	"fmt"
	redis "gopkg.in/redis.v4"
	"strconv"
)

func main() {
	fmt.Println("vim-go")
}

func croutinePopRedisTidData(c chan int, client *redis.Client, i int) {
	fmt.Println("[notice] pop splitPid")
	redisStr := (*client).LPop("splitPid").Val()
	if redisStr == "" {
		fmt.Println("[notice] got nothing")
		c <- 1
		return
	}

	tid, err := strconv.Atoi(redisStr)

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

func testLlen(client *redis.Client) {
	len := (*client).LLen(redisQueueName).Val()
	if int(len) > numForOneLoop {
		taskNum = numForOneLoop
	} else {
		taskNum = int(len)
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

func connect(conn string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("[Error] redis connect error")
	}
	return client
}
