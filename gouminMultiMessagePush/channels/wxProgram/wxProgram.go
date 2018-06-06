package wxProgram

import (
	"net/http"
	"fmt"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"github.com/donnie4w/go-logger/logger"
	"encoding/json"
	"errors"
	redis "gopkg.in/redis.v4"
	//"gouminGitlab/common/weixin/accessToken"
	"time"
)

type Task struct {
	AppId string
	Secret string
	TaskJson string  //发送的内容
	AccessToken string
	//Redisconn *redis.ClusterClient
}

const _sendUrl = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?"
const accesstoken_key = "card_access_token_"
const pushKey  = "weixinPush"

const _tokenUrl = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&"


/**
实例化
 */
func NewTask(appid string,jobStr string,accessToken string) (t *Task){
	var tR Task

	tR.AppId = appid
	tR.TaskJson = jobStr
	tR.AccessToken = accessToken
	return &tR
}
/**
发送请求
 */
func (p *Task) SendRequest() error {
	if p.AppId == "" {
		logger.Error("[appid empty] ", p.AppId, p.TaskJson)
		return nil
	}

	//请求微信
	err := p.requestWeixin(p.AccessToken,p.TaskJson)
	if err != nil {
		logger.Error("[request weixin fail] error: ", err, " -appid:",p.AppId, " -jobStr:",p.TaskJson)
		return err
	}
	return nil
}
/**
do request weixin
 */
func (p *Task)requestWeixin(accesstoken string,messqge string) error{
	var target = _sendUrl+"access_token="+accesstoken
	var h http.Header = make(http.Header)
	abuyun := p.getAbuyun()
	statusCode, _, body, err := abuyun.SendPostRequest(target,h,messqge,true)
	if err != nil {
		return err
	}
	if statusCode == 200 {
		var result map[string]interface{}
		if err:=json.Unmarshal([]byte(body),&result);err==nil{
			errcode := result["errcode"]
			if errcode != float64(0) {
				//if accesstoken 失效
				if errcode == float64(40001) { //40001:accesstoken 失效导致请求微信失败
					//删除缓存,并把该失败任务重新添加到任务中
					//cacheErr := p.delAccessTokenCache(redisConn)
					//logger.Info("update cache",cacheErr)
				}
				return errors.New(body)
			}
		}
		logger.Info("request weixin success",body)
	}
	return nil
}
/**
accesstoken失效，检查缓存是否过期，如果没有过期，删除缓存，并重建重建任务
 */
func (p *Task)delAccessTokenCache(redisConn *redis.ClusterClient) error{
	accesstokenKey := accesstoken_key+p.AppId
	cache := redisConn.Get(accesstokenKey).Val()
	if cache != "" {
		delerror := redisConn.Del(accesstokenKey).Err()
		if delerror == nil {
			time.Sleep(1*time.Second)
			newTask := p.AppId+"|"+p.TaskJson
			err := redisConn.LPush(pushKey,newTask).Err()
			if err != nil {
				return err
			}
		}

	}
	return  nil
}

/**
get secret form config by appid
 */
func GetAppSecret(appId string) string{
	configs := make(map[string]string)
	secret := configs[appId]
	if secret != "" {
		return secret
	}
	return ""
}

func (p *Task) getAbuyun() *abuyunHttpClient.AbuyunProxy {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy("", "", "")

	if abuyun == nil {
		fmt.Println("create abuyun error")
		return nil
	}
	return abuyun
}