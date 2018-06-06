package main

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/jackson198608/goProject/common/tools"
	"github.com/jackson198608/goProject/weixinPush/task"
	"gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v4"
	//"gouminGitlab/common/weixin/accessToken/accesstokenManager"
	"gouminGitlab/common/weixin/accessToken"
	"strings"
	"sync"
)

var c Config = Config{
	"127.0.0.1:6379", //redis info
	1,                //thread num
	"weixinPush",     //queuename
	"appSecret",      // app and secret of every program
}

//var weixinAccessTokens * accesstokenManager.Manager

var w sync.Once

func init() {
	initAppids()
	loadConfig()
}

func initAppids() {
	appids := make(map[string]string)
	a := strings.Split(appSecret, ",")
	for i := 0; i < len(a); i++ {
		appid := strings.Split(a[i], ":")[0]
		secret := strings.Split(a[i], ":")[1]
		appids[appid] = secret
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

func jobFuc(job string, redisConn *redis.ClusterClient, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {

	t, err := task.NewTask(job, mysqlConns, mgoConns)
	if err != nil {
		return err
	}

	token := getToken(t.Appid, c.appSecret, redisConn)

	err = t.Do(redisConn, token)
	if err != nil {
		return err
	}
	return err
}

func getToken(appid string, appSecret string, redisConn *redis.ClusterClient) string {
	//using appid to find the redis key of this appid's accessToken
	key := appid
	var token string

	//check it the token exist
	if(!redisConn.Exists(key).Val()){
		token=w.Do(generateToken(appid,redisConn))
	} else{
		token=redisConn.Get(key).Val()
	}

	return token
}

func generateToken(appid string,redisConn *redis.ClusterClient)) string {
	//@todo
	//using appid to search var appids to find the secret
	//using appid and secret to get token and set back into redis
	return token
}
