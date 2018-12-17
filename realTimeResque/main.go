package main

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/jackson198608/goProject/common/tools"
	"github.com/jackson198608/goProject/realTimeResque/task"
	"gopkg.in/redis.v4"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
)

var c Config = Config{
	"127.0.0.1:6379", //redis info
	1,                //thread num
	"imageCompressTaskList", //queue name
	"127.0.0.1"}            //php server ip

func init() {
	loadConfig()
}

func main() {
	var mongoConnInfo []string
	var mysqlInfo []string

	redisInfo := tools.FormatRedisOption(c.redisConn)
	logger.Info("start work")
	r, err := redisEngine.NewRedisEngine(c.queueName, &redisInfo, mongoConnInfo, mysqlInfo, c.coroutinNum, 1, jobFuc, c.phpServerIp)
	if err != nil {
		logger.Error("[NewRedisEngine] ", err)
	}

	err = r.Do()
	if err != nil {
		logger.Error("[redisEngine Do] ", err)
	}

}

func jobFuc(job string, redisConn *redis.ClusterClient, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {
	t, err := task.NewTask(job, taskarg[0])
	if err != nil {
		return err
	}
	err = t.Do()
	if err != nil {
		return err
	}
	return err
}
