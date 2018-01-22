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
	//get all active user from hashmap
	f.pushPersons(f.activeUserData)
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
	tableNameX := getTableNum(person)
	c := f.mongoConn[0].DB("RecommendData").C(tableNameX)
	if f.jsonData.Action == 0 {
		logger.Info("recommend data push uid is ", person, ", type is ", f.jsonData.Type, ", infoid is ", f.jsonData.Infoid)
		err := f.insertPerson(c, person)
		if err != nil {
			return err
		}
	} else if f.jsonData.Action == -1 {
		//删除数据
		logger.Info(" recommend data remove by uid is ", person, ", type is ", f.jsonData.Type, ", infoid is ", f.jsonData.Infoid)
		err := f.removePerson(c, person)
		if err != nil {
			return err
		}
	}
	return nil
}

func getTableNum(person int) string {
	tableNumX := person % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "user_recommend_" + strconv.Itoa(tableNumX) //用户推荐表
	return tableNameX
}

func (f *RecommendAllPersons) insertPerson(c *mgo.Collection, person int) error {
	//新增数据
	var data RecommendData.UserRecommendX
	data = RecommendData.UserRecommendX{bson.NewObjectId(),
		f.jsonData.Pid,
		f.jsonData.Uid,
		person,
		f.jsonData.Type,
		f.jsonData.Infoid,
		f.jsonData.Title,
		f.jsonData.Description,
		f.jsonData.Created,
		f.jsonData.Images,
		f.jsonData.Imagenums,
		f.jsonData.Tag,
		f.jsonData.Tags,
		f.jsonData.QstType,
		f.jsonData.AdType,
		f.jsonData.AdUrl,
		f.jsonData.Channel,
		f.jsonData.Rauth}
	err := c.Insert(&data) //插入数据
	if err != nil {
		return err
	}
	return nil
}

func (f *RecommendAllPersons) removePerson(c *mgo.Collection, person int) error {
	//删除数据
	_, err := c.RemoveAll(bson.M{"type": f.jsonData.Type, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid, "channel": f.jsonData.Channel})
	if err != nil {
		return err
	}
	return nil
}
