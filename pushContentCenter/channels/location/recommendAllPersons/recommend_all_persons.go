package recommendAllPersons

import (
	"errors"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/RecommendData"
	"strconv"
)

type RecommendAllPersons struct {
	mysqlXorm      []*xorm.Engine
	mongoConn      []*mgo.Session
	jsonData       *job.RecommendJsonColumn
	activeUserData *map[int]bool
}

const count = 1000

func NewRecommendAllPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.RecommendJsonColumn, activeUserData *map[int]bool) *RecommendAllPersons {
	if (mongoConn == nil) || (jsonData == nil) {
		return nil
	}

	f := new(RecommendAllPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jsonData = jsonData
	f.activeUserData = activeUserData

	return f
}

func (f *RecommendAllPersons) Do() error {
	if f.jsonData.Action == 0 {
		//get all active user from hashmap
		err := f.pushPersons(f.activeUserData)
		if err != nil {
			return err
		}
	} else if f.jsonData.Action == -1 {
		err := f.removeInfoByTables()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RecommendAllPersons) pushPersons(persons *map[int]bool) error {
	if persons == nil {
		return errors.New("push to all active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	for k := range *persons {
		err := f.pushPerson(k)
		if err != nil {
			for i := 0; i < 5; i++ {
				err = f.pushPerson(k)
				if err == nil {
					break
				}
			}
		}
	}
	return nil
}

func (f *RecommendAllPersons) pushPerson(person int) error {
	ur := elasticsearch.NewUserRecommendX(f.esConn, f.jsonData)
	err := ur.Create(person)
	if err != nil {
		return err
	}
	return nil
}

func (f *RecommendAllPersons) removeInfoByTables() error {
	for i := 1; i < 101; i++ {
		tableNameX := "user_recommend_" + strconv.Itoa(i)
		c := f.mongoConn[0].DB("RecommendData").C(tableNameX)
		err := f.removeInfo(c)
		logger.Info(" recommend data remove by connection is ", tableNameX, ", channel is ", f.jsonData.Channel, ", type is ", f.jsonData.Type, ", infoid is ", f.jsonData.Infoid)
		if err != nil {
			for n := 0; n < 5; n++ {
				err = f.removeInfo(c)
				logger.Info("[try next]", n, " recommend data remove by connection is ", tableNameX, ", channel is ", f.jsonData.Channel, ", type is ", f.jsonData.Type, ", infoid is ", f.jsonData.Infoid)
				if err == nil {
					break
				}
			}
		}
	}
	return nil
}

func (f *RecommendAllPersons) removeInfo(c *mgo.Collection) error {
	//删除数据
	_, err := c.RemoveAll(bson.M{"type": f.jsonData.Type, "infoid": f.jsonData.Infoid, "channel": f.jsonData.Channel})
	if err != nil {
		return err
	}
	return nil
}
