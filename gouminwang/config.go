package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn          string
	dbName         string
	dbAuth         string
	redisConn      string
	redisQueueName string
	saveDir        string
	imgUrl         string
	threadNum      int
	logPath        string
	startUrl       string
	startUrlTag    string
}

func loadConfig(args []string) {
	var config *jconfig.Config

	if len(args) >= 2 {
		config = jconfig.LoadConfig(args[1])
	} else {
		config = jconfig.LoadConfig("/etc/gouminwang.json")
	}
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.redisConn = config.GetString("redisConn")
	c.redisQueueName = config.GetString("redisQueueName")
	c.saveDir = config.GetString("saveDir")
	c.imgUrl = config.GetString("imgUrl")
	c.threadNum = config.GetInt("threadNum")
	c.logPath = config.GetString("logPath")
	c.startUrl = config.GetString("startUrl")
	c.startUrlTag = config.GetString("startUrlTag")
}
