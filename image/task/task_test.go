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
	//jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200,\"callback\":\"compressnotify/updatepromotiongoods\",\"args\":\"1\",\"watermark\":\"/Users/Snow/img/1.png\",\"gravity\":\"northwest\",\"x\":0,\"y\":0}|all"
	jobStr   := "{\"path\": \"/Users/Snow/img/IMG_03001.JPG\",\"width\": 800,\"height\": 800,\"callback\": \"compressnotify/updatepromotiongoods\",\"args\": {\"p_id\": 6726,\"pg_id\": 42029,\"goods_id\": \"5244\",\"sku_id\": \"19421\",\"module\": 3,\"url\": \"/goodsSku/day_181106/20181106_f097c48.jpg\",\"width\": 800,\"size\": 3},\"watermark\": \"/Users/Snow/img/1.png\",\"gravity\": \"northwest\",\"x\": 0,\"y\": 0,\"afterPath\":\"/Users/Snow/img/IMG_03001_2222.JPG\"}|all"
	f, _ := NewTask(jobStr, taskarg)
	fmt.Println(f.Do())
}
