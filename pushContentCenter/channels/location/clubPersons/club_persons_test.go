package clubPersons

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	mgo "gopkg.in/mgo.v2"

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

func jsonData() *job.FocusJsonColumn {
	var jsonData job.FocusJsonColumn
	jsonData.Uid = 2060500
	jsonData.TypeId = 1
	jsonData.Created = "2017-10-23 22:54"
	jsonData.Tid = 0
	jsonData.Bid = 36
	jsonData.Infoid = 234567
	jsonData.Title = "相同俱乐部推送title"
	jsonData.Content = "相同俱乐部推送正文正文"
	jsonData.Forum = "36club"
	jsonData.Imagenums = 0
	jsonData.Tag = 0
	jsonData.Qsttype = 0
	jsonData.Fid = 34
	jsonData.Source = 2
	jsonData.Status = -1
	jsonData.Action = -1
	return &jsonData
}

func TestGetPersons(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()

	f := NewClubPersons(mysqlXorm, mongoConn, jsonData)
	fmt.Println(f.getPersons(2))
}

func TestTryPushPerson(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()

	f := NewClubPersons(mysqlXorm, mongoConn, jsonData)
	fmt.Println(f.tryPushPerson(881050, 6))
}

func TestPushPerson(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()

	f := NewClubPersons(mysqlXorm, mongoConn, jsonData)
	fmt.Println(f.pushPerson(881050))
}

func TestPushPersons(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()
	var persons = []int{2060500, 2060400}

	f := NewClubPersons(mysqlXorm, mongoConn, jsonData)
	fmt.Println(f.pushPersons(persons))
}

func TestGetPersonPageNum(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()

	f := NewClubPersons(mysqlXorm, mongoConn, jsonData)
	fmt.Println(f.getPersonPageNum())
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()
	f := NewClubPersons(mysqlXorm, mongoConn, jsonData)
	fmt.Println(f.Do())
}
