package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn       string
	dbName      string
	dbAuth      string
	redisConn   string
	coroutinNum int
	queueName   string
	mongoConn   string
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/pushContentCenterConfig.json")
	// config := jconfig.LoadConfig("/tmp/pushContentCenterConfig.json")
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.redisConn = config.GetString("redisConn")
	c.coroutinNum = config.GetInt("coroutinNum")
	c.queueName = config.GetString("queueName")
	c.mongoConn = config.GetString("mongoConn")
}
