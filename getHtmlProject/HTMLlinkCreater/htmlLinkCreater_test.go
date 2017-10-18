package HTMLlinkCreater

import "fmt"

var channelTable map[string]string = map[string]string{
	"ask":    "ask_question",
	"thread": "pre_forum_thread",
}

type HtmlLinkCreater struct {
	channelType   string
	jobType       int
	lastIndexTime string
	tableName     string
	tableLastId   int64
	xorm          *xorm
}

func NewHtmlLinkCreater(channelType string, jobType int, lastIndexTime string) *HtmlLinkCreater {
	h := new(HtmlLinkCreater)
	if h == nil {
		return nil
	}

	//check param

	//pass prams to object

	//get tableName and lastId
	h.tableName = h.getTableNameFromChannel()

	//new xorm instance

	return h
}

func (h *HtmlLinkCreater) Create() error {
	buidAll()
	updateAll()
}

//
func (h *HtmlLinkCreater) buidAll() error {
	//for startid:=0,endId:=100; startid>h.tableLastId ;startid=startid+100,endId=endId+100{
	h.getIdsFromTableByRange()
	//}
}

func (h *HtmlLinkCreater) updateAll() error {
	//for startid:=0,endId:=100; startid>h.tableLastId ;startid=startid+100,endId=endId+100{

	//}
}

//
func (h *HtmlLinkCreater) getIdsFromTableByRange(startid int, endId int) []int {
	h.xorm.where().find(&pre_forum_thread.Pre_fomrum_thread)

}

//get tableName by oject's channelType from common var channelType
func (h *HtmlLinkCreater) getTableNameFromChannel() {

}

// get last id of table
func (h *HtmlLinkCreater) getLastIdFromTable() {

}
