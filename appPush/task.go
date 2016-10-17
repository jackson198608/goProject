package appPush

import (
	"strconv"
	"strings"
)

type Task struct {
	phoneType   int //o for ios ,1 for android
	DeviceToken string
	TaskJson    string
}

func NewTask(redisString string) (t *Task) {
	var tR Task
	result := strings.Split(redisString, "|")
	tR.phoneType, _ = strconv.Atoi(result[0])
	tR.DeviceToken = result[1]
	tR.TaskJson = result[2]
	return &tR
}
