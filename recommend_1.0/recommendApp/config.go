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
	elkDsn      string
}

func loadConfig() {
	//@todo change online path
	config := jconfig.LoadConfig("/etc/userRecommendConfig.json")
	// config := jconfig.LoadConfig("//Users/Snow/recommendUserConfig.json")
	c.dbDsn = config.GetString("dbDsn")
	c.dbName = config.GetString("dbName")
	c.dbAuth = config.GetString("dbAuth")
	c.redisConn = config.GetString("redisConn")
	c.coroutinNum = config.GetInt("coroutinNum")
	c.queueName = config.GetString("queueName")
	c.mongoConn = config.GetString("mongoConn")
	c.elkDsn = config.GetString("elkDsn")
}
