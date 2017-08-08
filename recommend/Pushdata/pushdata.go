package Pushdata

import (
	// "encoding/json"
	// "bufio"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/recommend/mysql"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"database/sql"
	"time"
	"strconv"
	// "github.com/frustra/bbcode"
	"reflect"
	// "math/rand"
)

type RecommendNew struct {
	db      *sql.DB
	session *mgo.Session
}

type UserRecommend struct {
	Id        bson.ObjectId "_id"
	Uuid        int        "uuid"
	Uid        	int        "uid"
	ContentType	int        "content_type"
	ContentId	int        "content_id"
	Score		int        "score"
	Time		string     "time"
}

type RecommendRecord struct {
	Id        		bson.ObjectId   "_id"
	UseRecommendId  string          "user_recommend_id"
	TypeId    		int           	"type"
	Uid       		int           	"uid"
	Bid       		int           	"bid"
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

func changeCreated(i int) string {
	
	// startTime := time.Now().Unix();
	tm := time.Unix(time.Now().Unix(), 0)
	today := tm.Format("2006-01-02 15:04:05")
	return today
	// if i<8 {
	// 	today = rand.Intn(100)
	// 	return today
	// }else if i>=8 && i<20 {
		
	// }else if i>=20 && i<30 {
		
	// }else if i>=30 && i<38 {
		
	// }else{

	// }
}
func (e *RecommendNew) PushActiveUserRecommendTask(uid int) error {
	session := e.session //主库存储
	
	userRecommends := GetUserRecommendData(uid, e.session)
	if len(userRecommends) == 0 {
		logger.Info("userRecommends arr is empty")
		return nil
	}

	for i:=0; i<len(userRecommends); i++ {
		var Uid int
		var TypeId int
		var Title string
		var Content string
		var Forum string
		var Created string
		var Images string
		var Imagenums int
		Qsttype := 0
		Bid := 0
		Tag := 0
		Source := 1
		IsRead := 0

		Fuid := userRecommends[i].Uid
		contentType := userRecommends[i].ContentType
		contentId := userRecommends[i].ContentId

		tableNumX := Fuid % 100
		if tableNumX == 0 {
			tableNumX = 100
		}
		tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表

		fmt.Println(Fuid)
		fmt.Println(tableNameX)

		c := session.DB("FansData").C(tableNameX)
		if  contentType==1 {
			postData := mysql.GetPostData(contentId, e.db)
			Tag = mysql.GetTagData(contentId, e.db)

			for _, post := range postData {
				Title = post.Subject;
				Content = post.Message;
				Forum = post.Name;
				Uid = post.Authorid;
			}
			Images = "";
			Imagenums = 0;
			TypeId = 9;
			Created = changeCreated(i);

			m := EventLogX{bson.NewObjectId(), TypeId, Uid, Fuid, Created, contentId, 1, 0, Bid, Content, Title, Imagenums, Images, Forum, Tag, Qsttype, IsRead, Source}
			err := c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
				return err
			}
			logger.Info("slave FansData mongodb push fans data ", m)
		
		}else if contentType==6 {
			videoData := mysql.GetVideoData(contentId, e.db)
			for _, video := range videoData {
				Images = post.Thumb;
				Imagenums = 0;
				if Images!="" {
					Imagenums = 1;
				}
				Content = post.Content;
				Forum = post.Name;
				Uid = post.Authorid;
			}
			TypeId = 15;
			// m := EventLogX{bson.NewObjectId(), TypeId, Uid, Fuid, Created, contentId, 1, 0, Bid, Content, Title, Imagenums, Images, Forum, Tag, Qsttype, IsRead, Source}
			// err := c.Insert(&m) //插入数据
			// if err != nil {
			// 	logger.Info("mongodb insert fans data", err, c)
			// 	return err
			// }
			// logger.Info("slave FansData mongodb push fans data ", m)
		
		}
		
	}
	return nil
}

//获取用户推荐数据
func GetUserRecommendData(uid int, session *mgo.Session) []*UserRecommend {
	var user []*UserRecommend
	c := session.DB("BigData").C("user_recommend")
	
	// tm := time.Unix(time.Now().Unix(), 0)
	// today := tm.Format("20060102")
	// fmt.Println(reflect.TypeOf(today))
	// err := c.Find(&bson.M{"uid": uid,"time":{"$gt":today}}).Sort(bson.M{"score": "desc"}).All(&user)
	err := c.Find(&bson.M{"uid": uid}).All(&user)
	if err != nil {
		panic(err)
	}
	return user
}

//获取根据犬种推荐数据
func GetDogRecommendData(uid int, bid int, session *mgo.Session) []*UserRecommend {
	var user []*UserRecommend
	// c := session.DB("RecommendData").C("recommend_by_dog_or_age")

	// err := c.Find(&bson.M{"uid": uid,"time":{"$gt":today}}).Sort(bson.M{"score": "desc"}).All(&user)
	// if err != nil {
	// 	panic(err)
	// }
	return user
}

//获取犬种推荐数据记录
func GetRecommendRecordLastId(uid int, bid int, session *mgo.Session) string {
	var lastId string
	c := session.DB("RecommendData").C("user_recommend_record")
	err := c.Find(&bson.M{"uid": uid,"type":2,"bid": bid}).One(&lastId)
	if err != nil {
		panic(err)
	}
	return lastId
}

//更新犬种推荐数据记录last id
func updateRecommendRecordLastId(newLastId string, uid int, bid int, session *mgo.Session) string {
	var lastId string
	c := session.DB("RecommendData").C("user_recommend_record")
	err := c.Find(&bson.M{"uid": uid,"type":2,"bid": bid}).One(&lastId)
	if err != nil {
		panic(err)
	}
	if lastId!="" {
		selector := bson.M{"uid": uid,"bid": bid,"type": "2"}
		data := bson.M{"$set": bson.M{"UseRecommendId": newLastId}}
		err := c.Update(selector, data)
		if err != nil {
		    panic(err)
		}
		logger.Info("mongodb update user_recommend_record data:", newLastId, uid, bid)
	}else{
		m1 := RecommendRecord{bson.NewObjectId(), newLastId, 2, uid, bid}
	    err = c.Insert(&m1)
	    if err != nil {
	        panic(err)
	    }
	    logger.Info("mongodb inster user_recommend_record data:", newLastId, uid, bid)
	}
	return lastId
}

//全部活跃用户
func GetAllActiveUsers() []int {
	var user []int
	mongoConn := "192.168.86.68:27017"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		logger.Error("[error] connect mongodb err")
		return user
	}

	c := session.DB("ActiveUser").C("active_user")
	err = c.Find(nil).Distinct("uid", &user)
	if err != nil {
		panic(err)
	}
	return user
}
