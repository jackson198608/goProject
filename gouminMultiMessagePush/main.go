package main

import (
	//"github.com/jackson198608/gotest/appPush"
	//"github.com/jackson198608/gotest/inMongo"
	"stathat.com/c/jconfig"
)

//const
const numForOneLoop = 1000

//define the config var
var c Config = Config{"127.0.0.1:6379", "127.0.0.1:27017"}

//define the tasks array for each loop
var tasks [numForOneLoop]redisData
var taskNum int = 0

func loadConfig() {
	config := jconfig.LoadConfig("/etc/msgConfig.json")
	c.redisConn = config.GetString("redisConn")
	c.mongoConn = config.GetString("mongoConn")
}

func main() {
	loadConfig()
	//testCreateTestData()

	//main loop
	loadDataFromRedis()
	insertMongo()
	//pushMesssage()
	//clearTaskData()

	//testPrintTasks()
}
