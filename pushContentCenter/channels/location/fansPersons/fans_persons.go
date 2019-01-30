package fansPersons

import (
	"errors"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
	"gouminGitlab/common/orm/mysql/new_dog123"
	// "reflect"
	"strconv"
	"gouminGitlab/common/orm/elasticsearch"
)

type FansPersons struct {
	mysqlXorm      []*xorm.Engine //@todo to be []
	mongoConn      []*mgo.Session //@todo to be []
	jsonData       *job.FocusJsonColumn
	activeUserData *map[int]bool
	nodes []string
}

const count = 1000

func NewFansPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, activeUserData *map[int]bool,nodes []string) *FansPersons {
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
	f.nodes = nodes

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
	//fmt.Println(active_user)

	var endId int
	elx := elasticsearch.NewEventLogX(f.nodes, f.jsonData)
	for _, person := range persons {
		//check key in actice user
		_, ok := active_user[person.FollowId]
		if ok {
			err := elx.PushPerson(person.FollowId)
			if err != nil {
				for i := 0; i < 5; i++ {
					err := elx.PushPerson(person.FollowId)
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

