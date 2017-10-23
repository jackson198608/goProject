package focus

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/allPersons"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/breedPersons"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/clubPersons"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/fansPersons"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/FansData"
	"gouminGitlab/common/orm/mysql/new_dog123"
	"math"
	"strconv"
	"strings"
)

type Focus struct {
	mysqlXorm *xorm.Engine //@todo to be []
	mongoConn *mgo.Session //@todo to be []
	jobstr    string
	jsonData  *jsonColumn
}

//json column
type jsonColumn struct {
	TypeId    int
	Uid       int
	Created   string
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

const count = 1000

func NewFocus(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jobStr string) *Focus {
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
		return nil
	}
	f.jsonData = jsonColumn

	return f

}

//TypeId = 1 bbs, push fans and club active persons
//TypeId = 6 video, push fans active persons
//TypeId = 8 bbs, push fans and breed active persons
//TypeId = 9 recommend bbs, push all active persons
//TypeId = 15 recommend video, push all active persons
func (f *Focus) Do() error {
	if f.jsonData.TypeId == 1 {
		fp := fansPersons.NewFansPersons(f.mysqlXorm, f.mongoConn, f.formatData(), f.jsonData.Status, f.jsonData.Uid)
		err := fp.Do()
		if err != nil {
			return err
		}
		cp := clubPersons.NewClubPersons(f.mysqlXorm, f.mongoConn, f.formatData(), f.jsonData.Status, f.jsonData.Fid)
		err = cp.Do()
		if err != nil {
			return err
		}
	} else if f.jsonData.TypeId == 6 {
		fp := fansPersons.NewFansPersons(f.mysqlXorm, f.mongoConn, f.formatData(), f.jsonData.Status, f.jsonData.Uid)
		err := fp.Do()
		if err != nil {
			return err
		}
	} else if f.jsonData.TypeId == 8 {
		bp := breedPersons.NewBreedPersons(f.mysqlXorm, f.mongoConn, f.formatData(), f.jsonData.Status, f.jsonData.Bid)
		err := bp.Do()
		if err != nil {
			return err
		}
	} else if ((f.jsonData.TypeId == 9) || (f.jsonData.TypeId == 15)) && (f.jsonData.Source) == 1 {
		ap := allPersons.NewAllPersons(f.mysqlXorm, f.mongoConn, f.formatData(), f.jsonData.Status)
		err := ap.Do()
		if err != nil {
			return err
		}
	} else {
		ap := allPersons.NewAllPersons(f.mysqlXorm, f.mongoConn, f.formatData(), f.jsonData.Status)
		err := ap.Do()
		if err != nil {
			return err
		}
	}
	return nil
}

//change json colum to object private member
func (f *Focus) parseJson() (*jsonColumn, error) {
	var jsonC jsonColumn
	js, err := simplejson.NewJson([]byte(f.jobstr))
	if err != nil {
		return &jsonC, err
	}

	jsonC.Uid, _ = js.Get("uid").Int()
	jsonC.TypeId, _ = js.Get("event_type").Int()
	jsonC.Created, _ = js.Get("time").String()
	jsonC.Tid, _ = js.Get("tid").Int()
	jsonC.Bid, _ = js.Get("event_info").Get("bid").Int()
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
