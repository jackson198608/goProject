package fansPersons

import (
	"errors"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/FansData"
	"gouminGitlab/common/orm/mysql/new_dog123"
	// "reflect"
	"strconv"
)

type FansPersons struct {
	mysqlXorm      []*xorm.Engine //@todo to be []
	mongoConn      []*mgo.Session //@todo to be []
	jsonData       *job.FocusJsonColumn
	activeUserData *map[int]bool
}

const count = 1000

func NewFansPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, activeUserData *map[int]bool) *FansPersons {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jsonData == nil) {
		return nil
	}

	f := new(FansPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jsonData = jsonData
	f.activeUserData = activeUserData

	return f
}

func (f *FansPersons) Do() error {
	startId := 0
	for {
		//获取粉丝用户
		currentPersionList := f.getPersons(startId)
		if currentPersionList == nil {
			return nil
		}
		endId, err := f.pushPersons(currentPersionList)
		startId = endId
		if err != nil {
			return err
		}
		if len(*currentPersionList) < count {
			break
		}
	}
	return nil
}

func (f *FansPersons) pushPersons(follows *[]new_dog123.Follow) (int, error) {
	if follows == nil {
		return 0, errors.New("push to fans active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	active_user := *f.activeUserData
	persons := *follows

	var endId int
	for _, person := range persons {
		//check key in actice user
		_, ok := active_user[person.FollowId]
		if ok {
			err := f.pushPerson(person.FollowId)
			if err != nil {
				for i := 0; i < 5; i++ {
					err := f.pushPerson(person.FollowId)
					if err == nil {
						break
					}
				}
			}
			endId = person.Id
		}
	}

	return endId, nil
}

func getTableNum(person int) string {
	tableNumX := person % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	return tableNameX
}

func (f *FansPersons) pushPerson(person int) error {
	tableNameX := getTableNum(person)
	c := f.mongoConn[0].DB("FansData").C(tableNameX)
	if f.jsonData.Action == 0 {
		// fmt.Println("insert" + strconv.Itoa(person))
		err := f.insertPerson(c, person)
		if err != nil {
			return err
		}
	} else if f.jsonData.Action == 1 {
		//修改数据
		// fmt.Println("update" + strconv.Itoa(person))
		err := f.updatePerson(c, person)
		if err != nil {
			return err
		}
	} else if f.jsonData.Action == -1 {
		//删除数据
		// fmt.Println("remove" + strconv.Itoa(person))
		err := f.removePerson(c, person)
		if err != nil {
			return err
		}
	}
	return nil
}

//get fans persons by uid
func (f *FansPersons) getPersons(startId int) *[]new_dog123.Follow {
	// var persons []int
	var follows []new_dog123.Follow
	err := f.mysqlXorm[0].Where("user_id=? and id>? and fans_active=1", f.jsonData.Uid, startId).Asc("id").Limit(count).Find(&follows)
	if err != nil {
		return nil
	}

	return &follows
}

func (f *FansPersons) insertPerson(c *mgo.Collection, person int) error {
	//新增数据
	var data FansData.EventLog
	data = FansData.EventLog{bson.NewObjectId(),
		f.jsonData.TypeId,
		f.jsonData.Uid,
		person,
		f.jsonData.Created,
		f.jsonData.Infoid,
		f.jsonData.Status,
		f.jsonData.Tid,
		f.jsonData.Bid,
		f.jsonData.Content,
		f.jsonData.Title,
		f.jsonData.Imagenums,
		f.jsonData.Forum,
		f.jsonData.Tag,
		f.jsonData.Qsttype,
		f.jsonData.Source}
	err := c.Insert(&data) //插入数据
	if err != nil {
		return err
	}
	return nil
}

func (f *FansPersons) updatePerson(c *mgo.Collection, person int) error {
	_, err := c.UpdateAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid}, bson.M{"$set": bson.M{"status": f.jsonData.Status}})
	if err != nil {
		return err
	}
	return nil
}

func (f *FansPersons) removePerson(c *mgo.Collection, person int) error {
	_, err := c.RemoveAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid, "tid": f.jsonData.Tid})
	if err != nil {
		return err
	}
	return nil
}
