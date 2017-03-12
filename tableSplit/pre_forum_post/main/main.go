package main

import (
	"github.com/jackson198608/goProject/redisLoopTask"
)

const dbDsn = "210.14.154.198:3306"
const dbName = "new_dog123"
const dbAuth = "dog123:dog123"
const numloops = 20
const lastTid = 2731136250
const firstTid = 0
const redisConn = "127.0.0.1:6379"

func main() {
	r := redisLoopTask.NewRedisEngine("movePost", redisConn, "", 0, 100, dbAuth, dbDsn, dbName)
	tids := getTask(10)
	r.PushTaskData(tids)
	r.Loop()

}
