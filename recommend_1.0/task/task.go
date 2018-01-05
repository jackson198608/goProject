package task

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/recommend_1.0/Club/club"
	"github.com/jackson198608/goProject/recommend_1.0/User/user"
	mgo "gopkg.in/mgo.v2"
	"strings"
)

type Task struct {
	Raw       string         //the data get from redis queue
	MysqlXorm []*xorm.Engine //mysql single instance
	MongoConn []*mgo.Session //mongo single instance
	Uid       int            //private member parse from raw
}

//job: redisQueue pop string
//taskarg: mongoHost,mongoDatabase,mongoReplicaSetName
func NewTask(raw string, mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session) (*Task, error) {
	//check prams
	if (raw == "") || (mysqlXorm == nil) || (mongoConn == nil) {
		return nil, errors.New("params can not be null")
	}

	t := new(Task)
	if t == nil {
		return nil, errors.New("there is no space to create struct")
	}

	//pass params
	t.Raw = raw
	t.MysqlXorm = mysqlXorm
	t.MongoConn = mongoConn

	return t, nil

}

// public interface for task
// if you have New channles you need to add logic here
func (t *Task) Do() error {
	err := t.ChannelFocus()
	if err != nil {
		return err
	} else {
		return nil
	}
}

// focus channel's invoke function
func (t *Task) user() error {
	u := user.NewUser(t.MysqlXorm, t.MongoConn, t.Uid, t.Raw[0])
	err := u.Do()
	if err != nil {
		return err
	}
	return nil
}
