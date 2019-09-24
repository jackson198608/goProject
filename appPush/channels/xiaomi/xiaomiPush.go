package xiaomi

import (
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

var push_url = "https://api.xmpush.xiaomi.com/v3/message/regid"
var secret = "i4ktMRyMLHE70rT6X/fU2A=="
var package_name = "com.goumin.forum"
var timeout time.Duration = 5

func Init(t time.Duration) {
	timeout = t
}

type Worker struct {
	token     string
	jsonStr   string
	redisConn *redis.ClusterClient
}

type PushResponse struct {
	Result  string `json:"result"`
	Reason  string `json:"reason"`
	TraceId string `json:"trace_id"`
	Code    int    `json:"code"`
	Desc    string `json:"description"`
}

func NewPush(token string, jsonStr string, redisConn *redis.ClusterClient) (w *Worker) {
	//init the worker
	var wR Worker
	wR.token = token
	wR.jsonStr = jsonStr
	wR.redisConn = redisConn
	return &wR
}

func (w Worker) AndroidXiaomiPush() (result bool) {
	err, commonInfo := ParsePushInfo(w.jsonStr)
	if err != nil {
		log.Error(err)
		return false
	}

	v := url.Values{}
	v.Set("payload", w.jsonStr)
	v.Set("title", commonInfo.Title)
	v.Set("description", commonInfo.Content)
	v.Set("pass_through", strconv.Itoa(0))
	v.Set("notify_type", strconv.Itoa(-1))
	v.Set("restricted_package_name", package_name)
	v.Set("notify_id", strconv.Itoa(2))
	v.Set("extra.notify_foreground", strconv.Itoa(1))
	v.Set("registration_id", w.token)
	params := strings.NewReader(v.Encode())
	req, err := http.NewRequest("POST", push_url, params)
	req.Header.Set("Authorization", "key="+secret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("X-PUSH-SDK-VERSION", "PHP_SDK_V2.2.21")
	req.Header.Set("X-PUSH-HOST-LIST", strconv.FormatBool(true))
	req.Header.Set("Expect", "")
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

	err, response := parseResponse(string(body))
	if err != nil {
		log.Error(err)
		return false
	}
	if response.Code != 0 {
		log.Info("[error] response error :", response.Desc)
		return false
	}
	log.Info("[success] xiaomi push response: ", string(body))
	return true
}

func parseResponse(response string) (error, PushResponse) {
	var arr PushResponse
	err := json.Unmarshal([]byte(response), &arr)
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
