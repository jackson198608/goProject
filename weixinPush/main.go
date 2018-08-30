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
	//"gouminGitlab/common/weixin/accessToken"
	"encoding/json"
	"fmt"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var c Config = Config{
	"127.0.0.1:6379", //redis info
	1,                //thread num
	"weixinPush",     //queuename
	"appSecret",      // app and secret of every program
	"workTime",       //工作时间，该时间段内取任务
}

//var weixinAccessTokens * accesstokenManager.Manager
const _tokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&"
const accesstoken_key = "card_access_token_"

var w sync.Once

var appids map[string]string

func init() {
	loadConfig()
	initAppids()
}

func initAppids() {
	appids = make(map[string]string)
	a := strings.Split(c.appSecret, ",")
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

	//check if it is work time
	for {
		if isWorkTime() {
			break
		}
		time.Sleep(10 * time.Minute)
	}

	t, err := task.NewTask(job, mysqlConns, mgoConns)
	if err != nil {
		return err
	}

	token := getToken(t.Appid, redisConn)

	err = t.Do(redisConn, token)
	if err != nil {
		return err
	}

	return err
}

func isWorkTime() bool {
	//该时间段内不发送任务
	work := c.workTime

	rawSlice := []byte(work)
	rawLen := len(rawSlice)
	lastIndex := strings.LastIndex(work, "|")
	start, err := strconv.Atoi(string(rawSlice[0:lastIndex]))
	if err != nil {
		start = int(8)
	}
	end, err1 := strconv.Atoi(string(rawSlice[lastIndex+1 : rawLen]))
	if err1 != nil {
		end = int(22)
	}
	wtime := time.Now().Hour()
	if wtime >= start && wtime < end {
		return true
	}

	return false
}

func getToken(appid string, redisConn *redis.ClusterClient) string {
	//using appid to find the redis key of this appid's accessToken
	key := accesstoken_key + appid
	var token string

	//check it the token exist
	token = redisConn.Get(key).Val()
	if token == "" {
		w.Do(func() {
			token = generateToken(appid, redisConn)
		})
	}

	return token
}

func generateToken(appid string, redisConn *redis.ClusterClient) string {
	//@todo
	//using appid to search var appids to find the secret
	secret := appids[appid]
	//using appid and secret to get token and set back into redis
	target := _tokenUrl + "appid=" + appid + "&secret=" + secret
	logger.Error("[generateToken] target:", target)
	var h http.Header = make(http.Header)
	abuyun := getAbuyun()
	defer abuyun.Close()
	statusCode, _, body, err := abuyun.SendRequest(target, h, "", true)
	if err != nil {
		return err.Error()
	}
	if statusCode == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(body), &result); err == nil { //解析json，获取accesstoken
			token := fmt.Sprintf("%s", result["access_token"]) //access_token默认2个小时的有效期
			//更新缓存
			key := accesstoken_key + appid
			redisConn.Set(key, token, time.Duration((7200-20)*time.Second))
			return token
		}
	}
	return ""
}

func getAbuyun() *abuyunHttpClient.AbuyunProxy {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy("", "", "")

	if abuyun == nil {
		fmt.Println("create abuyun error")
		return nil
	}
	return abuyun
}
