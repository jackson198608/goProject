package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn       string
	dbName      string
	dbName1      string
	dbAuth      string
	redisConn   string
	coroutinNum int
	queueName   string
	mongoConn   string
	elkNodes string
}

func loadConfig() {
	//@todo change online path
	config := jconfig.LoadConfig("/etc/pushContentCenterConfig.json")
	//config := jconfig.LoadConfig("/Users/Snow/Work/go/config/pushContentCenterConfig.json")
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbName1 = config.GetString("dbName1")
	c.dbAuth = config.GetString("dbAuth")
	c.redisConn = config.GetString("redisConn")
	c.coroutinNum = config.GetInt("coroutinNum")
	c.queueName = config.GetString("queueName")
	c.mongoConn = config.GetString("mongoConn")
	c.elkNodes = config.GetString("elkNodes")
}
