package main

import (
	"stathat.com/c/jconfig"
	"time"
)

type Config struct {
	httpTimeOut time.Duration
	currentNum  int
	redisConn   string
	//mongoConn   string
	elasticConn string
	log string
}

type redisData struct {
	pushStr   string
	insertStr string
	times     int
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/msgConfig.json")
	c.httpTimeOut = time.Duration(config.GetInt("httpTimeOut"))
	c.currentNum = config.GetInt("currentNum")
	c.redisConn = config.GetString("redisConn")
	//c.mongoConn = config.GetString("mongoConn")
	c.elasticConn = config.GetString("elasticConn")
	c.log = config.GetString("log")
}
