package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn         string
	dbName        string
	dbAuth        string
	numloops      int
	redisConn     string
	queueName     string
	logFile       string
	logLevel      int
	mongoConn     string
	mongoDb       string
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/recommendPositionConfig.json")
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.numloops = config.GetInt("numloops")
	c.redisConn = config.GetString("redisConn")
	c.queueName = config.GetString("queueName")
	c.logFile = config.GetString("logFile")
	c.logLevel = config.GetInt("logLevel")
	c.mongoConn = config.GetString("mongoConn")
	c.mongoDb = config.GetString("mongoDb")
}
