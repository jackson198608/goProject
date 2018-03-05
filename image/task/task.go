package task

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"github.com/jackson198608/goProject/image/compress"
	"net/http"
)

type Task struct {
	Raw         string //the data get from redis queue
	phpServerIp string
	jsonData    *JsonColumn
}

//json column
type JsonColumn struct {
	imgaePath     string
	width         int
	height        int
	callbackRoute string
	args          string
}

//job: redisQueue pop string
func NewTask(raw string, taskarg []string) (*Task, error) {
	//check prams
	if raw == "" {
		return nil, errors.New("params can not be null")
	}

	t := new(Task)
	if t == nil {
		return nil, errors.New("there is no space to create struct")
	}

	//pass params
	t.Raw = raw
	jsonColumn, err := t.parseJson()
	if err != nil {
		return t, nil
	}
	t.jsonData = jsonColumn

	t.phpServerIp = taskarg[0]

	return t, nil

}

func (t *Task) setAbuyun() *abuyunHttpClient.AbuyunProxy {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy("", "", "")

	if abuyun == nil {
		logger.Error("create abuyun error")
		return nil
	}
	return abuyun
}

// public interface for task
// If the compression is successful, the callback PHP
func (t *Task) Do() error {
	c := compress.NewCompress(t.jsonData.imgaePath, t.jsonData.width, t.jsonData.height)
	err := c.Do()
	if err == nil {
		err = t.callbackPhp()
	}
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) callbackPhp() error {
	abuyun := t.setAbuyun()
	targetUrl := "http://" + t.phpServerIp + "/" + t.jsonData.callbackRoute
	logger.Info("targetUrl", targetUrl)
	var h http.Header = make(http.Header)
	h.Set("HOST", "dev.lingdang.goumintest.com") //@todo change to online domain
	statusCode, _, body, err := abuyun.SendPostRequest(targetUrl, h, t.jsonData.args, true)

	if err != nil {
		logger.Error("http request error", err, "; task is ", t.Raw)
		return err
	}
	if statusCode == 200 {
		if body == "fail" {
			return errors.New("callback php fail ; task is " + t.Raw)
		} else if body == "sucess" {
			logger.Error("callback php sucess ; task is ", t.Raw)
		}
		return nil
	}
	return err
}

//change json colum to object private member
func (t *Task) parseJson() (*JsonColumn, error) {
	var jsonC JsonColumn
	js, err := simplejson.NewJson([]byte(t.Raw))
	if err != nil {
		return &jsonC, err
	}

	jsonC.imgaePath, _ = js.Get("path").String()
	jsonC.width, _ = js.Get("width").Int()
	jsonC.height, _ = js.Get("height").Int()
	jsonC.callbackRoute, _ = js.Get("callback").String()
	jsonC.args, _ = js.Get("args").String()
	return &jsonC, nil
}
