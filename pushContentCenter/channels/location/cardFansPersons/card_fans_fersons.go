package cardFansPersons

import (
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"strconv"
	"gouminGitlab/common/orm/mongo/FansData"
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
	"gouminGitlab/common/orm/mysql/card"
	"fmt"
)

type CardFansPersons struct {
	mysqlXorm      []*xorm.Engine //@todo to be []
	mongoConn      []*mgo.Session //@todo to be []
	jsonData       *job.FocusJsonColumn
	activeUserData *map[int]bool
}

const count = 1000

func NewCardFansPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, activeUserData *map[int]bool) *CardFansPersons {
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

	return f
}

func (f *CardFansPersons) Do() error {
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

func (f *CardFansPersons) pushPersons(follows *[]card.HaremCard) (int, error) {
	if follows == nil {
		return 0, errors.New("push to fans active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}
	active_user := *f.activeUserData
	persons := *follows

	var endId int
	for _, person := range persons {
		//check key in actice user
		_, ok := active_user[person.Uid]
		if ok {
			err := f.pushPerson(person.Uid)
			if err != nil {
				for i := 0; i < 5; i++ {
					err := f.pushPerson(person.Uid)
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

func getTableNum(person int) string {
	tableNumX := person % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	return tableNameX
}

func (f *CardFansPersons) pushPerson(person int) error {
	tableNameX := getTableNum(person)
	c := f.mongoConn[0].DB("FansData").C(tableNameX)
	if f.jsonData.Action == 0 {
		//fmt.Println("insert" + strconv.Itoa(person))
		err := f.insertPerson(c, person)
		if err != nil {
			return err
		}
	} else if f.jsonData.Action == 1 {
		//修改数据
		//fmt.Println("update" + strconv.Itoa(person))
		err := f.updatePerson(c, person)
		if err != nil {
			return err
		}
	} else if f.jsonData.Action == -1 {
		//删除数据
		//fmt.Println("remove" + strconv.Itoa(person))
		err := f.removePerson(c, person)
		if err != nil {
			return err
		}
	}
	return nil
}

//get fans persons by pet_id
func (f *CardFansPersons) getPersons(startId int) *[]card.HaremCard {
	// var persons []int
	var follows []card.HaremCard
	err := f.mysqlXorm[1].Where("pet_id=? and type=? and id>?", f.jsonData.PetId, f.jsonData.PetType, startId).Asc("id").Limit(count).Find(&follows)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &follows
}

func (f *CardFansPersons) insertPerson(c *mgo.Collection, person int) error {
	//新增数据
	var data FansData.EventLog
	data = FansData.EventLog{bson.NewObjectId(),
		f.jsonData.TypeId,
		f.jsonData.Uid,
		person,
		f.jsonData.Created,
		f.jsonData.Infoid,
		f.jsonData.Status,
		f.jsonData.Tid,
		f.jsonData.Bid,
		f.jsonData.Content,
		f.jsonData.Title,
		f.jsonData.Imagenums,
		f.jsonData.ImageInfo,
		f.jsonData.Forum,
		f.jsonData.Tag,
		f.jsonData.Qsttype,
		f.jsonData.Source,
		f.jsonData.PetId,
		f.jsonData.PetType,
		f.jsonData.VideoUrl}
	err := c.Insert(&data) //插入数据
	if err != nil {
		return err
	}
	return nil
}

func (f *CardFansPersons) updatePerson(c *mgo.Collection, person int) error {
	_, err := c.UpdateAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid}, bson.M{"$set": bson.M{"status": f.jsonData.Status}})
	if err != nil {
		return err
	}
	return nil
}

func (f *CardFansPersons) removePerson(c *mgo.Collection, person int) error {
	_, err := c.RemoveAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid, "tid": f.jsonData.Tid})
	if err != nil {
		return err
	}
	return nil
}
