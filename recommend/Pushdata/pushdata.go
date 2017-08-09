package Pushdata

import (
	// "encoding/json"
	// "bufio"
	"database/sql"
	"fmt"
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
	UseRecommendId bson.ObjectId        "user_recommend_id"
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

func GetNextTimestamp(second int) int64 {
	t := time.Now()
	tm := t.Unix()
	currenttime := tm + int64(second)
	return currenttime
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func changeCreated(i int) string {
	if i < 8 {
		startTime := time.Now().Unix()
		endTime := GetNextTimestamp(3600*5 + 60*59)
		currenttime := RandInt64(startTime, endTime)
		showTime := time.Unix(currenttime, 0).Format("2006-01-02 15:04:05")
		return showTime
	} else if i >= 8 && i < 20 {
		startTime := GetNextTimestamp(3600 * 6)
		endTime := GetNextTimestamp(3600*11 + 60*59)
		currenttime := RandInt64(startTime, endTime)
		showTime := time.Unix(currenttime, 0).Format("2006-01-02 15:04:05")
		return showTime

	} else if i >= 20 && i < 30 {
		startTime := GetNextTimestamp(3600 * 12)
		endTime := GetNextTimestamp(3600*13 + 60*59)
		currenttime := RandInt64(startTime, endTime)
		showTime := time.Unix(currenttime, 0).Format("2006-01-02 15:04:05")
		return showTime
	} else if i >= 30 && i < 38 {
		startTime := GetNextTimestamp(3600 * 14)
		endTime := GetNextTimestamp(3600*17 + 60*59)
		currenttime := RandInt64(startTime, endTime)
		showTime := time.Unix(currenttime, 0).Format("2006-01-02 15:04:05")
		return showTime
	} else {
		startTime := GetNextTimestamp(3600 * 18)
		endTime := GetNextTimestamp(3600*23 + 60*59)
		currenttime := RandInt64(startTime, endTime)
		showTime := time.Unix(currenttime, 0).Format("2006-01-02 15:04:05")
		return showTime
	}
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
	logger.Info("CheckThreadIsExist is ", isExists)
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
			logger.Info("user recommend post data is empty, tid is ", contentId)
		}
	}
	return nil
}

func saveVideoToFansData(i int, Fuid int, contentId int, db *sql.DB, session *mgo.Session) error {
	Qsttype := 0
	Bid := 0
	Tag := 0
	Source := 1
	IsRead := 0
	Imagenums := 0
	Images := ""
	Content := ""
	Forum := ""
	TypeId := 15

	tableNumX := Fuid % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	c := session.DB("FansData").C(tableNameX)

	// contentId = 48428
	videoData := mysql.GetVideoData(contentId, db)
	fmt.Println("get video data vid is ", contentId)
	if len(videoData) > 0 {
		Images = videoData[0].Thumb
		if Images != "" {
			Images = "http://c1.cdn.goumin.com/diary/" + videoData[0].Thumb
			Imagenums = 1
		}
		Title := videoData[0].Content
		Uid := videoData[0].Uid
		Created := changeCreated(i)
		m := EventLogX{bson.NewObjectId(), TypeId, Uid, Fuid, Created, contentId, 1, 0, Bid, Content, Title, Imagenums, Images, Forum, Tag, Qsttype, IsRead, Source}
		err := c.Insert(&m) //插入数据
		if err != nil {
			logger.Error("mongodb insert fans data", err, c)
			return err
		}
		logger.Info("slave FansData mongodb push fans data ", m)
	}
	logger.Info("video data is empty, video id is ", contentId)
	return nil
}

func savePetBreedRecommendData(i int, Fuid int, dogRecommend *DogRecommend, db *sql.DB, session *mgo.Session) error{
	tableNumX := Fuid % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	c := session.DB("FansData").C(tableNameX)
	Created := changeCreated(i)
	TypeId := dogRecommend.TypeId
	Source := 1
	if dogRecommend.TypeId==1 {
		TypeId = 9
	}else if dogRecommend.TypeId==6 {
		TypeId = 15
	}
	m := EventLogX{bson.NewObjectId(), TypeId, dogRecommend.Uid, Fuid, Created, dogRecommend.Infoid, 1, 0, dogRecommend.Bid, dogRecommend.Content, dogRecommend.Title, dogRecommend.Imagenums, dogRecommend.Images, dogRecommend.Forum, dogRecommend.Tag, dogRecommend.Qsttype, dogRecommend.IsRead, Source}
	err := c.Insert(&m) //插入数据
	if err != nil {
		logger.Error("mongodb insert fans data", err, c)
		return err
	}
	logger.Info("slave FansData mongodb push fans data ", m)
	return nil
}
func GetPetBreedRecommendData(uid int, db *sql.DB, session *mgo.Session) error {
	mongoConn := "192.168.86.104:27017"
	session,_ = mgo.Dial(mongoConn)
	
	Breed := mysql.GetPetBreed(uid, db);
	if len(Breed)>0 {
		for i := 0; i < len(Breed); i++ {
			bid := Breed[i].Bid
			dogRecommend := GetDogRecommendData(uid, bid, session);
			logger.Info("dogRecommend len ", len(dogRecommend))
			if len(dogRecommend)>0 {
				var newlastId bson.ObjectId
				for d := 0; d < len(dogRecommend); d++ {
					savePetBreedRecommendData(d, uid, dogRecommend[d], db, session);
					newlastId = dogRecommend[d].Id;
				}
				updateRecommendRecordLastId(newlastId, uid, bid, session)
			}
			logger.Error("get dog recommend data is empty, uid is ", uid, " bid is ",bid)
			return nil
		}
	}
	logger.Error("user breed is empty, uid is ", uid)
	return nil
}

