package huawei

import (
	"bytes"
	"encoding/json"
	"github.com/jackson198608/goProject/ActiveRecord"
	log "github.com/thinkboy/log4go"
	"gopkg.in/redis.v4"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var timeout time.Duration = 5

type Worker struct {
	DeviceToken string
	TaskJson    string
	redisConn   *redis.ClusterClient
}

type HuaweiPushParams struct {
	AccessToken string `json:"access_token"`
	Mark        string `json:"mark"`
	Uid         int    `json:"uid"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type HuaweiTokenInfo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func Init(t time.Duration) {
	timeout = t
}

func NewPush(token string, jsonStr string, redisConn *redis.ClusterClient) (w *Worker) {
	//init the worker
	var wR Worker
	wR.DeviceToken = token
	wR.TaskJson = jsonStr
	wR.redisConn = redisConn
	return &wR
}

func (w Worker) AndroidHuaweiPush() (result bool) {
	token := w.getHuaweiToken()
	params := w.getHuaweiPushParams(w.TaskJson, token)
	url := "https://api.push.hicloud.com/pushsend.do?nsp_ctx=%7B%22ver%22%3A%221%22%2C+%22appId%22%3A%2210207663%22%7D"
	req, err := http.NewRequest("POST", url, params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Close = true

	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	ioutil.ReadAll(resp.Body)
	log.Info("[success] android response Body:", string(body))
	return true
}

func (w Worker) getHuaweiToken() (token string) {
	tokenKey := "huawei_access_token"
	token = w.redisConn.Get(tokenKey).Val()
	if token != "" {
		return token
	}
	tokenUrl := "https://login.cloud.huawei.com/oauth2/v2/token"
	postStr := "grant_type=client_credentials&client_secret=0tth2426x40wfye7nedlkh3o4gfwrp7i&client_id=10207663"
	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBuffer([]byte(postStr)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Close = true
	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err, tokenInfo := ParseToken(string(body))
	if err != nil {
		log.Error(err)
		return ""
	}
	w.redisConn.Set(tokenKey, tokenInfo.AccessToken, time.Duration((3600-300)*time.Second))
	log.Info("[success] android response huaweiGetToken:", string(tokenInfo.AccessToken))
	return tokenInfo.AccessToken
}

func (w Worker) getHuaweiPushParams(jsonStr string, token string) (params *strings.Reader) {
	err, commonInfo := ParsePushInfo(jsonStr)
	if err != nil {
		log.Error(err)
		return params
	}
	second := int(time.Now().Unix())
	payload := `{"hps":{"msg":{"action":{"param":{"appPkgName":"com.goumin.forum","intent":"intent:\/\/com.goumin.forum\/push_detail?message=what#Intent;scheme=myscheme;launchFlags=0x10000000;end"},"type":1},"type":3,"body":{"title":"` + commonInfo.Title + `","content":"` + commonInfo.Content + `"}},"ext":{"biTag":"Trump","icon":"","customize":[{"type":` + strconv.Itoa(commonInfo.Type) + `,"mark":"` + commonInfo.Mark + `","title":"` + commonInfo.Title + `","content":"` + commonInfo.Content + `","uid":` + strconv.Itoa(commonInfo.Uid) + `}]}}}`
	v := url.Values{}
	v.Set("access_token", token)
	v.Set("nsp_svc", "openpush.message.api.send")
	v.Set("nsp_ts", strconv.Itoa(second))
	v.Set("device_token_list", "[\""+w.DeviceToken+"\"]")
	v.Set("payload", payload)
	params = strings.NewReader(v.Encode())
	return params
}

func ParsePushInfo(jsonStr string) (error, ActiveRecord.PushInfo) {
	var arr ActiveRecord.PushInfo
	err := json.Unmarshal([]byte(jsonStr), &arr)
	if err != nil {
		return err, arr
	}
	return nil, arr
}

//解析参数
func ParseToken(redisStr string) (error, HuaweiTokenInfo) {
	var arr HuaweiTokenInfo
	err := json.Unmarshal([]byte(redisStr), &arr)
	if err != nil {
		return err, arr
	}
	return nil, arr
}
