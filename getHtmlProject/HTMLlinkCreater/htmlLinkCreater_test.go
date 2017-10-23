package HTMLlinkCreater

import (
	"fmt"
	"testing"
)

const logLevel = 1
const jobType = "asksave"
const dbAuth = "dog123:dog123"
const dbDsn = "192.168.86.193:3307"
const dbName = "new_dog123"
const sdbAuth = "dog123:dog123"
const sdbDsn = "192.168.86.193:3307"
const sdbName = "process"

func TestCreate(t *testing.T) {
	var h *HtmlLinkCreater = NewHtmlLinkCreater(logLevel, jobType, dbAuth, dbDsn, dbName, sdbAuth, sdbDsn, sdbName)
	startId := 0
	endId := 50
	page := 1
	cat := "update"
	lastdate := ""
	str := h.Create(startId, endId, page, cat, lastdate)
	fmt.Println(str)
}

func TestGetMaxId(t *testing.T) {
	var h *HtmlLinkCreater = NewHtmlLinkCreater(logLevel, jobType, dbAuth, dbDsn, dbName, sdbAuth, sdbDsn, sdbName)
	str := h.GetMaxId()
	fmt.Println(str)
}

func TestGetProcessLastDate(t *testing.T) {
	var h *HtmlLinkCreater = NewHtmlLinkCreater(logLevel, jobType, dbAuth, dbDsn, dbName, sdbAuth, sdbDsn, sdbName)
	str := h.GetProcessLastDate()
	fmt.Println(str)
}
