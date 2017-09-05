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
}

func NewTask(loggerLevel int, redisStr string, db *sql.DB, session *mgo.Session, taskNewArgs []string) *Task {
	if loggerLevel < 0 {
		loggerLevel = 0
	}
	logger.SetLevel(logger.LEVEL(loggerLevel))
	redisArr := strings.Split(redisStr, "|")
	var id int
	var typeid int = 0
	if len(redisArr) == 2 {
		if redisArr[1] == "1" { //thread
			id, _ = strconv.Atoi(redisArr[0])
			typeid, _ = strconv.Atoi(redisArr[1])
		} else {
			// id, _ = strconv.Atoi(redisArr[0])
		}
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
	// t.templateType = templateType
	return t

}

func (t *Task) Do() {
	m := NewInfo(t.loggerLevel, t.id, t.typeid, t.db, t.session, t.taskNewArgs)
	if m != nil {
		if t.id > 0 && t.typeid == 0 {
			logger.Info("export event to mongo")
			m.CreateThreadHtmlContent(t.id)
		}
	}
}
