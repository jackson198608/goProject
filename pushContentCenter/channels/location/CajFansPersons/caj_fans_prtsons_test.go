package CajFansPersons

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
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
	dataSourceName := dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil
	}
	var engineAry []*xorm.Engine
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
	mongoConn := "192.168.86.80:27017"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		fmt.Println("[error] connect mongodb err")
		return nil, nil, nil
	}

	var nodes []string
	nodes = append(nodes, "http://192.168.86.230:9200")
	nodes = append(nodes, "http://192.168.86.231:9200")
	r, _ := elasticsearchBase.NewClient(nodes)
	esConn, _ := r.Run()

	var sessionAry []*mgo.Session
	sessionAry = append(sessionAry, session)
	return engineAry, sessionAry, esConn
	// return engine, session
}

func jsonData() *job.FocusJsonColumn {
	var tag []string
	tag = append(tag, "26")
	tag = append(tag, "27")
	var jsonData job.FocusJsonColumn
	jsonData.Uid = 2265027
	jsonData.TypeId = 36
	jsonData.Created = "2019-03-07 22:54:20"
	jsonData.Tid = 0
	jsonData.Bid = 36
	jsonData.Infoid = 2435
	jsonData.Title = ""
	jsonData.Content = ""
	jsonData.Forum = ""
	jsonData.Imagenums = 0
	jsonData.Tag = 0
	jsonData.Qsttype = 0
	jsonData.Fid = 0
	jsonData.Source = 2
	jsonData.Status = 1
	jsonData.Action = -1
	jsonData.AdoptId = 2435
	jsonData.PetName = "宠物名称1"
	jsonData.PetAge = "61"
	jsonData.PetBreed = 1
	jsonData.PetGender = 7
	jsonData.PetSpecies = "拉布拉多"
	jsonData.Province = "北京"
	jsonData.City = "北京市"
	jsonData.County = "昌平区"
	jsonData.Reason = "9"
	jsonData.Image = "aaaaaa"
	jsonData.PetImmunity = 16
	jsonData.PetExpelling = 18
	jsonData.PetSterilization = 19
	jsonData.PetStatus = 19
	jsonData.AdoptStatus = 1
	jsonData.PetIntroduction = "20"
	jsonData.UserIdentity = 1
	jsonData.AdoptTag = tag
	return &jsonData
}

func TestDo(t *testing.T) {
	mysqlXorm, mongoConn, esConn := testConn()
	jsonData := jsonData()

	f := NewCajFansPersons(mysqlXorm, mongoConn, jsonData, esConn)
	fmt.Println(f.Do())
}
