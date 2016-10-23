package main

import (
	"stathat.com/c/jconfig"
	"time"
)

type Config struct {
	httpTimeOut time.Duration
	currentNum  int
	redisConn   string
	mongoConn   string
}

type redisData struct {
	pushStr   string
	insertStr string
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/msgConfig.json")
	c.httpTimeOut = time.Duration(config.GetInt("httpTimeOut"))
	c.currentNum = config.GetInt("currentNum")
	c.redisConn = config.GetString("redisConn")
	c.mongoConn = config.GetString("mongoConn")
	numForOneLoop = c.currentNum
	timeout = c.httpTimeOut
}
