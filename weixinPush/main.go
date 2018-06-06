package main

import (
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/tools"
	redis "gopkg.in/redis.v4"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
	"github.com/jackson198608/goProject/weixinPush/task"
	//"gouminGitlab/common/weixin/accessToken/accesstokenManager"
	"sync"
	"gouminGitlab/common/weixin/accessToken"
	"strings"
)


var c Config = Config{
	"127.0.0.1:6379",       //redis info
	1,                      //thread num
	"weixinPush",    //queuename
	"appSecret",  // app and secret of every program
	}

//var weixinAccessTokens * accesstokenManager.Manager

var w sync.Once


func init() {
	loadConfig()
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

	token := getToken(t.Appid,c.appSecret,redisConn)

	err = t.Do(redisConn,token)
	if err != nil {
		return err
	}
	return err
}

func getToken(appid string,appSecret string,redisConn *redis.ClusterClient) string {
	//获取配置里面的appsecret(appid、secret)
	appids := make(map[string]string)
	a := strings.Split(appSecret,",")
	for i:=0; i<len(a);i++  {
		appid := strings.Split(a[i],":")[0]
		secret := strings.Split(a[i],":")[1]
		appids[appid] = secret
	}
	secret := appids[appid]
	var token = ""
	i := 0
	for i <= 3 {
		w.Do(func() {
			accesstoken := accessToken.NewAccessToken(appid, secret, redisConn)
			token = accesstoken.GetToken(redisConn)
		})
		if token != "" {
			break
		}
		i++
	}
	return token
}

