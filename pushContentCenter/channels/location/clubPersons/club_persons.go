package clubPersons

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/ActiveUser"
	"gouminGitlab/common/orm/mongo/FansData"
	"strconv"
)

type ClubPersons struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jsonData  *job.FocusJsonColumn
	fid       int
}

const count = 1000

func NewClubPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn) *ClubPersons {
	if (mongoConn == nil) || (jsonData == nil) {
		return nil
	}

	f := new(ClubPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jsonData = jsonData
	f.fid = f.jsonData.Fid

	return f
}

func (f *ClubPersons) Do() error {
	var startId bson.ObjectId
	startId = bson.ObjectId("000000000000")
	for {
		currentPersionList := f.getPersons(startId)
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

func (f *ClubPersons) pushPersons(ActiveUser *[]ActiveUser.ActiveForumUser) (bson.ObjectId, error) {
	if ActiveUser == nil {
		return bson.NewObjectId(), errors.New("push to club active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}

	var endId bson.ObjectId
	persons := *ActiveUser
	for _, person := range persons {
		// fmt.Println(person.Uid)
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

func (f *ClubPersons) pushPerson(person int) error {
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

func getTableNum(person int) string {
	tableNumX := person % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	return tableNameX
}

//获取相同俱乐部的活跃用户
func (f *ClubPersons) getPersons(startId bson.ObjectId) *[]ActiveUser.ActiveForumUser {
	var result []ActiveUser.ActiveForumUser

	c := f.mongoConn[0].DB("ActiveUser").C("active_forum_user")
	err := c.Find(&bson.M{"forum_id": f.jsonData.Fid, "_id": bson.M{"$gt": startId}}).
		Limit(count).
		All(&result)
	if err != nil {
		// panic(err)
		return &result
	}
	return &result
}

func (f *ClubPersons) insertPerson(c *mgo.Collection, person int) error {
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

func (f *ClubPersons) updatePerson(c *mgo.Collection, person int) error {
	//修改数据
	_, err := c.UpdateAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid}, bson.M{"$set": bson.M{"status": f.jsonData.Status}})
	if err != nil {
		return err
	}
	return nil
}

func (f *ClubPersons) removePerson(c *mgo.Collection, person int) error {
	//删除数据
	_, err := c.RemoveAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid, "tid": f.jsonData.Tid})
	if err != nil {
		return err
	}
	return nil
}
