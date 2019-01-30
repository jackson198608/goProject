package allPersons

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
	"gouminGitlab/common/orm/elasticsearch"
	"encoding/json"
)

func testConn() ([]*xorm.Engine, []*mgo.Session) {
	dbAuth := "dog123:dog123"
	dbDsn := "192.168.86.194:3307"
	// dbDsn := "210.14.154.117:33068"
	dbName := "new_dog123"
	dataSourceName := dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	mongoConn := "192.168.86.80:27017"
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
	jsonData.Created = "2017-10-23 22:54:11"
	jsonData.Tid = 0
	jsonData.Bid = 36
	jsonData.Infoid = 234567
	jsonData.Title = "所有活跃用户推送title"
	jsonData.Content = "所有活跃用户推送正文正文"
	jsonData.Forum = "36club"
	jsonData.Imagenums = 0
	jsonData.Tag = 0
	jsonData.Qsttype = 0
	jsonData.Fid = 0
	jsonData.Source = 2
	jsonData.Status =-1
	jsonData.Action = -1
	return &jsonData
}

var m map[int]bool

func init() {
	var nodes []string
	nodes = append(nodes, "http://192.168.86.230:9200")
	nodes = append(nodes, "http://192.168.86.231:9200")
	m = make(map[int]bool)
	er := elasticsearch.NewUser(nodes)
	from := 0
	size := 1000
	rst := er.SearchAllActiveUser(from, size)
	if rst.Hits.TotalHits >0 {
		for _,hit := range rst.Hits.Hits{
			var userinfo elasticsearch.UserData
			err := json.Unmarshal(*hit.Source, &userinfo) //另外一种取数据的方法
			if err != nil {
				fmt.Println("Deserialization failed")
			}
			m[userinfo.Id] = true
		}
	}
}

func TestDo(t *testing.T) {
	var nodes []string
	nodes = append(nodes, "http://192.168.86.230:9200")
	nodes = append(nodes, "http://192.168.86.231:9200")
	mysqlXorm, mongoConn := testConn()
	jsonData := jsonData()
	f := NewAllPersons(mysqlXorm, mongoConn, jsonData, &m,nodes)
	fmt.Println(f.Do())
}
