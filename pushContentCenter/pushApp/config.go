package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn       string
	dbDsn3       string
	dbName      string
	dbName1     string
	dbName2     string
	dbName3     string
	dbName4     string
	dbName5     string
	dbAuth      string
	redisConn   string
	coroutinNum int
	queueName   string
	//mongoConn   string
	elkNodes string
	log     string
}

func loadConfig() {
	//@todo change online path
	config := jconfig.LoadConfig("/etc/pushContentCenterConfig.json")
	//config := jconfig.LoadConfig("/Users/Snow/Work/go/config/pushContentCenterConfig.json")
	c.dbDsn = config.GetString("dbDsn")
	c.dbDsn3 = config.GetString("dbDsn3")
	c.dbName = config.GetString("dbName")
	c.dbName1 = config.GetString("dbName1")
	c.dbName2 = config.GetString("dbName2")
	c.dbName3 = config.GetString("dbName3")
	c.dbName4 = config.GetString("dbName4")
	c.dbName5 = config.GetString("dbName5")
	c.dbAuth = config.GetString("dbAuth")
	c.redisConn = config.GetString("redisConn")
	c.coroutinNum = config.GetInt("coroutinNum")
	c.queueName = config.GetString("queueName")
	//c.mongoConn = ""//config.GetString("mongoConn")
	c.elkNodes = config.GetString("elkNodes")
	c.log = config.GetString("log")
}
