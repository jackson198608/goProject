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
	"github.com/olivere/elastic"
)

type AllPersons struct {
	mysqlXorm      []*xorm.Engine
	mongoConn      []*mgo.Session
	jsonData       *job.FocusJsonColumn
	esConn  *elastic.Client
}

const count = 100

func NewAllPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, esConn *elastic.Client) *AllPersons {
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
	f.esConn = esConn
	return f
}

func (f *AllPersons) Do() error {
	//get all active user from hashmap
	er := elasticsearch.NewUser(f.esConn)
	from := 0
	i :=1
	for {
		var uids []int
		rst := er.SearchAllActiveUser(from, count)
		total := rst.Hits.TotalHits
		if total> 0 {
			for _, hit := range rst.Hits.Hits {
				uid,_ := strconv.Atoi(hit.Id)
				uids = append(uids, uid)
			}
		}
		f.pushPersons(uids)
		i++
		from = (i-1)*count
		if int(total) < from {
			break
		}
	}
	return nil
}

func (f *AllPersons) pushPersons(persons []int) error {
	if persons == nil {
		return errors.New("push to all active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	elx := elasticsearch.NewEventLogX(f.esConn, f.jsonData)
	for _,uid := range persons {
		err := elx.PushPerson(uid)
		if err != nil {
			for i := 0; i < 5; i++ {
				err := elx.PushPerson(uid)
				if err == nil {
					break
				}
			}
		}
	}
	return nil
}
