package content

import (
	// "errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/RecommendData"
	"testing"
	"time"
)

const dbAuth = "dog123:dog123"
const dbDsn = "192.168.86.193:3307"
const dbName = "new_dog123"
const mongoConn = "192.168.86.192:27017" //"192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019"
const elkDsn = "210.14.154.117:8986"     //"192.168.86.5:9200"

func testConn() ([]*xorm.Engine, []*mgo.Session) {
	dbAuth := "dog123:dog123"
	dbDsn := "192.168.0.110:3306" //"210.14.154.117:33068"
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
	// return engine, session
	var engineAry []*xorm.Engine
	engineAry = append(engineAry, engine)
	var sessionAry []*mgo.Session
	sessionAry = append(sessionAry, session)
	return engineAry, sessionAry

}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	uid := "2060500"
	// c := NewUser(mysqlXorm, mongoConn, uid, "210.14.154.117:8986")
	c := NewContent(mysqlXorm, mongoConn, uid)
	fmt.Println(c.Do())
}

func TestGetContent(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	uid := "2060500"
	// c := NewUser(mysqlXorm, mongoConn, uid, "210.14.154.117:8986")
	c := NewContent(mysqlXorm, mongoConn, uid)
	fmt.Println(c.getContents())
}

func TestInsert(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	uid := "2060500"
	created := int(time.Now().Unix())
	// c := NewUser(mysqlXorm, mongoConn, uid, "210.14.154.117:8986")
	c := NewContent(mysqlXorm, mongoConn, uid)

	mongoConn1 := "192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019"
	session, err := mgo.Dial(mongoConn1)
	if err != nil {
		fmt.Println("[error] connect mongodb err")
	}

	mc := session.DB("RecommendData").C("user_recommend_100")
	fmt.Println("start time:")
	fmt.Println(created)
	var content RecommendData.ContentChannel
	content.Id = bson.NewObjectId()
	content.Uid = 1499668
	content.Title = "林肯小朋友之冲出游泳馆，奔向大海边"
	content.Description = ""
	content.Channel = 1
	content.Tag = 0
	content.Type = 1
	content.Rauth = ""
	content.QstType = 0
	content.AdUrl = ""
	content.AdType = 0
	content.Images = "[\"http:\\/\\/dev.img.goumintest.com\\/bbs\\/201602\\/10\\/201602100753384632.jpg\",\"http:\\/\\/dev.img.goumintest.com\\/bbs\\/201602\\/09\\/201602091809523708.jpg\"]"
	content.ImageNum = 2
	content.Tags = ""
	content.Pid = 59603004
	content.Created = created
	content.Infoid = 3344140

	for i := 0; i < 10000; i++ {
		fmt.Println(c.insertPerson(mc, 2060500, content))
	}
	fmt.Println("end time:")
	fmt.Println(time.Now().Unix())

}
