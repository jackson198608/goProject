package main

import (
	"fmt"
	"github.com/jackson198608/goProject/getHtmlProject/HTMLlinkCreater"
	"github.com/jackson198608/goProject/getHtmlProject/RedisEngine"
	"strconv"
	"strings"
)

func idToUrl(jobType string, idstr []string) []string {
	var urls []string
	for _, v := range idstr {
		vArr := strings.Split(v, "|")
		if len(vArr) < 2 {
			break
		}
		id := vArr[0]
		pages, _ := strconv.Atoi(vArr[1])
		for page := 0; page <= pages; page++ {
			if page == 0 {
				page = 1
			}
			var url string = ""
			if jobType == "asksave" {
				url = c.domain + id + ".html?twig|" + id
				if page > 1 {
					url = c.domain + id + "-" + strconv.Itoa(page) + ".html?twig|" + id
				}
			}
			if jobType == "threadsave" {
				url = c.domain + "thread-" + id + "-" + strconv.Itoa(page) + "-1.html|" + id
			}
			urls = append(urls, url)
		}
	}
	return urls
}

func saveHtmlUrl(jobType string, cat string) {
	r := RedisEngine.NewEngine(c.logLevel, c.queueName, c.redisConn, "", c.numloops, c.dbAuth, c.dbDsn, c.dbName)
	h := HTMLlinkCreater.NewHtmlLinkCreater(c.logLevel, jobType, c.dbAuth, c.dbDsn, c.dbName, c.sdbAuth, c.sdbDsn, c.sdbName)
	page := 1
	intIdStart, _ := strconv.Atoi(c.tidStart)
	startId := intIdStart
	endId := startId + 1000
	maxId := h.GetMaxId()
	lastdate := ""
	if cat == "update" {
		lastdate = h.GetProcessLastDate()
	}
	fmt.Println(maxId)
	for {
		var ids []string
		ids = h.Create(startId, endId, page, cat, lastdate)
		idstr := idToUrl(jobType, ids)
		r.PushTaskData(idstr)
		if cat == "update" {
			if len(ids) == 0 {
				break
			}
		} else {
			if startId > maxId {
				break
			}
			startId += offset
			endId += offset
		}
		page++
	}
	h.SaveProcessLastdate()
}

func createHtmlByUrl(jobType string) {
	r := RedisEngine.NewEngine(c.logLevel, c.queueName, c.redisConn, jobType, c.numloops, c.saveDir, c.host)
	r.Loop()
}
