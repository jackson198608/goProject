package task

import (
	"database/sql"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

func TestDo(t *testing.T) {
	dbName := "test_dz2"
	db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/"+dbName+"?charset=utf8")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	mongoConn := "192.168.86.68:27017"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		logger.Error("[error] connect mongodb err")
		return
	}
	defer session.Close()
	// redisStr := "1934|-1"
	// redisStr := "1934|0"
	// redisStr := "1934|1"
	redisStr := "1935|2"
	task := NewTask(0, redisStr, db, session)
	if task != nil {
		task.Do()
	}
}
