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
	TypeId  int           "type"
	Uid     int           "uid"
	Fuid    int           "fuid" //fans id
	Info    string        "info"
	Created string        "created"
	Infoid  int           "infoid"
	Status  int           "status"
	Tid     int           "tid"
}

type EventLogNew struct {
	Id      bson.ObjectId "_id"
	TypeId  int           "type"
	Uid     int           "uid"
	Info    string        "info"
	Created string        "created"
	Infoid  int           "infoid"
	Status  int           "status"
	Tid     int           "tid"
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
	// tableName := "event_log_" + strconv.Itoa(tableNum)
	tableName := "event_log" //动态表
	c := session.DB("EventLog").C(tableName)

	m1 := EventLogNew{bson.NewObjectId(), event.typeId, event.uid, event.info, event.created, event.infoid, event.status, event.tid}
	logger.Info(m1)
	//判断数据是否存在
	eventIsExist := checkMongoIsExist(c, event, event.uid)
	// fmt.Println(eventIsExist)
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
		m := EventLogX{bson.NewObjectId(), event.typeId, event.uid, ar.follow_id, event.info, event.created, event.infoid, event.status, event.tid}
		eventIsExist := checkMongoIsExist(c, event, ar.follow_id)
		if eventIsExist == false {
			// err = c.Insert(&m) //插入数据
			if err != nil {
				logger.Info("mongodb insert fans data", err, c)
			}
		}
		// logger.Info(m)
		n1, _ := w.Seek(0, os.SEEK_END)
		lineStr := fmt.Sprintf("%s", m)
		_, err = w.WriteAt([]byte(lineStr+"\n"), n1)
		if err != nil {
			logger.Info("mongodb write data", err, c)
		}
	}

}

func checkMongoIsExist(c *mgo.Collection, event *EventLog, fuid int) bool {
	ms := []EventLogX{}
	err1 := c.Find(&bson.M{"uid": event.uid, "fuid": fuid, "created": event.created, "infoid": event.infoid}).All(&ms)

	if err1 != nil {
		logger.Info("mongodb find data", err1, c)
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
