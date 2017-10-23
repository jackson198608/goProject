package club

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

func testConn() (*xorm.Engine, *mgo.Session) {
	dbAuth := "dog123:dog123"
	dbDsn := "210.14.154.117:33068"
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
	return engine, session
}

func TestGetClubs(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060500,\"type\":1,\"typeid\":2,\"subject\":\"subject\",\"message\":\"message\",\"image_num\":\"image_num\",\"lastpost\":2,\"fid\":36,\"lastposter\":\"0ssss\",\"status\":1,\"displayorder\":1,\"disgest\":1,\"qst_type\":0,\"created\":1508469600}|1|0"
	c := NewClub(mysqlXorm, mongoConn, jobStr)
	fmt.Println(c.getClubs())
}

func TestPushData(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060500,\"type\":1,\"typeid\":2,\"subject\":\"subject\",\"message\":\"message\",\"image_num\":\"image_num\",\"lastpost\":2,\"fid\":36,\"lastposter\":\"0ssss\",\"status\":1,\"displayorder\":1,\"disgest\":1,\"qst_type\":0,\"created\":1508469600}|1|0"
	c := NewClub(mysqlXorm, mongoConn, jobStr)
	fmt.Println(c.pushData(34))
}

func TestTryPushClub(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060500,\"type\":1,\"typeid\":2,\"subject\":\"subject\",\"message\":\"message\",\"image_num\":\"image_num\",\"lastpost\":2,\"fid\":36,\"lastposter\":\"0ssss\",\"status\":1,\"displayorder\":1,\"disgest\":1,\"qst_type\":0,\"created\":1508469600}|1|0"
	c := NewClub(mysqlXorm, mongoConn, jobStr)
	fmt.Println(c.tryPushClub(34, 1))
}

func TestPushClub(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060501,\"type\":1,\"infoid\":1234567,\"typeid\":2,\"subject\":\"subject\",\"message\":\"message\",\"image_num\":\"image_num\",\"lastpost\":2,\"fid\":36,\"lastposter\":\"0ssss\",\"status\":1,\"displayorder\":1,\"disgest\":1,\"qst_type\":0,\"created\":1508469600}|1|0"
	c := NewClub(mysqlXorm, mongoConn, jobStr)
	fmt.Println(c.pushClub(34))
}

func TestPushClubs(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060501,\"type\":1,\"infoid\":1234567,\"typeid\":2,\"subject\":\"subject\",\"message\":\"message\",\"image_num\":\"image_num\",\"lastpost\":2,\"fid\":36,\"lastposter\":\"0ssss\",\"status\":1,\"displayorder\":1,\"disgest\":1,\"qst_type\":0,\"created\":1508469600}|1|0"
	c := NewClub(mysqlXorm, mongoConn, jobStr)

	var clubs = []int{34, 36, 44, 52}

	fmt.Println(c.pushClubs(clubs))
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060501,\"type\":1,\"infoid\":1234567,\"typeid\":2,\"subject\":\"subject\",\"message\":\"message\",\"image_num\":\"image_num\",\"lastpost\":2,\"fid\":\"\",\"lastposter\":\"0ssss\",\"status\":1,\"displayorder\":1,\"disgest\":1,\"qst_type\":0,\"created\":1508469600}|1|0"
	c := NewClub(mysqlXorm, mongoConn, jobStr)

	fmt.Println(c.Do())
}
