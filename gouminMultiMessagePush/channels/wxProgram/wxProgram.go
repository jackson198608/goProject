package wxProgram

import (
	"net/http"
	"fmt"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"github.com/donnie4w/go-logger/logger"
	"encoding/json"
	"errors"
)

type Task struct {
	AppId string
	TaskJson string  //发送的内容
	AccessToken string
}

const _sendUrl = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?"


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
				return errors.New(body)
			}
		}
		logger.Info("request weixin success",body)
	}
	return nil
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