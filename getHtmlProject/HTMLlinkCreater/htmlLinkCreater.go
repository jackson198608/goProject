package HTMLlinkCreater

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/getHtmlProject/ActiveRecord"
	"math"
	"strconv"
	"time"
)

var channelTable map[string]string = map[string]string{
	"ask":    "ask_question",
	"thread": "pre_forum_thread",
}

type HtmlLinkCreater struct {
	jobType     string
	dbAuth      string
	dbDsn       string
	dbName      string
	sdbAuth     string
	sdbDsn      string
	sdbName     string
	tableName   string
	tableLastId int64
	engine      *xorm.Engine
	sengine     *xorm.Engine
	// xorm        *xorm
}

// type PreForumThread struct {
// 	Tid          int
// 	Displayorder int
// 	Posttableid  int
// }

// type AskQuestion struct {
// 	Id     int
// 	AnsNum int
// }

func NewHtmlLinkCreater(logLevel int, jobType string, dbAuth string, dbDsn string, dbName string, sdbAuth string, sdbDsn string, sdbName string) *HtmlLinkCreater {
	logger.SetLevel(logger.LEVEL(logLevel))
	h := new(HtmlLinkCreater)
	if h == nil {
		return nil
	}
	h.dbAuth = dbAuth
	h.dbDsn = dbDsn
	h.dbName = dbName
	h.sdbAuth = sdbAuth
	h.sdbDsn = sdbDsn
	h.sdbName = sdbName
	//check param

	//pass prams to object

	//get tableName and lastId
	h.tableName = h.getTableNameFromChannel()
	h.jobType = jobType
	h.engine = h.setEngine()
	h.sengine = h.setSEngine()
	//new xorm instance

	return h
}

func (h *HtmlLinkCreater) setEngine() *xorm.Engine {
	if h.jobType == "asksave" {
		h.dbName = "ask"
	}
	dataSourceName := h.dbAuth + "@tcp(" + h.dbDsn + ")/" + h.dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return engine
}

func (h *HtmlLinkCreater) setSEngine() *xorm.Engine {
	dataSourceName := h.sdbAuth + "@tcp(" + h.sdbDsn + ")/" + h.sdbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return engine
}

func (h *HtmlLinkCreater) Create(startId int, endId int, page int, cat string, lastdate string) []string {
	var ids []string
	if cat == "update" {
		ids = h.updateAll(page, lastdate)
	} else {
		ids = h.buildAll(startId, endId)
	}
	return ids
}

//get all data
func (h *HtmlLinkCreater) buildAll(startId int, endId int) []string {
	var ids []string
	if h.jobType == "threadsave" {
		ids = h.getThreadTask(startId, endId)
	}
	if h.jobType == "asksave" {
		ids = h.getAskTask(startId, endId)
	}
	return ids
}

func (h *HtmlLinkCreater) getAskTask(startId int, endId int) []string {
	var ids []string
	var asks []ActiveRecord.AskQuestion
	err := h.engine.Where("is_hide=? and id>? and id<=?", 1, startId, endId).Cols("id", "ans_num").Asc("id").Find(&asks)
	if err != nil {
		return nil
	}
	for _, v := range asks {
		totalpages := int(math.Ceil(float64(v.AnsNum) / float64(5))) //page总数
		s := strconv.Itoa(v.Id) + "|" + strconv.Itoa(totalpages)
		ids = append(ids, s)
	}
	return ids
}

func (h *HtmlLinkCreater) getThreadTask(startId int, endId int) []string {
	var ids []string
	var thread []ActiveRecord.PreForumThread
	err := h.engine.In("displayorder", 0, 1).Where("tid>? and tid<=?", startId, endId).Cols("tid", "posttableid").Asc("tid").Find(&thread)
	if err != nil {
		return nil
	}
	for _, v := range thread {
		count := h.getPostCount(v.Posttableid)
		totalpages := int(math.Ceil(float64(count) / float64(5))) //page总数
		s := strconv.Itoa(v.Tid) + "|" + strconv.Itoa(totalpages)
		ids = append(ids, s)
	}
	return ids
}

//get update data
func (h *HtmlLinkCreater) updateAll(page int, lastdate string) []string {
	var ids []string
	if h.jobType == "asksave" {
		ids = h.getUpdateAskData(page, lastdate)
	}
	if h.jobType == "threadsave" {
		ids = h.getUpdateThreadData(page, lastdate)
	}
	return ids
}

func (h *HtmlLinkCreater) getUpdateAskData(page int, lastdate string) []string {
	var limit int = 1000
	offset := (page - 1) * limit
	var ids []string
	if lastdate == "" {
		return ids
	}
	var asks []ActiveRecord.AskQuestion
	err := h.engine.Where("is_hide=?", 1).And("created>=? or end_date>=?", lastdate, lastdate).Cols("id", "ans_num").Asc("id").Limit(limit, offset).Find(&asks)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, v := range asks {
		totalpages := int(math.Ceil(float64(v.AnsNum) / float64(5))) //page总数
		s := strconv.Itoa(v.Id) + "|" + strconv.Itoa(totalpages)
		ids = append(ids, s)
	}
	return ids
}

