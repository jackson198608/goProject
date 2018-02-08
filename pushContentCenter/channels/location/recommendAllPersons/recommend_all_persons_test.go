package recommendAllPersons

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	mgo "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	// "gouminGitlab/common/orm/mongo/FansData"
	// "reflect"
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

func jsonData() *job.RecommendJsonColumn {
	var jsonData job.RecommendJsonColumn
	jsonData.Uid = 2060500
	jsonData.Type = 1
	jsonData.Created = 1516598106
	jsonData.Infoid = 3399348
	jsonData.Title = "所有活跃用户推送title"
	jsonData.Description = "所有活跃用户推送正文正文"
	jsonData.Images = ""
	jsonData.Imagenums = 0
	jsonData.Tag = 0
	jsonData.Tags = ""
	jsonData.QstType = 0
	jsonData.AdType = 0
	jsonData.AdUrl = ""
	jsonData.Channel = 1
	jsonData.Action = 0
	return &jsonData
}

var m map[int]bool

func init() {
	m = make(map[int]bool)

	mongoConn := "192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		// return m
	}

	var uids []int
	c := session.DB("ActiveUser").C("active_user")
	err = c.Find(nil).Distinct("uid", &uids)
	if err != nil {
		// panic(err)
		// return m
	}
	for i := 0; i < len(uids); i++ {
		m[uids[i]] = true
	}
}

func TestPushPerson(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()
	f := NewRecommendAllPersons(mysqlXorm, mongoConn, jsonData, &m)
	fmt.Println(f.pushPerson(881050))
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()
	f := NewRecommendAllPersons(mysqlXorm, mongoConn, jsonData, &m)
	fmt.Println(f.Do())
}

func TestRemoveInfoByTables(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()
	f := NewRecommendAllPersons(mysqlXorm, mongoConn, jsonData, &m)
	fmt.Println(f.removeInfoByTables())
}

func TestRemoveinfo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()
	tableNameX := "user_recommend_100"
	c := mongoConn[0].DB("RecommendData").C(tableNameX)
	f := NewRecommendAllPersons(mysqlXorm, mongoConn, jsonData, &m)
	fmt.Println(f.removeInfo(c))
}
