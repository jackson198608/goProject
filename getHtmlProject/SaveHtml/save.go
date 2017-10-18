package SaveHtml

import (
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

const proxyServer = "http-pro.abuyun.com:9010"
const proxyUser = "HK71T41EZ21304GP"
const proxyPasswd = "75FE0C4E23EEA0E7"

// const proxyServer = ""
// const proxyUser = ""
// const proxyPasswd = ""

type HtmlInfo struct {
	id        int
	url       string
	queueName string
	saveDir   string
	tidStart  string
	tidEnd    string
	domain    string
	client    *redis.Client
}

func NewHtml(logLevel int, queueName string, id int, url string, taskNewArgs []string, client *redis.Client) *HtmlInfo {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(HtmlInfo)
	e.id = id
	e.url = url
	e.queueName = queueName
	e.saveDir = taskNewArgs[3]
	e.tidStart = taskNewArgs[4]
	e.tidEnd = taskNewArgs[5]
	e.domain = taskNewArgs[6]
	e.client = client
	return e
}

func (e *HtmlInfo) CreateHtmlByUrl() {
	statusCode, _, body, err := e.changeIpByAbuyun()
	if err != nil {
		logger.Error("change ip abuyun error", err)
		return
	}
	if statusCode == 200 {
		urlname := e.saveFileName()
		status := e.saveContentToHtml(urlname, body)
		if status == true {
			logger.Info("save content to html: ", urlname)
		}
	} else {
		str := e.url + "|" + strconv.Itoa(e.id)
		result := (*e.client).LPush(e.queueName, str).Val()
		fmt.Println("resave redis ", str, result)
		logger.Info("resave redis: ", str, result)
	}
}

func (e *HtmlInfo) changeIpByAbuyun() (int, *http.Header, string, error) {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy(proxyServer,
		proxyUser,
		proxyPasswd)

	logger.Info("begin the test", e.id)

	if abuyun == nil {
		logger.Error("create abuyun error")
	}
	var h http.Header = make(http.Header)
	h.Set("a", "1")
	statusCode, responseHeader, body, err := abuyun.SendRequest(e.url, h, true)
	return statusCode, responseHeader, body, err
}

func (e *HtmlInfo) saveFileName() string {
	filename := ""
	dir := ""
	if e.id < 1000 {
		dir = ""
	} else {
		n4 := e.id % 10               //个位数
		n3 := (e.id - n4) % 100       //十位数
		n2 := (e.id - n4 - n3) % 1000 //百位数
		dir = strconv.Itoa(n2/100) + "/" + strconv.Itoa(n3/10) + "/" + strconv.Itoa(n4) + "/"
	}
	urlstr := strings.Split(e.url, "/")
	strlen := len(urlstr)
	if strlen >= 1 {
		filename = dir + urlstr[strlen-1]
	}
	return filename
}

func (e *HtmlInfo) saveContentToHtml(urlname string, content string) bool {
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
