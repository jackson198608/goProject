package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"net/http"
	"os"
	"path"
	"strconv"
)

const proxyServer = ""
const proxyUser = ""
const proxyPasswd = ""

func (e *InfoNew) CreateHtmlByUrl(id int, pages int, jobType string) {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy(proxyServer,
		proxyUser,
		proxyPasswd)

	logger.Info("begin the test", id)

	if abuyun == nil {
		logger.Error("create abuyun error")
		return
	}
	for page := 1; page <= pages; page++ {
		targetUrl := e.getTargetUrl(id, page, jobType)
		var h http.Header = make(http.Header)
		h.Set("a", "1")
		statusCode, responseHeader, body, err := abuyun.SendRequest(targetUrl, h, true)
		fmt.Println(statusCode)
		fmt.Println(responseHeader)
		fmt.Println(body)
		fmt.Println(err)
		urlname := e.saveFilename(id, page, jobType)
		e.saveContentToHtml(urlname, body)
	}

}

type InfoNew struct {
	db           *sql.DB
	id           int
	templateType string
	templatefile string
	saveDir      string
	tidStart     string
	tidEnd       string
	domain       string
}

func NewInfo(logLevel int, id int, db *sql.DB, taskNewArgs []string) *InfoNew {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(InfoNew)
	e.db = db
	e.id = id
	e.saveDir = taskNewArgs[3]
	e.tidStart = taskNewArgs[4]
	e.tidEnd = taskNewArgs[5]
	e.domain = taskNewArgs[6]
	return e
}

func (e *InfoNew) getTargetUrl(id int, page int, jobType string) string {
	var url string = ""
	if jobType == "ask" {
		url = e.domain + strconv.Itoa(id) + ".html"
	}
	if jobType == "threadsave" {
		url = e.domain + "thread-" + strconv.Itoa(id) + "-" + strconv.Itoa(page) + "-1.html"
	}
	return url
}

func (e *InfoNew) saveFilename(id int, page int, jobType string) string {
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
	if jobType == "ask" {
		filename = e.saveDir + dir + strconv.Itoa(id) + ".html"
	}
	if jobType == "thread" {
		filename = e.saveDir + dir + "thread-" + strconv.Itoa(id) + "-" + strconv.Itoa(page) + "-1.html"
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
