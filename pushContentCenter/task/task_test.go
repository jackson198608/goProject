package task

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/tools"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

const dbAuth = "dog123:dog123"
const dbDsn = "192.168.86.193:3307"
const dbName = "new_dog123"
const mongoConn = "192.168.86.192:27017" //"192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019"

func newtask() (*Task, error) {
	//getXormEngine
	connStr := tools.GetMysqlDsn(dbAuth, dbDsn, dbName)
	engine, err := xorm.NewEngine("mysql", connStr)
	if err != nil {
		return nil, err
	}

	engines := []*xorm.Engine{engine}

	//get mongo session
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		return nil, err
	}

	sessions := []*mgo.Session{session}

	// jobStr := "{\"uid\":881050,\"event_type\":1,\"event_info\":{\"title\":\"subject\",\"focus content\":\" focus  message\",\"image_num\":\"image_num\",\"forum\":\"金毛俱乐部\",\"tag\":0,\"source\":1,\"fid\":36,\"bid\":34},\"tid\":0,\"status\":1,\"action\":0,\"time\":\"2017-10-23 10:54:00\"}"
	jobStr := "{\"uid\":2060501,\"type\":6,\"infoid\":2234567,\"typeid\":2,\"subject\":\"subject\",\"message\":\" push to 36 message\",\"image_num\":\"image_num\",\"lastpost\":2,\"fid\":\"37,38,77\",\"lastposter\":\"0ssss\",\"status\":1,\"displayorder\":1,\"disgest\":1,\"qst_type\":0,\"created\":1508469600,\"action\":0}"

	t, err := NewTask(jobStr+"|club", engines, sessions)
	return t, err

}

func TestNewTask(t *testing.T) {

	task, err := newtask()
	if task == nil {
		t.Log("task create error", err)
		t.Fail()
	}

}

func TestDoTask(t *testing.T) {
	task, err := newtask()
	if task == nil {
		t.Log("task create error", err)
		t.Fail()
	}

	err = task.Do()
	if err != nil {
		t.Log("task do error", err)
		t.Fail()
	}

	closetask(task)

}

func TestChannelFocus(t *testing.T) {
	task, err := newtask()
	if task == nil {
		t.Log("task create error", err)
		t.Fail()
	}

	err = task.ChannelFocus()
	if err != nil {
		t.Log("task do error", err)
		t.Fail()
	}

	closetask(task)
}

func TestChannelClub(t *testing.T) {
	task, err := newtask()
	if task == nil {
		t.Log("task create error", err)
		t.Fail()
	}

	err = task.ChannelClub()
	if err != nil {
		t.Log("task do error", err)
		t.Fail()
	}

	closetask(task)

}
func closetask(t *Task) {
	for _, v := range t.MysqlXorm {
		v.Close()
	}
	for _, v := range t.MongoConn {
		v.Close()
	}
}
