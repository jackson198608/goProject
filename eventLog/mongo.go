package main

import (
	// "encoding/json"
	// "bufio"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "io"
	"os"
	// "reflect"
	"strconv"
)

type EventLogX struct {
	Id     bson.ObjectId "_id"
	TypeId int           "type"
	Uid    int           "uid"
	Fuid   int           "fuid" //fans id
	// Info    string        "info"
	Created string "created"
	Infoid  int    "infoid"
	Status  int    "status"
	Tid     int    "tid"
}

type EventLogNew struct {
	Id     bson.ObjectId "_id"
	TypeId int           "type"
	Uid    int           "uid"
	// Info    string        "info"
	Created string "created"
	Infoid  int    "infoid"
	Status  int    "status"
	Tid     int    "tid"
}

func SaveMongoEventLog(event *EventLog, fans []*Follow, w *os.File) {
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		logger.Info("mongodb connect fail", err)
		return
	}
	defer session.Close()

	tableName := "event_log" //动态表
	c := session.DB(c.mongoDb).C(tableName)
	m1 := EventLogNew{bson.NewObjectId(), event.typeId, event.uid, event.created, event.infoid, event.status, event.tid}
	logger.Info(m1)
	//判断数据是否存在
	eventIsExist := checkEventLogIsExist(c, event)
	if eventIsExist == false {
		err = c.Insert(&m1) //插入数据
		if err != nil {
			logger.Info("mongo insert one data error:", err)
		}
	}
	// fmt.Println("type:", reflect.TypeOf(c))
	lineStr1 := fmt.Sprintf("%s", m1)
	// 查找文件末尾的偏移量
	n, _ := w.Seek(0, os.SEEK_END)
	// 从末尾的偏移量开始写入内容
	_, err = w.WriteAt([]byte(lineStr1+"\n"), n)
	for _, ar := range fans {
		tableNum1 := ar.follow_id % 100
		if tableNum1 == 0 {
			tableNum1 = 100
		}
		tableName1 := "event_log_" + strconv.Itoa(tableNum1) //粉丝表
		c := session.DB("EventLog").C(tableName1)
		m := EventLogX{bson.NewObjectId(), event.typeId, event.uid, ar.follow_id, event.created, event.infoid, event.status, event.tid}
		eventIsExist := checkFansDataIsExist(c, event, ar.follow_id)
		if eventIsExist == false {
			err = c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
			}
		}
		n1, _ := w.Seek(0, os.SEEK_END)
		lineStr := fmt.Sprintf("%s", m)
		_, err = w.WriteAt([]byte(lineStr+"\n"), n1)
		if err != nil {
			logger.Info("mongodb write data", err, c)
		}
	}

}

func PushFansEventLog(event *EventLogNew, fans []*Follow) error {
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		logger.Info("mongodb connect fail", err)
		return err
	}
	defer session.Close()
	for _, ar := range fans {
		tableNumX := ar.follow_id % 100
		if tableNumX == 0 {
			tableNumX = 100
		}
		tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
		c := session.DB("EventLog").C(tableNameX)
		m := EventLogX{bson.NewObjectId(), event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, event.Status, event.Tid}
		eventIsExist := checkMongoFansDataIsExist(c, event, ar.follow_id)
		if eventIsExist == false && event.Status == 1 {
			err = c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
				return err
			}
		}
	}
	return nil
}

func UpdateMongoEventLogStatus(event *EventLogNew, fans []*Follow, status string) {
	//-1隐藏,0:删除,1显示,2动态推送给粉丝
	if status == "-1" {
		HideOrShowEventLog(event, fans, -1)
	}
	if status == "0" {
		HideOrShowEventLog(event, fans, 0)
	}
	if status == "1" {
		HideOrShowEventLog(event, fans, 1)
	}
	if status == "2" {
		PushFansEventLog(event, fans)
	}
}

