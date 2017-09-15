package main

import (
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
	"strconv"
	"strings"
)

type Task struct {
	loggerLevel  int
	id           int
	typeid       int
	db           *sql.DB
	loopNum      int
	session      *mgo.Session
	templateType string
	taskNewArgs  []string
	relateAsk    string
	jobType      string
}

func NewTask(loggerLevel int, redisStr string, db *sql.DB, session *mgo.Session, taskNewArgs []string, relateAsk string, jobType string) *Task {
	if loggerLevel < 0 {
		loggerLevel = 0
	}
	logger.SetLevel(logger.LEVEL(loggerLevel))
	redisArr := strings.Split(redisStr, "|")
	var id int
	var typeid int = 0
	if len(redisArr) == 2 {
		id, _ = strconv.Atoi(redisArr[0])
		typeid, _ = strconv.Atoi(redisArr[1])
	}
	if len(redisArr) == 1 {
		id, _ = strconv.Atoi(redisStr)
	}
	t := new(Task)
	t.id = id
	t.typeid = typeid
	t.db = db
	t.session = session
	t.taskNewArgs = taskNewArgs
	t.relateAsk = relateAsk
	t.jobType = jobType
	return t

}

func (t *Task) Do() {
	if t.jobType == "thread" {
		m := NewInfo(t.loggerLevel, t.id, t.typeid, t.db, t.session, t.taskNewArgs)
		if m != nil {
			if t.id > 0 && t.typeid == 0 {
				logger.Info("export thread to miphtml")
				m.CreateThreadHtmlContent(t.id, t.relateAsk)
			}

		}
	}
	if t.jobType == "ask" {
		m := NewAskInfo(t.loggerLevel, t.id, t.db, t.session, t.taskNewArgs)
		if m != nil {
			if t.id > 0 {
				logger.Info("export ask to miphtml")
				m.CreateAskHtmlContent(t.id)
			}
		}
	}
}
