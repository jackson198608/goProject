package task

import (
	"errors"
	"github.com/bitly/go-simplejson"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"net/http"
)

type Task struct {
	Raw         string //the data get from redis queue
	phpServerIp string
	jsonData    *JsonColumn
}

//json column
type JsonColumn struct {
	callback string
	args     string
}

//job: redisQueue pop string
func NewTask(raw string, phpServerIp string) (*Task, error) {
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
	t.phpServerIp = phpServerIp

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
	err := t.callbackPhp()
	if err != nil {
		for i := 0; i < 5; i++ {
			err = t.callbackPhp()
			if err == nil {
				break
			}
		}
	}
	return err
}

func (t *Task) callbackPhp() error {
	abuyun := t.setAbuyun()
	targetUrl := "http://" + t.phpServerIp + "/" + t.jsonData.callback
	var h http.Header = make(http.Header)
	h.Set("HOST", "lingdang.goumin.com") //@todo change to online domain
	statusCode, _, body, err := abuyun.SendPostRequest(targetUrl, h, t.jsonData.args, true)

	if err != nil {
		logger.Error("http request error", err, "; task is ", t.Raw)
		return err
	}
	if statusCode == 200 {
		if body == "fail" {
			return errors.New("callback php fail ; task is " + t.Raw)
		} else if body == "success" {
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

	jsonC.callback, _ = js.Get("callback").String()
	jsonC.args, _ = js.Get("args").String()
	return &jsonC, nil
}
