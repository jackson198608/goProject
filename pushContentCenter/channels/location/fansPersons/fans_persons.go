package fansPersons

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/FansData"
	"gouminGitlab/common/orm/mysql/new_dog123"
	"math"
	"strconv"
	// "strings"
)

type FansPersons struct {
	mysqlXorm []*xorm.Engine //@todo to be []
	mongoConn []*mgo.Session //@todo to be []
	jobData   *FansData.EventLog
	status    int
	uid       int
}

const count = 1000

func NewFansPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jobData *FansData.EventLog, status int, uid int) *FansPersons {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jobData == nil) {
		return nil
	}

	f := new(FansPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm[0]
	f.mongoConn = mongoConn[0]
	f.jobData = jobData
	f.status = status
	f.uid = uid

	return f
}

func (f *FansPersons) Do() error {
	initialId := f.getPersonFirstId()
	page := f.getPersonPageNum()

	for i := 1; i <= page; i++ {
		initialId = initialId + (i-1)*count
		startId, endId := f.getIdRange(initialId)
		currentPersionList := f.getPersons(startId, endId)
		err := f.pushPersons(currentPersionList)
		if err != nil {
			return err
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
	err := f.mysqlXorm.Where("user_id=? and fans_active=1", f.uid).Asc("id").Limit(1).Find(&follows)
	if err != nil {
		fmt.Println(err)
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
	err := f.mysqlXorm.Where("user_id=? and fans_active=1", f.uid).Desc("id").Limit(1).Find(&follows)
	if err != nil {
		fmt.Println(err)
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
	page := int(math.Ceil(float64(startId-endId) / float64(count)))
	return page
}

func (f *FansPersons) pushPersons(persons []int) error {
	if persons == nil {
		return errors.New("you have no person to push " + f.jobData.Uid)
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
	c := f.mongoConn.DB("FansData").C(tableNameX)

	if f.status == 1 {
		//新增数据
		// m := f.jobData
		err := c.Insert(&f.jobData) //插入数据
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FansPersons) tryPushPerson(person int, num int) error {
	if num > 5 {
		return errors.New("Attempting to push has failed 5 times: " + f.jobstr + "; person is " + strconv.Itoa(person))
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
	err := f.mysqlXorm.Where("user_id=? and id>=? and id<=? and fans_active=1", f.uid, startId, endId).Cols("follow_id").Find(&follows)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, v := range follows {
		persons = append(persons, v.FollowId)
	}
	uids := f.getFansActivePersons(persons)
	return uids
}

//@todo make active_user to be hash struct
func (f *Focus) getFansActivePersons(persons []int) []int {
	var uids []int
	c := f.mongoConn.DB("ActiveUser").C("active_user")
	err := c.Find(&bson.M{"uid": bson.M{"$in": persons}}).Distinct("uid", &uids)
	if err != nil {
		panic(err)
	}
	return uids
}
