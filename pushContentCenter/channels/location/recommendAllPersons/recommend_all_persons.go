package recommendAllPersons

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	mgo "gopkg.in/mgo.v2"
	"strconv"
	"gouminGitlab/common/orm/elasticsearch"
	"github.com/olivere/elastic"
)

type RecommendAllPersons struct {
	mysqlXorm      []*xorm.Engine
	mongoConn      []*mgo.Session
	jsonData       *job.RecommendJsonColumn
	esConn *elastic.Client
}

const count = 1000

func NewRecommendAllPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, esConn *elastic.Client,jsonData *job.RecommendJsonColumn) *RecommendAllPersons {
	if (mongoConn == nil) || (jsonData == nil) {
		return nil
	}

	f := new(RecommendAllPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jsonData = jsonData
	f.esConn = esConn

	return f
}

func (f *RecommendAllPersons) Do() error {
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
		if len(uids)>0 {
			if f.jsonData.Action == 0 {
				//get all active user from hashmap
				err := f.pushPersons(uids)
				if err != nil {
					return err
				}
			} else if f.jsonData.Action == -1 {
				err := f.removeInfoByTables()
				if err != nil {
					return err
				}
			}
		}
		i++
		from = (i-1)*count
		if int(total) < from {
			break
		}
	}
	return nil
}

func (f *RecommendAllPersons) pushPersons(persons []int) error {
	if persons == nil {
		return errors.New("push to all active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	for k := range persons {
		err := f.pushPerson(k)
		if err != nil {
			for i := 0; i < 5; i++ {
				err = f.pushPerson(k)
				if err == nil {
					break
				}
			}
		}
	}
	return nil
}

func (f *RecommendAllPersons) pushPerson(person int) error {
	ur := elasticsearch.NewUserRecommendX(f.esConn, f.jsonData)
	err := ur.Create(person)
	if err != nil {
		return err
	}
	return nil
}

func (f *RecommendAllPersons) removeInfoByTables() error {
	ur := elasticsearch.NewUserRecommendX(f.esConn, f.jsonData)
	err := ur.Remove()
	if err != nil {
		return err
	}
	return nil
}

