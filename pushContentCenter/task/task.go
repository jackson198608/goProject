package task

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/club"
	"github.com/jackson198608/goProject/pushContentCenter/channels/focus"
	mgo "gopkg.in/mgo.v2"
)

type Task struct {
	raw       string
	mysqlXorm *xorm.Engine
	mongoConn *mgo.Session
	jobstr    string
	jobType   string
}

//job: redisQueue pop string
//taskarg: mongoHost,mongoDatabase,mongoReplicaSetName
func NewTask(raw string, mysqlXorm *xorm.Engine, mongoConn *mgo.Session) *task {
	//@todo check prams

	t = new(Task)
	if t == nil {
		return nil
	}

	//@todo pass param

	//@todo create private member
	jobStr, jobType, err := t.parseRaw()
	if err != nil {
		return nil
	}
	//@todo check return detail

	t.jobstr = jobStr
	t.jobType = jobType

	return t

}

func (t *Task) Do() error {

	return nil
}

//return:
//         jobstr
//	       type
//         trytimes
//		   error
func (t *Tasl) parseRaw() (string, string, error) {

	return "", "", 0, nil
}
