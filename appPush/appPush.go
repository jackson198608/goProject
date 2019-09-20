package appPush

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/thinkboy/log4go"
	"io/ioutil"
	"net/http"
	"time"
	"gopkg.in/redis.v4"
	"github.com/jackson198608/goProject/appPush/channels/oppo"
	"github.com/jackson198608/goProject/appPush/channels/huawei"
	"github.com/jackson198608/goProject/appPush/channels/vivo"
	"github.com/jackson198608/goProject/appPush/channels/meizu"
	"github.com/jackson198608/goProject/appPush/channels/xiaomi"
)

//gloabl variables
var timeout time.Duration = 5

//mobPush推送
var mobKey = "2b14bb6c10bac"
var mobSign = "172cca337b3b9ea4cf4b250bd7c773e0"

type Worker struct {
	t *Task
	redisConn *redis.ClusterClient
}

type modPushRes struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Res    string `json:"res"`
}


func Init(t time.Duration) {
	timeout = t
}




func NewWorker(t *Task, redisConn *redis.ClusterClient) (w *Worker) {
	//init the worker
	var wR Worker
	wR.t = t
	wR.redisConn = redisConn
	return &wR
}

func (w Worker) Push(p12bytes []byte) (result bool) {
	phoneType := w.t.phoneType
	if phoneType == 0 {
		result = w.iosPush(p12bytes)
	} else if phoneType == 1 {
		//result = w.androidPush()

		//result = w.androidPushMob()
		result = true
	} else if phoneType == 2 {
		fmt.Println("[request] huawei phoneType:2")
		hw := huawei.NewPush(w.t.DeviceToken, w.t.TaskJson, w.redisConn)
		result = hw.AndroidHuaweiPush()
	} else if phoneType == 3 {
		fmt.Println("[request] xiaomi phoneType:3")
		xm := xiaomi.NewPush(w.t.DeviceToken, w.t.TaskJson, w.redisConn)
		result = xm.AndroidXiaomiPush()
	} else if phoneType == 4 {
		fmt.Println("[request] oppo phoneType:4")
		op := oppo.NewPush(w.t.DeviceToken, w.t.TaskJson, w.redisConn)
		result = op.AndroidOppoPush()
	} else if phoneType == 5 {
		fmt.Println("[request] vivo phoneType:5")
		vv := vivo.NewPush(w.t.DeviceToken, w.t.TaskJson, w.redisConn)
		result = vv.AndroidVivoPush()
	} else if phoneType == 6 {
		fmt.Println("[request] meizu phoneType:6")
		mz := meizu.NewPush(w.t.DeviceToken, w.t.TaskJson, w.redisConn)
		result = mz.AndroidMeizuPush()
	} else {
		//wx program
		result = w.wxProgramPush()
	}
	return result
}

func (w Worker) wxProgramPush() (result bool) {
	//调用channel 请求微信
	return true
}

func (w Worker) iosPush(p12bytes []byte) (result bool) {
	//cert, pemErr := certificate.FromPemFile("/etc/pro-lingdang.pem", "gouminwang")
	cert, pemErr := certificate.FromPemBytes(p12bytes, "gouminwang")
	if pemErr != nil {
		//fmt.Println("[Error]Cert Error:", pemErr)
		result = false
		return result
	}

	notification := &apns.Notification{}
	notification.DeviceToken = w.t.DeviceToken
	notification.Topic = "com.goumin.bell"
	//notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below
	notification.Payload = w.t.TaskJson // See Payload section below

	client := apns.NewClient(cert).Production()
	res, err := client.Push(notification)
	//_, err := client.Push(notification)
	log4go.Info("res:", res)


	if err != nil {
		//fmt.Println("Error:", err)
		log4go.Error("Error:", err)
		result = false
		return result
	}

	log4go.Info("APNs ID:", res.ApnsID)
	//fmt.Println("APNs ID:", res.ApnsID)
	return true
}

func (w Worker) androidPush() (result bool) {
	//fmt.Println("[notice]androidPush")
	url := "http://sdk.open.api.igexin.com/apiex.htm"

	var jsonStr []byte = []byte(w.t.TaskJson)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "GeTui PHP/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true

	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		result = false
		return result
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	ioutil.ReadAll(resp.Body)
	fmt.Println("[notice] android response Body:", string(body))
	return true
}

//android push 新的推送方法
func (w Worker) androidPushMob() (result bool) {
	url := "http://api.push.mob.com/v2/push"

	var jsonStr []byte = []byte(w.t.TaskJson)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//sign加密
	info := w.t.TaskJson + mobSign
	secret := md5Str(info)

	req.Header.Set("key", mobKey)
	req.Header.Set("sign", secret)
	req.Close = true

	timeout := time.Duration(timeout * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		result = false
		return result
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	ioutil.ReadAll(resp.Body)

	var a modPushRes
	json.Unmarshal([]byte(string(body)), &a)
	if a.Status != 200 {
		log4go.Debug("android push fail status:", a.Status, " error:", a.Error, " taskJson:", w.t.TaskJson)
	}
	log4go.Info("[notice] android response Body:", string(body))

	return true
}

//md5加密
func md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}


