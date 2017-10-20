package focus

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mysql/new_dog123"
	"strings"
)

type Focus struct {
	mysqlXorm *xorm.Engine
	mongoConn *mgo.Session
	jobstr    string
	jsonData  *jsonColumn
}

//json column
type jsonColumn struct {
	TypeId    int
	Uid       int
	Created   int
	Infoid    int
	Status    int
	Tid       int
	Bid       int
	Fid       int
	Content   string
	Title     string
	Imagenums int
	Forum     string
	Tag       int
	Qsttype   int
	Source    int
}

const count = 100

func NewFocus(mysqlXorm *xorm.Engine, mongoConn *mgo.Session, jobStr string) *Focus {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jobStr == "") {
		return nil
	}

	f := new(Focus)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jobstr = jobStr

	//@todo pass params
	jsonColumn, err := f.parseJson()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	f.jsonData = jsonColumn

	return f

}

func (f *Focus) Do() error {
	page := f.getPersionsPageNum()
	if page <= 0 {
		return nil
	}

	for i := 1; i <= page; i++ {
		currentPersionList := f.getPersons(page)
		f.pushPersons(currentPersionList)

	}
	return nil
}

//change json colum to object private member
func (f *Focus) parseJson() (*jsonColumn, error) {
	var jsonC jsonColumn

	jobs := strings.Split(f.jobstr, "|")
	if len(jobs) <= 1 {
		return &jsonC, errors.New("you have no job")
	}

	jsonStr := jobs[0]
	js, err := simplejson.NewJson([]byte(jsonStr))
	if err != nil {
		return &jsonC, err
	}

	jsonC.Uid, _ = js.Get("uid").Int()
	jsonC.TypeId, _ = js.Get("event_type").Int()
	jsonC.Created, _ = js.Get("time").Int()
	jsonC.Tid, _ = js.Get("tid").Int()
	jsonC.Bid, _ = js.Get("bid").Int()
	jsonC.Infoid, _ = js.Get("event_info").Get("infoid").Int()
	jsonC.Title, _ = js.Get("event_info").Get("title").String()
	jsonC.Content, _ = js.Get("event_info").Get("content").String()
	jsonC.Forum, _ = js.Get("event_info").Get("forum").String()
	jsonC.Imagenums, _ = js.Get("event_info").Get("image_num").Int()
	jsonC.Tag, _ = js.Get("event_info").Get("tag").Int()
	jsonC.Qsttype, _ = js.Get("event_info").Get("qst_type").Int()
	jsonC.Fid, _ = js.Get("event_info").Get("fid").Int()
	jsonC.Source, _ = js.Get("source").Int()
	jsonC.Status, _ = js.Get("status").Int()

	return &jsonC, nil
}

func (f *Focus) pushPersons(persons []int) error {
	if persons == nil {
		return errors.New("you have no person to push " + f.jobstr)
	}

	for _, person := range persons {
		err := f.pushPerson(person)
		if err != nil {
			//@todo if err times < 5 ,just print log
			//      if err times > 5 ,return err
		}
	}
	return nil
}

func (f *Focus) pushPerson(person int) error {

	return nil
}

//@todo how to remove duplicate uid from to lists
func (f *Focus) getPersons(page int) []int {
	var uid []int
	return uid
}

func (f *Focus) getPersionsPageNum() int {

	return 0
}

//获取相同犬种的活跃用户
func (f *Focus) getBreedPersons(page int) []int {
	var uid []int
	return uid
}

//获取相同俱乐部的活跃用户
func (f *Focus) getClubPersons(page int) []int {
	var uids []int
	fid := f.jsonData.Fid
	c := f.mongoConn.DB("ActiveUser").C("active_forum_user")
	err := c.Find(&bson.M{"forum_id": fid}).Distinct("uid", &uids)
	if err != nil {
		panic(err)
	}
	return uids
}

//获取活跃粉丝用户
func (f *Focus) getFansPersons(page int) []int {
	var uids []int
	uid := f.jsonData.Uid
	var follows []new_dog123.Follow
	err := f.mysqlXorm.Where("user_id=? and fans_active=1", uid).Cols("follow_id").Limit(count, (page-1)*count).Find(&follows)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, v := range follows {
		uids = append(uids, v.FollowId)
	}
	return uids
}

//获取所有活跃用户
func (f *Focus) getActivePersons(page int) []int {
	var uids []int
	c := f.mongoConn.DB("ActiveUser").C("active_user")
	err := c.Find(nil).Distinct("uid", &uids)
	if err != nil {
		panic(err)
	}
	return uids
}
