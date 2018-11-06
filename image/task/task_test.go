package task

import (
	"fmt"
	"testing"
)

func TestParseJson(t *testing.T) {
	var taskarg []string
	taskarg = append(taskarg, "127.0.0.1")
	taskarg = append(taskarg, "")
	jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200,\"callback\":\"compressnotify/dog\",\"args\":\"1\",\"watermark\":\"/Users/Snow/img/1.png\"}"
	c, _ := NewTask(jobStr, taskarg)
	fmt.Println(c.parseJson())
}

func TestCallbackPhp(t *testing.T) {
	var taskarg []string
	taskarg = append(taskarg, "127.0.0.1")
	jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200,\"callback\":\"compressnotify/dog\",\"args\":\"1\"}"
	f, _ := NewTask(jobStr, taskarg)
	//path := "forum/201802/12/1518412896316_200.png"
	fmt.Println(f.callbackPhp())
}

func TestDo(t *testing.T) {
	var taskarg []string
	taskarg = append(taskarg, "127.0.0.1")
	taskarg = append(taskarg, "")
	jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200,\"callback\":\"compressnotify/dog\",\"args\":\"1\",\"watermark\":\"/Users/Snow/img/1.png\",\"gravity\":\"northwest\",\"x\":0,\"y\":0}|composite"
	f, _ := NewTask(jobStr, taskarg)
	fmt.Println(f.Do())
}
