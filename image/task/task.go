package task

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"github.com/jackson198608/goProject/image/composite"
	"github.com/jackson198608/goProject/image/compress"
	"net/http"
	"strings"
	"os"
)

type Task struct {
	Raw         string //the data get from redis queue
	phpServerIp string
	waterPath   string
	jsonData    *JsonColumn
	JobType     string
	Jobstr      string
}

//json column
type JsonColumn struct {
	imgaePath   string
	width       int
	height      int
	callback    string
	args        string
	watermark   string
	gravityType string
	offsetX     int
	offsetY     int
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
	t.parseRaw()

	jsonColumn, err := t.parseJson()
	if err != nil {
		return t, nil
	}
	t.jsonData = jsonColumn

	t.phpServerIp = taskarg[0]
	t.waterPath = taskarg[1]

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
	switch t.JobType {
	case "all":
		err := t.channelAll()
		if err != nil {
			return err
		} else {
			return nil
		}
		break
	case "compress":
		_, err := t.channelCompress()
		if err != nil {
			return err
		} else {
			return nil
		}
		break
	case "composite":
		watermarkPath := t.watermarkImage()
		err := t.channelComposite(t.jsonData.imgaePath, watermarkPath, t.jsonData.gravityType, t.jsonData.offsetX, t.jsonData.offsetY)
		if err != nil {
			return err
		} else {
			return nil
		}
		break
	}
	return nil
}

func (t *Task) channelAll() error {
	path, err := t.channelCompress()
	fmt.Println(err)
	if err == nil {
		watermarkPath := t.watermarkImage()
		err = t.channelComposite(path, watermarkPath, t.jsonData.gravityType, t.jsonData.offsetX, t.jsonData.offsetY)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (t *Task) channelCompress() (string, error) {
	c := compress.NewCompress(t.jsonData.imgaePath, t.jsonData.width, t.jsonData.height)
	path, err := c.Do()
	if err == nil {
		if t.jsonData.callback != "" {
			err = t.callback()
		}
	}
	if err != nil {
		return path, err
	}
	return path, nil
}

func (t *Task) channelComposite(path string, watermarkPath string, gravityType string, offsetX int, offsetY int) error {
	cp := composite.NewComposite(path, watermarkPath, gravityType, offsetX, offsetY)
	compositeErr := cp.Do()
	if compositeErr != nil {
		for i := 0; i < 5; i++ {
			compositeErr := cp.Do()
			if compositeErr == nil {
				break
			}
		}
	}
	return nil
}

func (t *Task) watermarkImage() string {
	if t.jsonData.watermark != "" {
		if (Exist(t.jsonData.watermark)) {
			return t.jsonData.watermark
		}
		return t.waterPath + t.jsonData.watermark
	} else {
		if t.jsonData.width <= 220 {
			return t.waterPath + "220.png"
		} else if t.jsonData.width > 220 && t.jsonData.width <= 340 {
			return t.waterPath + "340.png"
		} else {
			return t.waterPath + "720.png"
		}
	}
}

func (t *Task) callback() error {
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
	js, err := simplejson.NewJson([]byte(t.Jobstr))
	if err != nil {
		return &jsonC, err
	}

	jsonC.imgaePath, _ = js.Get("path").String()
	jsonC.width, _ = js.Get("width").Int()
	jsonC.height, _ = js.Get("height").Int()
	jsonC.callback, _ = js.Get("callback").String()
	jsonC.args, _ = js.Get("args").String()
	jsonC.watermark, _ = js.Get("watermark").String()
	jsonC.gravityType, _ = js.Get("gravity").String()
	jsonC.offsetX, _ = js.Get("x").Int()
	jsonC.offsetY, _ = js.Get("y").Int()
	return &jsonC, nil
}

// this function parase raw to judge jobstr and job type
// sep string : '|'
//return:
//         jobstr
//	       type
//		   error
func (t *Task) parseRaw() error {
	rawSlice := []byte(t.Raw)
	lastIndex := strings.LastIndex(t.Raw, "|")
	if lastIndex > 1 {
		rawLen := len(rawSlice)
		t.Jobstr = string(rawSlice[0:lastIndex])
		t.JobType = string(rawSlice[lastIndex+1 : rawLen])
	} else {
		t.Jobstr = t.Raw
		t.JobType = "compress"
	}
	return nil
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
