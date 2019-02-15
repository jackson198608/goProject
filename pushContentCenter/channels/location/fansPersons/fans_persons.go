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
	"fmt"
)

type FansPersons struct {
	mysqlXorm      []*xorm.Engine //@todo to be []
	mongoConn      []*mgo.Session //@todo to be []
	jsonData       *job.FocusJsonColumn
	esConn  *elastic.Client
}

const count = 100

func NewFansPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, esConn *elastic.Client) *FansPersons {
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
	f.esConn = esConn

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
	persons := *follows

	var endId int
	elx := elasticsearch.NewEventLogX(f.esConn, f.jsonData)
	var active_user map[int]bool

	if f.jsonData.Action != -1 {
		active_user = f.getActiveUserByUids(follows)
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

/**
获取活跃用户的粉丝
 */
func (f *FansPersons) getActiveUserByUids(follows *[]new_dog123.Follow) map[int]bool {
	var m map[int]bool
	m = make(map[int]bool)
	er := elasticsearch.NewUser(f.esConn)
	var uids []int
	persons := *follows
	for _, person := range persons {
		uids = append(uids, person.FollowId)
	}
	rst := er.SearchActiveUserByUids(uids, 0, count)
	total := rst.Hits.TotalHits
	if total> 0 {
		for _, hit := range rst.Hits.Hits {
			uid,_ := strconv.Atoi(hit.Id)
			m[uid] = true
		}
	}
	fmt.Println(m)
	return m
}

