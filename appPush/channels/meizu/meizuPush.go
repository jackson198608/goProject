package meizu

import (
	"crypto/md5"
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

var push_url = "http://server-api-push.meizu.com/garcia/api/server/push/varnished/pushByPushId"
var appId = "124063"
var appSercet = "44851959ea724e948fac6fcc2d288dd2"
var timeout time.Duration = 5

func Init(t time.Duration) {
	timeout = t
}

type Worker struct {
	token     string
	jsonStr   string
	redisConn *redis.ClusterClient
}

func NewPush(token string, jsonStr string, redisConn *redis.ClusterClient) (w *Worker) {
	//init the worker
	var wR Worker
	wR.token = token
	wR.jsonStr = jsonStr
	wR.redisConn = redisConn
	return &wR
}

func (w Worker) AndroidMeizuPush() (result bool) {
	msg := w.pushMessage()
	if msg == "" {
		return false
	}
	sign := w.getPushSign(msg)
	log.Info("[success] meizu sign :", sign)
	v := url.Values{}
	v.Set("appId", appId)
	v.Set("messageJson", msg)
	v.Set("pushIds", w.token)
	v.Set("sign", sign)

	params := strings.NewReader(v.Encode())
	req, err := http.NewRequest("POST", push_url, params)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("sign", sign)
	req.Close = true

	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Info("[error] meizu client :", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Info("[success] meizu response ", string(body))
	return true
}

func (w Worker) pushMessage() (message string) {
	err, commonInfo := ParsePushInfo(w.jsonStr)
	if err != nil {
		return ""
	}
	message = `{"noticeBarInfo":{"noticeBarType":0,"title":"` + commonInfo.Title + `","content":"` + commonInfo.Content + `"},"noticeExpandInfo":{"noticeExpandType":0,"noticeExpandContent":"content"},"clickTypeInfo":{"clickType":3,"customAttribute":"{\"type\":` + strconv.Itoa(commonInfo.Type) + `,\"mark\":\"` + commonInfo.Mark + `\",\"title\":\"` + commonInfo.Title + `\",\"content\":\"` + commonInfo.Content + `\",\"uid\":` + strconv.Itoa(commonInfo.Uid) + `}"},"pushTimeInfo":{"offLine":1,"validTime":24},"advanceInfo":{"suspend":1,"clearNoticeBar":1,"notificationType":{"vibrate":1,"lights":1,"sound":1}}}`
	return message
}

func (w Worker) getPushSign(msgJson string) (sign string) {
	str := "appId=" + appId + "messageJson=" + msgJson + "pushIds=" + w.token + appSercet
	ctx := md5.New()
	ctx.Write([]byte(str))
	sign = hex.EncodeToString(ctx.Sum(nil))
	return sign
}

func ParsePushInfo(jsonStr string) (error, ActiveRecord.PushInfo) {
	var arr ActiveRecord.PushInfo
	err := json.Unmarshal([]byte(jsonStr), &arr)
	if err != nil {
		return err, arr
	}
	return nil, arr
}
