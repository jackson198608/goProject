package breedPersons

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/ActiveUser"
	"gouminGitlab/common/orm/mongo/FansData"
	"strconv"
)

const count = 1000

type BreedPersons struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jsonData  *job.FocusJsonColumn
}

func NewBreedPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn) *BreedPersons {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jsonData == nil) {
		return nil
	}

	f := new(BreedPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jsonData = jsonData

	return f
}

func (f *BreedPersons) Do() error {
	var startId bson.ObjectId
	startId = bson.ObjectId("000000000000")
	for {
		currentPersionList := f.getPersons(startId)
		if currentPersionList == nil {
			return nil
		}
		endId, err := f.pushPersons(currentPersionList)
		if err != nil {
			return err
		}
		startId = endId
		if len(*currentPersionList) < count {
			break
		}
	}
	return nil
}

func (f *BreedPersons) pushPersons(ActiveUser *[]ActiveUser.ActiveBreedUser) (bson.ObjectId, error) {
	if ActiveUser == nil {
		return bson.NewObjectId(), errors.New("push to breed active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	var endId bson.ObjectId
	persons := *ActiveUser
	for _, person := range persons {
		err := f.pushPerson(person.Uid)
		if err != nil {
			for i := 0; i < 5; i++ {
				err := f.pushPerson(person.Uid)
				if err == nil {
					break
				}
			}
		}
		endId = person.Id
	}
	return endId, nil
}

func (f *BreedPersons) pushPerson(person int) error {
	tableNameX := getTableNum(person)
	c := f.mongoConn[0].DB("FansData").C(tableNameX)

	if f.jsonData.Action == 0 {
		//fmt.Println("insert" + strconv.Itoa(person))
		err := f.insertPerson(c, person)
		if err != nil {
			return err
		}
	} else if f.jsonData.Action == 1 {
		//修改数据
		//fmt.Println("update" + strconv.Itoa(person))
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

func (f *BreedPersons) insertPerson(c *mgo.Collection, person int) error {
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
		f.jsonData.ImageInfo,
		f.jsonData.Forum,
		f.jsonData.Tag,
		f.jsonData.Qsttype,
		f.jsonData.Source,
		f.jsonData.PetId,
		f.jsonData.PetType,
		f.jsonData.VideoUrl,
		f.jsonData.IsVideo}
	err := c.Insert(&data) //插入数据
	if err != nil {
		return err
	}
	return nil
}

func getTableNum(person int) string {
	tableNumX := person % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	return tableNameX
}

//获取相同俱乐部的活跃用户
func (f *BreedPersons) getPersons(startId bson.ObjectId) *[]ActiveUser.ActiveBreedUser {
	var result []ActiveUser.ActiveBreedUser
	Bid := f.jsonData.Bid
	if Bid == 0 {
		return nil
	}
	c := f.mongoConn[0].DB("ActiveUser").C("active_breed_user")
	err := c.Find(&bson.M{"breed_id": Bid, "_id": bson.M{"$gt": startId}}).
		Limit(count).
		All(&result)
	if err != nil {
		return nil
	}
	return &result
}

func (f *BreedPersons) updatePerson(c *mgo.Collection, person int) error {
	//修改数据
	_, err := c.UpdateAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "infoid": f.jsonData.Infoid}, bson.M{"$set": bson.M{"status": f.jsonData.Status, "created": f.jsonData.Created}})
	if err != nil {
		return err
	}
	return nil
}

func (f *BreedPersons) removePerson(c *mgo.Collection, person int) error {
	//删除数据
	_, err := c.RemoveAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "infoid": f.jsonData.Infoid, "tid": f.jsonData.Tid})
	if err != nil {
		return err
	}
	return nil
}
