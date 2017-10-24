package fansPersons

import (
	"errors"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/FansData"
	"gouminGitlab/common/orm/mysql/new_dog123"
	"math"
	// "reflect"
	"strconv"
	// "strings"
)

type FansPersons struct {
	mysqlXorm      []*xorm.Engine //@todo to be []
	mongoConn      []*mgo.Session //@todo to be []
	jsonData       *job.FocusJsonColumn
	activeUserData *map[int]bool
	uid            int
}

const count = 1000

func NewFansPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, activeUserData *map[int]bool) *FansPersons {
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
	f.uid = f.jsonData.Uid
	f.activeUserData = activeUserData

	return f
}

func (f *FansPersons) Do() error {
	initialId := f.getPersonFirstId()
	page := f.getPersonPageNum()
	for i := 1; i <= page; i++ {
		if i != 1 {
			initialId += count
		}
		startId, endId := f.getIdRange(initialId)
		// fmt.Println(startId)
		// fmt.Println(endId)
		currentPersionList := f.getPersons(startId, endId)
		// fmt.Println(currentPersionList)
		if len(currentPersionList) > 0 {
			err := f.pushPersons(currentPersionList)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//return id range
func (f *FansPersons) getIdRange(startId int) (int, int) {
	endId := startId + count
	maxId := f.getPersonLastId()
	if endId > maxId {
		endId = maxId
	}
	return startId, endId
}

//获取活跃用户第一个ID
func (f *FansPersons) getPersonFirstId() int {
	var follows []new_dog123.Follow
	err := f.mysqlXorm[0].Where("user_id=? and fans_active=1", f.uid).Asc("id").Limit(1).Find(&follows)
	if err != nil {
		return 0
	}

	if len(follows) == 0 {
		return 0
	}

	id := follows[0].Id - 1
	return id
}

//获取活跃用户最后一个ID
func (f *FansPersons) getPersonLastId() int {
	var follows []new_dog123.Follow
	err := f.mysqlXorm[0].Where("user_id=? and fans_active=1", f.uid).Desc("id").Limit(1).Find(&follows)
	if err != nil {
		return 0
	}
	if len(follows) == 0 {
		return 0
	}

	id := follows[0].Id

	return id
}

func (f *FansPersons) getPersonPageNum() int {
	startId := f.getPersonFirstId()
	endId := f.getPersonLastId()
	page := int(math.Ceil(float64(endId-startId) / float64(count)))
	return page
}

func (f *FansPersons) pushPersons(persons []int) error {
	if persons == nil {
		return errors.New("push to fans active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}

	for _, person := range persons {
		err := f.pushPerson(person)
		if err != nil {
			f.tryPushPerson(person, 1)
		}
	}
	return nil
}

func getTableNum(person int) string {
	tableNumX := person % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	return tableNameX
}

func (f *FansPersons) pushPerson(person int) error {
	tableNameX := getTableNum(person)
	c := f.mongoConn[0].DB("FansData").C(tableNameX)
	if f.jsonData.Action == 0 {
		// fmt.Println("insert" + strconv.Itoa(person))

		err := f.insertPerson(c, person)
		if err != nil {
			return err
		}
	} else if f.jsonData.Action == 1 {
		// } else if (f.jsonData.Action == 1) && (f.checkDataIsExist(person)) {
		//修改数据
		// fmt.Println("update" + strconv.Itoa(person))
		err := f.updatePerson(c, person)
		if err != nil {
			return err
		}
		// } else if (f.jsonData.Action == -1) && (f.checkDataIsExist(person)) {
	} else if f.jsonData.Action == -1 {
		//删除数据
		// fmt.Println("remove" + strconv.Itoa(person))

		err := f.removePerson(c, person)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FansPersons) tryPushPerson(person int, num int) error {
	if num > 5 {
		return errors.New("push to fans active user : Attempting to push has failed 5 times; infoid is " + strconv.Itoa(f.jsonData.Infoid) + "; person is " + strconv.Itoa(person))
	}
	err := f.pushPerson(person)
	if err != nil {
		f.tryPushPerson(person, num+1)
	}
	return nil
}

//获取活跃粉丝用户
func (f *FansPersons) getPersons(startId int, endId int) []int {
	var persons []int
	var follows []new_dog123.Follow
	err := f.mysqlXorm[0].Where("user_id=? and id>? and id<=? and fans_active=1", f.uid, startId, endId).Cols("follow_id").Find(&follows)
	if err != nil {
		return nil
	}
	for _, v := range follows {
		persons = append(persons, v.FollowId)
	}
	uids := f.getFansActivePersons(persons)
	return uids
}

//get active user by hashmap
func (f *FansPersons) getFansActivePersons(persons []int) []int {
	var uids []int
	active_user := *f.activeUserData
	for i := 0; i < len(persons); i++ {
		//check key is exists
		_, ok := active_user[persons[i]]
		if ok {
			uids = append(uids, persons[i])
		}
	}
	return uids
}

//检查mongo中是否存在该条数据
func (f *FansPersons) checkDataIsExist(person int) bool {
	var ms []FansData.EventLog
	tableNameX := getTableNum(person)
	c := f.mongoConn[0].DB("FansData").C(tableNameX)
	err1 := c.Find(&bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid, "tid": f.jsonData.Tid}).All(&ms)

	if err1 != nil {
		return false
	}
	if len(ms) == 0 {
		return false
	}
	return true
}

func (f *FansPersons) insertPerson(c *mgo.Collection, person int) error {
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
		f.jsonData.Forum,
		f.jsonData.Tag,
		f.jsonData.Qsttype,
		f.jsonData.Source}
	err := c.Insert(&data) //插入数据
	if err != nil {
		return err
	}
	return nil
}

func (f *FansPersons) updatePerson(c *mgo.Collection, person int) error {
	_, err := c.UpdateAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid}, bson.M{"$set": bson.M{"status": f.jsonData.Status}})
	if err != nil {
		return err
	}
	return nil
}

func (f *FansPersons) removePerson(c *mgo.Collection, person int) error {
	_, err := c.RemoveAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid, "tid": f.jsonData.Tid})
	if err != nil {
		return err
	}
	return nil
}
