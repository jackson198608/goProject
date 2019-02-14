package content

import (
	// "errors"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gouminGitlab/common/orm/mongo/RecommendData"
	"strconv"
	// "strings"
	"math/rand"
	"time"
	"github.com/olivere/elastic"
)

type Content struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	Uid       int
	esConn    *elastic.Client
}

func NewContent(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, uid string, esConn *elastic.Client) *Content {
	if (mysqlXorm == nil) || (mongoConn == nil) || (uid == "") {
		return nil
	}
	logger.Info("start recommend: ", uid)
	c := new(Content)
	if c == nil {
		return nil
	}
	c.mysqlXorm = mysqlXorm
	c.mongoConn = mongoConn
	c.Uid, _ = strconv.Atoi(uid)
	c.esConn = esConn
	return c
}

func (c *Content) Do() error {
	contents, err := c.getContents()
	if err != nil {
		return err
	}
	tableNameX := getTableNum(c.Uid)
	mc := c.mongoConn[0].DB("RecommendData").C(tableNameX)
	for i, _ := range contents {
		content := contents[i]
		// logger.Info("recommend data push uid is ", c.Uid, ", type is ", c.jsonData.Type, ", infoid is ", c.jsonData.Infoid)
		fmt.Println(*content)
		err := c.insertPerson(mc, c.Uid, *content)
		if err != nil {
			return err
		}
	}
	return nil
}

//生成随机数
func (c *Content) getRand() int {
	seed := time.Now().Unix() + int64(c.Uid)
	rand.Seed(seed)
	rangeNum := rand.Intn(10)
	if rangeNum == 0 {
		rangeNum = 1
	}
	max := rangeNum * 10000
	randid := rand.Intn(max)

	return randid
}

func (c *Content) getContents() ([]*RecommendData.ContentChannel, error) {
	var data []*RecommendData.ContentChannel
	mc := c.mongoConn[0].DB("RecommendData").C("content_channel")
	randid := c.getRand()
	fmt.Println(randid)
	err := mc.Find(&bson.M{"stick": 0, "channel": 1, "randid": bson.M{"$gt": randid}}).Sort("created").Limit(5).All(&data)
	if len(data) == 0 {
		err = mc.Find(&bson.M{"stick": 0, "channel": 1, "randid": bson.M{"$lt": randid}}).Sort("created").Limit(5).All(&data)
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getTableNum(person int) string {
	tableNumX := person % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "user_recommend_" + strconv.Itoa(tableNumX) //用户推荐表
	return tableNameX
}

func (c *Content) insertPerson(mc *mgo.Collection, person int, recmmendData RecommendData.ContentChannel) error {
	//新增数据
	created := int(time.Now().Unix())

	var data RecommendData.UserRecommendX
	data = RecommendData.UserRecommendX{bson.NewObjectId(),
		recmmendData.Pid,
		recmmendData.Uid,
		person,
		recmmendData.Type,
		recmmendData.Infoid,
		recmmendData.Title,
		recmmendData.Description,
		created,
		recmmendData.Images,
		recmmendData.ImageNum,
		recmmendData.Tag,
		recmmendData.Tags,
		recmmendData.QstType,
		recmmendData.AdType,
		recmmendData.AdUrl,
		recmmendData.Channel,
		recmmendData.VideoUrl,
		recmmendData.Duration}
	err := mc.Insert(&data) //插入数据
	if err != nil {
		return err
	}
	return nil
}
