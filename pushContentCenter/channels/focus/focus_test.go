package focus

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

func testConn() ([]*xorm.Engine, []*mgo.Session) {
	dbAuth := "dog123:dog123"
	dbDsn := "192.168.86.193:3307"
	// dbDsn := "210.14.154.117:33068"
	dbName := "new_dog123"
	dataSourceName := dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	mongoConn := "192.168.86.192:27017"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		fmt.Println("[error] connect mongodb err")
		return nil, nil
	}
	var engineAry []*xorm.Engine
	engineAry = append(engineAry, engine)
	var sessionAry []*mgo.Session
	sessionAry = append(sessionAry, session)
	return engineAry, sessionAry
	// return engine, session
}

func TestParseJson(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36},\"tid\":0,\"status\":1,\"time\":1508469600}"
	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	f.parseJson()
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":881050,\"event_type\":1,\"event_info\":{\"title\":\"subject\",\"focus content\":\" focus  message\",\"image_num\":\"image_num\",\"forum\":\"金毛俱乐部\",\"tag\":0,\"source\":1,\"fid\":36,\"bid\":34},\"tid\":0,\"status\":1,\"action\":0,\"time\":\"2017-10-23 10:54:00\"}"
	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	fmt.Println(f.Do())
}
