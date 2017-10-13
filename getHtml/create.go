package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	redis "gopkg.in/redis.v4"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

const proxyServer = ""
const proxyUser = ""
const proxyPasswd = ""

type InfoNew struct {
	db       *sql.DB
	id       int
	saveDir  string
	tidStart string
	tidEnd   string
	domain   string
	client   *redis.Client
}

func NewInfo(logLevel int, id int, db *sql.DB, taskNewArgs []string, client *redis.Client) *InfoNew {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(InfoNew)
	e.db = db
	e.id = id
	e.saveDir = taskNewArgs[3]
	e.tidStart = taskNewArgs[4]
	e.tidEnd = taskNewArgs[5]
	e.domain = taskNewArgs[6]
	e.client = client
	return e
}

func (e *InfoNew) CreateHtmlByUrl(id int, targetUrl string, jobType string) {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy(proxyServer,
		proxyUser,
		proxyPasswd)

	logger.Info("begin the test", id)

	if abuyun == nil {
		logger.Error("create abuyun error")
		return
	}
	var h http.Header = make(http.Header)
	h.Set("a", "1")
	statusCode, responseHeader, body, err := abuyun.SendRequest(targetUrl, h, true)
	fmt.Println(statusCode, id)
	fmt.Println(responseHeader)
	// fmt.Println(body)
	fmt.Println(err)
	if statusCode == 200 {
		urlname := e.newSaveFilename(targetUrl, id)
		status := e.saveContentToHtml(urlname, body)
		if status == true {
			logger.Info("save content to html: ", urlname)
		}
	} else {
		fmt.Println("resave id and pages to redis")
		str := targetUrl + "|" + strconv.Itoa(id)
		result := (*e.client).LPush(c.queueName, str).Val()
		fmt.Println("resave redis ", str, result)
	}

}

func (e *InfoNew) newSaveFilename(url string, id int) string {
	filename := ""
	dir := ""
	if id < 1000 {
		dir = ""
	} else {
		n4 := id % 10               //个位数
		n3 := (id - n4) % 100       //十位数
		n2 := (id - n4 - n3) % 1000 //百位数
		dir = strconv.Itoa(n2/100) + "/" + strconv.Itoa(n3/10) + "/" + strconv.Itoa(n4) + "/"
	}
	urlstr := strings.Split(url, "/")
	strlen := len(urlstr)
	if strlen >= 1 {
		filename = dir + urlstr[strlen-1]
	}
	return filename
}

func (e *InfoNew) saveContentToHtml(urlname string, content string) bool {
	var filename = e.saveDir + urlname
	var f *os.File
	var err1 error
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	fullDirPath := path.Dir(filename)
	err1 = os.MkdirAll(fullDirPath, 0777)
	check(err1)
	f, err1 = os.Create(filename) //创建文件
	check(err1)
	n, err1 := f.WriteString(content) //写入文件(字符串)
	check(err1)
	fmt.Printf("写入 %d 个字节n", n)
	return true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

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
				url = c.domain + id + ".html|" + id
				if page > 1 {
					url = c.domain + id + "-" + strconv.Itoa(page) + ".html|" + id
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
