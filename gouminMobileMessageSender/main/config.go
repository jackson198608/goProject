package main

import (
	"github.com/donnie4w/go-logger/logger"
	"stathat.com/c/jconfig"
)

type Config struct {
	numloops  int
	redisConn string
	queueName string
	logLevel  logger.LEVEL
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/sendPhoneMessage.json")
	c.numloops = config.GetInt("numloops")
	c.redisConn = config.GetString("redisConn")
	c.queueName = config.GetString("queueName")
	c.logLevel = logger.LEVEL(config.GetInt("logLevel"))
}
