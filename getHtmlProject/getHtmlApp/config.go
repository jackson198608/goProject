package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn     string
	dbName    string
	dbAuth    string
	sdbDsn    string //5.86
	sdbName   string //5.86
	sdbAuth   string //5.86
	numloops  int
	redisConn string
	queueName string
	saveDir   string
	logFile   string
	logLevel  int
	tidStart  string
	tidEnd    string
	domain    string
	host      string
}

func loadConfig(args []string) {
	var config *jconfig.Config

	if len(args) >= 3 {
		config = jconfig.LoadConfig(args[2])
	} else {
		config = jconfig.LoadConfig("/etc/configask.json")
	}
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.sdbDsn = config.GetString("sdbDsn")
	c.sdbName = config.GetString("sdbName")
	c.sdbAuth = config.GetString("sdbAuth")
	c.numloops = config.GetInt("numloops")
	c.redisConn = config.GetString("redisConn")
	c.queueName = config.GetString("queueName")
	c.saveDir = config.GetString("saveDir")
	c.logFile = config.GetString("logFile")
	c.logLevel = config.GetInt("logLevel")
	c.tidStart = config.GetString("tidStart")
	c.tidEnd = config.GetString("tidEnd")
	c.domain = config.GetString("domain")
	c.host = config.GetString("host")
}
