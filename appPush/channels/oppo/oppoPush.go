package oppo

import (
	"crypto/sha256"
	"encoding/hex"
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

var auth_url = "https://api.push.oppomobile.com/server/v1/auth"
var save_message_content_url = "https://api.push.oppomobile.com/server/v1/message/notification/save_message_content"
var broadcast_url = "https://api.push.oppomobile.com/server/v1/message/notification/broadcast"
var unicast_batch_url = "https://api.push.oppomobile.com/server/v1/message/notification/unicast_batch"
var APP_KEY = "7jGL33e7EMsC4S488s4O8OwS8"
var MasterSecret = "87874c1285Eb2C2E3f3ed0159cDd06F7"
var timeout time.Duration = 5

func Init(t time.Duration) {
	timeout = t
}

type Worker struct {
	token     string
	jsonStr   string
	redisConn *redis.ClusterClient
}

type OppoTokenInfo struct {
	Code    int      `json:"code"`
	Data    DataInfo `json:"data"`
	Message string   `json:"message"`
}

type DataInfo struct {
	AuthToken  string `json:"auth_token"`
	CreateTime int    `json:"create_time"`
}

type MessageInfo struct {
	Code    int         `json:"code"`
	MsgData MsgDataInfo `json:"data"`
	Message string      `json:"message"`
}

type MsgDataInfo struct {
	MessageId string `json:"message_id"`
}

func NewPush(token string, jsonStr string, redisConn *redis.ClusterClient) (w *Worker) {
	var wR Worker
	wR.token = token
	wR.jsonStr = jsonStr
	wR.redisConn = redisConn
	return &wR
}

func (w Worker) AndroidOppoPush() (result bool) {
	auth_token := w.getAuth()
	if auth_token == "" {
		return false
	}
	log.Info("[success] oppo auth_token is :", auth_token)
	message_id := w.getPushMessageId(auth_token)
	v := url.Values{}
	v.Set("message_id", message_id)
	v.Set("auth_token", auth_token)
	v.Set("target_type", strconv.Itoa(2))
	v.Set("target_value", w.token)
	params := strings.NewReader(v.Encode())
	req, err := http.NewRequest("POST", broadcast_url, params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Close = true
	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Info("[error] oppo client response :", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Info("[success] oppo push message response :", string(body))
	return true
}

func (w Worker) getPushMessageId(auth_token string) (message_id string) {
	err, commonInfo := ParsePushInfo(w.jsonStr)
	if err != nil {
		log.Info("[error] oppo parse :", err)
		return ""
	}

	v := url.Values{}
	v.Set("title", commonInfo.Title)
	v.Set("content", commonInfo.Content)
	v.Set("click_action_type", strconv.Itoa(1))
	v.Set("click_action_activity", "com.goumin.forum.push.oppo.internal")
	v.Set("action_parameters", w.jsonStr)
	v.Set("auth_token", auth_token)
	params := strings.NewReader(v.Encode())

	req, err := http.NewRequest("POST", save_message_content_url, params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Close = true
	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Info("[error] oppo client ", err)
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Info("[message] oppo get message id :", string(body))

	err, msgInfo := w.ParseMessageInfo(string(body))
	if err != nil {
		log.Error(err)
		return ""
	}
	log.Info("[success] oppo android response oppoMessageId:", msgInfo.MsgData.MessageId)
	return msgInfo.MsgData.MessageId
}

func (w Worker) getAuth() (token string) {
	msectime := int(time.Now().UnixNano() / 1e6)
	hashStr := APP_KEY + strconv.Itoa(msectime) + MasterSecret
	//使用sha256哈希函数
	h := sha256.New()
	h.Write([]byte(hashStr))
	sum := h.Sum(nil)
	//由于是十六进制表示，因此需要转换
	sign := hex.EncodeToString(sum)
	v := url.Values{}
	v.Set("app_key", APP_KEY)
	v.Set("sign", string(sign[:]))
	v.Set("timestamp", strconv.Itoa(msectime))
	params := strings.NewReader(v.Encode())
	req, err := http.NewRequest("POST", auth_url, params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Close = true
	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Info("[error] oppo client :", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Info("[notice] oppo android response oppoGetToken:", string(body))

	err, tokenInfo := w.ParseToken(string(body))
	if err != nil {
		log.Info("[error] oppo parse token :", err)
		return ""
	}
	log.Info("[success] oppo android response oppoGetToken:", string(tokenInfo.Message))
	return tokenInfo.Data.AuthToken
}

func ParsePushInfo(jsonStr string) (error, ActiveRecord.PushInfo) {
	var arr ActiveRecord.PushInfo
	err := json.Unmarshal([]byte(jsonStr), &arr)
	if err != nil {
		return err, arr
	}
	return nil, arr
}

func (w Worker) ParseMessageInfo(msgStr string) (error, MessageInfo) {
	var arr MessageInfo
	err := json.Unmarshal([]byte(msgStr), &arr)
	if err != nil {
		return err, arr
	}
	return nil, arr
}

//解析参数
func (w Worker) ParseToken(redisStr string) (error, OppoTokenInfo) {
	var arr OppoTokenInfo
	err := json.Unmarshal([]byte(redisStr), &arr)
	if err != nil {
		return err, arr
	}
	return nil, arr
}
