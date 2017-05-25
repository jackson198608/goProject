package main

import (
	// "encoding/json"
	// "bufio"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "io"
	"database/sql"
	"os"
	// "reflect"
	"strconv"
)

type EventLogX struct {
	// Oid    bson.ObjectId "_id"
	Id     int "_id"
	TypeId int "type"
	Uid    int "uid"
	Fuid   int "fuid" //fans id
	// Info    string        "info"
	Created string "created"
	Infoid  int    "infoid"
	Status  int    "status"
	Tid     int    "tid"
}

type EventLogNew struct {
	// Oid bson.ObjectId "_id"
	Id int "_id"
	// Id     int           "id"
	TypeId int "type"
	Uid    int "uid"
	// Info    string        "info"
	Created string "created"
	Infoid  int    "infoid"
	Status  int    "status"
	Tid     int    "tid"
}

func SaveMongoEventLog(event *EventLog, fans []*Follow, session *mgo.Session) {
	// session, err := mgo.Dial(c.mongoConn)
	// if err != nil {
	// 	logger.Info("mongodb connect fail", err)
	// 	return
	// }
	// defer session.Close()

	if event.uid > 0 {
		tableName := "event_log" //动态表
		x := session.DB(c.mongoDb).C(tableName)
		Id := createAutoIncrementId(session, "")
		// Id := 0
		// m1 := EventLogNew{bson.NewObjectId(), Id, event.typeId, event.uid, event.created, event.infoid, event.status, event.tid}
		m1 := EventLogNew{Id, event.typeId, event.uid, event.created, event.infoid, event.status, event.tid}
		//判断数据是否存在
		// eventIsExist := checkEventLogIsExist(x, event)
		// if eventIsExist == false {
		// logger.Info(x)
		logger.Info(m1)
		err := x.Insert(&m1) //插入数据
		if err != nil {
			logger.Info("mongo insert one data error:", err)
		}
	}
	// }
	/*fansLimit, _ := strconv.Atoi(c.fansLimit)
	if fansLimit > 0 && len(fans) > fansLimit {
		if event.created > c.dateLimit {
			saveFansEventLog(fans, session, event)
		}
	} else if fansLimit == 0 && c.dateLimit > "0" {
		if event.created > c.dateLimit {
			saveFansEventLog(fans, session, event)
		}
	} else {
		saveFansEventLog(fans, session, event)
	}*/

}

func saveFansEventLog(fans []*Follow, session *mgo.Session, event *EventLog) {
	for _, ar := range fans {
		tableNum1 := ar.follow_id % 100
		if tableNum1 == 0 {
			tableNum1 = 100
		}
		tableName1 := "event_log_" + strconv.Itoa(tableNum1) //粉丝表
		x := session.DB("EventLog").C(tableName1)
		eventIsExist := checkFansDataIsExist(x, event, ar.follow_id)
		if eventIsExist == false {
			IdX := 0
			// IdX := createAutoIncrementId(session, strconv.Itoa(tableNum1))
			// m := EventLogX{bson.NewObjectId(), IdX, event.typeId, event.uid, ar.follow_id, event.created, event.infoid, event.status, event.tid}
			m := EventLogX{IdX, event.typeId, event.uid, ar.follow_id, event.created, event.infoid, event.status, event.tid}
			logger.Info(m)
			// err := x.Insert(&m) //插入数据
			// if err != nil {
			// 	logger.Info("mongodb insert fans data", err, x)
			// }
		}
	}
}

func PushFansEventLog(event *EventLogNew, fans []*Follow, session *mgo.Session) error {
	for _, ar := range fans {
		tableNumX := ar.follow_id % 100
		if tableNumX == 0 {
			tableNumX = 100
		}
		tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
		c := session.DB("EventLog").C(tableNameX)
		eventIsExist := checkMongoFansDataIsExist(c, event, ar.follow_id)
		if eventIsExist == false && event.Status == 1 {
			IdX := createAutoIncrementId(session, strconv.Itoa(tableNumX))
			// m := EventLogX{bson.NewObjectId(), IdX, event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, event.Status, event.Tid}
			m := EventLogX{IdX, event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, event.Status, event.Tid}
			err := c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
				return err
			}
		}
	}
	return nil
}

