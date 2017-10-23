package redisEngine

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/tools"
	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v4"
	"testing"
)

const dbAuth = "dog123:dog123"
const dbDsn = "192.168.86.193:3307"
const dbName = "new_dog123"
const mongoConn = "192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019"

func newtask() (*RedisEngine, error) {

	redisInfo := redis.Options{
		Addr: "127.0.0.1:6379",
	}
	//getXormEngine
	connStr := tools.GetMysqlDsn(dbAuth, dbDsn, dbName)
	conns := []string{connStr}

	//get mongo session
	mgos := []string{mongoConn}

	r := NewRedisEngine(redisInfo, conns, mgos, 3, jobFunc)

}

func jobFunc(job string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {

	fmt.Println("this is jobFunc")
	return nil
}

func TestNewTask(t *testing.T) {
	r, err := newtask()
	if err != nil {
		t.Log("create task error")
	}
}
