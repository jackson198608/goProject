package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn     string
	dbName    string
	dbAuth    string
	numloops  int
	lastTid   int
	firstTid  int
	redisConn string
	queueName string
	logFile   string
	logLevel  int
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/movePostConfig.json")
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.numloops = config.GetInt("numloops")
	c.lastTid = config.GetInt("lastTid")
	c.firstTid = config.GetInt("firstTid")
	c.redisConn = config.GetString("redisConn")
	c.queueName = config.GetString("queueName")
	c.logFile = config.GetString("logFile")
	c.logLevel = config.GetInt("logLevel")
}
