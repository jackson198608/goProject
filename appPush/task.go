package iosPush

import (
	"strings"
)

type Task struct {
	DeviceToken string
	TaskJson    string
}

func NewTask(redisString string) (t *Task) {
	var tR Task
	result := strings.Split(redisString, "|")
	tR.DeviceToken = result[0]
	tR.TaskJson = result[1]
	return &tR
}