func (e *RecommendNew) PushActiveUserRecommendTask(uid int) error {
	session := e.session //主库存储

	userRecommends := GetUserRecommendData(uid, e.session)
	if len(userRecommends) == 0 {
		logger.Info("userRecommends arr is empty")
		return nil
	}
	for i := 0; i < len(userRecommends); i++ {
		Fuid := userRecommends[i].Uid
		contentType := userRecommends[i].ContentType
		contentId := userRecommends[i].ContentId

		if contentType == 1 {
			savePostToFansData(i, Fuid, contentId, e.db, session)
		} else if contentType == 6 {
			saveVideoToFansData(i, Fuid, contentId, e.db, session)
		} else {
			logger.Info("user recommend data type is ", contentType)
		}
	}

	GetPetBreedRecommendData(uid, e.db, e.session);
	return nil
}

//获取用户推荐数据
func GetUserRecommendData(uid int, session *mgo.Session) []*UserRecommend {
	var user []*UserRecommend
	c := session.DB("BigData").C("user_recommend")

	today := time.Unix(time.Now().Unix(), 0).Format("20060102")
	today = ""
	err := c.Find(&bson.M{"uid": uid,"time":today}).Sort("score desc").All(&user)
	if err != nil {
		panic(err)
	}
	return user
}

//获取根据犬种推荐数据
func GetDogRecommendData(uid int, bid int, session *mgo.Session) []*DogRecommend {
	var data []*DogRecommend
	c := session.DB("RecommendData").C("recommend_by_dog_or_age")
	bids := []int{0, bid}
	RecommendRecord := GetRecommendRecordLastId(uid, bid, session)

	if len(RecommendRecord)>0 {
		lastId := RecommendRecord[0].UseRecommendId
		err := c.Find(bson.M{"bid":bson.M{"$in": bids},"_id":bson.M{"$gt": lastId}}).Limit(50).All(&data)
		if err != nil {
			logger.Error(" get recommend record mongodb find data", err, c)
		}
	}else{
		err := c.Find(bson.M{"bid":bson.M{"$in": bids}}).Limit(50).All(&data)
		if err != nil {
			logger.Error(" get recommend record mongodb find data", err, c)
		}
	}
	
	return data
}

//获取上次犬种推荐数据记录ID
func GetRecommendRecordLastId(uid int, bid int, session *mgo.Session) []*RecommendRecord {
	var rr []*RecommendRecord
	c := session.DB("RecommendData").C("user_recommend_record")
	fmt.Println(uid)
	fmt.Println(bid)
	err := c.Find(bson.M{"uid": uid, "type": 2, "bid": bid}).All(&rr)
	if err != nil {
		logger.Error(" get recommend record lastId mongodb find data", err, c)
		return rr
	}
	return rr
}

//更新犬种推荐数据记录last id
func updateRecommendRecordLastId(newLastId bson.ObjectId, uid int, bid int, session *mgo.Session) error {
	logger.Info("in update RecommendRecord")
	rr := []RecommendRecord{}
	c := session.DB("RecommendData").C("user_recommend_record")
	err := c.Find(bson.M{"uid": uid, "type": 2, "bid": bid}).All(&rr)

	if err != nil {
		logger.Error("mongodb find data", err, c)
	}
	logger.Info(" RecommendRecord len ", len(rr))
	if len(rr)>0 {

		selector := bson.M{"uid": uid, "bid": bid, "type": "2"}
		data := bson.M{"$set": bson.M{"UseRecommendId": newLastId}}
		err := c.Update(selector, data)
		if err != nil {
			logger.Error("update recommend record to mongodb error ", err, c)
			return nil
		}
		logger.Info("mongodb update user_recommend_record data:", newLastId, uid, bid)
	} else {
		m1 := RecommendRecord{bson.NewObjectId(), newLastId, int(2), int(uid), int(bid)}

		err = c.Insert(&m1)
		if err != nil {
			logger.Error("save recommend record to mongodb error ", err, c)
			return nil
		}
		logger.Info("mongodb inster user_recommend_record data:", newLastId, uid, bid)
	}
	return nil
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
