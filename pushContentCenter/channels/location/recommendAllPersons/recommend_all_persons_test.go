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
	"gouminGitlab/common/orm/elasticsearchBase"
	"github.com/olivere/elastic"
)

func testConn() ([]*xorm.Engine, []*mgo.Session, *elastic.Client) {
	dbAuth := "dog123:dog123"
	dbDsn := "192.168.86.193:3307"
	// dbDsn := "210.14.154.117:33068"
	dbName := "new_dog123"
	dataSourceName := dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil, nil,nil
	}

	mongoConn := "192.168.86.80:27017"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		fmt.Println("[error] connect mongodb err")
		return nil, nil,nil
	}

	var nodes []string
	nodes = append(nodes, "http://192.168.86.230:9200")
	nodes = append(nodes, "http://192.168.86.231:9200")
	r,_ := elasticsearchBase.NewClient(nodes)
	esConn,_ :=r.Run()

	var engineAry []*xorm.Engine
	engineAry = append(engineAry, engine)
	var sessionAry []*mgo.Session
	sessionAry = append(sessionAry, session)
	return engineAry, sessionAry,esConn
	// return engine, session
}

func jsonData() *job.RecommendJsonColumn {
	var jsonData job.RecommendJsonColumn
	jsonData.Uid = 2060500
	jsonData.Type = 1
	jsonData.Created = 1516598106
	jsonData.Infoid = 3399349
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
	jsonData.Action = -1
	return &jsonData
}

func TestPushPerson(t *testing.T) {
	mysqlXorm, mongoConn,esConn := testConn()
	jsonData := jsonData()
	f := NewRecommendAllPersons(mysqlXorm, mongoConn,esConn, jsonData)
	fmt.Println(f.pushPerson(881050))
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn,esConn := testConn()
	jsonData := jsonData()
	f := NewRecommendAllPersons(mysqlXorm, mongoConn,esConn, jsonData)
	fmt.Println(f.Do())
}

func TestRemoveInfoByTables(t *testing.T) {
	mysqlXorm, mongoConn,esConn := testConn()
	jsonData := jsonData()
	f := NewRecommendAllPersons(mysqlXorm, mongoConn,esConn, jsonData)
	fmt.Println(f.removeInfoByTables())
}
