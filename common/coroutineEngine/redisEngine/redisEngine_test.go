package redisEngine

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/tools"
	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v4"
	"testing"
)

const dbAuth = "root:goumin123"
const dbDsn = "127.0.0.1:3306"
const dbName = "test"
const mongoConn = "127.0.0.1:27017"

func newtask() (*RedisEngine, error) {

	redisInfo := redis.Options{
		Addr: "127.0.0.1:6379",
	}
	//getXormEngine
	connStr := tools.GetMysqlDsn(dbAuth, dbDsn, dbName)
	conns := []string{connStr}

	//get mongo session
	mgos := []string{mongoConn}

	r, err := NewRedisEngine("test", &redisInfo, mgos, conns, 3, 1, jobFunc)
	if err != nil {
		return nil, err
	}
	return r, nil

}

func newtaskWithEmptyInfo() (*RedisEngine, error) {

	redisInfo := redis.Options{
		Addr: "127.0.0.1:6379",
	}
	//getXormEngine
	conns := []string{}

	//get mongo session
	mgos := []string{}

	r, err := NewRedisEngine("test", &redisInfo, mgos, conns, 3, jobFunc)
	if err != nil {
		return nil, err
	}
	return r, nil

}

func jobFunc(job string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {

	fmt.Println("this is jobFunc", job)
	return errors.New("job func fail")
}

func TestNewTask(t *testing.T) {
	_, err := newtask()
	if err != nil {
		t.Log("create task error")
	}

}

func TestDo(t *testing.T) {
	r, err := newtask()
	if err != nil {
		t.Log("create task error")
	}

	err = r.Do()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

}

func TestDoWithEmptyParams(t *testing.T) {
	r, err := newtaskWithEmptyInfo()
	if err != nil {
		t.Log("create task error")
		t.Fail()
	}

	err = r.Do()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

}

func TestParseRaw(t *testing.T) {
	r, err := newtask()
	if err != nil {
		t.Log("create task error")
	}

	raw, tryTimes, err := r.parseRaw("te_s_d")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log(raw)
	t.Log(tryTimes)
}
