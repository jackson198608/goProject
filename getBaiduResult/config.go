package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn     string
	dbName    string
	dbAuth    string
	numloops  int
	redisConn string
	queueName string
	logLevel  int
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/baiduRankConfig.json")
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.numloops = config.GetInt("numloops")
	c.redisConn = config.GetString("redisConn")
	c.queueName = config.GetString("queueName")
	c.logLevel = config.GetInt("logLevel")
}
