package task

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/gouminMobileMessageSender/Montnets"
	"strconv"
	"strings"
)

type Task struct {
	loggerLevel logger.LEVEL
	phone       string
	message     string
}

/*
	redis str format:    A|B
	A for phone
	B for message want to send
*/
func NewTask(loggerLevel logger.LEVEL, redisStr string) *Task {
	if loggerLevel < 0 {
		loggerLevel = 0
	}
	logger.SetLevel(loggerLevel)

	redisArr := strings.Split(redisStr, "|")
	if len(redisArr) != 2 {
		logger.Error("error redis str ", redisStr)
		return nil
	}

	phoneStr := redisArr[0]
	messageStr := redisArr[1]

	if !checkPhone(phoneStr) {
		return nil
	}

	if len(messageStr) == 0 {
		logger.Error("can not send empty message")
		return nil
	}

	t := new(Task)
	t.loggerLevel = loggerLevel
	t.phone = phoneStr
	t.message = messageStr
	return t

}

func checkPhone(phoneStr string) bool {
	_, err := strconv.Atoi(phoneStr)
	if err != nil {
		logger.Error("phone is not valid ", err)
		return false
	}
	return true
}

func (t *Task) Do() {
	m := Montnets.NewMontnets(t.loggerLevel, t.phone, t.message)
	if m != nil {
		m.Send()
	}
}
