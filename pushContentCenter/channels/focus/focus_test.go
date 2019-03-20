package focus

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"testing"
	"gouminGitlab/common/orm/elasticsearchBase"
	"github.com/olivere/elastic"
)

func testConn() ([]*xorm.Engine, []*mgo.Session, *elastic.Client) {
	dbAuth := "dog123:dog123"
	dbDsn := "192.168.86.194:3307"
	// dbDsn := "210.14.154.117:33068"
	dbName := "new_dog123"
	dbName1 := "card"
	dbName2 := "adoption"
	var engineAry []*xorm.Engine
	dataSourceName := dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil, nil,nil
	}
	engineAry = append(engineAry, engine)

	if dbName1!="" {
		dataSourceName1 := dbAuth + "@tcp(" + dbDsn + ")/" + dbName1 + "?charset=utf8mb4"
		engine1, err := xorm.NewEngine("mysql", dataSourceName1)
		if err != nil {
			fmt.Println(err)
			return nil, nil,nil
		}
		engineAry = append(engineAry, engine1)
	}

	if dbName2!="" {
		dataSourceName2 := dbAuth + "@tcp(" + dbDsn + ")/" + dbName2 + "?charset=utf8mb4"
		engine2, err := xorm.NewEngine("mysql", dataSourceName2)
		if err != nil {
			fmt.Println(err)
			return nil, nil,nil
		}
		engineAry = append(engineAry, engine2)
	}

	//mongoConn := "192.168.86.80:27017"
	//session, err := mgo.Dial(mongoConn)
	//if err != nil {
	//	fmt.Println("[error] connect mongodb err")
	//	return nil, nil,nil
	//}

	var nodes []string
	nodes = append(nodes, "http://192.168.86.230:9200")
	nodes = append(nodes, "http://192.168.86.231:9200")
	r,_ := elasticsearchBase.NewClient(nodes)
	esConn,_ :=r.Run()

	var sessionAry []*mgo.Session
	//sessionAry = append(sessionAry, session)
	//Init()
	return engineAry, sessionAry, esConn
	// return engine, session
}

func TestParseJson(t *testing.T) {
	var nodes []string
	nodes = append(nodes, "http://192.168.86.230:9200")
	nodes = append(nodes, "http://192.168.86.231:9200")
	mysqlXorm, mongoConn, esConn := testConn()
	//jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36},\"tid\":0,\"status\":1,\"time\":1508469600}"
	jobStr := "{\"uid\":2060500,\"event_type\":30,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36,\"pet_id\":71,\"pet_type\":1},\"tid\":0,\"status\":1,\"time\":1508469600}"
	f := NewFocus(mysqlXorm, mongoConn, jobStr, esConn)
	fmt.Println(f.parseJson())
}

func TestDo0(t *testing.T) {
	mysqlXorm, mongoConn,esConn := testConn()

	//jobStr :="{\"uid\":2417364,\"infoid\":56921,\"event_type\":30,\"event_info\":{\"content\":\"\\u6d4b\\u9996\\u9875\",\"images\":\"7916\",\"source\":11},\"pet_id\":6351,\"pet_type\":2,\"status\":1,\"time\":\"2018-07-18 11:54:38\",\"action\":0,\"is_video\":0}|focus"
	jobStr := "{\"uid\":881050,\"event_type\":8,\"event_info\":{\"title\":\"subject\",\"focus content\":\" focus  message\",\"image_num\":\"image_num\",\"forum\":\"金毛俱乐部\",\"tag\":0,\"source\":1,\"fid\":36,\"bid\":60},\"tid\":0,\"status\":1,\"action\":0,\"time\":\"2017-10-23 10:54:00\"}"
	//jobStr := "{\"uid\":2060500,\"event_type\":30,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"1\",\"images\":\"http://img1.goumin.com\",\"video_url\":\"http://video.goumin.com\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36},\"tid\":0,\"status\":1,\"pet_id\":71,\"pet_type\":1,\"time\":\"2017-10-23 10:54:00\"}"
	//jobStr := "{\"uid\":2060500,\"event_type\":2,\"event_info\":{\"title\":\"subject\",\"content\":\"message\",\"image_num\":\"image_num\",\"forum\":\"forum->name\",\"tag\":0,\"source\":2,\"fid\":36},\"tid\":0,\"status\":1,\"time\":1508469600}"

	f := NewFocus(mysqlXorm, mongoConn, jobStr, esConn)
	fmt.Println(f.Do())
}

func TestDo1(t *testing.T) {
	mysqlXorm, mongoConn,esConn := testConn()

	jobStr := "{\"uid\":881050,\"event_type\":1,\"event_info\":{\"title\":\"subject\",\"event type 1focus content\":\" focus  message\",\"image_num\":\"image_num\",\"forum\":\"金毛俱乐部\",\"tag\":0,\"source\":1,\"fid\":36,\"bid\":34},\"tid\":0,\"status\":1,\"action\":0,\"time\":\"2017-10-23 10:54:00\"}"
	f := NewFocus(mysqlXorm, mongoConn, jobStr, esConn)
	fmt.Println(f.Do())
}

func TestDo2(t *testing.T) {
	mysqlXorm, mongoConn,esConn := testConn()

	jobStr := "{\"uid\":881050,\"event_type\":18,\"event_info\":{\"title\":\"subject\",\"event type 1focus content\":\" focus  message\",\"image_num\":\"image_num\",\"forum\":\"金毛俱乐部\",\"tag\":0,\"source\":1,\"fid\":36,\"bid\":34},\"infoid\":123,\"tid\":0,\"status\":1,\"action\":0,\"time\":\"2017-10-23 10:54:00\"}"
	f := NewFocus(mysqlXorm, mongoConn, jobStr, esConn)
	fmt.Println(f.Do())
}


func TestDo3(t *testing.T) {
	mysqlXorm, mongoConn,esConn := testConn()

	jobStr := "{\"time\":1550387388,\"event_info\":{\"adopt_id\":2435,\"pet_name\":\"5\",\"pet_age\":\"6\",\"pet_breed\":1,\"pet_gender\":7,\"pet_species\":\"\u6bd4\u683c\",\"province\":\"11\",\"city\":\"12\",\"county\":\"13\",\"reason\":\"9\",\"image\":\"\",\"is_video\":1,\"pet_immunity\":16,\"pet_expelling\":17,\"pet_sterilization\":18,\"pet_status\":19,\"pet_agenum\":15,\"adopt_status\":0,\"pet_introduction\":\"20\",\"user_identity\":0,\"adopt_tag\":[\"26\",\"27\"]},\"uid\":2265027,\"event_type\":36,\"infoid\":2435,\"status\":1,\"action\":1}|focus"
	f := NewFocus(mysqlXorm, mongoConn, jobStr, esConn)
	fmt.Println(f.Do())
}
