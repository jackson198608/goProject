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
	"github.com/olivere/elastic"
	log "github.com/thinkboy/log4go"
)

type FansPersons struct {
	mysqlXorm      []*xorm.Engine //@todo to be []
	mongoConn      []*mgo.Session //@todo to be []
	jsonData       *job.FocusJsonColumn
	esConn  *elastic.Client
}

const count = 100

func NewFansPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, esConn *elastic.Client) *FansPersons {
	if (mysqlXorm == nil) || (jsonData == nil) || (esConn ==nil){
		return nil
	}

	f := new(FansPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jsonData = jsonData
	f.esConn = esConn

	return f
}

func (f *FansPersons) Do() error {
	startId := 0
	for {
		//获取粉丝用户
		currentPersionList,err := f.getPersons(startId)
		if err != nil {
			return err
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
	persons := *follows

	var endId int
	elx,err := elasticsearch.NewEventLogX(f.esConn, f.jsonData)
	if err !=nil {
		return 0, err
	}
	var active_user map[int]bool

	if f.jsonData.Action != -1 {
		active_user,err = f.getActiveUserByUids(follows)
		if err !=nil {
			return 0, err
		}
	}
	for _, person := range persons {
		ok := false
		//如果是删除操作，则所有粉丝用户全部更新
		if f.jsonData.Action == -1  {
			ok = true
		}else {
			//check key in actice user
			_, ok = active_user[person.FollowId]
		}
		if ok {
			err := elx.PushPerson(person.FollowId)
			if err != nil {
				for i := 0; i < 5; i++ {
					log.Info("push fans ", person.FollowId, " try ", i, " by ",f.jsonData)
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
func (f *FansPersons) getPersons(startId int) (*[]new_dog123.Follow, error) {
	// var persons []int
	var follows []new_dog123.Follow
	err := f.mysqlXorm[0].Where("user_id=? and id>? and fans_active=1", f.jsonData.Uid, startId).Asc("id").Limit(count).Find(&follows)
	if err != nil {
		return nil,err
	}

	return &follows,nil
}

/**
获取活跃用户的粉丝
 */
func (f *FansPersons) getActiveUserByUids(follows *[]new_dog123.Follow) (map[int]bool,error) {

	er,err := elasticsearch.NewUserInfo(f.esConn)
	if err!=nil {
		return nil,err
	}
	var uids []int
	persons := *follows
	for _, person := range persons {
		uids = append(uids, person.FollowId)
	}
	rst,err := er.GetActiveUserInfoByUids(uids, 0, count)
	if err!=nil {
		return nil,err
	}
	return rst,nil
}



