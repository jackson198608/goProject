package appPush

import (
	"bytes"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//gloabl variables
var timeout time.Duration = 5

type Worker struct {
	t *Task
}

func Init(t time.Duration) {
	timeout = t
}

func NewWorker(t *Task) (w *Worker) {
	//init the worker
	var wR Worker
	wR.t = t
	return &wR
}

func (w Worker) Push(p12bytes []byte) (result bool) {
	phoneType := w.t.phoneType
	if phoneType == 0 {
		result = w.iosPush(p12bytes)
	} else {
		result = w.androidPush()
	}
	return result
}

func (w Worker) iosPush(p12bytes []byte) (result bool) {
	//cert, pemErr := certificate.FromPemFile("/etc/pro-lingdang.pem", "gouminwang")
	cert, pemErr := certificate.FromPemBytes(p12bytes, "gouminwang")
	if pemErr != nil {
		log.Println("[Error]Cert Error:", pemErr)
		result = false
		return result
	}

	notification := &apns.Notification{}
	notification.DeviceToken = w.t.DeviceToken
	notification.Topic = "com.goumin.bell"
	//notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below
	notification.Payload = w.t.TaskJson // See Payload section below

	client := apns.NewClient(cert).Production()
	//res, err := client.Push(notification)
	_, err := client.Push(notification)

	if err != nil {
		//log.Println("Error:", err)
		result = false
		return result
	}

	//log.Println("APNs ID:", res.ApnsID)
	return true
}

func (w Worker) androidPush() (result bool) {
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

	//body, _ := ioutil.ReadAll(resp.Body)
	ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	return true
}
