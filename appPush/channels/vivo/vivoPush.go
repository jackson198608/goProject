package vivo

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/jackson198608/goProject/ActiveRecord"
	log "github.com/thinkboy/log4go"
	"gopkg.in/redis.v4"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var app_key = "b2466161-ef62-414f-bb10-1b6150978001"
var app_id = "15312"
var app_secret = "bef72815-a082-4f93-a857-52e758bd5b74"
var push_url = "https://api-push.vivo.com.cn/message/send"
var auth_url = "https://api-push.vivo.com.cn/message/auth"
var timeout time.Duration = 5

func Init(t time.Duration) {
	timeout = t
}

type Worker struct {
	token     string
	jsonStr   string
	redisConn *redis.ClusterClient
}

type AuthInfo struct {
	Result    int    `json:"result"`
	Desc      string `json:"desc"`
	AuthToken string `json:"authToken"`
}

type ResponseInfo struct {
	Result int    `json:"result"`
	Desc   string `json:"desc"`
	TaskId string `json:"taskId"`
}

func NewPush(token string, jsonStr string, redisConn *redis.ClusterClient) (w *Worker) {
	var wR Worker
	wR.token = token
	wR.jsonStr = jsonStr
	wR.redisConn = redisConn
	return &wR
}

func (w Worker) AndroidVivoPush() (result bool) {
	authToken := w.getAuth()
	if authToken == "" {
		return false
	}
	log.Info("[success] vivo authToken : ", authToken)
	params := w.getPushParams()
	var jsonStr []byte = []byte(params)
	req, err := http.NewRequest("POST", push_url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authToken", authToken)
	req.Close = true
	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("[error] vivo push request :", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Info("[messge] vivo android response vivoPushMsg:", string(body))
	//{"result":0,"desc":"请求成功","taskId":"624304170434887680"}
	err, response := w.parseResponseInfo(string(body))
	if err != nil {
		log.Error("[error] vivo parase response info :", err)
		return false
	}
	if response.Result != 0 {
		log.Info("[error] vivo response error info :", response.Desc)
		return false
	}
	log.Info("[success] vivo android response vivoMessageId: ", response.Desc)
	return true
}

func (w Worker) getPushParams() (params string) {
	err, commonInfo := ParsePushInfo(w.jsonStr)
	if err != nil {
		return ""
	}
	msectime := int(time.Now().UnixNano() / 1e6)
	params = `{"regId":"` + w.token + `","notifyType":4,"title":"` + commonInfo.Title + `","content":"` + commonInfo.Content + `","skipType":3,"skipContent":"start","clientCustomMap":` + w.jsonStr + `,"requestId":"` + strconv.Itoa(msectime) + `"}`
	return params
}

func (w Worker) getAuth() (auth_token string) {
	msectime := int(time.Now().UnixNano() / 1e6)
	sign := w.getSignStr(msectime)
	params := `{"appId":"` + app_id + `","appKey":"` + app_key + `","timestamp":` + strconv.Itoa(msectime) + `,"sign":"` + sign + `"}`
	var jsonStr []byte = []byte(params)
	req, err := http.NewRequest("POST", auth_url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("[error] auth request :", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	err, tokenInfo := w.parseAuthInfo(string(body))
	if err != nil {
		log.Error(err)
		return ""
	}
	log.Info("[success] android response vivoGetToken:", string(tokenInfo.AuthToken))
	return tokenInfo.AuthToken
}

func (w Worker) getSignStr(msectime int) (sign string) {
	str := app_id + app_key + strconv.Itoa(msectime) + app_secret
	ctx := md5.New()
	ctx.Write([]byte(str))
	sign = hex.EncodeToString(ctx.Sum(nil))
	return sign
}

func (w Worker) parseAuthInfo(jsonStr string) (error, AuthInfo) {
	var arr AuthInfo
	err := json.Unmarshal([]byte(jsonStr), &arr)
	if err != nil {
		return err, arr
	}
	return nil, arr
}

func (w Worker) parseResponseInfo(jsonStr string) (error, ResponseInfo) {
	var arr ResponseInfo
	err := json.Unmarshal([]byte(jsonStr), &arr)
	if err != nil {
		return err, arr
	}
	return nil, arr
}

func ParsePushInfo(jsonStr string) (error, ActiveRecord.PushInfo) {
	var arr ActiveRecord.PushInfo
	err := json.Unmarshal([]byte(jsonStr), &arr)
	if err != nil {
		return err, arr
	}
	return nil, arr
}
