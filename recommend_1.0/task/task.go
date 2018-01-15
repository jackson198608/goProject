package task

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/recommend_1.0/User"
	mgo "gopkg.in/mgo.v2"
	// "strings"
)

type Task struct {
	Uid       string         //the data get from redis queue
	MysqlXorm []*xorm.Engine //mysql single instance
	MongoConn []*mgo.Session //mongo single instance
	elkDsn    string
}

//job: redisQueue pop string
//taskarg: mongoHost,mongoDatabase,mongoReplicaSetName
func NewTask(raw string, mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, taskarg []string) (*Task, error) {
	//check prams
	if (raw == "") || (mysqlXorm == nil) || (mongoConn == nil) {
		return nil, errors.New("params can not be null")
	}

	t := new(Task)
	if t == nil {
		return nil, errors.New("there is no space to create struct")
	}

	//pass params
	t.Uid = raw
	t.MysqlXorm = mysqlXorm
	t.MongoConn = mongoConn
	t.elkDsn = taskarg[0]

	return t, nil
}

func (t *Task) Do() error {
	u := user.NewUser(t.MysqlXorm, t.MongoConn, t.Uid, t.elkDsn)
	err := u.Do()
	if err != nil {
		return err
	}
	return nil
}
