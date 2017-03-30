package main

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/redisLoopTask"
	"time"
)

var c Config = Config{
	100,
	"127.0.0.1:6379",
	"movePost",
	0}

func do() {
	r := redisLoopTask.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops)
	for {
		r.Loop()
		time.Sleep(5)

	}
}

func Init() {
	loadConfig()
	logger.SetLevel(c.logLevel)

}

func main() {
	Init()
	do()
}
