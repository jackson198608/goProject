package CajFansPersons

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
	"strconv"
	"gouminGitlab/common/orm/elasticsearch"
	"github.com/olivere/elastic"
	"gouminGitlab/common/orm/mysql/adoption"
	"fmt"
)

type CajFansPersons struct {
	mysqlXorm []*xorm.Engine //@todo to be []
	jsonData  *job.FocusJsonColumn
	esConn    *elastic.Client
}

const count = 100

func NewCajFansPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, esConn *elastic.Client) *CajFansPersons {
	if (mysqlXorm == nil) || (jsonData == nil) || (esConn == nil) {
		return nil
	}

	f := new(CajFansPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.jsonData = jsonData
	f.esConn = esConn

	return f
}

func (f *CajFansPersons) Do() error {
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

func (f *CajFansPersons) pushPersons(follows *[]adoption.UserFollow) (int, error) {
	if follows == nil {
		return 0, errors.New("push to caj fans active user : you have no person to push " + strconv.Itoa(f.jsonData.AdoptId))
	}
	persons := *follows

	var endId int
	elx, err := elasticsearch.NewEventLogX(f.esConn, f.jsonData)
	if err != nil {
		return 0, err
	}

	for _, person := range persons {
		err := elx.PushPerson(person.Uid)
		if err != nil {
			for i := 0; i < 5; i++ {
				err := elx.PushPerson(person.Uid)
				if err == nil {
					break
				}else{
					fmt.Println("caj push person fail, person is ", person.Uid, "infoid is ",f.jsonData.Infoid," try times is ",i)
				}
			}
		}
		endId = person.Id
	}

	return endId, nil
}

//get fans persons by uid
func (f *CajFansPersons) getPersons(startId int) *[]adoption.UserFollow {
	// var persons []int
	var follows []adoption.UserFollow
	err := f.mysqlXorm[2].Where("fuid=? and id>?", f.jsonData.Uid, startId).Asc("id").Limit(count).Find(&follows)
	if err != nil {
		return nil
	}

	return &follows
}
