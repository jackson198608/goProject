package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type InfoAsk struct {
	db           *sql.DB
	id           int
	session      *mgo.Session
	templateType string
	templatefile string
	saveDir      string
	tidStart     string
	tidEnd       string
	domain       string
}

func NewAskInfo(logLevel int, id int, db *sql.DB, session *mgo.Session, taskNewArgs []string) *InfoAsk {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(InfoAsk)
	e.db = db
	e.id = id
	e.session = session
	e.templateType = taskNewArgs[3]
	e.templatefile = taskNewArgs[4]
	e.saveDir = taskNewArgs[5]
	e.tidStart = taskNewArgs[6]
	e.tidEnd = taskNewArgs[7]
	e.domain = taskNewArgs[8]
	return e
}

func (e *InfoAsk) CreateAskHtmlContent(id int, relateDefaultThread string) error {
	question := LoadQuestionById(id, e.db)
	if question.Id <= 0 || question == nil {
		logger.Info("ask_question is not exist id=", id)
		return nil
	}
	fmt.Println(question.Images)
	//相关帖子 eg:tid=12
	relateThread := threadByAsk(question.Id, e.db, e.session)
	if relateThread == "" {
		relateThread = relateDefaultThread
	}
	//相关问答 eg:tid=12
	relateAsk := e.askByAsk(question.Id, question.Pid, e.db, e.session)

	//相关犬种 eg:tid=4682521
	relateDogs := dogsByAsk(question.Id, question.Pid, e.db)
	answers := LoadAnswersById(question.Id, e.db)
	e.groupContentToSaveAskHtml(question, relateThread, relateAsk, relateDogs, answers)
	return nil
}
func getTypeName(typeid int) string {
	if typeid == 1 {
		return "医疗"
	}
	if typeid == 1 {
		return "养护"
	}
	if typeid == 1 {
		return "训练"
	}
	if typeid == 1 {
		return "综合"
	}
	return "医疗"
}

func (e *InfoAsk) groupContentToSaveAskHtml(question *Question, relateThread string, relateAsk string, relateDogs string, answers []*Answer) {
	var title = question.Subject + "_狗民知道_狗民网"
	var keyword = question.Subject + " , " + getTypeName(question.Typeid)
	var description = substr(question.Content, 40)
	html := e.getH5TemplateHtml()
	if html == "" {
		logger.Error("template file not found")
		return
	}
	html = strings.Replace(html, "cmsTypeName", getTypeName(question.Typeid), -1)
	html = strings.Replace(html, "cmsTitle", title, -1)
	html = strings.Replace(html, "cmsSubject", question.Subject, -1)
	html = strings.Replace(html, "cmsRelateThread", relateThread, -1)
	html = strings.Replace(html, "cmsRelateAsk", relateAsk, -1)
	html = strings.Replace(html, "cmsRelateDogs", relateDogs, -1)
	html = strings.Replace(html, "cmsKeywords", keyword, -1)
	html = strings.Replace(html, "cmsDescription", description, -1)
	questionContent := filterContent(question.Content)
	questionContent = findface(questionContent)
	html = strings.Replace(html, "cmsQuestionContent", questionContent, -1)
	images := imageHtml(question.Images)
	html = strings.Replace(html, "cmsQuestionImages", images, -1)
	content := ""
	filename := createAskFileName(question.Id)
	len := len(answers) //总数
	// count := 20                                                 //每页条数
	// totalpages := int(math.Ceil(float64(len) / float64(count))) //page总数
	// for i := 1; i <= totalpages; i++ {
	// start := (i - 1) * count
	// end := start + count
	// if end > len {
	// 	end = len
	// }
	// pagepost := answers[start:end]
	if len > 0 {
		for _, v := range answers {
			message := filterContent(v.Content)
			userinfo := LoadUserinfoByUid(v.Uid, e.db)
			if userinfo == nil {
				continue
			}
			content += "<li><p class=\"text\">" + message + "</p><mip-img src=\"" + userinfo.Avatar + "\"></mip-img><span class=\"time\">回答者：<em class=\"name\"><span>" + userinfo.Nickname + "</span>-" + userinfo.Grouptitle + "</em><em class=\"num\">" + v.Created + "</em></span>"
			comment := getCommentListHtml(v.Id, e.db)
			content += comment
		}
	}
	// cmsPage := cmsPage(totalpages, question.Id)
	oldUrl := "http://m.goumin.com/ask/" + strconv.Itoa(question.Id) + ".html"
	html = strings.Replace(html, "cmsCanical", oldUrl, -1)
	content = findface(content)
	html = strings.Replace(html, "cmsAnswers", content, -1)
	status := e.saveContentToHtml(filename, html)
	if status == true {
		logger.Info("save content to ask-mip-html: ", filename)
	} else {
		logger.Error("[error] to ask-mip-html error ")
	}
	// }
}

