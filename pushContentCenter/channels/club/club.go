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
	TypeId    int
	Uid       int
	Fid       int
	Created   int
	Infoid    int
	Status    int
	Content   string
	Title     string
	Imagenums int
}

const count = 100

func NewClub(mysqlXorm *xorm.Engine, mongoConn *mgo.Session, jobStr string) *Club {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jobStr == "") {
		return nil
	}

	f := new(Club)
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

func (c *Club) Do() error {
	currentPersionList := f.getClubs()
	f.pushClubs(currentPersionList)
	return nil
}

//change json colum to object private member
func (c *Club) parseJson() (*jsonColumn, error) {
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
	jsonC.TypeId, _ = js.Get("type").Int()
	jsonC.Created, _ = js.Get("time").Int()
	jsonC.Infoid, _ = js.Get("infoid").Int()
	jsonC.Title, _ = js.Get("title").String()
	jsonC.Content, _ = js.Get("content").String()
	jsonC.Imagenums, _ = js.Get("image_num").Int()
	jsonC.Fid, _ = js.Get("fid").Int()
	jsonC.Status, _ = js.Get("status").Int()

	return &jsonC, nil
}

func (c *Club) pushClubs(clubs []int) error {
	if clubs == nil {
		return errors.New("you have no club to push " + f.jobstr)
	}

	for _, club := range persons {
		err := f.pushClub(club)
		if err != nil {
			//@todo if err times < 5 ,just print log
			//      if err times > 5 ,return err
		}
	}
	return nil
}

func (c *Club) pushClub(club int) error {

	return nil
}

//@todo how to remove duplicate uid from to lists
func (c *Club) getClubs() []int {
	var club_id []int
	return club_id
}
