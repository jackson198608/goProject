package task

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/recommend_1.0/User"
	"github.com/jackson198608/goProject/recommend_1.0/content"
	mgo "gopkg.in/mgo.v2"
	"strings"
)

type Task struct {
	Raw       string         //the data get from redis queue
	MysqlXorm []*xorm.Engine //mysql single instance
	MongoConn []*mgo.Session //mongo single instance
	elkDsn    string
	Jobstr    string //private member parse from raw
	JobType   string //private membe parse from raw jobType: focus|club
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
	t.Raw = raw
	t.MysqlXorm = mysqlXorm
	t.MongoConn = mongoConn

	//create private member
	err := t.parseRaw()
	if err != nil {
		return nil, errors.New("raw format error ,can not find jobstr and jobtype " + raw)
	}

	return t, nil
}

func (t *Task) Do() error {
	switch t.JobType {
	case "follow":
		t.ChannelFollow()
	case "content":
		t.ChannelContent()
	}
	return nil
}

// follow channel's invoke function
func (t *Task) ChannelFollow() error {
	u := user.NewUser(t.MysqlXorm, t.MongoConn, t.Jobstr, t.elkDsn)
	err := u.Do()
	if err != nil {
		return err
	}
	return nil
}

// content channel's invoke function
func (t *Task) ChannelContent() error {
	c := content.NewContent(t.MysqlXorm, t.MongoConn, t.Jobstr)
	err := c.Do()
	if err != nil {
		return err
	}
	return nil
}

// this function parase raw to judge jobstr and job type
// sep string : '|'
//return:
//         jobstr
//	       type
//		   error
func (t *Task) parseRaw() error {
	rawSlice := []byte(t.Raw)
	rawLen := len(rawSlice)
	lastIndex := strings.LastIndex(t.Raw, "|")

	t.Jobstr = string(rawSlice[0:lastIndex])
	t.JobType = string(rawSlice[lastIndex+1 : rawLen])

	return nil

}
