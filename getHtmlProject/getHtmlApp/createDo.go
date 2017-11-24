package main

import (
	"errors"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"github.com/jackson198608/goProject/getHtmlProject/HTMLlinkCreater"
	"github.com/jackson198608/goProject/getHtmlProject/Redis"
	"github.com/jackson198608/goProject/getHtmlProject/Task"
	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v4"
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
				url = c.domain + "thread-" + id + "-" + strconv.Itoa(page) + "-1.html?twig|" + id
			}
			if jobType == "bbsindexsave" {
				url = c.domain + "index-area-" + strconv.Itoa(page) + ".html?twig|" + "1"
			}
			if jobType == "forumsave" {
				url = c.domain + "forum-" + id + "-" + strconv.Itoa(page) + ".html?twig|" + id
			}
			urls = append(urls, url)
		}
	}
	return urls
}

const proxyServer = "http-pro.abuyun.com:9010"
const proxyUser = "HK71T41EZ21304GP"
const proxyPasswd = "75FE0C4E23EEA0E7"

// const proxyServer = ""
// const proxyUser = ""
// const proxyPasswd = ""

func setAbuyun() *abuyunHttpClient.AbuyunProxy {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy(proxyServer, proxyUser, proxyPasswd)
	return abuyun
}

func connect() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.redisConn,
		Password: "", // no password set
		DB:       0,  //e.db use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.New("[Error] redis connect error")
	}
	return client, nil
}

func saveHtmlUrl(jobType string, cat string) {
	client, err := connect()
	if err != nil {
		logger.Error("[Redis connect error]", err)
		return
	}
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
		Redis.PushTaskData(client, c.queueName, idstr)
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
	r := getDoRedisEngine(jobType)
	err := r.Do()
	if err != nil {
		logger.Error("[redisEngine Do] ", err)
	}
}

func jobFunc(redisStr string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {
	var abuyun *abuyunHttpClient.AbuyunProxy
	if taskarg[2] == "1" {
		abuyun = setAbuyun()
	}
	t, err := Task.NewTask(c.logLevel, c.queueName, redisStr, taskarg, abuyun)
	if err != nil {
		logger.Error("[NewTask]", err)
	}
	err = t.Do()
	if err != nil {
		logger.Error("[task Do]", err)
	}
	return err
}

func getDoRedisEngine(jobType string) *redisEngine.RedisEngine {
	var mongoConnInfo []string
	var mysqlInfo []string
	mysqlInfo = append(mysqlInfo, c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")

	redisInfo := redis.Options{
		Addr: c.redisConn,
	}
	r, err := redisEngine.NewRedisEngine(c.queueName, &redisInfo, mongoConnInfo, mysqlInfo, c.numloops, jobFunc, c.saveDir, c.host, c.is_abuyun, jobType)
	if err != nil {
		logger.Error("[NewRedisEngine] ", err)
	}
	return r
}
