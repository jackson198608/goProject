package HTMLlinkCreater

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/ActiveRecord"
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
	if h.jobType == "bbsindexsave" {
		ids = h.getBbsIndexTask(startId, endId)
	}
	if h.jobType == "forumsave" {
		ids = h.getForumTask(startId, endId)
	}
	return ids
}

func (h *HtmlLinkCreater) getBbsIndexTask(startId int, endId int) []string {
	var ids []string
	forum := new(ActiveRecord.PreForumForum)
	count, err := h.engine.In("fup", 76, 78).Where("status=?", 1).Count(forum)
	if err != nil {
		logger.Error("get forum data ", err)
		return nil
	}
	counts := int(count)
	s := getIdAndPages(1, counts, 20)
	ids = append(ids, s)
	return ids
}

func (h *HtmlLinkCreater) getForumTask(startId int, endId int) []string {
	var ids []string
	var forum []ActiveRecord.PreForumForum
	err := h.engine.In("fup", 76, 78).Where("status=?", 1).Cols("fid").Find(&forum)
	if err != nil {
		logger.Error("get forum data ", err)
		return nil
	}
	for _, v := range forum {
		count := h.getForumThreadCount(v.Fid)
		fmt.Println(count)
		s := getIdAndPages(v.Fid, count, 20)
		ids = append(ids, s)
	}
	return ids
}

func (h *HtmlLinkCreater) getForumThreadCount(fid int) int {
	//SELECT `pre_forum_thread`.* FROM `pre_forum_thread` INNER JOIN `pre_common_member` ON `pre_forum_thread`.`authorid` = `pre_common_member`.`uid` INNER JOIN `pre_ucenter_members` ON `pre_forum_thread`.`authorid` = `pre_ucenter_members`.`uid` WHERE fid=120 and displayorder>=0 and pre_common_member.groupid!=4 ORDER BY `lastpost` DESC LIMIT 20
	thread := new(ActiveRecord.PreForumThread)
	count, err := h.engine.Table("pre_forum_thread").Join("left", "pre_common_member", "pre_forum_thread.authorid=pre_common_member.uid").Where("groupid!=? and displayorder>=? and fid=?", 4, 0, fid).Count(thread)
	if err != nil {
		logger.Error("get forum_thread count ", err)
		return 0
	}
	return int(count)
}

func (h *HtmlLinkCreater) getAskTask(startId int, endId int) []string {
	var ids []string
	var asks []ActiveRecord.AskQuestion
	err := h.engine.Where("is_hide=? and id>? and id<=?", 1, startId, endId).Cols("id", "ans_num").Asc("id").Find(&asks)
	if err != nil {
		return nil
	}
	for _, v := range asks {
		s := getIdAndPages(v.Id, v.AnsNum, 20)
		ids = append(ids, s)
	}
	return ids
}

func (h *HtmlLinkCreater) getThreadTask(startId int, endId int) []string {
	var ids []string
	var thread []ActiveRecord.PreForumThread
	err := h.engine.In("displayorder", 0, 1).NotIn("fid", 62, 74, 7, 71).Where("tid>? and tid<=?", startId, endId).Cols("tid", "posttableid").Asc("tid").Find(&thread)
	if err != nil {
		logger.Error("get thread data ", err)
		return nil
	}
	for _, v := range thread {
		count := h.getPostCount(v.Posttableid, v.Tid)
		s := getIdAndPages(v.Tid, count, 20)
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
		s := getIdAndPages(v.Id, v.AnsNum, 20)
		ids = append(ids, s)
	}
	return ids
}

func formatDateToTime(lastdate string) string {
	timeLayout := "2006-01-02 15:04:05"                           //转化所需模板
	loc, _ := time.LoadLocation("Local")                          //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, lastdate, loc) //使用模板在对应时区转化为time.time类型
	dateint := int(theTime.Unix())
	time := strconv.Itoa(dateint)
	return time
}

func (h *HtmlLinkCreater) getUpdateThreadData(page int, lastdate string) []string {
	var ids []string
	if lastdate == "" {
		return ids
	}
	var limit int = 1000
	offset := (page - 1) * limit
	date := formatDateToTime(lastdate)
	var thread []ActiveRecord.PreForumThread
	err := h.engine.In("displayorder", 0, 1).Where("lastpost>=?", date).Cols("tid", "posttableid").Asc("tid").Limit(limit, offset).Find(&thread)
	if err != nil {
		return nil
	}
	for _, v := range thread {
		count := h.getPostCount(v.Posttableid, v.Tid)
		s := getIdAndPages(v.Tid, count, 20)
		ids = append(ids, s)
	}
	return ids
}

func getIdAndPages(id int, count int, per int) string {
	totalpages := int(math.Ceil(float64(count) / float64(per))) //page总数
	s := strconv.Itoa(id) + "|" + strconv.Itoa(totalpages)
	return s
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

func (h *HtmlLinkCreater) getPostCount(posttableid int, tid int) int {
	tablename := "pre_forum_post"
	if posttableid > 0 {
		tablename = "pre_forum_post_" + strconv.Itoa(posttableid)
	}
	var post = new(ActiveRecord.PreForumPostX)
	count, err := h.engine.Table(tablename).Where("invisible=?", 0).And("tid=?", tid).Count(post)
	if err != nil {
		fmt.Println(err)
		logger.Error("get pre_forum_post count error", err)
		return 0
	}
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
	if has == true {
		return exec.Id
	}
	return 0
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
	has, err := h.sengine.Where("process_name=?", processName).And("data_source=?", h.jobType).Cols("lastdate").Desc("id").Get(&exec)
	if err != nil {
		logger.Error("get ask maxId error", err)
		return ""
	}
	if has == true {
		return exec.Lastdate
	}
	return ""
}
