package clubPersons

import (
	"errors"
	// "fmt"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/ActiveUser"
	"gouminGitlab/common/orm/mongo/FansData"
	"math"
	"strconv"
)

type ClubPersons struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jsonData  *job.FocusJsonColumn
	fid       int
}

const count = 2 //1000

func NewClubPersons(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn) *ClubPersons {
	if (mongoConn == nil) || (jsonData == nil) {
		return nil
	}

	f := new(ClubPersons)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jsonData = jsonData
	f.fid = f.jsonData.Fid

	return f
}

func (f *ClubPersons) Do() error {
	page := f.getPersonPageNum()
	for i := 1; i <= page; i++ {
		currentPersionList := f.getPersons(i)
		// fmt.Println(currentPersionList)
		err := f.pushPersons(currentPersionList)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *ClubPersons) pushPersons(persons []int) error {
	if persons == nil {
		return errors.New("push to club active user : you have no person to push " + strconv.Itoa(f.jsonData.Infoid))
	}

	for _, person := range persons {
		err := f.pushPerson(person)
		if err != nil {
			f.tryPushPerson(person, 1)
		}
	}
	return nil
}

func (f *ClubPersons) tryPushPerson(person int, num int) error {
	if num > 5 {
		return errors.New("push to club active user : Attempting to push has failed 5 times; infoid is " + strconv.Itoa(f.jsonData.Infoid) + "; person is " + strconv.Itoa(person))
	}
	err := f.pushPerson(person)
	if err != nil {
		f.tryPushPerson(person, num+1)
	}
	return nil
}

func (f *ClubPersons) pushPerson(person int) error {
	tableNameX := getTableNum(person)
	c := f.mongoConn[0].DB("FansData").C(tableNameX)
	if f.jsonData.Action == 0 {
		// fmt.Println("insert" + strconv.Itoa(person))
		err := f.insertPerson(c, person)
		if err != nil {
			return err
		}
		// } else if (f.jsonData.Action == 1) && (f.checkDataIsExist(person)) {
	} else if f.jsonData.Action == 1 {
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

func getTableNum(person int) string {
	tableNumX := person % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	return tableNameX
}

// Get the same club user data page number
func (f *ClubPersons) getPersonPageNum() int {
	c := f.mongoConn[0].DB("ActiveUser").C("active_forum_user")
	countNum, err := c.Find(&bson.M{"forum_id": f.fid}).Count()
	if err != nil {
		panic(err)
	}
	page := int(math.Ceil(float64(countNum) / float64(count)))

	return page
}

//获取相同俱乐部的活跃用户
//@todo 使用id范围分页查询
func (f *ClubPersons) getPersons(page int) []int {
	var uids []int
	var result []ActiveUser.ActiveForumUser

	c := f.mongoConn[0].DB("ActiveUser").C("active_forum_user")
	err := c.Find(&bson.M{"forum_id": f.fid}).
		Select(bson.M{"uid": 1}).
		Skip((page - 1) * count).
		Limit(count).
		All(&result)
	if err != nil {
		panic(err)
		return uids
	}
	for _, v := range result {
		uids = append(uids, v.Uid)
	}
	return uids
}

//检查mongo中是否存在该条数据
func (f *ClubPersons) checkDataIsExist(person int) bool {
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

func (f *ClubPersons) insertPerson(c *mgo.Collection, person int) error {
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

func (f *ClubPersons) updatePerson(c *mgo.Collection, person int) error {
	//修改数据
	_, err := c.UpdateAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid}, bson.M{"$set": bson.M{"status": f.jsonData.Status}})
	if err != nil {
		return err
	}
	return nil
}

func (f *ClubPersons) removePerson(c *mgo.Collection, person int) error {
	//删除数据
	_, err := c.RemoveAll(bson.M{"type": f.jsonData.TypeId, "uid": f.jsonData.Uid, "fuid": person, "created": f.jsonData.Created, "infoid": f.jsonData.Infoid, "tid": f.jsonData.Tid})
	if err != nil {
		return err
	}
	return nil
}
