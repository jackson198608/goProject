package Pushdata

import (
	// "encoding/json"
	// "bufio"
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/recomPosition/mysql"
	mgo "gopkg.in/mgo.v2"
	// "github.com/bitly/go-simplejson"
	"gopkg.in/mgo.v2/bson"
	// "strconv"
	// "time"
	// "github.com/frustra/bbcode"
	// "reflect"	
	// "math/rand"
	"encoding/json"
)

type RecommendNew struct {
	db      *sql.DB
	session *mgo.Session
}

type RecommendPosition struct {
	Id             bson.ObjectId "_id"
	Uid 		   int 			 "uid"
	TypeId         int           "type"
	Info           string        "info"
}

func RecommendUser(logLevel int, db *sql.DB, session *mgo.Session) *RecommendNew {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(RecommendNew)
	e.db = db
	e.session = session //主库
	// e.slave = slave     //从库
	return e
}

func saveData(Uid int, TypeId int, Info string, db *sql.DB, session *mgo.Session) error {
	tableName := "recommend_position"
	c := session.DB("RecommendData").C(tableName)
	m := RecommendPosition{bson.NewObjectId(),  Uid, TypeId, Info}
	err := c.Insert(&m) //插入数据
	if err != nil {
		logger.Error("mongodb insert fans data", err, c)
		return err
	}
	return nil
}

func (e *RecommendNew) PushRecommendTask(uid int) error {
	db := e.db
	session := e.session
	logger.Info("push uid ",uid)
	// recommend forum
	isExists := CheckTypeIsExists(uid, 1, session)
	if isExists==0 {
		FidsData :=  mysql.GetFids(uid, db)
		RecommendForum :=  mysql.GetClubsInfo(FidsData, db)
		ForumJson ,_ := json.Marshal(RecommendForum)
		forumStr := string(ForumJson)
		saveData(uid, 1, forumStr, db, session)
	}
	
	// recommend user
	Uids := mysql.GetUids(uid, db)
	RecommendUser := mysql.GetUserInfoByUids(uid, Uids, db)
	UserJson,_ := json.Marshal(RecommendUser)
	userStr := string(UserJson)
	err := saveData(uid, 2, userStr, db, session)

	// recommend goods
	Pet := mysql.GetPetInfoByUid(uid, db)
	var species int
	for _, v := range Pet {
		species = v.DogSpecies
	}
	speciesName := mysql.GetSpeciesnameBySpeciesid(species, db)
	RecommendGoods := mysql.GetGoods(speciesName)
	GoodsJson ,_ := json.Marshal(RecommendGoods)
	goodsStr := string(GoodsJson)
	err = saveData(uid, 3, goodsStr, db, session)

	// recommend ad
	RecomendAD := mysql.GetAdInfo(db)
	AdJson ,_ := json.Marshal(RecomendAD)
	adStr := string(AdJson)
	err = saveData(uid, 4, adStr, db, session)

	return err
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

func CheckTypeIsExists(uid int, typeId int, session *mgo.Session) int {
	var rp []*RecommendPosition
	c := session.DB("RecommendData").C("recommend_position")
	err := c.Find(bson.M{"uid": uid, "type": typeId}).All(&rp)

	if err != nil {
		logger.Error("mongodb find data", err, c)
		return 1
	}
	if len(rp)==0 {
		logger.Info("uid ",uid, " type ",typeId, "is not exists")
		return 0
	}
	logger.Info("uid ",uid, " type ",typeId, "is exists")
	return 1
}
