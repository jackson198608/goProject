package Pushdata

import (
	// "encoding/json"
	// "bufio"
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/recommend/mysql"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
	// "github.com/frustra/bbcode"
	// "reflect"
	"math/rand"
)

type RecommendNew struct {
	db      *sql.DB
	session *mgo.Session
}

type DogRecommend struct {
	Id        bson.ObjectId "_id"
	TypeId    int           "type"
	Uid       int           "uid"
	Fuid      int           "fuid" //fans id
	Created   string        "created"
	Infoid    int           "infoid"
	Status    int           "status"
	Tid       int           "tid"
	Bid       int           "bid"
	Content   string        "content"
	Title     string        "title"
	Imagenums int           "image_num"
	Images    string        "images"
	Forum     string        "forum"
	Tag       int           "tag"
	Qsttype   int           "qst_type"
	IsRead    int           "is_read"
	Source    int           "source"
}

type UserRecommend struct {
	Id          bson.ObjectId "_id"
	Uuid        int           "uuid"
	Uid         int           "uid"
	ContentType int           "content_type"
	ContentId   int           "content_id"
	Score       int           "score"
	Time        string        "time"
}

type RecommendRecord struct {
	Id             bson.ObjectId "_id"
	UseRecommendId bson.ObjectId "user_recommend_id"
	TypeId         int           "type"
	Uid            int           "uid"
	Bid            int           "bid"
}

type EventLogX struct {
	Id        bson.ObjectId "_id"
	TypeId    int           "type"
	Uid       int           "uid"
	Fuid      int           "fuid" //fans id
	Created   string        "created"
	Infoid    int           "infoid"
	Status    int           "status"
	Tid       int           "tid"
	Bid       int           "bid"
	Content   string        "content"
	Title     string        "title"
	Imagenums int           "image_num"
	Images    string        "images"
	Forum     string        "forum"
	Tag       int           "tag"
	Qsttype   int           "qst_type"
	IsRead    int           "is_read"
	Source    int           "source"
}

func RecommendUser(logLevel int, db *sql.DB, session *mgo.Session) *RecommendNew {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(RecommendNew)
	e.db = db
	e.session = session //主库
	// e.slave = slave     //从库
	return e
}

func savePostToFansData(i int, Fuid int, contentId int, db *sql.DB, session *mgo.Session) error {
	Qsttype := 0
	Bid := 0
	Tag := 0
	Source := 1
	IsRead := 0
	tableNumX := Fuid % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	c := session.DB("FansData").C(tableNameX)

	isExists := mysql.CheckThreadIsExist(contentId, db)
	logger.Info(contentId, " CheckThreadIsExist is ", isExists)
	if isExists != 0 {
		Images := ""
		Imagenums := 0
		TypeId := 9
		Created := changeCreated(i)

		postData := mysql.GetPostData(contentId, db)
		Tag = mysql.GetTagData(contentId, db)
		if len(postData) > 0 {
			Title := postData[0].Subject
			Content := postData[0].Message
			Content = show_substr(Content, 100)

			Forum := postData[0].Name
			Uid := postData[0].Authorid

			m := EventLogX{bson.NewObjectId(), TypeId, Uid, Fuid, Created, contentId, 1, 0, Bid, Content, Title, Imagenums, Images, Forum, Tag, Qsttype, IsRead, Source}
			err := c.Insert(&m) //插入数据
			if err != nil {
				logger.Error("mongodb insert fans data", err, c)
				return err
			}
			logger.Info("slave FansData mongodb push fans data ", m)
		} else {
			logger.Error("user recommend post data is empty, tid is ", contentId)
		}
	}
	return nil
}

func (e *RecommendNew) PushActiveUserRecommendTask(uid int) error {
	// userRecommends := GetUserRecommendData(uid, e.session)
	// if len(userRecommends) == 0 {
	// 	logger.Info("userRecommends arr is empty")
	// 	return nil
	// }
	// logger.Info(uid, " user Recommends arr is ", userRecommends)
	// for i := 0; i < len(userRecommends); i++ {
	// 	Fuid := userRecommends[i].Uid
	// 	contentType := userRecommends[i].ContentType
	// 	contentId := userRecommends[i].ContentId

	// 	if contentType == 1 {
	// 		savePostToFansData(i, Fuid, contentId, e.db, e.session)
	// 	} else if contentType == 6 {
	// 		saveVideoToFansData(i, Fuid, contentId, e.db, e.session)
	// 	} else {
	// 		logger.Info("user recommend data type is ", contentType)
	// 	}
	// }
	// return nil

}

//全部活跃用户
func GetAllActiveUsers(mongoConn string) []int {
	var user []int
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		logger.Error("[error] connect mongodb err")
		return user
	}

	defer session.Close()

	c := session.DB("ActiveUser").C("active_user")
	err = c.Find(nil).Distinct("uid", &user)
	if err != nil {
		panic(err)
	}
	return user
}

//---------------------------------------------
//推荐关注的人
func RecommendUsers() []*RecommendUser {

}

//最近10分钟登录的用户
func GetLastLoginUids(mongoConn string) []int {
	var user []int
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		logger.Error("[error] connect mongodb err")
		return user
	}
	defer session.Close()
	c := session.DB("UserAction").C("user_login")
	logintime := time.Now().Add(-time.minute * 10)
	err = c.Find(bson.M{"created": bson.M{"$gt": today}}).All(&user)
	if err != nil {
		panic(err)
	}
	return user
}
