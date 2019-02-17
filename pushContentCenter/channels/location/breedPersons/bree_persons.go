package breedPersons

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
	"strconv"
	"gouminGitlab/common/orm/elasticsearch"
	"github.com/olivere/elastic"
)

const count = 100

type BreedPersons struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jsonData  *job.FocusJsonColumn
	esConn  *elastic.Client
}

func NewBreedPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn,esConn *elastic.Client) *BreedPersons {
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
	f.esConn = esConn

	return f
}

func (f *BreedPersons) Do() error {
	from := 0
	i :=1
	for {
		currentPersionList := f.getPersonsByElk(from)
		if currentPersionList == nil {
			return nil
		}
		total := len(currentPersionList)
		err := f.pushPersons(currentPersionList)
		if err != nil {
			return err
		}
		i++
		from = (i-1)*count
		if int(total) < count {
			break
		}
	}
	return nil
}

func (f *BreedPersons) pushPersons(ActiveUser []int) error {
	if ActiveUser == nil {
		return errors.New("push to breed active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	elx := elasticsearch.NewEventLogX(f.esConn, f.jsonData)
	for _, person := range ActiveUser {
		err := elx.PushPerson(person)
		if err != nil {
			for i := 0; i < 5; i++ {
				err := elx.PushPerson(person)
				if err == nil {
					break
				}
			}
		}
	}
	return nil
}

func (f *BreedPersons) getPersonsByElk(from int) []int {
	Bid := f.jsonData.Bid
	if Bid == 0 {
		return nil
	}
	var activeBreedUsers []int
	er := elasticsearch.NewUser(f.esConn)
	rst := er.SearchAllActiveUserByBreed(Bid, f.jsonData.Uid, from, count)
	total := rst.Hits.TotalHits
	if total > 0 {
		for _, hit := range rst.Hits.Hits {
			id, _ := strconv.Atoi(hit.Id)
			activeBreedUsers = append(activeBreedUsers, id)
		}
	}
	return activeBreedUsers
}
