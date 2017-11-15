package club

import (
	"errors"
	// "fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/ClubData"
	"gouminGitlab/common/orm/mysql/new_dog123"
	// "reflect"
	"strconv"
	"strings"
)

type Club struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jobstr    string
	jsonData  *jsonColumn
}

//json column
type jsonColumn struct {
	Fid          string //push to the club, if fid is 0, to all clubs
	Infoid       int
	Uid          int
	Type         int
	TypeId       int
	Title        string
	Content      string
	Imagenums    int
	Created      int
	Lastpost     int
	Lastposter   string
	Status       int
	Displayorder int
	Digest       int
	Qsttype      int
	ThreadStatus int
	Cover        int
	Closed       int
	Highlight    int
	Sortid       int
	Recommends   int
	Special      int
	Replies      int
	Isgroup      int
	Price        int
	Heats        int
	Action       int
}

func NewClub(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jobStr string) *Club {
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
	err := c.pushClubs(currentClubList)

	if err != nil {
		return err
	}
	return nil
}

//change json colum to object private member
func (c *Club) parseJson() (*jsonColumn, error) {
	var jsonC jsonColumn
	js, err := simplejson.NewJson([]byte(c.jobstr))
	if err != nil {
		return &jsonC, err
	}

	jsonC.Fid, _ = js.Get("fid").String()
	jsonC.Infoid, _ = js.Get("infoid").Int()
	jsonC.Uid, _ = js.Get("uid").Int()
	jsonC.Type, _ = js.Get("type").Int()
	jsonC.TypeId, _ = js.Get("typeid").Int()
	jsonC.Title, _ = js.Get("subject").String()
	jsonC.Content, _ = js.Get("message").String()
	jsonC.Imagenums, _ = js.Get("image_num").Int()
	jsonC.Created, _ = js.Get("created").Int()
	jsonC.Lastpost, _ = js.Get("lastpost").Int()
	jsonC.Lastposter, _ = js.Get("lastposter").String()
	jsonC.Status, _ = js.Get("status").Int()
	jsonC.Displayorder, _ = js.Get("displayorder").Int()
	jsonC.Digest, _ = js.Get("digest").Int()
	jsonC.Qsttype, _ = js.Get("qst_type").Int()
	jsonC.ThreadStatus, _ = js.Get("thread_status").Int()
	jsonC.Cover, _ = js.Get("cover").Int()
	jsonC.Closed, _ = js.Get("closed").Int()
	jsonC.Highlight, _ = js.Get("highlight").Int()
	jsonC.Sortid, _ = js.Get("sortid").Int()
	jsonC.Recommends, _ = js.Get("recommend").Int()
	jsonC.Special, _ = js.Get("special").Int()
	jsonC.Replies, _ = js.Get("replies").Int()
	jsonC.Isgroup, _ = js.Get("isgroup").Int()
	jsonC.Price, _ = js.Get("price").Int()
	jsonC.Heats, _ = js.Get("heats").Int()
	jsonC.Action, _ = js.Get("action").Int()

	return &jsonC, nil
}

func (c *Club) pushClubs(clubs []int) error {
	if clubs == nil {
		return errors.New("you have no club to push " + c.jobstr)
	}

	for _, club := range clubs {
		err := c.pushClub(club)
		if err != nil {
			for i := 0; i < 5; i++ {
				err := c.pushClub(club)
				if err == nil {
					break
				}
			}
		}
	}
	return nil
}

func (c *Club) pushClub(club int) error {
	tableNameX := "forum_content_" + strconv.Itoa(club)
	mc := c.mongoConn[0].DB("ClubData").C(tableNameX)
	if c.jsonData.Action == 0 {
		// fmt.Println("insert" + strconv.Itoa(club))

		err := c.insertClub(mc)
		if err != nil {
			return err
		}
	} else if c.jsonData.Action == 1 {
		//修改数据状态
		fmt.Println("update" + strconv.Itoa(club))
		err := c.updateClub(mc)
		if err != nil {
			return err
		}
	} else if c.jsonData.Action == -1 {
		//删除数据
		// fmt.Println("remove" + strconv.Itoa(club))
		err := c.removeClub(mc)
		if err != nil {
			return err
		}
	}
	return nil
}

//@todo how to remove duplicate uid from to lists
func (c *Club) getClubs() []int {
	var cluds []int
	var forums []new_dog123.PreForumForum
	err := c.mysqlXorm[0].Where("status=1 and fup!=0").Cols("fid").Find(&forums)
	if err != nil {
		return nil
	}
	for _, v := range forums {
		cluds = append(cluds, v.Fid)
	}
	return cluds
}

func (c *Club) insertClub(mc *mgo.Collection) error {
	//新增数据
	var data ClubData.ClubX
	data = ClubData.ClubX{bson.NewObjectId(),
		c.jsonData.Infoid,
		c.jsonData.Uid,
		c.jsonData.Type,
		c.jsonData.TypeId,
		c.jsonData.Title,
		c.jsonData.Content,
		c.jsonData.Imagenums,
		c.jsonData.Created,
		c.jsonData.Lastpost,
		c.jsonData.Lastposter,
		c.jsonData.Status,
		c.jsonData.Displayorder,
		c.jsonData.Digest,
		c.jsonData.Qsttype,
		c.jsonData.ThreadStatus,
		c.jsonData.Cover,
		c.jsonData.Closed,
		c.jsonData.Highlight,
		c.jsonData.Sortid,
		c.jsonData.Recommends,
		c.jsonData.Special,
		c.jsonData.Replies,
		c.jsonData.Isgroup,
		c.jsonData.Price,
		c.jsonData.Heats}
	err := mc.Insert(&data) //插入数据
	if err != nil {
		return err
	}
	return nil
}

//update thread
func (c *Club) updateClub(mc *mgo.Collection) error {
	_, err := mc.UpdateAll(bson.M{"type": c.jsonData.Type, "infoid": c.jsonData.Infoid},
		bson.M{"$set": bson.M{"status": c.jsonData.Status,
			"digest":        c.jsonData.Digest,
			"displayorder":  c.jsonData.Displayorder,
			"qst_type":      c.jsonData.Qsttype,
			"cover":         c.jsonData.Cover,
			"closed":        c.jsonData.Closed,
			"highlight":     c.jsonData.Highlight,
			"sortid":        c.jsonData.Sortid,
			"isgroup":       c.jsonData.Isgroup,
			"price":         c.jsonData.Price,
			"thread_status": c.jsonData.ThreadStatus,
			"lastpost":      c.jsonData.Lastpost,
			"lastposter":    c.jsonData.Lastposter,
			"replies":       c.jsonData.Replies,
			"heats":         c.jsonData.Heats}})
	if err != nil {
		return err
	}
	return nil
}

func (c *Club) removeClub(mc *mgo.Collection) error {
	_, err := mc.RemoveAll(bson.M{"type": c.jsonData.Type, "uid": c.jsonData.Uid, "infoid": c.jsonData.Infoid})
	if err != nil {
		return err
	}
	return nil
}