func getCommentListHtml(ans_id int, db *sql.DB) string {
	var s string = ""
	comments := LoadCommentsById(ans_id, db)
	if len(comments) > 0 {
		for _, v := range comments {
			content := filterContent(v.Content)
			userinfo := LoadUserinfoByUid(v.Uid, db)
			s += "<li><div class=\"com-head\"><div class=\"cavatar\"><a href=\"http://i.goumin.com/user/" + strconv.Itoa(v.Uid) + "\"><img src=\"" + userinfo.Avatar + "\" alt=\"\"></a></div><div class=\"c-head\"><p><span class=\"c-user\"><a href=\"http://i.goumin.com/user/" + strconv.Itoa(v.Uid) + "\">" + userinfo.Nickname + "</a></span>"
			if v.Typeid == 2 {
				replyuser := LoadUserinfoByUid(v.Replyuid, db)
				s += "<span>回复</span><span class=\"c-user\">" + replyuser.Nickname + "</span>"
			}
			s += "</p></div></div><div class=\"com-con\">" + content + "</div><div class=\"com-bottom\"><p><span class=\"com-time\">" + v.Created + "</span></p></div></li>"
		}
	}
	return s
}

func imageHtml(images string) string {
	s := ""
	if images != "" {
		imgArr := strings.Split(images, ",")
		if len(imgArr) == 0 {
			for _, img := range imgArr {
				s += "<mip-img src=\"" + img + "\" ></mip-img>"
			}
		}
	}
	return s
}

func (e *InfoAsk) getH5TemplateHtml() string {
	html := ""
	templatefile := e.templatefile
	if checkFileIsExist(templatefile) {
		fi, err := os.Open(templatefile)
		if err != nil {
			check(err)
		}
		defer fi.Close()
		fd, err := ioutil.ReadAll(fi)
		html = string(fd)
	} else {
		logger.Error("template file not found")
		return ""
	}
	return html
}

func (e *InfoAsk) saveContentToHtml(urlname string, content string) bool {
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

func substr(s string, l int) string {
	s = filterContent(s)
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	s = re.ReplaceAllString(s, "")
	re, _ = regexp.Compile("\\s")
	s = re.ReplaceAllString(s, "")
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}

		if sl+rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}

func createAskFileName(id int) string {
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
	filename = dir + strconv.Itoa(id) + ".html"
	return filename
}

