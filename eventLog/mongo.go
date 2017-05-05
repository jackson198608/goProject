package main

import (
	// "encoding/json"
	// "bufio"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"os"
	// "reflect"
	"strconv"
)

type EventLogX struct {
	Id      bson.ObjectId "_id"
	TypeId  int
	Uid     int
	Fuid    int //fans id
	Info    string
	Created string
	Infoid  int
	Status  int
	Tid     int
}

func SaveMongoEventLog(event *EventLog, fans []*Follow, w *os.File) {
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		logger.Info("mongodb connect fail", err)
		return
	}
	defer session.Close()

	tableNum := event.uid % 100
	if tableNum == 0 {
		tableNum = 100
	}
	tableName := "event_log_" + strconv.Itoa(tableNum)
	c := session.DB("EventLog").C(tableName)

	m1 := EventLogX{bson.NewObjectId(), event.typeId, event.uid, event.uid, event.info, event.created, event.infoid, event.status, event.tid}
	logger.Info(m1)
	// err = c.Insert(&m1)
	//判断数据是否存在
	eventIsExist := checkMongoIsExist(c, event, event.uid)
	fmt.Println(eventIsExist)
	if eventIsExist == false {
		err = c.Insert(&m1)
		if err != nil {
			logger.Info("mongo insert error:", err)
		}
	}
	// fmt.Println("type:", reflect.TypeOf(c))
	lineStr1 := fmt.Sprintf("%s", m1)
	// /*fmt.Fprintln(w, lineStr1)*/
	// 查找文件末尾的偏移量
	n, _ := w.Seek(0, os.SEEK_END)
	// 从末尾的偏移量开始写入内容
	_, err = w.WriteAt([]byte(lineStr1+"\n"), n)
	for _, ar := range fans {
		tableNum1 := ar.follow_id % 100
		if tableNum1 == 0 {
			tableNum1 = 100
		}
		tableName1 := "event_log_" + strconv.Itoa(tableNum1)
		fmt.Println(tableName1)
		c := session.DB("EventLog").C(tableName1)
		m := EventLogX{bson.NewObjectId(), event.typeId, event.uid, ar.follow_id, event.info, event.created, event.infoid, event.status, event.tid}
		// err = c.Insert(&m)
		eventIsExist := checkMongoIsExist(c, event, ar.follow_id)
		fmt.Println(eventIsExist)
		if eventIsExist == false {
			err = c.Insert(&m)
			logger.Info("mongo insert error:", err)
		}
		logger.Info(m)
		if err != nil {
			logger.Info("mongodb insert data", err, c)
		}
		n1, _ := w.Seek(0, os.SEEK_END)
		lineStr := fmt.Sprintf("%s", m)
		_, err = w.WriteAt([]byte(lineStr+"\n"), n1)
		if err != nil {
			logger.Info("mongodb insert data", err, c)
		}
		// fmt.Fprintln(w, lineStr)
		// writeFile(fmt.Sprintln(m))
	}
	// w.Flush()
	// fmt.Println("type:", reflect.TypeOf(w))
}

func checkMongoIsExist(c *mgo.Collection, event *EventLog, fuid int) bool {
	ms := []EventLogX{}
	err1 := c.Find(&bson.M{"uid": event.uid, "fuid": fuid, "created": event.created, "infoid": event.infoid}).All(&ms)

	if err1 != nil {
		logger.Info("mongodb insert data", err1, c)
		return false
	}
	if len(ms) == 0 {
		return false
	}
	return true
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func writeFile(m string) {
	var err1 error
	var f *os.File
	// var m string = "ddd"
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	check(err1)
	n, err1 := io.WriteString(f, m) //写入文件(字符串)
	check(err1)
	fmt.Printf("写入 %d 个字节n", n)
}

func check(e error) {
	if e != nil {
		logger.Info("check file error", e)
	}
}

// func NewEvent(logLevel int, dbFlavor string, id int64) *EventLog {
// 	u := new(EventLog)
// 	logger.SetLevel(logger.LEVEL(logLevel))

// 	// u.isSplit = isSplit
// 	// if id > 0 {
// 	// 	u.id = id
// 	// }

// 	if id > 0 {
// 		u = u.LoadById()
// 	}

// 	u.logLevel = logLevel

// 	return u
// }

// func (p *EventLog) LoadById() *EventLog {
// 	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
// 	if err != nil {
// 		logger.Error("[error] connect db err")
// 	}
// 	defer db.Close()
// 	tableName := "event_log"
// 	rows, err := db.Query("select id,type as typeId,uid,info,created,infoid,status,tid from `" + tableName + "` where id=" + strconv.Itoa(int(p.id)) + "")
// 	defer rows.Close()
// 	if err != nil {
// 		logger.Error("[error] check event_log sql prepare error: ", err)
// 		return nil
// 	}
// 	for rows.Next() {
// 		var row = new(EventLog)
// 		rows.Scan(&row.id, &row.typeId, &row.uid, &row.info, &row.created, &row.infoid, &row.status, &row.tid)
// 		return row
// 	}
// 	return &EventLog{}
// }

// func (p *EventLog) MoveToSplit() bool {
// 	session, err := mgo.Dial(c.mongoConn)
// 	if err != nil {
// 		return
// 	}
// 	defer session.Close()

// 	// Optional. Switch the session to a monotonic behavior.
// 	session.SetMode(mgo.Monotonic, true)

// 	//create channel
// 	c := make(chan int, taskNum)

// 	for i := 0; i < taskNum; i++ {
// 		go doInMongo(c, session, i)
// 	}
// 	for i := 0; i < taskNum; i++ {
// 		<-c
// 	}
// }

// func doInMongo(c chan int, session *mgo.Session, i int) {
// 	mongoStr := tasks[i].insertStr
// 	if mongoStr == "0" {
// 		c <- 1
// 		return
// 	}
// 	t := NewMongoTask(mongoStr)
// 	w := NewWorker(t)
// 	w.Insert(session)

// 	c <- 1
// }

// func insertMongo() {
// 	//init session
// 	session, err := mgo.Dial(c.mongoConn)
// 	if err != nil {
// 		return
// 	}
// 	defer session.Close()

// 	// Optional. Switch the session to a monotonic behavior.
// 	session.SetMode(mgo.Monotonic, true)

// 	//create channel
// 	c := make(chan int, taskNum)

// 	for i := 0; i < taskNum; i++ {
// 		go doInMongo(c, session, i)
// 	}
// 	for i := 0; i < taskNum; i++ {
// 		<-c
// 	}
// }

// type Worker struct {
// 	t *Task
// }

// func NewWorker(t *Task) (w *Worker) {
// 	//init the worker
// 	var wR Worker
// 	wR.t = t
// 	return &wR
// }

// func (w Worker) Insert(session *mgo.Session) {
// 	//convert json string to struct
// 	var m row
// 	if err := json.Unmarshal([]byte(w.t.columData), &m); err != nil {
// 		//fmt.Println("[error] mongo json error", err, w.t.columData)
// 		return
// 	}

// 	//get the table name
// 	tableNumber := strconv.Itoa(m.Uid % 1000)
// 	tableName := "message_push_record_" + tableNumber

// 	//create mongo session
// 	c := session.DB("MessageCenter").C(tableName)

// 	err := c.Insert(&m)
// 	if err != nil {
// 		//fmt.Println("[Error]insert into mongo error", err)
// 		return
// 	}
// }