func UpdateMongoEventLogStatus(event *EventLogNew, fans []*Follow, status string, session *mgo.Session) {
	//-1隐藏,0:删除,1显示,2动态推送给粉丝,3取消关注
	if status == "-1" {
		HideOrShowEventLog(event, fans, session, -1)
	}
	if status == "0" {
		HideOrShowEventLog(event, fans, session, 0)
	}
	if status == "1" {
		HideOrShowEventLog(event, fans, session, 1)
	}
	if status == "2" {
		PushFansEventLog(event, fans, session)
	}
}

func RemoveFansEventLog(fuid string, uid string, session *mgo.Session) error {
	uidN, _ := strconv.Atoi(uid)
	fuidN, _ := strconv.Atoi(fuid)
	tableNumX := fuidN % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	c := session.DB("EventLog").C(tableNameX)
	_, err := c.RemoveAll(bson.M{"uid": uidN, "fuid": fuidN}) //取消关注删除数据
	if err != nil {
		logger.Info("mongodb insert fans data", err, c)
		return err
	}

	return nil
}

func HideOrShowEventLog(event *EventLogNew, fans []*Follow, session *mgo.Session, status int) error {
	tableName := "event_log" //动态表
	c := session.DB(c.mongoDb).C(tableName)
	//判断数据是否存在
	eventIsExist := checkMongoEventLogIsExist(c, event)
	if eventIsExist == true {
		err := c.Update(bson.M{"_id": event.Id}, bson.M{"$set": bson.M{"status": status}}) //插入数据
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
				err := c.Update(bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": ar.follow_id, "created": event.Created, "infoid": event.Infoid}, bson.M{"$set": bson.M{"status": status}}) //插入数据
				if err != nil {
					logger.Info("mongodb insert fans data", err, c)
					return err
				}
			}
			if status == 0 {
				err := c.Remove(bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": ar.follow_id, "created": event.Created, "infoid": event.Infoid, "tid": event.Tid}) //插入数据
				if err != nil {
					logger.Info("mongodb insert fans data", err, c)
					return err
				}
			}
		}
		if eventIsExist == false && status == 1 {
			IdX := createAutoIncrementId(session, strconv.Itoa(tableNumX))
			// m := EventLogX{bson.NewObjectId(), IdX, event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, status, event.Tid}
			m := EventLogX{IdX, event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, status, event.Tid}
			err := c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
				return err
			}
		}
	}
	return nil
}

func createAutoIncrementId(session *mgo.Session, tableNum string) int {
	// session, err := mgo.Dial(c.mongoConn)
	// if err != nil {
	// 	logger.Info("mongodb connect fail", err)
	// 	// return nil
	// }
	// defer session.Close()
	c := session.DB(c.mongoDb).C("ids" + tableNum)
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"id": 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	doc := struct{ Id int }{}
	_, err := c.Find(bson.M{"_id": 0}).Apply(change, &doc)
	if err != nil {
		logger.Info("get counter failed:", err)
	}
	return doc.Id
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

func LoadMongoById(Id string, session *mgo.Session) *EventLogNew {
	// objectId := bson.ObjectIdHex(Id)
	objectId, _ := strconv.Atoi(Id)
	event := new(EventLogNew)
	c := session.DB(c.mongoDb).C("event_log")
	// query := func(c *mgo.Collection) error {
	c.FindId(objectId).One(&event)
	// return c.Find(bson.M{"_id": objectId}).One(&event)
	// }
	// witchCollection("event_log", query)
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

func pushEventToFansTask(fans string, uid string, count string, session *mgo.Session, db *sql.DB) {
	page := 0
	for {
		countNum, _ := strconv.Atoi(count)
		fansNum, _ := strconv.Atoi(fans)
		if fansNum >= c.fansLimit {
			offset := page * c.numloops
			if offset >= c.eventLimit {
				break
			}
		}
		if fansNum < c.fansLimit && countNum > c.pushLimit {
			offset := page * c.numloops
			if offset >= c.pushLimit {
				break
			}
		}
		user_id, _ := strconv.Atoi(uid)
		datas := getMysqlData(fansNum, user_id, countNum, page, db)
		if len(datas) == 0 {
			break
		}
		if datas == nil {
			break
		}
		fansData := GetFansData(user_id, db)
		for _, event := range datas {
			// fmt.Println(event.created)
			// logger.Info("event", ar)
			// fmt.Println(reflect.TypeOf(ar))
			saveFansEventLog(fansData, session, event)
		}
		page++
	}
}
