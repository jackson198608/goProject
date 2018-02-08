package recommend

import (
	// "fmt"
	"github.com/bitly/go-simplejson"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/recommendAllPersons"
	mgo "gopkg.in/mgo.v2"
)

const (
	mongoConn = "192.168.5.22:27017,192.168.5.26:27017,192.168.5.200:27017"
	// mongoConn = "192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019" //@todo change online dsn
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
	// fmt.Println("in recommend")
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

	logger.Info("[new recommend] jobStr is ", f.jobstr)

	//@todo pass params
	jsonColumn, err := f.parseJson()
	if err != nil {
		return nil
	}
	f.jsonData = jsonColumn
	return f
}

func (f *Recommend) Do() error {
	//推送所有用户
	// fmt.Println(f.jsonData.RecommendType)
	if f.jsonData.RecommendType == "all" {
		logger.Info("[recommend do] jsonData is ", f.jsonData)
		ap := recommendAllPersons.NewRecommendAllPersons(f.mysqlXorm, f.mongoConn, f.jsonData, &m)
		err := ap.Do()
		if err != nil {
			return err
		}
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

	jsonC.Uid, _ = js.Get("uid").Int()       //发布内容uid
	jsonC.Ruid, _ = js.Get("ruid").Int()     //推荐用户uid
	jsonC.Infoid, _ = js.Get("infoid").Int() //内容ID
	jsonC.Pid, _ = js.Get("pid").Int()       //内容ID
	jsonC.Type, _ = js.Get("type").Int()     //内容类型 1帖子 6视频 8问答 13广告 18 宠家号文章 19宠家号视频
	jsonC.Tag, _ = js.Get("tag").Int()       //热门话题ID
	jsonC.Tags, _ = js.Get("tags").String()  //标签
	jsonC.QstType, _ = js.Get("qst_type").Int()
	jsonC.AdType, _ = js.Get("ad_type").Int()
	jsonC.AdUrl, _ = js.Get("ad_url").String()
	jsonC.Title, _ = js.Get("title").String()
	jsonC.Description, _ = js.Get("description").String()
	jsonC.Images, _ = js.Get("images").String()
	jsonC.Imagenums, _ = js.Get("image_num").Int()
	jsonC.VideoUrl, _ = js.Get("video_url").String() //认证信息
	jsonC.Created, _ = js.Get("created").Int()
	jsonC.Action, _ = js.Get("action").Int()                   //行为 -1 删除 0 插入
	jsonC.Channel, _ = js.Get("channel").Int()                 //展示渠道 1精选 2视频 3游记 4宠家号
	jsonC.Duration, _ = js.Get("duration").Int()               //视频时长
	jsonC.RecommendType, _ = js.Get("recommend_type").String() //推送方式 all 全部用户

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
