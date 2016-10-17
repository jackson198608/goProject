package appPush

import (
	"bytes"
	"fmt"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"io/ioutil"
	"log"
	"net/http"
)

type Worker struct {
	t *Task
}

func NewWorker(t *Task) (w *Worker) {
	//init the worker
	var wR Worker
	wR.t = t
	return &wR
}

func (w Worker) Push() {
	phoneType := w.t.phoneType
	if phoneType == 0 {
		w.iosPush()
	} else {
		w.androidPush()
	}
}

func (w Worker) iosPush() {
	cert, pemErr := certificate.FromPemFile("/etc/pro-lingdang.pem", "gouminwang")
	if pemErr != nil {
		log.Println("Cert Error:", pemErr)
	}

	notification := &apns.Notification{}
	notification.DeviceToken = w.t.DeviceToken
	notification.Topic = "com.goumin.bell"
	//notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below
	notification.Payload = w.t.TaskJson // See Payload section below

	client := apns.NewClient(cert).Production()
	res, err := client.Push(notification)

	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("APNs ID:", res.ApnsID)
}

func (w Worker) androidPush() {
	url := "http://sdk.open.api.igexin.com/apiex.htm"
	fmt.Println("URL:>", url)

	var jsonStr []byte = []byte(w.t.TaskJson)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "GeTui PHP/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
