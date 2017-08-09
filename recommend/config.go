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
	queueName1    string
	logFile       string
	logLevel      int
	mongoConn     string
	mongoDb       string
	mongoConn1    string
	mongoDb1      string
	pushLimit     string
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/recommendConfig.json")
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.numloops = config.GetInt("numloops")
	c.redisConn = config.GetString("redisConn")
	c.queueName = config.GetString("queueName")
	c.queueName1 = config.GetString("queueName1")
	c.logFile = config.GetString("logFile")
	c.logLevel = config.GetInt("logLevel")
	c.mongoConn = config.GetString("mongoConn")
	c.mongoDb = config.GetString("mongoDb")
	c.mongoConn1 = config.GetString("mongoConn1")
	c.mongoDb1 = config.GetString("mongoDb1")
	c.pushLimit = config.GetString("pushLimit")
}
