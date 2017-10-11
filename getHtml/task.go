package main

import (
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	redis "gopkg.in/redis.v4"
	"strconv"
	"strings"
)

type Task struct {
	loggerLevel       int
	id                int
	pages             int
	db                *sql.DB
	loopNum           int
	templateType      string
	taskNewArgs       []string
	relateDefaultData string
	jobType           string
	client            *redis.Client
}

func NewTask(loggerLevel int, redisStr string, db *sql.DB, taskNewArgs []string, jobType string, client *redis.Client) *Task {
	if loggerLevel < 0 {
		loggerLevel = 0
	}
	logger.SetLevel(logger.LEVEL(loggerLevel))
	redisArr := strings.Split(redisStr, "|")
	var id int
	var pages int = 1
	if len(redisArr) == 2 {
		id, _ = strconv.Atoi(redisArr[0])
		pages, _ = strconv.Atoi(redisArr[1])
	}
	if len(redisArr) == 1 {
		id, _ = strconv.Atoi(redisStr)
	}
	t := new(Task)
	t.id = id
	t.pages = pages
	t.db = db
	t.taskNewArgs = taskNewArgs
	t.jobType = jobType
	t.client = client
	return t

}

func (t *Task) Do() {
	m := NewInfo(t.loggerLevel, t.id, t.db, t.taskNewArgs, t.client)
	if m != nil {
		if t.id > 0 {
			logger.Info("export thread to miphtml")
			m.CreateHtmlByUrl(t.id, t.pages, t.jobType)
		}

	}
}
