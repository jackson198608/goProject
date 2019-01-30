package user

import (
	// "errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"testing"
	"strings"
)

const dbAuth = "dog123:dog123"
const dbDsn = "192.168.86.194:3307"
const dbName = "new_dog123"
const mongoConn ="192.168.86.80:27017,192.168.86.81:27017,192.168.86.82:27017" //"192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019"
const elkDsn = "http://192.168.86.231:9200,http://192.168.86.230:9200"     //"192.168.86.5:9200"

func testConn() ([]*xorm.Engine, []*mgo.Session) {
	dbAuth := "dog123:dog123"
	dbDsn := "192.168.86.194:3307" //"210.14.154.117:33068"
	dbName := "new_dog123"
	dataSourceName := dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	// mongoConn := "192.168.86.192:27017"
	// session, err := mgo.Dial(mongoConn)
	// if err != nil {
	// 	fmt.Println("[error] connect mongodb err")
	// 	return nil, nil
	// }

	Host := []string{
		"192.168.86.80:27017",
		"192.168.86.81:27018",
		"192.168.86.82:27019",
	}
	const (
		Username       = ""
		Password       = ""
		Database       = ""
		ReplicaSetName = "goumin"
	)
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:          Host,
		Username:       Username,
		Password:       Password,
		Database:       Database,
		ReplicaSetName: ReplicaSetName,
	})
	if err != nil {
		return nil, nil
	}
	// return engine, session
	var engineAry []*xorm.Engine
	engineAry = append(engineAry, engine)
	var sessionAry []*mgo.Session
	sessionAry = append(sessionAry, session)
	return engineAry, sessionAry

}

func TestGetMyData(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	uid := "2060500"
	elkNodes := strings.SplitN(elkDsn, ",", -1)
	// c := NewUser(mysqlXorm, mongoConn, uid, "210.14.154.117:8986")
	c := NewUser(mysqlXorm, mongoConn, uid, elkNodes)
	fmt.Println(c.getMyData())
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	uid := "2060500"
	elkNodes := strings.SplitN(elkDsn, ",", -1)
	// c := NewUser(mysqlXorm, mongoConn, uid, "210.14.154.117:8986")
	c := NewUser(mysqlXorm, mongoConn, uid, elkNodes)
	fmt.Println(c.Do())
}

func TestRecommendUserBySpecies(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	uid := "2060500"
	elkNodes := strings.SplitN(elkDsn, ",", -1)
	c := NewUser(mysqlXorm, mongoConn, uid, elkNodes)
	fmt.Println(c.recommendUserBySpecies(0, 5))
}

func TestFollowClubs(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	uid := "2060500"
	elkNodes := strings.SplitN(elkDsn, ",", -1)
	c := NewUser(mysqlXorm, mongoConn, uid, elkNodes)
	fmt.Println(c.followClubs())
}

func TestRecommendClubBySpecies(t *testing.T) {
	mysqlXorm, mongoConn := testConn()
	uid := "2060500"
	elkNodes := strings.SplitN(elkDsn, ",", -1)
	c := NewUser(mysqlXorm, mongoConn, uid, elkNodes)
	c.getMyData() //获取我的数据
	fmt.Println(c.recommendClubBySpecies())
}

//{"size" : 5,"query": {"query_string":{"query":"\"法国斗牛\",\"金毛\"","fields":["pets"]}},"filter" : {"bool":{"must_not":{"term":{"id":2060500}}}},"sort": { "lastlogintime": { "order": "desc" }}}
