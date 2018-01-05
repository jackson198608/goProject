package club

import (
	"errors"
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
	t, err := NewTask(2060500, engines, sessions, "192.168.86.5:9200")
	return t, err
}

func TestNewUser(t *testing.T) {

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

func TestGetUser(t *testing.T) {

}
