package main

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/redisLoopTask"
	"time"
)

var c Config = Config{
	100,
	"127.0.0.1:6379",
	"MallSendPhone",
	0}

func do() {
	r := redisLoopTask.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops)
	for {
		logger.Info("going to do one loop")
		r.Loop()

		logger.Info("going to sleep for 5 second")
		time.Sleep(5 * time.Second)

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
