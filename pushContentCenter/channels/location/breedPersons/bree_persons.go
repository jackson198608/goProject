package breedPersons

import (
	"error"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/FansData"
	"gouminGitlab/common/orm/mysql/new_dog123"
	"math"
	"strconv"
	"strings"
)

type BreedPersons struct {
	mysqlXorm []*xorm.Engine //@todo to be []
	mongoConn []*mgo.Session //@todo to be []
	jobData   *FansData.EventLog
	status    int
	bid       int
}

const count = 1000

func NewBreedPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jobData *FansData.EventLog, status int, bid int) *BreedPersons {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jobData == nil) {
		return nil
	}

	f := new(BreedPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jobData = jobData
	f.status = status
	f.bid = bid

	return f
}

func (f *BreedPersons) Do() error {
	page := f.getPersonPageNum()

	for i := 1; i <= page; i++ {
		currentPersionList := f.getPersons(i)
		err := f.pushPersons(currentPersionList)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *BreedPersons) pushPersons(persons []int) error {
	if persons == nil {
		return errors.New("you have no person to push " + f.jobstr)
	}

	for _, person := range persons {
		err := f.pushPerson(person)
		if err != nil {
			f.tryPushPerson(person, 1)
		}
	}
	return nil
}

func (f *BreedPersons) tryPushPerson(person int, num int) error {
	if num > 5 {
		return errors.New("Attempting to push has failed 5 times: " + f.jobstr + "; person is " + strconv.Itoa(person))
	}
	err := f.pushPerson(person)
	if err != nil {
		f.tryPushPerson(person, num+1)
	}
	return nil
}

func (f *BreedPersons) pushPerson(person int) error {
	tableNameX := getTableNum(person)
	c := f.mongoConn.DB("FansData").C(tableNameX)
	if f.status == 1 {
		//数据展示
		err := c.Insert(&f.jobData) //插入数据
		if err != nil {
			return err
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

// Get the same club user data page number
func (f *BreedPersons) getPersonPageNum() int {
	Bid := f.jsonData.Bid
	if Bid == 0 {
		return 0
	}

	c := f.mongoConn.DB("ActiveUser").C("active_breed_user")
	countNum, err := c.Find(&bson.M{"breed_id": Bid}).Count()
	if err != nil {
		panic(err)
	}
	page := int(math.Ceil(float64(countNum) / float64(count)))
	return page
}

//获取相同俱乐部的活跃用户
//@todo 使用id范围分页查询
func (f *BreedPersons) getPersons(page int) []int {
	var uids []int
	Bid := f.jsonData.Bid
	if Bid == 0 {
		return uids
	}

	c := f.mongoConn.DB("ActiveUser").C("active_breed_user")
	err := c.Find(&bson.M{"breed_id": Bid}).
		Skip((page-1)*count).
		Limit(count).
		Distinct("uid", &uids)
	if err != nil {
		panic(err)
	}
	return uids
}
