package recommend

import (
	// "fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/recommendAllPersons"
	mgo "gopkg.in/mgo.v2"
)

const (
	mongoConn = "192.168.5.22:27017,192.168.5.26:27017,192.168.5.200:27017"
	//mongoConn = "192.168.86.192:27017" //@todo change online dsn
)

var m map[int]bool

type Recommend struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jobstr    string
	jsonData  *job.RecommendJsonColumn
}

func init() {
	m = make(map[int]bool)
	m = loadActiveUserToMap()
}

func NewRecommend(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jobStr string) *Recommend {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jobStr == "") {
		return nil
	}

	f := new(Recommend)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jobstr = jobStr

	//@todo pass params
	jsonColumn, err := f.parseJson()
	if err != nil {
		return nil
	}
	f.jsonData = jsonColumn
	return f
}

//TypeId = 15 recommend, push all active persons
func (f *Recommend) Do() error {

	ap := recommendAllPersons.NewAllPersons(f.mysqlXorm, f.mongoConn, f.jsonData, &m)
	err := ap.Do()
	if err != nil {
		return err
	}

	return nil
}

//change json colum to object private member
func (f *Recommend) parseJson() (*job.RecommendJsonColumn, error) {
	var jsonC job.RecommendJsonColumn
	js, err := simplejson.NewJson([]byte(f.jobstr))
	if err != nil {
		return &jsonC, err
	}

	jsonC.Uid, _ = js.Get("uid").Int()
	jsonC.Created, _ = js.Get("create").String()
	jsonC.Infoid, _ = js.Get("infoid").Int()
	jsonC.Type, _ = js.Get("type").Int()
	jsonC.Title, _ = js.Get("title").String()
	jsonC.Description, _ = js.Get("description").String()
	jsonC.Images, _ = js.Get("image").String()
	jsonC.Rauth, _ = js.Get("rauth").String()
	jsonC.Imagenums, _ = js.Get("image_num").Int()
	jsonC.QstType, _ = js.Get("qst_type").Int()
	jsonC.AdType, _ = js.Get("ad_type").Int()
	jsonC.AdUrl, _ = js.Get("ad_url").String()
	jsonC.Rauth, _ = js.Get("rauth").String()
	jsonC.Action, _ = js.Get("action").Int() //行为 -1 删除 0 插入 1 修改

	return &jsonC, nil
}

func loadActiveUserToMap() map[int]bool {
	var m map[int]bool
	m = make(map[int]bool)
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		return m
	}
	defer session.Close()

	var uids []int
	c := session.DB("ActiveUser").C("active_user")
	err = c.Find(nil).Distinct("uid", &uids)
	if err != nil {
		// panic(err)
		return m
	}
	for i := 0; i < len(uids); i++ {
		m[uids[i]] = true
	}
	return m
}
