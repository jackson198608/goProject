package clubPersons

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/ActiveUser"
	"strconv"
	"gouminGitlab/common/orm/elasticsearch"
)

type ClubPersons struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jsonData  *job.FocusJsonColumn
	nodes []string
}

//已废弃

const count = 1000

func NewClubPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, nodes []string) *ClubPersons {
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
	f.nodes = nodes

	return f
}

func (f *ClubPersons) Do() error {
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

func (f *ClubPersons) pushPersons(ActiveUser *[]ActiveUser.ActiveForumUser) (bson.ObjectId, error) {
	if ActiveUser == nil {
		return bson.NewObjectId(), errors.New("push to club active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}

	var endId bson.ObjectId
	persons := *ActiveUser
	elx := elasticsearch.NewEventLogX(f.nodes, f.jsonData)
	for _, person := range persons {
		// fmt.Println(person.Uid)
		err := elx.PushPerson(person.Uid)
		if err != nil {
			for i := 0; i < 5; i++ {
				err := elx.PushPerson(person.Uid)
				if err == nil {
					break
				}
			}
		}
		endId = person.Id
	}
	return endId, nil
}

//获取相同俱乐部的活跃用户
//@todo 活跃用户需要改为从elk获取  未测试
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



