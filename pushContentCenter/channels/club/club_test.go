package club

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

func testConn() ([]*xorm.Engine, []*mgo.Session) {
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
	// return engine, session
	var engineAry []*xorm.Engine
	engineAry = append(engineAry, engine)
	var sessionAry []*mgo.Session
	sessionAry = append(sessionAry, session)
	return engineAry, sessionAry

}

func TestGetClubs(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060500,\"type\":1,\"typeid\":2,\"subject\":\"subject\",\"message\":\"message\",\"image_num\":\"image_num\",\"lastpost\":2,\"fid\":36,\"lastposter\":\"0ssss\",\"status\":1,\"displayorder\":1,\"digest\":1,\"qst_type\":0,\"created\":1508469600}"
	c := NewClub(mysqlXorm, mongoConn, jobStr)
	fmt.Println(c.getClubs())
}

func TestPushClub(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060501,\"heats\":1,\"type\":6,\"infoid\":2234567,\"typeid\":2,\"subject\":\"subject\",\"message\":\" push to 36 message\",\"image_num\":\"image_num\",\"lastpost\":881051,\"fid\":\"37,38,77\",\"lastposter\":\"123asad232344\",\"status\":1,\"displayorder\":1,\"digest\":1,\"qst_type\":0,\"created\":1508469600,\"action\":1,\"replies\":1,\"price\":1,\"isgroup\":1,\"special\":0,\"recommends\":3,\"sortid\":12,\"highlight\":1,\"closed\":1,\"cover\":2,\"thread_status\":256}"
	c := NewClub(mysqlXorm, mongoConn, jobStr)
	fmt.Println(c.pushClub(34))
}

func TestPushClubs(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060501,\"heats\":1,\"type\":6,\"infoid\":2234567,\"typeid\":2,\"subject\":\"subject\",\"message\":\" push to 36 message\",\"image_num\":\"image_num\",\"lastpost\":881051,\"fid\":\"37,38,77\",\"lastposter\":\"123asad2323455\",\"status\":1,\"displayorder\":1,\"digest\":1,\"qst_type\":0,\"created\":1508469600,\"action\":1,\"replies\":1,\"price\":1,\"isgroup\":1,\"special\":0,\"recommends\":3,\"sortid\":12,\"highlight\":1,\"closed\":1,\"cover\":2,\"thread_status\":256}"
	c := NewClub(mysqlXorm, mongoConn, jobStr)

	var clubs = []int{34, 36, 44, 52}

	fmt.Println(c.pushClubs(clubs))
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060500,\"heats\":1,\"type\":6,\"infoid\":2234567,\"typeid\":2,\"subject\":\"asdasdasda\",\"message\":\" push to 36 message\",\"lastpost\":1510716086,\"fid\":\"37\",\"lastposter\":\"12我去玩若323578\",\"displayorder\":1,\"digest\":1,\"qst_type\":0,\"created\":1508469600,\"action\":-1,\"replies\":1,\"price\":1,\"isgroup\":1,\"special\":0,\"recommends\":3,\"sortid\":12,\"highlight\":1,\"closed\":1,\"cover\":2,\"thread_status\":256}"
	// jobStr := "{\"type\":6,\"infoid\":2234567,\"fid\":\"37\",\"heats\":2,\"typeid\":22,\"message\":\" push to 321wsa6 message\",\"subject\":\"asdsadsddddddddd\",\"action\":1,\"image_num\":2,\"lastpost\":1508469900,\"lastposter\":\"12我去玩若32asdas3578\",\"displayorder\":2,\"digest\":2,\"qst_type\":1,\"replies\":0,\"price\":12,\"isgroup\":0,\"special\":1,\"recommends\":0,\"sortid\":10,\"highlight\":2,\"closed\":0,\"cover\":1,\"thread_status\":128,\"status\":1}"
	c := NewClub(mysqlXorm, mongoConn, jobStr)
	fmt.Println(c.Do())
}

func TestAddBsonMap(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060501,\"heats\":1,\"type\":6,\"infoid\":2234567,\"typeid\":2,\"message\":\" push to 36 message\",\"lastpost\":1510716086,\"fid\":\"37\",\"lastposter\":\"12我去玩若323578\",\"displayorder\":1,\"digest\":1,\"qst_type\":0,\"created\":1508469600,\"action\":1,\"replies\":1,\"price\":1,\"isgroup\":1,\"special\":0,\"recommends\":3,\"sortid\":12,\"highlight\":1,\"closed\":1,\"cover\":2,\"thread_status\":256}"
	c := NewClub(mysqlXorm, mongoConn, jobStr)

	query := c.updateBsonMap()
	fmt.Println(query)
}
