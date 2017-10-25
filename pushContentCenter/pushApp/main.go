package main

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/jackson198608/goProject/pushContentCenter/task"
	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v4"
	"reflect"
)

var c Config = Config{
	"192.168.86.193:3307",  //mysql dsn
	"new_dog123",           //mysql dbName
	"dog123:dog123",        //mysqldbAuth
	"127.0.0.1:6379",       //redis info
	1,                      //thread num
	"pushContentCenter",    //queuename
	"192.168.86.192:27017"} // mongo

func init() {
	loadConfig()
}

func main() {
	var mongoConnInfo []string
	mongoConnInfo = append(mongoConnInfo, c.mongoConn)
	var mysqlInfo []string
	mysqlInfo = append(mysqlInfo, c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")

	redisInfo := redis.Options{
		Addr: c.redisConn,
	}
	logger.Info("start work")
	r, err := redisEngine.NewRedisEngine(c.queueName, &redisInfo, mongoConnInfo, mysqlInfo, c.coroutinNum, jobFuc)
	if err != nil {
		logger.Error("[NewRedisEngine] ", err)
	}

	err = r.Do()
	if err != nil {
		logger.Error("[redisEngine Do] ", err)
	}
}

func jobFuc(job string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {
	t, err := task.NewTask(job, mysqlConns, mgoConns)
	if err != nil {
		logger.Error("[NewTask]", err)
	}
	fmt.Println(mysqlConns)
	err = t.Do()
	if err != nil {
		logger.Error("[task Do]", err)
	}
	return err
}
