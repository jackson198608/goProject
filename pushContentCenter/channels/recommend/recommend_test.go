package recommend

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

	mongoConn := "192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019"
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
	jobStr := "{\"uid\":2060500,\"recommend_type\":\"all\",\"type\":1,\"title\":\"subject\",\"description\":\"aaaaaa\",\"image_num\":0,\"image_num\":\"\",\"tags\":\"\",\"tag\":0,\"created\":1508469600,\"channel\":1}"
	f := NewRecommend(mysqlXorm, mongoConn, jobStr)
	f.parseJson()
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":2060500,\"recommend_type\":\"all\",\"type\":1,\"title\":\"subject\",\"description\":\"aaaaaa\",\"image_num\":0,\"image_num\":\"\",\"tags\":\"\",\"tag\":0,\"created\":1508469600,\"channel\":1}"
	f := NewRecommend(mysqlXorm, mongoConn, jobStr)
	fmt.Println(f.Do())
}