func HideOrShowEventLog(event *EventLogNew, fans []*Follow, status int) error {
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		logger.Info("mongodb connect fail", err)
		return err
	}
	defer session.Close()
	tableName := "event_log" //动态表
	c := session.DB(c.mongoDb).C(tableName)
	//判断数据是否存在
	eventIsExist := checkMongoEventLogIsExist(c, event)
	if eventIsExist == true {
		err = c.Update(bson.M{"_id": event.Id}, bson.M{"$set": bson.M{"status": status}}) //插入数据
		if err != nil {
			logger.Info("mongo insert one data error:", err)
			return err
		}
	}
	for _, ar := range fans {
		tableNumX := ar.follow_id % 100
		if tableNumX == 0 {
			tableNumX = 100
		}
		tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
		c := session.DB("EventLog").C(tableNameX)
		eventIsExist := checkMongoFansDataIsExist(c, event, ar.follow_id)
		if eventIsExist == true {
			if status == -1 || status == 1 {
				err = c.Update(bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": ar.follow_id, "created": event.Created, "infoid": event.Infoid}, bson.M{"$set": bson.M{"status": status}}) //插入数据
				if err != nil {
					logger.Info("mongodb insert fans data", err, c)
					return err
				}
			}
			if status == 0 {
				err = c.Remove(bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": ar.follow_id, "created": event.Created, "infoid": event.Infoid, "tid": event.Tid}) //插入数据
				if err != nil {
					logger.Info("mongodb insert fans data", err, c)
					return err
				}
			}
		}
		if eventIsExist == false && status == 1 {
			m := EventLogX{bson.NewObjectId(), event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, status, event.Tid}
			err = c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
				return err
			}
		}
	}
	return nil
}

//根据mysql中的数据检查mongo中是否存在该条数据
func checkEventLogIsExist(c *mgo.Collection, event *EventLog) bool {
	ms := []EventLogNew{}
	err1 := c.Find(&bson.M{"type": event.typeId, "uid": event.uid, "created": event.created, "infoid": event.infoid, "status": event.status}).All(&ms)
	if err1 != nil {
		logger.Info("mongodb find data", err1, c)
		return false
	}
	if len(ms) == 0 {
		return false
	}
	return true
}

//检查mongo中是否存在该条数据
func checkMongoEventLogIsExist(c *mgo.Collection, event *EventLogNew) bool {
	ms := []EventLogNew{}
	err1 := c.Find(&bson.M{"type": event.TypeId, "uid": event.Uid, "created": event.Created, "infoid": event.Infoid, "status": event.Status}).All(&ms)
	if err1 != nil {
		logger.Info("mongodb find data", err1, c)
		return false
	}
	if len(ms) == 0 {
		return false
	}
	return true
}

//检查mongo中是否存在该条fans数据
func checkMongoFansDataIsExist(c *mgo.Collection, event *EventLogNew, fuid int) bool {
	ms := []EventLogX{}
	err1 := c.Find(&bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": fuid, "created": event.Created, "infoid": event.Infoid, "tid": event.Tid}).All(&ms)

	if err1 != nil {
		logger.Info("mongodb find data", err1, c)
		return false
	}
	if len(ms) == 0 {
		return false
	}
	return true
}

//根据mysql中的数据检查mongo中是否存在该条fans数据
func checkFansDataIsExist(c *mgo.Collection, event *EventLog, fuid int) bool {
	ms := []EventLogX{}
	err1 := c.Find(&bson.M{"type": event.typeId, "uid": event.uid, "fuid": fuid, "created": event.created, "infoid": event.infoid, "status": event.status}).All(&ms)

	if err1 != nil {
		logger.Info("mongodb find data", err1, c)
		return false
	}
	if len(ms) == 0 {
		return false
	}
	return true
}

func LoadMongoById(Id string) *EventLogNew {
	objectId := bson.ObjectIdHex(Id)
	event := new(EventLogNew)
	query := func(c *mgo.Collection) error {
		return c.FindId(objectId).One(&event)
		// return c.Find(bson.M{"_id": objectId}).One(&event)
	}
	witchCollection("event_log", query)
	return event
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		logger.Info("mongodb connect fail", err)
		return err
	}
	defer session.Close()
	c := session.DB(c.mongoDb).C(collection)
	return s(c)
}

func check(e error) {
	if e != nil {
		logger.Info("check file error", e)
	}
}