func (e *InfoAsk) cmsPage(totalpages int, id int, i int) string {
	cmsPage := ""
	if totalpages == 1 {
		cmsPage = ""
	}
	if totalpages == 2 {
		if i == 1 {
			cmsPage = "<a href=\"" + e.domain + strconv.Itoa(id) + "-2-1.html\">下一页</a> <a href=\"" + e.domain + strconv.Itoa(id) + "-2-1.html\">尾页</a>"
		} else {
			cmsPage = "<a href=\"" + e.domain + strconv.Itoa(id) + "-1-1.html\">首页</a> <a href=\"" + e.domain + strconv.Itoa(id) + "-1-1.html\">上一页</a>"
		}
	}
	if totalpages > 2 {
		if i == 1 {
			cmsPage = "<a href=\"" + e.domain + strconv.Itoa(id) + "-2-1.html\">下一页</a> <a href=\"" + e.domain + strconv.Itoa(id) + "-2-1.html\">尾页</a>"
		} else if i == totalpages {
			cmsPage = "<a href=\"" + e.domain + strconv.Itoa(id) + "-1-1.html\">首页</a> <a href=\"" + e.domain + strconv.Itoa(id) + "-" + strconv.Itoa(i-1) + "-1.html\">上一页</a>"
		} else {
			cmsPage = "<a href=\"" + e.domain + strconv.Itoa(id) + "-1-1.html\">首页</a><a href=\"" + e.domain + strconv.Itoa(id) + "-" + strconv.Itoa(i-1) + "-1.html\">上一页</a> <a href=\"" + e.domain + strconv.Itoa(id) + "-" + strconv.Itoa(i+1) + "-1.html\">下一页</a><a href=\"" + e.domain + strconv.Itoa(id) + "-" + strconv.Itoa(totalpages) + "-1.html\">尾页</a>"
		}
	}
	return cmsPage
}

func filterContent(content string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	content = re.ReplaceAllStringFunc(content, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	content = re.ReplaceAllString(content, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	content = re.ReplaceAllString(content, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	content = re.ReplaceAllString(content, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	content = re.ReplaceAllString(content, "\n")
	return content
}

func threadByAsk(id int, db *sql.DB, session *mgo.Session) string {
	threads := LoadRelateThreadByAsk(id, db, session)
	if threads == nil {
		logger.Info("relate thread by ask not found askid=", id)
		return ""
	}
	content := ""
	for _, v := range threads {
		if v.Views < 3000 {
			v.Views = rand.Intn(5000)
		}
		content += "<a href=\"" + mipBbsUrl + "thread-" + strconv.Itoa(v.Tid) + "-1-1.html\" class=\"relate-a\"><span class=\"subj\">" + v.Subject + "</span><span class=\"seenum\">" + strconv.Itoa(v.Views) + "浏览</span></a>"
	}
	return content
}

func (e *InfoAsk) askByAsk(id int, pid int, db *sql.DB, session *mgo.Session) string {
	asks := LoadRelateAskByAsk(id, pid, db, session)
	if asks == nil {
		logger.Error("relate ask data by ask not found askid=", id)
		return ""
	}
	content := ""
	for _, v := range asks {
		content += "<a href=\"" + e.domain + strconv.Itoa(v.Id) + ".html\" class=\"relate-a\"><span class=\"subj\">" + v.Subject + "</span><span class=\"seenum\">" + strconv.Itoa(v.Views) + "浏览</span></a>"
	}
	return content
}

func dogsByAsk(id int, pid int, db *sql.DB) string {
	dogs := LoadRelateDogByAsk(id, pid, db)
	if dogs == nil {
		logger.Error("relate dog data by ask not found askid=", id)
		return ""
	}
	content := ""
	for _, v := range dogs {
		content += "<li><a href=\"http://dog.m.goumin.com/pet/" + strconv.Itoa(v.Speid) + "\" class=\"relate-pet-a\"><mip-img  src=\"http://c1.cdn.goumin.com/cms" + v.Img + "\" alt=\"" + v.Spename + "\" class=\"relate-img\"></mip-img><span class=\"relate-pet-seenum\">" + v.Spename + "</span></a></li>"
	}
	return content
}

func DefaultDoctors(db *sql.DB) string {
	doctors := LoadHealthDoctor(db)
	if doctors == nil {
		return ""
	}
	var s string = ""
	for _, v := range doctors {
		s += "<dl><a href=\"http://a.app.qq.com/o/simple.jsp?pkgname=com.goumin.forum\" class=\"doctor-avatar ui-link\"><dt><mip-img src=\"" + v.Avatar + "\"><em>" + v.Name + "</em></dt><dd>" + v.Hospital + "</dd></a></dl>"
	}
	return s
}
