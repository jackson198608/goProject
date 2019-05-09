package robots

import (
	"errors"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"github.com/olivere/elastic"
	"strconv"
	"gouminGitlab/common/orm/elasticsearch"
	"gouminGitlab/common/orm/mysql/member"
	"github.com/donnie4w/go-logger/logger"
)

type Robots struct {
	mysqlXorm      []*xorm.Engine //@todo to be []
	mongoConn      []*mgo.Session //@todo to be []
	jsonData       *job.FocusJsonColumn
	esConn  *elastic.Client
}

const count = 100
func NewRobots(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, esConn *elastic.Client) *Robots {
	if (mysqlXorm == nil) || (jsonData == nil) || (esConn ==nil){
		return nil
	}

	r := new(Robots)
	if r == nil {
		return nil
	}

	r.mysqlXorm = mysqlXorm
	r.mongoConn = mongoConn
	r.jsonData = jsonData
	r.esConn = esConn

	return r
}

func (r *Robots) Do() error {
	startId := 0
	for {
		//获取机器人数据
		currentPersionList,err := r.getPersons(startId)
		if err != nil {
			return err
		}
		endId, err := r.pushPersons(currentPersionList)
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

func (r *Robots) pushPersons(robots *[]member.PublishUser) (int, error) {
	if robots == nil {
		return 0, errors.New("push to fans active user : you have no person to push " + strconv.Itoa(r.jsonData.Infoid))
	}
	persons := *robots

	r.jsonData.Channel = 1
	var endId int
	elx,err := elasticsearch.NewEventLogX(r.esConn, r.jsonData)
	if err !=nil {
		return 0, err
	}

	for _, person := range persons {
		err := elx.PushPerson(person.RealUid)
		if err != nil {
			for i := 0; i < 5; i++ {
				logger.Info("push fans ", person.RealUid, " try ", i, " by ",r.jsonData)
				err := elx.PushPerson(person.RealUid)
				if err == nil {
					break
				}
			}
		}
		endId = person.Id
	}

	return endId, nil
}

/**
	获取机器人uid
 */
//get fans persons by uid
func (r *Robots) getPersons(startId int) (*[]member.PublishUser, error){
	// var persons []int
	var uids []member.PublishUser
	err := r.mysqlXorm[3].Where("robot_uid =? and robot_nums>? and id>?", 0, 0, startId).Asc("id").Limit(count).Find(&uids)
	if err != nil {
		return &uids,err
	}

	return &uids,nil
}

