package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	redisConn   string
	coroutinNum int
	queueName   string
	phpServerIp string
}

func loadConfig() {
	//@todo change online path
	config := jconfig.LoadConfig("/etc/compressImageConfig.json")
	// config := jconfig.LoadConfig("/Users/Snow/compressImageConfig.json")
	c.redisConn = config.GetString("redisConn")
	c.coroutinNum = config.GetInt("coroutinNum")
	c.queueName = config.GetString("queueName")
	c.phpServerIp = config.GetString("phpServerIp")
}
