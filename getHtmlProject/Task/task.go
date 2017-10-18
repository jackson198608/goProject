package Task

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/getHtmlProject/SaveHtml"
	redis "gopkg.in/redis.v4"
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
	client      *redis.Client
}

func NewTask(loggerLevel int, queueName string, redisStr string, taskNewArgs []string, client *redis.Client) *Task {
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
		return nil
	}
	t := new(Task)
	t.id = id
	t.url = url
	t.queueName = queueName
	t.taskNewArgs = taskNewArgs
	t.client = client
	return t

}

func (t *Task) Do() {
	m := SaveHtml.NewHtml(t.loggerLevel, t.queueName, t.id, t.url, t.taskNewArgs, t.client)
	if m != nil {
		if t.id > 0 {
			logger.Info("export thread to threadHtmlUrl")
			m.CreateHtmlByUrl()
		}

	}
}
