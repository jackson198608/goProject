package allPersons

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
	//"gouminGitlab/common/orm/mongo/FansData"
	"strconv"
	"gouminGitlab/common/orm/elasticsearch"
)

type AllPersons struct {
	mysqlXorm      []*xorm.Engine
	mongoConn      []*mgo.Session
	jsonData       *job.FocusJsonColumn
	activeUserData *map[int]bool
	nodes []string
}

const count = 1000

func NewAllPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, activeUserData *map[int]bool, nodes []string) *AllPersons {
	if (mongoConn == nil) || (jsonData == nil) {
		return nil
	}

	f := new(AllPersons)
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

func (f *AllPersons) Do() error {
	//get all active user from hashmap
	f.pushPersons(f.activeUserData)
	return nil
}

func (f *AllPersons) pushPersons(persons *map[int]bool) error {
	if persons == nil {
		return errors.New("push to all active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	elx := elasticsearch.NewEventLogX(f.nodes, f.jsonData)
	for k := range *persons {
		err := elx.PushPerson(k)
		if err != nil {
			for i := 0; i < 5; i++ {
				err := elx.PushPerson(k)
				if err == nil {
					break
				}
			}
		}
	}
	return nil
}
