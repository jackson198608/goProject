package main

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	dbDsn         string
	dbName        string
	dbAuth        string
	numloops      int
	lastId        int
	firstId       int
	redisConn     string
	queueName     string
	logFile       string
	logLevel      int
	fansLimit     string
	dateLimit     string
	currentNum    string
	mongoConn     string
	mongoDb       string
	followFirstId int
	followLastId  int
	eventLimit    string
	redisStart    string
	redisEnd      string
	pushLimit     string
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/moveEventLogConfigNew.json")
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.numloops = config.GetInt("numloops")
	c.lastId = config.GetInt("lastId")
	c.firstId = config.GetInt("firstId")
	c.redisConn = config.GetString("redisConn")
	c.queueName = config.GetString("queueName")
	c.logFile = config.GetString("logFile")
	c.logLevel = config.GetInt("logLevel")
	c.fansLimit = config.GetString("fansLimit")
	c.dateLimit = config.GetString("dateLimit")
	c.currentNum = config.GetString("currentNum")
	c.mongoConn = config.GetString("mongoConn")
	c.mongoDb = config.GetString("mongoDb")
	c.followFirstId = config.GetInt("followFirstId")
	c.followLastId = config.GetInt("followLastId")
	c.eventLimit = config.GetString("eventLimit")
	c.redisStart = config.GetString("redisStart")
	c.redisEnd = config.GetString("redisEnd")
	c.pushLimit = config.GetString("pushLimit")
}
