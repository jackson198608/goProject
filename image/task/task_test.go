package task

import (
	"fmt"
	"testing"
)

func TestParseJson(t *testing.T) {
	var taskarg []string
	taskarg = append(taskarg, "127.0.0.1")
	jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200,\"callback\":\"notify/compress-notify\",\"args\":\"1\"}"
	c, _ := NewTask(jobStr, taskarg)
	fmt.Println(c.parseJson())
}

func TestCallbackPhp(t *testing.T) {
	var taskarg []string
	taskarg = append(taskarg, "127.0.0.1")
	jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200,\"callback\":\"notify/compress-notify\",\"args\":\"1\"}"
	f, _ := NewTask(jobStr, taskarg)
	fmt.Println(f.callbackPhp())
}

func TestDo(t *testing.T) {
	var taskarg []string
	taskarg = append(taskarg, "127.0.0.1")
	jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200,\"callback\":\"notify/compress-notify\",\"args\":\"1\"}"
	f, _ := NewTask(jobStr, taskarg)
	fmt.Println(f.Do())
}
