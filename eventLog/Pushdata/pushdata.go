package Pushdata

import (
	// "encoding/json"
	// "bufio"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/eventLog/mysql"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "io"
	"database/sql"
	// "os"
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

// t.oid, t, id, t.fuid, t.uid, t.status, t.event
type EventLogNew struct {
	db      *sql.DB
	session *mgo.Session
	oid     int
	// event   *EventLog
}

type EventLogLast struct {
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

type EventLog struct {
	id        int64
	typeId    int
	uid       int
	info      string
	created   string
	infoid    int
	status    int
	tid       int
	isSplit   bool
	logLevel  int
	postTable string
}

func NewEventLogNew(logLevel int, oid int, id int, db *sql.DB, session *mgo.Session) *EventLogNew {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(EventLogNew)
	e.db = db
	e.session = session
	// e.event = event
	e.oid = oid
	return e
}

func (e *EventLogNew) SaveMongoEventLog(oid int) error {
	event := mysql.LoadById(oid, e.db)
	// 判断数据是否存在
	tableName := "event_log" //动态表
	x := e.session.DB("EventLog").C(tableName)
	eventIsExist := checkEventLogIsExist(x, event)
	if event.Uid > 0 && eventIsExist == false {
		Id := createAutoIncrementId(e.session, "")
		// Id := 0
		// m1 := EventLogNew{bson.NewObjectId(), Id, event.typeId, event.uid, event.created, event.infoid, event.status, event.tid}
		m1 := EventLogLast{Id, event.TypeId, event.Uid, event.Created, event.Infoid, event.Status, event.Tid}
		err := x.Insert(&m1) //插入数据
		if err != nil {
			logger.Info("mongo insert one data error:", err)
		}
		logger.Info("mysql to mongo data ", m1)
	}
	return nil
}

//根据mysql中的数据检查mongo中是否存在该条数据
func checkEventLogIsExist(c *mgo.Collection, event *mysql.EventLog) bool {
	ms := []EventLogNew{}
	err1 := c.Find(&bson.M{"type": event.TypeId, "uid": event.Uid, "created": event.Created, "infoid": event.Infoid, "status": event.Status}).All(&ms)
	// logger.Info("check event_log data is exist")
	if err1 != nil {
		logger.Info("mongodb find data", err1, c)
		return false
	}
	if len(ms) == 0 {
		return false
	}
	return true
}

func LoadMongoById(objectId int, session *mgo.Session) *EventLogLast {
	event := new(EventLogLast)
	c := session.DB("EventLog").C("event_log")
	c.FindId(objectId).One(&event)
	return event
}

//更改推送给粉丝的动态数据的状态
func (e *EventLogNew) UpdateMongoEventLogStatus(id int, status string) error {
	event := LoadMongoById(id, e.session)
	fans := mysql.GetFansData(event.Uid, e.db)
	//-1隐藏,0:删除,1显示,2动态推送给粉丝,3取消关注
	if status == "-1" {
		e.HideOrShowEventLog(event, fans, -1)
	}
	if status == "0" {
		e.HideOrShowEventLog(event, fans, 0)
	}
	if status == "1" {
		e.HideOrShowEventLog(event, fans, 1)
	}
	if status == "2" {
		e.PushFansEventLog(event, fans)
	}
	return nil
}

func (e *EventLogNew) HideOrShowEventLog(event *EventLogLast, fans []int, status int) error {
	tableName := "event_log" //动态表
	session := e.session
	c := session.DB("EventLog").C(tableName)
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
		tableNumX := ar % 100
		if tableNumX == 0 {
			tableNumX = 100
		}
		tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
		c := session.DB("EventLog").C(tableNameX)
		eventIsExist := checkMongoFansDataIsExist(c, event, ar)
		if eventIsExist == true {
			if status == -1 || status == 1 {
				err := c.Update(bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": ar, "created": event.Created, "infoid": event.Infoid}, bson.M{"$set": bson.M{"status": status}}) //插入数据
				if err != nil {
					logger.Info("mongodb update fans data error ", err, c)
					return err
				}
			}
			if status == 0 {
				err := c.Remove(bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": ar, "created": event.Created, "infoid": event.Infoid, "tid": event.Tid}) //插入数据
				if err != nil {
					logger.Info("mongodb remove fans data error ", err, c)
					return err
				}
			}
		}
		if eventIsExist == false && status == 1 {
			IdX := createAutoIncrementId(session, strconv.Itoa(tableNumX))
			// m := EventLogX{bson.NewObjectId(), IdX, event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, status, event.Tid}
			m := EventLogX{IdX, event.TypeId, event.Uid, ar, event.Created, event.Infoid, status, event.Tid}
			err := c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data error ", err, c)
				return err
			}
			logger.Info("mongodb insert fans data:", m)
		}
	}
	return nil
}

func (e *EventLogNew) RemoveFansEventLog(fuid int, uid int) error {
	uidN := uid
	fuidN := fuid
	tableNumX := fuidN % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	session := e.session
	c := session.DB("EventLog").C(tableNameX)
	_, err := c.RemoveAll(bson.M{"uid": uidN, "fuid": fuidN}) //取消关注删除数据
	if err != nil {
		logger.Info("mongodb insert fans data", err, c)
		return err
	}
	logger.Info("remove event is success")
	return nil
}

func (e *EventLogNew) PushFansEventLog(event *EventLogLast, fans []int) error {
	session := e.session
	for _, ar := range fans {
		tableNumX := ar % 100
		if tableNumX == 0 {
			tableNumX = 100
		}
		tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
		c := session.DB("EventLog").C(tableNameX)
		eventIsExist := checkMongoFansDataIsExist(c, event, ar)
		if eventIsExist == false && event.Status == 1 {
			IdX := createAutoIncrementId(session, strconv.Itoa(tableNumX))
			// m := EventLogX{bson.NewObjectId(), IdX, event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, event.Status, event.Tid}
			m := EventLogX{IdX, event.TypeId, event.Uid, ar, event.Created, event.Infoid, event.Status, event.Tid}
			err := c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
				return err
			}
			logger.Info("mongodb push fans data ", m)
		}
	}
	return nil
}

//检查mongo中是否存在该条数据
func checkMongoEventLogIsExist(c *mgo.Collection, event *EventLogLast) bool {
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

func createAutoIncrementId(session *mgo.Session, tableNum string) int {
	c := session.DB("EventLog").C("ids" + tableNum)
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

//检查mongo中是否存在该条fans数据
func checkMongoFansDataIsExist(c *mgo.Collection, event *EventLogLast, fuid int) bool {
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

func (e *EventLogNew) PushEventToFansTask(fans string, user_id int, count string, numLoop int, fansLimit string, eventLimit string, pushLimit string, dateLimit string) {
	page := 0
	for {
		countNum, _ := strconv.Atoi(count)
		fansNum, _ := strconv.Atoi(fans)
		fansLimitNum, _ := strconv.Atoi(fansLimit)
		eventLimitNum, _ := strconv.Atoi(eventLimit)
		pushLimit, _ := strconv.Atoi(pushLimit)
		if fansNum >= fansLimitNum {
			offset := page * numLoop
			if offset >= eventLimitNum {
				break
			}
		}
		if fansNum < fansLimitNum && countNum > pushLimit {
			offset := page * numLoop
			if offset >= pushLimit {
				break
			}
		}
		// user_id, _ := strconv.Atoi(uid)
		datas := mysql.GetMysqlData(fansNum, user_id, countNum, page, e.db, numLoop, fansLimitNum, eventLimitNum, pushLimit, dateLimit)
		if len(datas) == 0 {
			break
		}
		if datas == nil {
			break
		}
		fansData := mysql.GetFansData(user_id, e.db)
		for _, event := range datas {
			// fmt.Println(reflect.TypeOf(ar))
			e.saveFansEventLog(fansData, event)
		}
		page++
	}
}

func (e *EventLogNew) saveFansEventLog(fans []int, event *mysql.EventLog) {
	session := e.session
	for _, ar := range fans {
		tableNum1 := ar % 100
		if tableNum1 == 0 {
			tableNum1 = 100
		}
		tableName1 := "event_log_" + strconv.Itoa(tableNum1) //粉丝表
		x := session.DB("EventLog").C(tableName1)
		eventIsExist := checkFansDataIsExist(x, event, ar)
		if eventIsExist == false {
			// IdX := 0
			IdX := createAutoIncrementId(session, strconv.Itoa(tableNum1))
			// m := EventLogX{bson.NewObjectId(), IdX, event.typeId, event.uid, ar.follow_id, event.created, event.infoid, event.status, event.tid}
			m := EventLogX{IdX, event.TypeId, event.Uid, ar, event.Created, event.Infoid, event.Status, event.Tid}
			err := x.Insert(&m) //插入数据
			logger.Info(m)
			if err != nil {
				logger.Info("mongodb insert fans data", err, x)
			}
		}
	}
}

//根据mysql中的数据检查mongo中是否存在该条fans数据
func checkFansDataIsExist(c *mgo.Collection, event *mysql.EventLog, fuid int) bool {
	ms := []EventLogX{}
	err1 := c.Find(&bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": fuid, "created": event.Created, "infoid": event.Infoid, "status": event.Status}).All(&ms)

	if err1 != nil {
		logger.Info("mongodb find data", err1, c)
		return false
	}
	if len(ms) == 0 {
		return false
	}
	return true
}
