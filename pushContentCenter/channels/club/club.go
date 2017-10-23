package club

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mysql/new_dog123"
	// "reflect"
	"strconv"
	"strings"
)

type Club struct {
	mysqlXorm *xorm.Engine
	mongoConn *mgo.Session
	jobstr    string
	jsonData  *jsonColumn
}

//json column
type jsonColumn struct {
	TypeId       int
	Uid          int
	Fid          string //push to the club, if fid is 0, to all clubs
	Created      int
	Infoid       int
	Status       int
	Type         int
	Content      string
	Title        string
	Imagenums    int
	Disgest      int
	Qsttype      int
	Lastpost     int
	Lastposter   string
	Displayorder int
}

type ClubX struct {
	Id           bson.ObjectId "_id"
	Type         int           "type"
	TypeId       int           "typeid" //俱乐部主题ID
	Infoid       int           "infoid"
	Uid          int           "uid"
	Created      int           "created"
	Status       int           "status"
	Content      string        "message"
	Title        string        "subject"
	Imagenums    int           "image_num"
	Displayorder int           "displayorder"
	Lastpost     int           "lastpost"
	Lastposter   string        "lastposter"
	Disgest      int           "disgest"
	Qsttype      int           "qst_type"
}

const count = 100

func NewClub(mysqlXorm *xorm.Engine, mongoConn *mgo.Session, jobStr string) *Club {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jobStr == "") {
		return nil
	}

	c := new(Club)
	if c == nil {
		return nil
	}

	c.mysqlXorm = mysqlXorm
	c.mongoConn = mongoConn
	c.jobstr = jobStr

	jsonColumn, err := c.parseJson()
	if err != nil {
		return nil
	}
	c.jsonData = jsonColumn

	return c
}

func (c *Club) Do() error {
	var currentClubList []int
	if c.jsonData.Fid == "" {
		return errors.New("you have no club to push " + c.jobstr)
	}

	fids := strings.Split(c.jsonData.Fid, ",")

	clubid, _ := strconv.Atoi(fids[0])
	if (len(fids) == 1) && (clubid == 0) {
		//推送到所有展示状态俱乐部
		currentClubList = c.getClubs()
	} else {
		//	推送到指定俱乐部
		for i := 0; i < len(fids); i++ {
			clubid, _ := strconv.Atoi(fids[i])
			currentClubList = append(currentClubList, clubid)
		}
	}
	// fmt.Println(currentClubList)
	err := c.pushClubs(currentClubList)

	if err != nil {
		return err
	}
	return nil
}

//change json colum to object private member
func (c *Club) parseJson() (*jsonColumn, error) {
	var jsonC jsonColumn

	jobs := strings.Split(c.jobstr, "|")
	if len(jobs) <= 1 {
		return &jsonC, errors.New("you have no job")
	}

	jsonStr := jobs[0]
	js, err := simplejson.NewJson([]byte(jsonStr))
	if err != nil {
		return &jsonC, err
	}

	jsonC.Uid, _ = js.Get("uid").Int()
	jsonC.Fid, _ = js.Get("fid").String()
	jsonC.Type, _ = js.Get("type").Int()
	jsonC.TypeId, _ = js.Get("typeid").Int()
	jsonC.Created, _ = js.Get("created").Int()
	jsonC.Infoid, _ = js.Get("infoid").Int()
	jsonC.Title, _ = js.Get("subject").String()
	jsonC.Content, _ = js.Get("message").String()
	jsonC.Imagenums, _ = js.Get("image_num").Int()
	jsonC.Lastpost, _ = js.Get("lastpost").Int()
	jsonC.Lastposter, _ = js.Get("lastposter").String()
	jsonC.Status, _ = js.Get("status").Int()
	jsonC.Displayorder, _ = js.Get("displayorder").Int()
	jsonC.Disgest, _ = js.Get("disgest").Int()
	jsonC.Qsttype, _ = js.Get("qst_type").Int()

	return &jsonC, nil
}

func (c *Club) pushClubs(clubs []int) error {
	if clubs == nil {
		return errors.New("you have no club to push " + c.jobstr)
	}

	for _, club := range clubs {
		err := c.pushClub(club)
		if err != nil {
			c.tryPushClub(club, 1)
		}
	}
	return nil
}

func (c *Club) pushClub(club int) error {
	tableNameX := "forum_content_" + strconv.Itoa(club)
	mc := c.mongoConn.DB("ClubData").C(tableNameX)
	if c.jsonData.Status == 1 {
		//数据展示
		data := c.pushData(club)
		err := mc.Insert(&data) //插入数据
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Club) tryPushClub(club int, num int) error {
	if num > 5 {
		return errors.New("Attempting to push has failed 5 times: " + c.jobstr + "; club is " + strconv.Itoa(club))
	}
	err := c.pushClub(club)
	if err != nil {
		c.tryPushClub(club, num+1)
	}
	return nil
}

func (c *Club) pushData(club int) *ClubX {
	var data ClubX
	data = ClubX{bson.NewObjectId(),
		c.jsonData.Type,
		c.jsonData.TypeId,
		c.jsonData.Infoid,
		c.jsonData.Uid,
		c.jsonData.Created,
		c.jsonData.Status,
		c.jsonData.Content,
		c.jsonData.Title,
		c.jsonData.Imagenums,
		c.jsonData.Displayorder,
		c.jsonData.Lastpost,
		c.jsonData.Lastposter,
		c.jsonData.Disgest,
		c.jsonData.Qsttype}
	return &data
}

//@todo how to remove duplicate uid from to lists
func (c *Club) getClubs() []int {
	var cluds []int
	var forums []new_dog123.PreForumForum
	err := c.mysqlXorm.Where("status=1 and fup!=0").Cols("fid").Find(&forums)
	if err != nil {
		return nil
	}
	for _, v := range forums {
		cluds = append(cluds, v.Fid)
	}
	return cluds
}