func (h *HtmlLinkCreater) getUpdateThreadData(page int, lastdate string) []string {
	var ids []string
	if lastdate == "" {
		return ids
	}
	var limit int = 1000
	offset := (page - 1) * limit
	timeLayout := "2006-01-02 15:04:05"                           //转化所需模板
	loc, _ := time.LoadLocation("Local")                          //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, lastdate, loc) //使用模板在对应时区转化为time.time类型
	dateint := int(theTime.Unix())
	date := strconv.Itoa(dateint)
	var thread []ActiveRecord.PreForumThread
	err := h.engine.In("displayorder", 0, 1).Where("lastpost>=?", date).Cols("tid", "posttableid").Asc("tid").Limit(limit, offset).Find(&thread)
	if err != nil {
		return nil
	}
	for _, v := range thread {
		count := h.getPostCount(v.Posttableid)
		totalpages := int(math.Ceil(float64(count) / float64(5))) //page总数
		s := strconv.Itoa(v.Tid) + "|" + strconv.Itoa(totalpages)
		ids = append(ids, s)
	}
	return ids
}

func (h *HtmlLinkCreater) GetMaxId() int {

	var maxid int = 0
	if h.jobType == "asksave" {
		maxid = h.getAskMaxId()
	}
	if h.jobType == "threadsave" {
		maxid = h.getThreadMaxId()
	}
	return maxid
}

func (h *HtmlLinkCreater) getAskMaxId() int {
	var ask ActiveRecord.AskQuestion
	_, err := h.engine.Where("is_hide=?", 1).Cols("id").Desc("id").Get(&ask)
	if err != nil {
		logger.Error("get ask maxId error", err)
		return 0
	}
	return ask.Id
}

func (h *HtmlLinkCreater) getThreadMaxId() int {
	var thread ActiveRecord.PreForumThread
	_, err := h.engine.Cols("tid").In("displayorder", 0, 1).Desc("tid").Get(&thread)
	if err != nil {
		logger.Error("get ask maxId error", err)
		return 0
	}
	return thread.Tid
}

//
func (h *HtmlLinkCreater) getIdsFromTableByRange() []int {
	var ids []int
	return ids
}

//get tableName by oject's channelType from common var channelType
func (h *HtmlLinkCreater) getTableNameFromChannel() string {
	var tablename string = ""
	if h.jobType == "asksave" {
		tablename = "ask_question"
	}
	if h.jobType == "thread" {
		tablename = "pre_forum_thread"
	}
	return tablename
}

// get last id of table
func (h *HtmlLinkCreater) getLastIdFromTable() {

}

func (h *HtmlLinkCreater) getPostCount(posttableid int) int {
	tablename := "pre_forum_post_" + strconv.Itoa(posttableid)
	var post = make(map[string]string)
	count, err := h.engine.Table(tablename).Where("invisible=?", 0).Count(&post)
	if err != nil {
		return 0
	}
	fmt.Println(post)
	return int(count)
}

const processName = "getHtmlApp"

func (h *HtmlLinkCreater) checkProcessExist() int {
	var exec ActiveRecord.ExecuteRecord
	has, err := h.sengine.Where("process_name=?", processName).And("data_source=?", h.jobType).Cols("id").Desc("id").Get(&exec)
	if err != nil {
		logger.Error("get ask maxId error", err)
		return 0
	}
	fmt.Println(has)
	return exec.Id
}

func (h *HtmlLinkCreater) insertProcessLastdate() {
	lastdate := time.Now().Format("2006-01-02 15:04:05")
	exec := ActiveRecord.ExecuteRecord{ProcessName: processName, DataSource: h.jobType, Lastdate: lastdate, Created: lastdate}
	_, err := h.sengine.Insert(&exec)
	if err != nil {
		logger.Error("insert process table execute_record error", err)
	}
}

func (h *HtmlLinkCreater) updateProcessLastdate(id int) {
	lastdate := time.Now().Format("2006-01-02 15:04:05")
	exec := ActiveRecord.ExecuteRecord{Lastdate: lastdate}
	_, err := h.sengine.Where("id=?", id).And("process_name=? and data_source=?", processName, h.jobType).Update(exec)
	if err != nil {
		logger.Error("update process table execute_record error", id, err)
	}
}

func (h *HtmlLinkCreater) SaveProcessLastdate() {
	id := h.checkProcessExist()
	if id > 0 {
		h.updateProcessLastdate(id)
	} else {
		h.insertProcessLastdate()
	}
}

func (h *HtmlLinkCreater) GetProcessLastDate() string {
	var exec ActiveRecord.ExecuteRecord
	_, err := h.sengine.Where("process_name=?", processName).And("data_source=?", h.jobType).Cols("lastdate").Desc("id").Get(&exec)
	if err != nil {
		logger.Error("get ask maxId error", err)
		return ""
	}
	return exec.Lastdate
}
