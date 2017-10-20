package focus

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

func TestParseJson(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36},\"tid\":0,\"status\":1,\"time\":1508469600}|1|0"
	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	f.parseJson()
}

func TestFansPersons(t *testing.T) {
	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36},\"tid\":0,\"status\":1,\"time\":1508469600}|1|0"

	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	f.getFansPersons(1, 10)
}

func TestClubPersons(t *testing.T) {
	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36},\"tid\":0,\"status\":1,\"time\":1508469600}|1|0"

	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	f.getClubPersons(1)
}

func TestGetClubPersonPageNum(t *testing.T) {
	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36,\"bid\":34},\"tid\":0,\"status\":1,\"time\":1508469600}|1|0"

	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	page := f.getClubPersonPageNum()
	fmt.Println(page)
}

func TestGetBreedPersonsPagNum(t *testing.T) {
	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36,\"bid\":34},\"tid\":0,\"status\":1,\"time\":1508469600}|1|0"

	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	page := f.getBreedPersonsPagNum()
	fmt.Println(page)
}

func TestBreedPersons(t *testing.T) {
	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36,\"bid\":34},\"tid\":0,\"status\":1,\"time\":1508469600}|1|0"

	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	f.getBreedPersons(1)
}

func TestFansPersonFirstId(t *testing.T) {

	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36,\"bid\":34},\"tid\":0,\"status\":1,\"time\":1508469600}|1|0"

	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	f.getFansPersonFirstId()
}

func TestFansPersonLastId(t *testing.T) {
	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36,\"bid\":34},\"tid\":0,\"status\":1,\"time\":1508469600}|1|0"

	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	f.getFansPersonLastId()
}

func TestGetPersons(t *testing.T) {
	mysqlXorm, mongoConn := testConn()

	jobStr := "{\"uid\":881050,\"event_type\":6,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36,\"bid\":34},\"tid\":0,\"status\":1,\"time\":1508469600}|1|0"

	f := NewFocus(mysqlXorm, mongoConn, jobStr)
	f.getPersons(0, 1, 999999999)
}
