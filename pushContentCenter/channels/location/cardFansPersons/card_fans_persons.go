package cardFansPersons

import (
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"strconv"
	"github.com/pkg/errors"
	"gouminGitlab/common/orm/mysql/card"
	"fmt"
	"gouminGitlab/common/orm/elasticsearch"
	"github.com/olivere/elastic"
)

type CardFansPersons struct {
	mysqlXorm      []*xorm.Engine //@todo to be []
	mongoConn      []*mgo.Session //@todo to be []
	jsonData       *job.FocusJsonColumn
	activeUserData *map[int]bool
	esConn *elastic.Client
}

const count = 1000

func NewCardFansPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, activeUserData *map[int]bool,esConn *elastic.Client) *CardFansPersons {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jsonData == nil) {
		return nil
	}

	f := new(CardFansPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jsonData = jsonData
	f.activeUserData = activeUserData
	f.esConn = esConn

	return f
}

func (f *CardFansPersons) Do() error {
	startId := 0

	f.pushMyself()
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

func (f *CardFansPersons) pushMyself() {
	//推送给自己
	elx := elasticsearch.NewEventLogX(f.esConn, f.jsonData)
	err := elx.PushPerson(f.jsonData.Uid)
	if err != nil {
		for i := 0; i < 5; i++ {
			err := elx.PushPerson(f.jsonData.Uid)
			if err == nil {
				break
			}
		}
	}
}

func (f *CardFansPersons) pushPersons(follows *[]card.HaremCard) (int, error) {
	if follows == nil {
		return 0, errors.New("push to fans active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	active_user := *f.activeUserData
	persons := *follows

	var endId int
	elx := elasticsearch.NewEventLogX(f.esConn, f.jsonData)
	for _, person := range persons {
		//check key in actice user
		var ok bool
		if f.jsonData.Action == -1 {
			ok = true
		} else {
			_, ok = active_user[person.Uid]
		}
		if ok {
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
	}
	return endId, nil
}

//get fans persons by pet_id
func (f *CardFansPersons) getPersons(startId int) *[]card.HaremCard {
	// var persons []int
	var follows []card.HaremCard
	err := f.mysqlXorm[1].Where("pet_id=? and type=? and id>?", f.jsonData.PetId, f.jsonData.PetType, startId).Asc("id").Limit(count).Find(&follows)
	if err != nil {
		fmt.Println(err)
		for i := 0; i < 5; i++ {
			err := f.mysqlXorm[1].Where("pet_id=? and type=? and id>?", f.jsonData.PetId, f.jsonData.PetType, startId).Asc("id").Limit(count).Find(&follows)
			if err == nil {
				break
			}
		}
	}
	return &follows
}

