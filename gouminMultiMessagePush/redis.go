package main

import (
	"fmt"
	redis "gopkg.in/redis.v4"
	"strings"
)

func connect(conn string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}

func testLlen(client *redis.Client) {
	len := (*client).LLen("MessageCenter").Val()
	fmt.Println("len:", len)
	if len > numForOneLoop {
		taskNum = numForOneLoop
	} else {
		taskNum = int(len)
	}
}

func croutinePopRedis(c chan int, client *redis.Client, i int) {
	redisStr := (*client).LPop("MessageCenter").Val()
	redisStrArr := strings.Split(redisStr, "^")
	tasks[i].pushStr = redisStrArr[0]
	tasks[i].insertStr = redisStrArr[1]

	c <- 1
}

func lopDataFromRedis(client *redis.Client) {
	c := make(chan int, taskNum)
	for i := 0; i < taskNum; i++ {
		go croutinePopRedis(c, client, i)
	}

	for i := 0; i < taskNum; i++ {
		<-c
	}
}

func loadDataFromRedis() {
	client := connect(c.redisConn)
	testLlen(client)
	fmt.Println(taskNum)
	lopDataFromRedis(client)
	client.Close()
}
