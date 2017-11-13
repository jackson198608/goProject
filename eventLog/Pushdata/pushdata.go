package Pushdata

import (
	// "encoding/json"
	// "bufio"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/eventLog/mysql"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "io"
	"database/sql"
	// "os"
	// "reflect"
	"strconv"
	"time"
)

type EventLogX struct {
	// Id         int    "_id"
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

// t.oid, t, id, t.fuid, t.uid, t.status, t.event
type EventLogNew struct {
	db      *sql.DB
	session *mgo.Session
	oid     int
	// slave   *mgo.Session
	// event   *EventLog
}

type EventLogLast struct {
	// Oid bson.ObjectId "_id"
	Id        int    "_id"
	TypeId    int    "type"
	Uid       int    "uid"
	Created   string "created"
	Infoid    int    "infoid"
	Status    int    "status"
	Tid       int    "tid"
	Bid       int    "bid"
	Content   string "content"
	Title     string "title"
	Imagenums int    "image_num"
	Images    string "images"
	Forum     string "forum"
	Tag       int    "tag"
	Qsttype   int    "qst_type"
	IsRead    int    "is_read"
	Source    int    "source"
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

type BreedActiveUser struct {
	Uid     int "uid"
	Breedid int "breed_id"
}

type ForumActiveUser struct {
	Uid     int "uid"
	Forumid int "forum_id"
}

func NewEventLogNew(logLevel int, oid int, id int, db *sql.DB, session *mgo.Session) *EventLogNew {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(EventLogNew)
	e.db = db
	e.session = session //主库
	// e.slave = slave     //从库
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
		m1 := EventLogLast{Id, event.TypeId, event.Uid, event.Created, event.Infoid, event.Status, event.Tid, 0, "", "", 0, "", "", 0, 0, 0, 0}
		err := x.Insert(&m1) //插入数据
		if err != nil {
			logger.Info("mongo insert one data error:", err)
		}
		logger.Info("mysql to master mongo data ", m1)
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

func LoadMongoById(objectId int, slave *mgo.Session) *EventLogLast {
	event := new(EventLogLast)
	c := slave.DB("EventLog").C("event_log")
	c.FindId(objectId).One(&event)
	return event
}

//更改推送给粉丝的动态数据的状态
func (e *EventLogNew) UpdateMongoEventLogStatus(id int, status string) error {
	//从库读取数据 e.slave
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

func (e *EventLogNew) HideOrShowEventLog(event *EventLogLast, fans []*mysql.Follow, status int) error {
	tableName := "event_log" //动态表
	session := e.session     //主库存储
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
	var allusers []int
	if event.TypeId == 1 { //1:帖子
		//俱乐部所有活跃用户 + 活跃粉丝用户
		allusers = MergeFansAndForumUsers(fans, event.Infoid, e.session, e.db)
	} else if event.TypeId == 8 { //8:问答
		//获取相同犬种的活跃用户
		allusers = GetBreedActiveUser(event.Bid, e.session)
	} else if (event.TypeId == 9 || event.TypeId == 15) && event.Source == 1 { //人工推荐
		//全部活跃用户
		allusers = GetAllActiveUsers(e.session)
	} else {
		for _, v := range fans {
			// allusers[k] = v.Follow_id
			allusers = append(allusers, v.Follow_id)
		}
	}
	for _, ar := range allusers {
		tableNumX := ar % 100
		if tableNumX == 0 {
			tableNumX = 100
		}
		tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
		c := session.DB("FansData").C(tableNameX)
		eventIsExist := checkMongoFansDataIsExist(c, event, ar)
		if eventIsExist == true {
			if status == -1 || status == 1 {
				err := c.Update(bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": ar, "created": event.Created, "infoid": event.Infoid}, bson.M{"$set": bson.M{"status": status}}) //插入数据
				if err != nil {
					logger.Info("mongodb update fans data error ", err, c)
					return err
				}
				logger.Info("mongodb update fans data:", event.TypeId, event.Uid, ar, event.Infoid)
			}
			if status == 0 {
				err := c.Remove(bson.M{"type": event.TypeId, "uid": event.Uid, "fuid": ar, "created": event.Created, "infoid": event.Infoid, "tid": event.Tid}) //插入数据
				if err != nil {
					logger.Info("mongodb remove fans data error ", err, c)
					return err
				}
				logger.Info("mongodb remove fans data:", event.TypeId, event.Uid, ar, event.Infoid)
			}
		}
		if eventIsExist == false && status == 1 {
			// IdX := createFansAutoIncrementId(session, strconv.Itoa(tableNumX))
			// m := EventLogX{bson.NewObjectId(), IdX, event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, status, event.Tid}
			m := EventLogX{bson.NewObjectId(), event.TypeId, event.Uid, ar, event.Created, event.Infoid, status, event.Tid, event.Bid, event.Content, event.Title, event.Imagenums, event.Images, event.Forum, event.Tag, event.Qsttype, event.IsRead, event.Source}
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
	c := session.DB("FansData").C(tableNameX)
	_, err := c.RemoveAll(bson.M{"uid": uidN, "fuid": fuidN}) //取消关注删除数据
	if err != nil {
		logger.Info("mongodb insert fans data", err, c)
		return err
	}
	logger.Info("remove event is success")
	return nil
}

func (e *EventLogNew) PushFansEventLogOld(event *EventLogLast, fans []*mysql.Follow) error {
	// slave := e.slave     //从库查询
	session := e.session //主库存储
	for _, ar := range fans {
		tableNumX := ar.Follow_id % 100
		if tableNumX == 0 {
			tableNumX = 100
		}
		tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
		c := session.DB("FansData").C(tableNameX)
		// eventIsExist := checkMongoFansDataIsExist(c, event, ar.Follow_id)
		// if eventIsExist == false && event.Status == 1 {
		if event.Status == 1 {
			// IdX := createFansAutoIncrementId(session, strconv.Itoa(tableNumX))
			// m := EventLogX{bson.NewObjectId(), IdX, event.TypeId, event.Uid, ar.follow_id, event.Created, event.Infoid, event.Status, event.Tid}
			m := EventLogX{bson.NewObjectId(), event.TypeId, event.Uid, ar.Follow_id, event.Created, event.Infoid, event.Status, event.Tid, 0, "", "", 0, "", "", 0, 0, 0, 0}
			err := c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
				return err
			}
			logger.Info("slave FansData mongodb push fans data ", m)
		}
	}
	return nil
}

func (e *EventLogNew) PushFansEventLog(event *EventLogLast, fans []*mysql.Follow) error {
	// slave := e.slave     //从库查询
	session := e.session //主库存储

	var allusers []int
	if event.TypeId == 1 { //1:帖子
		//俱乐部所有活跃用户 + 活跃粉丝用户
		allusers = MergeFansAndForumUsers(fans, event.Infoid, e.session, e.db)
	} else if event.TypeId == 8 { //8:问答
		//获取相同犬种的活跃用户
		allusers = GetBreedActiveUser(event.Bid, e.session)
	} else if (event.TypeId == 9 || event.TypeId == 15) && event.Source == 1 { //小编推荐
		//全部活跃用户
		allusers = GetAllActiveUsers(e.session)
	} else {
		for _, v := range fans {
			// allusers[k] = v.Follow_id
			allusers = append(allusers, v.Follow_id)
		}
	}
	if len(allusers) == 0 {
		logger.Info("allusers arr is empty")
		return nil
	}

	for _, ar := range allusers {
		tableNumX := ar % 100
		if tableNumX == 0 {
			tableNumX = 100
		}
		tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
		c := session.DB("FansData").C(tableNameX)
		if event.Status == 1 {
			m := EventLogX{bson.NewObjectId(), event.TypeId, event.Uid, ar, event.Created, event.Infoid, event.Status, event.Tid, event.Bid, event.Content, event.Title, event.Imagenums, event.Images, event.Forum, event.Tag, event.Qsttype, event.IsRead, event.Source}
			err := c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
				return err
			}
			logger.Info("slave FansData mongodb push fans data ", m)
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

func createFansAutoIncrementId(session *mgo.Session, tableNum string) int {
	c := session.DB("FansData").C("ids" + tableNum)
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
		limit := 500
		countNum, _ := strconv.Atoi(count)
		fansNum, _ := strconv.Atoi(fans)
		fansLimitNum, _ := strconv.Atoi(fansLimit)
		eventLimitNum, _ := strconv.Atoi(eventLimit)
		pushLimit, _ := strconv.Atoi(pushLimit)
		offset := page * limit
		offlim := offset + limit
		if fansNum >= fansLimitNum {
			if limit > eventLimitNum {
				limit = eventLimitNum
			} else {
				if offlim > eventLimitNum && eventLimitNum > offset {
					limit = eventLimitNum - offset
				}
			}
			// fmt.Println(page, limit, offset)
			if offset >= eventLimitNum {
				break
			}
		}
		if fansNum < fansLimitNum && countNum > pushLimit {
			if limit > pushLimit {
				limit = pushLimit
			} else {
				if offlim > pushLimit && pushLimit > offset {
					limit = pushLimit - offset
				}
			}
			if offset >= pushLimit {
				break
			}
		}
		// user_id, _ := strconv.Atoi(uid)
		// datas := mysql.GetMysqlData(fansNum, user_id, countNum, page, e.db, numLoop, fansLimitNum, eventLimitNum, pushLimit, dateLimit)
		if fansNum <= 0 {
			break
		}
		datas := e.getMongoData(fansNum, user_id, countNum, page, limit, offset, fansLimitNum, pushLimit, dateLimit)
		if len(datas) == 0 {
			break
		}
		if datas == nil {
			break
		}
		fansData := mysql.GetFansData(user_id, e.db)
		// fmt.Println(len(fansData))
		for _, event := range datas {
			// fmt.Println(event)
			e.saveFansEventLog(fansData, event)
		}
		page++
	}
}

/**
 * fans 粉丝数
 * uid  用户uid
 * eventCount 动态数
 * page  当前分页
 * fansLimit 粉丝数限制
 * eventLimit 动态数限制
 * pushLimit 推送条数限制
 * dateLimit 动态时间限制
 */
func (e *EventLogNew) getMongoData(fans int, uid int, eventCount int, page int, limit int, offset int, fansLimit int, pushLimit int, dateLimit string) []*EventLogLast {
	tableName := "event_log" //动态表
	c := e.session.DB("EventLog").C(tableName)
	ms := []*EventLogLast{}
	typeIds := [7]int{1, 2, 3, 4, 5, 6, 7} //我关注人的推送信息
	if fans > fansLimit {
		err := c.Find(&bson.M{"uid": uid, "type": bson.M{"$in": typeIds}, "status": 1}).Sort("-created").Limit(limit).Skip(offset).All(&ms)
		if err != nil {
			logger.Info("mongodb find data", err, c)
			return nil
		}
		// } else if fans < fansLimit && eventCount > pushLimit {
		// 	err = c.Find(&bson.M{"uid": uid, "type": bson.M{"$in": typeIds}, "status": 1,"created":bson.M{"$gt": dateLimit}}).Sort("-created").Limit(limit).Skip(offset).All(&ms)
	} else {
		err := c.Find(&bson.M{"uid": uid, "type": bson.M{"$in": typeIds}, "status": 1, "created": bson.M{"$gt": dateLimit}}).Sort("-created").Limit(limit).Skip(offset).All(&ms)
		if err != nil {
			logger.Info("mongodb find data", err, c)
			return nil
		}
	}
	if len(ms) == 0 {
		return nil
	}
	return ms
}

func (e *EventLogNew) saveFansEventLog(fans []*mysql.Follow, event *EventLogLast) {
	session := e.session
	for _, ar := range fans {
		tableNum1 := ar.Follow_id % 100
		if tableNum1 == 0 {
			tableNum1 = 100
		}
		tableName1 := "event_log_" + strconv.Itoa(tableNum1) //粉丝表
		x := session.DB("FansData").C(tableName1)
		eventIsExist := checkFansDataIsExist(x, event, ar.Follow_id)
		if eventIsExist == false {
			// IdX := 0
			// IdX := createFansAutoIncrementId(session, strconv.Itoa(tableNum1))
			// m := EventLogX{bson.NewObjectId(), IdX, event.typeId, event.uid, ar.follow_id, event.created, event.infoid, event.status, event.tid}
			m := EventLogX{bson.NewObjectId(), event.TypeId, event.Uid, ar.Follow_id, event.Created, event.Infoid, event.Status, event.Tid, 0, "", "", 0, "", "", 0, 0, 0, 0}
			err := x.Insert(&m) //插入数据
			logger.Info(m)
			if err != nil {
				logger.Info("mongodb insert fans data", err, x)
			}
		}
	}
}

//根据mysql中的数据检查mongo中是否存在该条fans数据
func checkFansDataIsExist(c *mgo.Collection, event *EventLogLast, fuid int) bool {
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

func (e *EventLogNew) RemoveEventToFansTask(fans_uid int, numloop int, eventLimit string, sleeptime string) {
	session := e.session
	tableNumX := fans_uid % 100
	if tableNumX == 0 {
		tableNumX = 100
	}
	tableNameX := "event_log_" + strconv.Itoa(tableNumX) //粉丝表
	c := session.DB("FansData").C(tableNameX)
	//1：帖子，2：回复，3：关注，4：评论，5：日志萌图，6：视频，7：商城圈，8：问答
	typeIds := [7]int{1, 2, 3, 4, 5, 6, 7} //我关注人的推送信息
	count, _ := c.Find(&bson.M{"fuid": fans_uid, "type": bson.M{"$in": typeIds}}).Count()
	eventLimitNum, _ := strconv.Atoi(eventLimit)
	logger.Info("mongodb fans event_log uid total nums", tableNameX, fans_uid, count)
	if count > eventLimitNum {
		removeNum := count - eventLimitNum
		logger.Info("mongodb need remove fans event_log data nums", tableNameX, fans_uid, removeNum)
		created := e.getRemoveCreatedValue(c, eventLimitNum, fans_uid, typeIds)
		if created != "" {
			c.RemoveAll(&bson.M{"type": bson.M{"$in": typeIds}, "fuid": fans_uid, "created": bson.M{"$lte": created}})
			logger.Info("mongodb already remove fans event_log data ", fans_uid, created)
			slptime, _ := strconv.Atoi(sleeptime)
			time.Sleep(time.Duration(slptime) * time.Millisecond)
		}
	}
}

func (e *EventLogNew) getRemoveCreatedValue(c *mgo.Collection, offset int, fuid int, typeIds [7]int) string {
	ms := []EventLogX{}
	err1 := c.Find(&bson.M{"fuid": fuid, "type": bson.M{"$in": typeIds}}).Sort("-created").Skip(offset).Limit(1).All(&ms)

	if err1 != nil {
		logger.Info("mongodb find data", err1, c)
		return ""
	}
	if len(ms) == 0 {
		return ""
	}
	created := ""
	for _, v := range ms {
		created = v.Created
	}
	return created
}

//获取相同犬种的活跃用户
func GetBreedActiveUser(Bid int, session *mgo.Session) []int {
	var user []int
	c := session.DB("ActiveUser").C("active_breed_user")
	err := c.Find(&bson.M{"breed_id": Bid}).Select(bson.M{"pet_id": 1}).Distinct("uid", &user)
	if err != nil {
		panic(err)
	}
	return user
}

//获取相同俱乐部的活跃用户
func GetForumActiveUsers(tid int, session *mgo.Session, db *sql.DB) []int {
	// Forumids := mysql.GetFollowForumIds(uid, db)
	Forumid := mysql.GetFollowForumId(tid, db)
	fmt.Println(Forumid)
	var user []int
	c := session.DB("ActiveUser").C("active_forum_user")
	// err := c.Find(&bson.M{"forum_id": bson.M{"$in": Forumids}}).Distinct("uid", &user)
	err := c.Find(&bson.M{"forum_id": Forumid}).Distinct("uid", &user)
	if err != nil {
		panic(err)
	}
	return user
}

//合并俱乐部和粉丝数据排重
func MergeFansAndForumUsers(fans []*mysql.Follow, tid int, session *mgo.Session, db *sql.DB) []int {
	forumuids := GetForumActiveUsers(tid, session, db)
	var users []int
	//所有加入该帖子俱乐部的活跃用户
	for _, f := range forumuids {
		users = append(users, f)
	}
	// 发帖用户的所有活跃粉丝
	for _, v := range fans {
		users = append(users, v.Follow_id)
	}
	// 数据排重
	allusers := Rm_duplicate(users)
	return allusers
}

// 数据排重
func Rm_duplicate(list []int) []int {
	var x []int = []int{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

//全部活跃用户
func GetAllActiveUsers(session *mgo.Session) []int {
	var user []int
	c := session.DB("ActiveUser").C("active_user")
	err := c.Find(nil).Distinct("uid", &user)
	if err != nil {
		panic(err)
	}
	return user
}
