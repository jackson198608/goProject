package main

import (
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/tools"
	redis "gopkg.in/redis.v4"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
	"github.com/jackson198608/goProject/weixinPush/task"
	"gouminGitlab/common/weixin/accessToken/accesstokenManager"
	"sync"
	"errors"
)

//var weixinAccessTokens * accesstokenManager.Manager

var c Config = Config{
	"127.0.0.1:6379",       //redis info
	1,                      //thread num
	"weixinPush",    //queuename
	"appSecret",  // app and secret of every program
	}

var weixinAccessTokens * accesstokenManager.Manager

var w sync.Once

//type wx struct {
//	once sync.Once
//	weixinAccessTokens * accesstokenManager.Manager
//}

func init() {
	loadConfig()
	initWeixinaccessToken()
}
func initWeixinaccessToken() {
	weixinAccessTokens = accesstokenManager.NewManager(c.appSecret)
	if weixinAccessTokens == nil{
		errors.New("weixinAccessTokens is empty")
	}
}

func main() {
	var mongoConnInfo []string
	var mysqlInfo []string

	redisInfo := tools.FormatRedisOption(c.redisConn)
	logger.Info("start work")
	r, err := redisEngine.NewRedisEngine(c.queueName, &redisInfo, mongoConnInfo, mysqlInfo, c.coroutinNum, 1, jobFuc)
	if err != nil {
		logger.Error("[NewRedisEngine] ", err)
	}

	err = r.Do()
	if err != nil {
		logger.Error("[redisEngine Do] ", err)
	}
}

func jobFuc(job string,redisConn *redis.ClusterClient, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {

	t, err := task.NewTask(job, mysqlConns, mgoConns)
	if err != nil {
		return err
	}
	w.Do(func() {
		weixinAccessTokens.GetTokens(redisConn)
	})

	err = t.Do(weixinAccessTokens)
	if err != nil {
		return err
	}
	return err
}

