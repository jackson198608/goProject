package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	redisConn   string
	coroutinNum int
	queueName   string
	appSecret  string
	workTime     string
}

func loadConfig() {
	//@todo change online path
	config := jconfig.LoadConfig("/etc/weixinPushConfig.json")
	c.redisConn = config.GetString("redisConn")
	c.coroutinNum = config.GetInt("coroutinNum")
	c.queueName = config.GetString("queueName")
	c.appSecret = config.GetString("appSecret")
	c.workTime = config.GetString("workTime")
}
