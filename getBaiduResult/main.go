package main

import (
	"encoding/csv"
	"fmt"
	"github.com/jackson198608/goProject/stayProcess"
	"github.com/jackson198608/gotest/go_spider/core/pipeline"
	"github.com/jackson198608/gotest/go_spider/core/spider"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var c Config = Config{
	"192.168.86.72:3309",
	"test_dz2",
	"root:goumintech",
	1, "192.168.86.68:6379", "keyword", 1, "/tmp/importkeyword.csv"}

var result *os.File

var logger *log.Logger

func onstart() {
	fileName := os.Args[2]
	if checkFileIsExist(fileName) {
		fmt.Println("file exist")
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			panic("can not open file")
		}
		logger = log.New(file, "", log.LstdFlags)
	} else {
		fmt.Println("new file")
		file, err := os.Create(fileName)
		if err != nil {
			os.Exit(1)
		}
		logger = log.New(file, "", log.LstdFlags)
	}

}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func getRedisData() {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName)
	var start int64 = 0
	var stop int64 = 9
	var limit int64 = 10
	for {
		keywords := r.GetKeywordData(c.queueName, start, stop)
		logger.Println("keyword redis data list", keywords)
		if len(keywords) == 0 {
			break
		}
		for i := 0; i < len(keywords); i++ {
			getRankList(keywords[i])
			_, _, IsExist := checkKeywordExist(keywords[i])
			if IsExist == false {
				fmt.Println("keyword save ", keywords[i])
				saveKeywordRankData(keywords[i], 101, "http://m.goumin.com", "m.goumin.com")
			}
		}
		start += limit
		stop += limit
	}
}

func getRedisDataNew(x chan int, i int) {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName)
	for {
		keyword := r.GetKeywordDataNew(c.queueName)
		logger.Println("keyword redis data list", keyword)
		if keyword == "" {
			logger.Println("keyword got nothing", c.queueName)
			x <- 1
			return
		}
		getRankList(keyword)
		_, _, IsExist := checkKeywordExist(keyword)
		if IsExist == false {
			fmt.Println("keyword save ", keyword)
			saveKeywordRankData(keyword, 101, "http://m.goumin.com", "m.goumin.com")
		}
	}
}

func loopThead() {
	logger.Println("do in oneloop taskNum", c.numloops)
	x := make(chan int, c.numloops)
	for i := 0; i < c.numloops; i++ {
		go getRedisDataNew(x, i)
	}

	for i := 0; i < c.numloops; i++ {
		<-x
	}
}

func importKeyword() {
	importfile := c.importFile
	// importfile := os.Args[3]
	// importfile := "/tmp/importkeyword.csv"
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName)
	file, err := os.Open(importfile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	if r.KeywordExist(c.queueName) > 0 {
		return
	}
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("&&&" + strings.Trim(record[0], " ") + "&&&") // record has the type []string
		logger.Println("save keyword to redis:", record[0])       // record has the type []string
		r.SaveKeywordRedis(c.queueName, strings.Trim(record[0], " "))
	}
}

func Init() {
	loadConfig()
}

func judgeMode(keyword string) (bool, int, string, int) {
	r := stayProcess.NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName)
	times := r.GetTimes(keyword)
	if times >= 5 {
		return false, times, "", 0
	}
	url, num := r.GetUrl(keyword)
	if url == "" {
		return false, times, "", 0
	}
	return true, times, url, num
}

func getRankList(keyword string) {
	// keyword := os.Args[2]
	// realKeyWord := keyword + "  site:bbs.goumin.com"  //pc
	realKeyWord := keyword //+ "  site:m.goumin.com" //h5
	tempUrl := "http://1.com/?wd=" + realKeyWord
	tempUrlP, _ := url.Parse(tempUrl)
	realKeyWordEncode := tempUrlP.Query().Encode()
	result.WriteString(realKeyWordEncode + "\n")
	startUrl := ""
	startUrlTag := ""
	IsExist, times, url, num := judgeMode(keyword)
	fmt.Println("isExist,times,url,num:", IsExist, times, url, num)
	if times > 5 {
		return
	}
	if IsExist == true && num > 1 {
		startUrl = url
		startUrlTag = "searchListNextPage|" + strconv.Itoa(num)
	} else {

		url1 := "https://m.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&"
		url2 := "&rsv_pq=e04aca4d0000cf55&rsv_t=dabboKw3o5qu1XZAgI43hhbjd2olBB3puS%2Fqgn7abC1zNtc%2BA4jzjem3%2BEI&rqlang=cn&rsv_enter=1&rsv_sug3=28&rsv_sug1=17&rsv_sug7=100&rsv_sug2=0&inputT=145353&rsv_sug4=146135"

		startUrlTag = "searchList"
		//startUrl := "https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=%E9%87%91%E6%AF%9B%20site%3Abbs.goumin.com&rsv_pq=e04aca4d0000cf55&rsv_t=dabboKw3o5qu1XZAgI43hhbjd2olBB3puS%2Fqgn7abC1zNtc%2BA4jzjem3%2BEI&rqlang=cn&rsv_enter=1&rsv_sug3=28&rsv_sug1=17&rsv_sug7=100&rsv_sug2=0&inputT=145353&rsv_sug4=146135"
		startUrl = url1 + realKeyWordEncode + url2
	}
	req := newRequest(startUrlTag, startUrl)
	logger.Println("keyword", keyword)
	logger.Println("search url", startUrl)
	spider.NewSpider(NewMyPageProcesser(), "getBaiduResult").
		AddRequest(req).
		AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		SetThreadnum(5).                            // Crawl request by three Coroutines
		Run()
}

func main() {
	Init()
	onstart()
	jobType := os.Args[1]
	switch jobType {
	case "baidu":
		logger.Println("Start get baidu keyword rank")
		// getRedisData()
		loopThead()
	case "keyword":
		logger.Println("Start import keyword to redis")
		importKeyword()
	default:

	}
}
