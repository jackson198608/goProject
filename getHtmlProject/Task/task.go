package Task

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"github.com/jackson198608/goProject/getHtmlProject/SaveHtml"
	// redis "gopkg.in/redis.v4"
	"errors"
	"strconv"
	"strings"
)

type Task struct {
	loggerLevel int
	id          int
	url         string
	queueName   string
	loopNum     int
	taskNewArgs []string
	abuyun      *abuyunHttpClient.AbuyunProxy
}

func NewTask(loggerLevel int, queueName string, redisStr string, taskNewArgs []string, abuyun *abuyunHttpClient.AbuyunProxy) (*Task, error) {
	if loggerLevel < 0 {
		loggerLevel = 0
	}
	logger.SetLevel(logger.LEVEL(loggerLevel))
	redisArr := strings.Split(redisStr, "|")
	var id int
	var url string
	if len(redisArr) == 2 {
		id, _ = strconv.Atoi(redisArr[1])
		url = redisArr[0]
	} else {
		return nil, errors.New("redis value is error")
	}
	t := new(Task)
	t.id = id
	t.url = url
	t.queueName = queueName
	t.taskNewArgs = taskNewArgs
	t.abuyun = abuyun
	return t, nil

}

func (t *Task) Do() error {
	m := SaveHtml.NewHtml(t.loggerLevel, t.queueName, t.id, t.url, t.taskNewArgs, t.abuyun)
	if m != nil {
		if t.id > 0 {
			logger.Info("export data to", t.queueName)
			err := m.CreateHtmlByUrl()
			if err != nil {
				return errors.New("save content error")
			}
		}
		return nil
	}
	return nil
}
