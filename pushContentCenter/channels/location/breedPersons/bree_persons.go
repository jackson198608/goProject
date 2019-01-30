package breedPersons

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
	"strconv"
	"gouminGitlab/common/orm/elasticsearch"
)

const count = 1000

type BreedPersons struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jsonData  *job.FocusJsonColumn
	nodes []string
}

func NewBreedPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn,nodes []string) *BreedPersons {
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
	f.nodes = nodes

	return f
}

func (f *BreedPersons) Do() error {
	currentPersionList := f.getPersonsByElk()
	if currentPersionList == nil {
		return nil
	}
	err := f.pushPersons(currentPersionList)
	if err != nil {
		return err
	}
	return nil
}

func (f *BreedPersons) pushPersons(ActiveUser []int) error {
	if ActiveUser == nil {
		return errors.New("push to breed active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	elx := elasticsearch.NewEventLogX(f.nodes, f.jsonData)
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

func (f *BreedPersons) getPersonsByElk() []int {
	Bid := f.jsonData.Bid
	if Bid == 0 {
		return nil
	}
	var activeBreedUsers []int
	from := 0
	size := 1000
	i :=1
	for {
		er := elasticsearch.NewUser(f.nodes)
		rst := er.SearchAllActiveUserByBreed(Bid, f.jsonData.Uid, from, size)
		total := rst.Hits.TotalHits
		if total > 0 {
			for _, hit := range rst.Hits.Hits {
				id, _ := strconv.Atoi(hit.Id)
				activeBreedUsers = append(activeBreedUsers, id)
			}
		}
		if int(total) < from {
			break
		}
		i++
		from = (i-1)*size
	}
	return activeBreedUsers
}
