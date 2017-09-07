package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type InfoNew struct {
	db           *sql.DB
	id           int
	typeid       int
	session      *mgo.Session
	templateType string
	templatefile string
	saveDir      string
	tidStart     string
	tidEnd       string
}

func NewInfo(logLevel int, id int, typeid int, db *sql.DB, session *mgo.Session, taskNewArgs []string) *InfoNew {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(InfoNew)
	e.db = db
	e.id = id
	e.typeid = typeid
	e.session = session
	e.templateType = taskNewArgs[3]
	e.templatefile = taskNewArgs[4]
	e.saveDir = taskNewArgs[5]
	e.tidStart = taskNewArgs[6]
	e.tidEnd = taskNewArgs[7]
	return e
}

func (e *InfoNew) CreateThreadHtmlContent(tid int) error {
	thread := LoadThreadByTid(tid, e.db)
	if thread.Tid <= 0 {
		logger.Info("thread is not exist")
		return nil
	}
	//相关帖子 eg:tid=12
	relateThread := relateThread(tid, thread.Fid, e.db, e.session)
	//相关问答 eg:tid=12
	relateAsk := relateAsk(tid, e.db, e.session)
	//相关犬种 eg:tid=4682521
	relateDogs := relateDogs(tid, e.db, e.session, e.templateType)

	posts := LoadPostsByTid(tid, thread.Posttableid, e.db)
	forum := LoadForumByFid(thread.Fid, thread.Typeid, e.db)
	firstpost := LoadFirstPostByTid(tid, thread.Posttableid, e.db)
	relatelink := LoadRelateLink(e.db)

	e.groupContentToSaveHtml(tid, e.templateType, thread, posts, forum, firstpost, relatelink, relateThread, relateAsk, relateDogs, e.db)
	// saveContentToHtml(content, tid, page)
	return nil
}

func (e *InfoNew) saveContentToHtml(urlname string, content string, tid int, page int) bool {
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

func relateThread(tid int, fid int, db *sql.DB, session *mgo.Session) string {
	threads := LoadRelateThread(tid, fid, db, session)
	content := ""
	for _, v := range threads {
		if v.Views < 3000 {
			v.Views = rand.Intn(5000)
		}
		content += "<a href=\"/bbs/thread-" + strconv.Itoa(v.Tid) + "-1-1.html\" class=\"relate-a\"><span class=\"subj\">" + v.Subject + "</span><span class=\"seenum\">" + strconv.Itoa(v.Views) + "浏览</span></a>"
	}
	return content
}

func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}

func relateAsk(tid int, db *sql.DB, session *mgo.Session) string {
	asks := LoadRelateAsk(tid, db, session)
	content := ""
	for _, v := range asks {
		content += "<a href=\"/ask/" + strconv.Itoa(v.Id) + ".html\" class=\"relate-a\"><span class=\"subj\">" + v.Subject + "</span><span class=\"seenum\">" + strconv.Itoa(v.Views) + "浏览</span></a>"
	}
	return content
}

func relateDogs(tid int, db *sql.DB, session *mgo.Session, templateType string) string {
	dogs := LoadRelateDog(tid, db, session)
	content := ""
	for k, v := range dogs {
		if k <= 5 {
			if templateType == "1" {
				content += "<li><a href=\"http://dog.m.goumin.com/pet/" + strconv.Itoa(v.Speid) + "\" class=\"relate-pet-a\"><img src=\"http://c1.cdn.goumin.com/cms" + v.Img + "\" alt=\"" + v.Spename + "\" class=\"relate-img\" /><span class=\"relate-pet-seenum\">" + v.Spename + "</span></a></li>"
			} else {
				content += "<li><a href=\"http://dog.m.goumin.com/pet/" + strconv.Itoa(v.Speid) + "\" class=\"relate-pet-a\"><mip-img  src=\"http://c1.cdn.goumin.com/cms" + v.Img + "\" alt=\"" + v.Spename + "\" class=\"relate-img\"></mip-img><span class=\"relate-pet-seenum\">" + v.Spename + "</span></a></li>"
			}
		}
	}
	return content
}

func (e *InfoNew) groupContentToSaveHtml(tid int, templateType string, thread *Thread, posts []*Post, forum *Forum, firstpost *Post, relatelink []*Relatelink, relateThread string, relateAsk string, relateDogs string, db *sql.DB) {
	var subject = thread.Subject
	var url = staticH5Url + "thread-" + strconv.Itoa(tid) + "-1-1.html"
	var threadUrl = "thread-" + strconv.Itoa(tid) + "-1-1.html"
	var forumUrl = "forum-" + strconv.Itoa(thread.Fid) + "-1.html"

	var subMessage = show_substr(firstpost.Message, 160)
	var keyword = forum.Name + "，" + subject + subMessage + " ..."
	var description = subMessage + " ... " + subject + " ,狗民网｜铃铛宠物App"
	var views int = 0
	html := e.getH5TemplateHtml()
	if html == "" {
		logger.Error("template file not found")
		return
	}
	html = strings.Replace(html, "cmsRand", strconv.Itoa(rand.Intn(3000)), -1)
	html = strings.Replace(html, "cmsViews", strconv.Itoa(views), -1)
	html = strings.Replace(html, "cmsSubject", subject, -1)
	html = strings.Replace(html, "cmsKeywords", keyword, -1)
	html = strings.Replace(html, "cmsLink", url, -1)
	html = strings.Replace(html, "cmsDescription", description, -1)
	html = strings.Replace(html, "cmsThreadUrl", threadUrl, -1)
	html = strings.Replace(html, "cmsForumUrl", forumUrl, -1)
	html = strings.Replace(html, "cmsForumName", forum.Name, -1)
	html = strings.Replace(html, "cmsTypeName", forum.Threadtype, -1)
	html = strings.Replace(html, "cmsRelateThread", relateThread, -1)
	html = strings.Replace(html, "cmsRelateAsk", relateAsk, -1)
	html = strings.Replace(html, "cmsRelateDogs", relateDogs, -1)
	len := len(posts)                                           //帖子楼层总数
	count := 20                                                 //每页条数
	totalpages := int(math.Ceil(float64(len) / float64(count))) //page总数
	for i := 1; i <= totalpages; i++ {
		var title = subject + " - 第" + strconv.Itoa(i) + "页 - " + forum.Name + " -  狗民社区-移动版"
		html = strings.Replace(html, "cmsTitle", title, -1)
		cmsPage := ""
		content := ""
		// dir := strconv.Itoa(tid % 1000)
		// dir := ""
		// if tid < 5000000 { //tid<5百万的数据，生成目录4/tid%1000/
		// 	dir = "4/" + strconv.Itoa(tid%1000)

		// } else { //tid>5百万，每增加1百万，生成目录/tid%1百万/tid%500个
		// 	dir = strconv.Itoa(tid/1000000) + "/" + strconv.Itoa(tid%500)
		// }
		// filename := dir + "/thread-" + strconv.Itoa(tid) + "-" + strconv.Itoa(i) + "-1.html"
		filename := createFileName(tid, i, 0)
		start := (i - 1) * count
		end := start + count
		if end > len {
			end = len
		}
		pagepost := posts[start:end]
		floor := count * (i - 1) //楼层
		var images []*AttachmentX
		var message string = ""
		for _, v := range pagepost {
			message = regexp_string(v.Message)
			floor++
			for _, lv := range relatelink {
				replace := "<a href='" + lv.Url + "'>" + lv.Name + "</a>"
				message = strings.Replace(message, lv.Name, replace, -1)
			}
			userinfo := LoadUserinfoByUid(v.Authorid, db)
			tm := time.Unix(int64(v.Dateline), 0)
			dateline := tm.Format("2006-01-02 15:04:05")
			images = LoadAttachmentByPid(tid, v.Pid, db)
			message = replaceImgOrAttach(message, subject, images)
			if templateType == "1" {
				content += "<div class=\"post-detail-a\"><div class=\"user-info\"><a href=\"javascript:;\"><img src=\"" + userinfo.Avatar + "\" alt=\"" + userinfo.Nickname + "\"><span class=\"info\"><em class=\"user-name\">" + userinfo.Nickname + "</em><em class=\"level\">" + userinfo.Grouptitle + "</em>"
				if v.First == 1 {
					content += "<em  class=\"identity\">楼主</em>"
				}
				content += "</span><span class=\"dataTime\">" + dateline + "</span><span class=\"floor\">" + strconv.Itoa(floor) + "楼</span></a></div><div class=\"post-detail-content\"><p>" + message + "</p></div></div>"
			} else {
				content += "<div class=\"post-detail-a\"><div class=\"user-info\"><a href=\"javascript:;\"><mip-img src=\"" + userinfo.Avatar + "\" class=\"user-avatar\"></mip-img><span class=\"info\"><em class=\"user-name\">" + userinfo.Nickname + "</em>"
				if v.First == 1 {
					content += "<em class=\"identity\">楼主</em>"
				}
				content += "</span><span class=\"dataTime\">" + userinfo.Grouptitle + "</span></a></div><div class=\"post-detail-content\"><div><p>" + message + "</p></div><div class=\"detail-date\"><span>" + strconv.Itoa(floor) + "楼</span><span>" + dateline + "</span></div></div></div>"
			}
		}
		if totalpages == 1 {
			cmsPage = ""
		}
		if totalpages == 2 {
			if i == 1 {
				cmsPage = "<a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-2-1.html\">下一页</a> <a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-2-1.html\">尾页</a>"
			} else {
				cmsPage = "<a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-1-1.html\">首页</a> <a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-1-1.html\">上一页</a>"
			}
		}
		if totalpages > 2 {
			if i == 1 {
				cmsPage = "<a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-2-1.html\">下一页</a> <a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-2-1.html\">尾页</a>"
			} else if i == totalpages {
				cmsPage = "<a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-1-1.html\">首页</a> <a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-" + strconv.Itoa(i-1) + "-1.html\">上一页</a>"
			} else {
				cmsPage = "<a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-1-1.html\">首页</a><a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-" + strconv.Itoa(i-1) + "-1.html\">上一页</a> <a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-" + strconv.Itoa(i+1) + "-1.html\">下一页</a><a href=\"/bbs/thread-" + strconv.Itoa(tid) + "-" + strconv.Itoa(totalpages) + "-1.html\">尾页</a>"
			}
		}
		htmlhtml := strings.Replace(html, "cmsPost", strconv.Itoa(len-1), -1)
		content = findface(content)
		htmlhtml = strings.Replace(htmlhtml, "cmsMessage", content, -1)
		htmlhtml = strings.Replace(htmlhtml, "cmsPage", cmsPage, -1)
		status := e.saveContentToHtml(filename, htmlhtml, tid, i)
		if status == true {
			logger.Info("save content to html: ", filename)
		} else {
			logger.Error("[error] to html error ")
		}
	}
}

func createFileName(tid int, page int, typeid int) string {
	filename := ""
	dir := ""
	if typeid == 1 {
		if tid < 5000000 { //tid<5百万的数据，生成目录4/tid%1000/
			dir = "4/" + strconv.Itoa(tid%1000)

		} else { //tid>5百万，每增加1百万，生成目录/tid%1百万/tid%500个
			dir = strconv.Itoa(tid/1000000) + "/" + strconv.Itoa(tid%500)
		}
	} else {
		n4 := tid % 10               //个位数
		n3 := (tid - n4) % 100       //十位数
		n2 := (tid - n4 - n3) % 1000 //百位数
		dir = strconv.Itoa(n2/100) + "/" + strconv.Itoa(n3/10) + "/" + strconv.Itoa(n4)
	}
	filename = dir + "/thread-" + strconv.Itoa(tid) + "-" + strconv.Itoa(page) + "-1.html"
	return filename
}

func replaceImgOrAttach(content string, subject string, images []*AttachmentX) string {
	re, _ := regexp.Compile("\\[img(.*?)\\](.*?)\\[\\/img\\]")
	content = re.ReplaceAllString(content, "<img alt=\""+subject+"\" src=\"$2\">")
	len := len(images)
	if len == 0 {
		return content
	}
	if strings.Contains(content, "[attach]") == true && strings.Contains(content, "[/attach]") == true {
		count := strings.Count(content, "[attach]")
		if count == 0 {
			return content
		}
		result := content
		reg := regexp.MustCompile("\\[attach\\](\\d+)\\[\\/attach\\]")
		str := reg.FindAllStringSubmatch(content, -1)
		for i := 0; i < count; i++ {
			result = changeMessage(result, i, images, str)
		}
		return result
	} else {
		for _, v := range images {
			if v.Thumb == "" {
				content += "<img class='imgs' src='" + bbsDomain + v.Attachment + "' />"
			} else {
				content += "<img class='imgs' src='" + bbsDomain + v.Thumb + "' />"
			}
		}
	}
	return content
}

func changeMessage(content string, i int, images []*AttachmentX, str [][]string) string {
	if len(str) == 0 || len(images) == 0 {
		return content
	}
	img := ""
	for _, v := range images {
		aid, _ := strconv.Atoi(str[i][1])
		if v.Aid == aid {
			if v.Thumb == "" {
				img = "<img class='imgs' src='" + bbsDomain + v.Attachment + "' />"
			} else {
				img = "<img class='imgs' src='" + bbsDomain + v.Thumb + "' />"
			}
		}
	}
	content = strings.Replace(content, str[i][0], img, -1)
	return content
}

func (e *InfoNew) getH5TemplateHtml() string {
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

func show_substr(s string, l int) string {
	s = regexp_string(s)
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

func regexp_string(content string) string {

	m := strings.Index(content, "[quote]")
	n := strings.Index(content, "[/quote]")
	if m >= 0 && n > 0 {
		substr := content[m:n]
		content = strings.Replace(content, substr+"[/quote]", "", -1)
	}
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	content = re.ReplaceAllStringFunc(content, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	content = re.ReplaceAllString(content, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	content = re.ReplaceAllString(content, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	// re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	// content = re.ReplaceAllString(content, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	content = re.ReplaceAllString(content, "\n")

	re, _ = regexp.Compile("\\[hr\\]")
	content = re.ReplaceAllString(content, "<hr>")

	re, _ = regexp.Compile("\\[size\\=(.*?)\\](.*?)")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[font\\=(.*?)\\](.*?)")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[color\\=(.*?)\\](.*?)")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[backcolor\\=(.*?)\\](.*?)")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/size\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/font\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/color\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/backcolor\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[list=(\\d+)\\]")
	content = re.ReplaceAllString(content, "<ul class='numeric'>")

	re, _ = regexp.Compile("\\[list\\]")
	content = re.ReplaceAllString(content, "<ul class='dot'>")

	re, _ = regexp.Compile("\\[/list\\]")
	content = re.ReplaceAllString(content, "</ul>")

	re, _ = regexp.Compile("\\[img(.*?)\\](.*?)\\[/img\\]")
	content = re.ReplaceAllString(content, "<div><img class='post_content_image' src='$2' /></div>")

	re, _ = regexp.Compile("\\[url=.*?goto=findpost&pid=\\d+&ptid=\\d+\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[url=(http://bbs.goumin.com/)?forum-\\d+-\\d+\\.html\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "$2")

	re, _ = regexp.Compile("\\[url=home.php\\?mod=space&uid=\\d+\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "$1")

	re, _ = regexp.Compile("\\[url=(http.*?)\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "<a href='$1'>$2</a>")

	re, _ = regexp.Compile("\\[url\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "<a href='$1'>$1</a>")

	re, _ = regexp.Compile("\\[url=(mailto:.*?)\\](.*?)\\[/url\\]")
	content = re.ReplaceAllString(content, "<a href='$1'>$2</a>")

	// re, _ = regexp.Compile("(http:\\/\\/bbs\\.goumin\\.com\\/thread-\\d+-\\d+-\\d+.html)")
	// content = re.ReplaceAllString(content, "<a href='$1'>$2</a>")

	//代码 code
	re, _ = regexp.Compile("\\[code\\](.*?)\\[/code\\]")
	content = re.ReplaceAllString(content, "<div class='code'><pre>$1</pre></div>")
	// 音视频 audio
	re, _ = regexp.Compile("\\[audio\\](.*?)\\[/audio\\]")
	content = re.ReplaceAllString(content, "<audio src='$1' controls></audio>")

	re, _ = regexp.Compile("\\[audio(.*?)\\](.*?)\\[/audio\\]")
	content = re.ReplaceAllString(content, "<audio src='$2' controls></audio>")

	// 视频 media
	re, _ = regexp.Compile("\\[media\\](.*?\\.mp4)\\[/media\\]")
	content = re.ReplaceAllString(content, "<video src='$1' controls></video>")

	re, _ = regexp.Compile("\\[media(.*?)\\](.*?)\\[/media\\]")
	content = re.ReplaceAllString(content, "<video src='$1' controls></video>")

	re, _ = regexp.Compile("\\[media.*?\\]http:\\/\\/v\\.youku\\.com\\/v_show\\/id_(.*?)\\.html.*?\\[/media\\]")
	content = re.ReplaceAllString(content, "<a class='post_content_link' href='http://v.youku.com/v_show/id_$1.html'>***优酷视频点击播放***</a>")

	re, _ = regexp.Compile("\\[media.*?\\]http:/\\/player\\.youku\\.com\\/player\\.php\\/sid\\/(.*?)\\/v\\.swf\\[/media\\]")
	content = re.ReplaceAllString(content, "<a class='post_content_link' href='http://v.youku.com/v_show/id_$1.html'>***优酷视频点击播放***</a>")

	re, _ = regexp.Compile("\\[media.*?\\]http:\\/\\/www\\.tudou\\.com\\/programs\\/view\\/(.*?)\\/\\[/media\\]")
	content = re.ReplaceAllString(content, "<a class='post_content_link' href='http://www.tudou.com/programs/view/$1/'>***土豆视频点击播放***</a>")

	re, _ = regexp.Compile("\\[media\\](.*?)\\[/media\\]")
	content = re.ReplaceAllString(content, "")

	// 表格 table td tr
	re, _ = regexp.Compile("\\[table.*?\\]")
	content = re.ReplaceAllString(content, "<table class='post_table'>")

	re, _ = regexp.Compile("\\[\\/table\\]")
	content = re.ReplaceAllString(content, "</table>")

	re, _ = regexp.Compile("\\[tr(=.*?)?\\]")
	content = re.ReplaceAllString(content, "<tr>")

	re, _ = regexp.Compile("\\[\\/tr\\]")
	content = re.ReplaceAllString(content, "</tr>")

	re, _ = regexp.Compile("\\[td(=.*?)?\\]")
	content = re.ReplaceAllString(content, "<td>")

	re, _ = regexp.Compile("\\[\\/td\\]")
	content = re.ReplaceAllString(content, "</td>")

	re, _ = regexp.Compile("\\[\\/td\\]")
	content = re.ReplaceAllString(content, "</td>")

	re, _ = regexp.Compile("\\[\\*\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[flash\\](.*?)\\[/flash\\]")
	content = re.ReplaceAllString(content, "<embed width=\"330\" height=\"240\" allownetworking=\"internal\" allowscriptaccess=\"never\" src='$1' quality=\"high\" bgcolor=\"#ffffff\" wmode=\"transparent\" allowfullscreen=\"true\" type=\"application/x-shockwave-flash\">")

	// re, _ = regexp.Compile("\r\n")
	// content = re.ReplaceAllString(content, "")

	// re, _ = regexp.Compile("\r")
	// content = re.ReplaceAllString(content, "")

	// re, _ = regexp.Compile("\n")
	// content = re.ReplaceAllString(content, "")

	// 左右对齐align，浮动 float
	re, _ = regexp.Compile("\\[align=(.*?)\\](.*?)\\[/align\\]")
	content = re.ReplaceAllString(content, "<div align='$1'><p>$2</p></div>")

	re, _ = regexp.Compile("\\[align.*?\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/align\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[float=(.*?)\\](.*?)\\[/float\\]")
	content = re.ReplaceAllString(content, "<div style='float:$1'>$2</div>")

	re, _ = regexp.Compile("\\[float.*?\\]")
	content = re.ReplaceAllString(content, "")

	// re, _ = regexp.Compile("\\[quote\\][\n\r]*(.+?)[\n\r]*\\[\\/quote\\]")
	// content = re.ReplaceAllString(content, "<div class='quote'><blockquote>“$1”</blockquote></div>")

	re, _ = regexp.Compile("\\[free\\](.*?)\\[/free\\]")
	content = re.ReplaceAllString(content, "<div class='quote'><blockquote>“$1”</blockquote></div>")

	re, _ = regexp.Compile("\\[hide\\](.*?)\\[/hide\\]")
	content = re.ReplaceAllString(content, "<div class='quote'><blockquote>“$1”</blockquote></div>")

	re, _ = regexp.Compile("\\[qq\\](.*?)\\[/qq\\]")
	content = re.ReplaceAllString(content, "QQ:$1")

	// i
	re, _ = regexp.Compile("\\[i.*?\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/i\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[u.*?\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/u\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[b.*?\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/b\\]")
	content = re.ReplaceAllString(content, "")

	re, _ = regexp.Compile("\\[/url\\]")
	content = re.ReplaceAllString(content, "")

	return strings.TrimSpace(content)
}

//替换表情标签
func findface(content string) string {
	content = strings.Replace(content, ":cm101:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-01.gif'>", -1)
	content = strings.Replace(content, ":cm102:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-02.gif'>", -1)
	content = strings.Replace(content, ":cm103:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-03.gif'>", -1)
	content = strings.Replace(content, ":cm104:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-04.gif'>", -1)
	content = strings.Replace(content, ":cm105:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-05.gif'>", -1)
	content = strings.Replace(content, ":cm106:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-06.gif'>", -1)
	content = strings.Replace(content, ":cm107:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-07.gif'>", -1)
	content = strings.Replace(content, ":cm108:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-08.gif'>", -1)
	content = strings.Replace(content, ":cm109:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-09.gif'>", -1)
	content = strings.Replace(content, ":cm110:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-10.gif'>", -1)
	content = strings.Replace(content, ":cm111:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-11.gif'>", -1)
	content = strings.Replace(content, ":cm112:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-12.gif'>", -1)
	content = strings.Replace(content, ":cm113:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-13.gif'>", -1)
	content = strings.Replace(content, ":cm114:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-14.gif'>", -1)
	content = strings.Replace(content, ":cm115:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-15.gif'>", -1)
	content = strings.Replace(content, ":cm116:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-16.gif'>", -1)
	content = strings.Replace(content, ":cm117:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-17.gif'>", -1)
	content = strings.Replace(content, ":cm118:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-18.gif'>", -1)
	content = strings.Replace(content, ":cm119:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-19.gif'>", -1)
	content = strings.Replace(content, ":cm120:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-20.gif'>", -1)
	content = strings.Replace(content, ":cm121:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-21.gif'>", -1)
	content = strings.Replace(content, ":cm122:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-22.gif'>", -1)
	content = strings.Replace(content, ":cm123:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-23.gif'>", -1)
	content = strings.Replace(content, ":cm124:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-24.gif'>", -1)
	content = strings.Replace(content, ":cm125:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-25.gif'>", -1)
	content = strings.Replace(content, ":cm126:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-26.gif'>", -1)
	content = strings.Replace(content, ":cm127:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-27.gif'>", -1)
	content = strings.Replace(content, ":cm128:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-28.gif'>", -1)
	content = strings.Replace(content, ":cm129:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-29.gif'>", -1)
	content = strings.Replace(content, ":cm130:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-30.gif'>", -1)
	content = strings.Replace(content, ":cm131:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-31.gif'>", -1)
	content = strings.Replace(content, ":cm132:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-32.gif'>", -1)
	content = strings.Replace(content, ":cm133:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-33.gif'>", -1)
	content = strings.Replace(content, ":cm134:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-34.gif'>", -1)
	content = strings.Replace(content, ":cm135:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-35.gif'>", -1)
	content = strings.Replace(content, ":cm136:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-36.gif'>", -1)
	content = strings.Replace(content, ":cm137:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-37.gif'>", -1)
	content = strings.Replace(content, ":cm138:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-38.gif'>", -1)
	content = strings.Replace(content, ":cm139:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-39.gif'>", -1)
	content = strings.Replace(content, ":cm140:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-40.gif'>", -1)
	content = strings.Replace(content, ":cm141:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-41.gif'>", -1)
	content = strings.Replace(content, ":cm142:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-42.gif'>", -1)
	content = strings.Replace(content, ":cm143:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-43.gif'>", -1)
	content = strings.Replace(content, ":cm144:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-44.gif'>", -1)
	content = strings.Replace(content, ":cm145:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-45.gif'>", -1)
	content = strings.Replace(content, ":cm146:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-46.gif'>", -1)
	content = strings.Replace(content, ":cm147:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-47.gif'>", -1)
	content = strings.Replace(content, ":cm148:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-48.gif'>", -1)
	content = strings.Replace(content, ":cm149:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-49.gif'>", -1)
	content = strings.Replace(content, ":cm150:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-50.gif'>", -1)
	content = strings.Replace(content, ":cm151:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-51.gif'>", -1)
	content = strings.Replace(content, ":cm152:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-52.gif'>", -1)
	content = strings.Replace(content, ":cm153:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-53.gif'>", -1)
	content = strings.Replace(content, ":cm154:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-54.gif'>", -1)
	content = strings.Replace(content, ":cm155:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-55.gif'>", -1)
	content = strings.Replace(content, ":cm156:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-56.gif'>", -1)
	content = strings.Replace(content, ":cm157:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-57.gif'>", -1)
	content = strings.Replace(content, ":cm158:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-58.gif'>", -1)
	content = strings.Replace(content, ":cm159:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-59.gif'>", -1)
	content = strings.Replace(content, ":cm160:", "<img src='http://bbs.goumin.com/static/image/smiley/default/emot/emot-60.gif'>", -1)
	content = strings.Replace(content, ":dog61:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/61.gif'>", -1)
	content = strings.Replace(content, ":dog62:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/62.gif'>", -1)
	content = strings.Replace(content, ":dog63:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/63.gif'>", -1)
	content = strings.Replace(content, ":dog64:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/64.gif'>", -1)
	content = strings.Replace(content, ":dog65:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/65.gif'>", -1)
	content = strings.Replace(content, ":dog66:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/66.gif'>", -1)
	content = strings.Replace(content, ":dog67:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/67.gif'>", -1)
	content = strings.Replace(content, ":dog68:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/68.gif'>", -1)
	content = strings.Replace(content, ":dog69:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/69.gif'>", -1)
	content = strings.Replace(content, ":dog70:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/70.gif'>", -1)
	content = strings.Replace(content, ":dog71:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/71.gif'>", -1)
	content = strings.Replace(content, ":dog72:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/72.gif'>", -1)
	content = strings.Replace(content, ":dog73:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/73.gif'>", -1)
	content = strings.Replace(content, ":dog74:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/74.gif'>", -1)
	content = strings.Replace(content, ":dog75:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/75.gif'>", -1)
	content = strings.Replace(content, ":dog76:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/76.gif'>", -1)
	content = strings.Replace(content, ":dog77:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/77.gif'>", -1)
	content = strings.Replace(content, ":dog78:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/78.gif'>", -1)
	content = strings.Replace(content, ":dog79:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/79.gif'>", -1)
	content = strings.Replace(content, ":dog80:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/80.gif'>", -1)
	content = strings.Replace(content, ":dog81:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/81.gif'>", -1)
	content = strings.Replace(content, ":dog82:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/82.gif'>", -1)
	content = strings.Replace(content, ":dog83:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/83.gif'>", -1)
	content = strings.Replace(content, ":dog84:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/84.gif'>", -1)
	content = strings.Replace(content, ":dog85:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/85.gif'>", -1)
	content = strings.Replace(content, ":dog86:", "<img src='http://bbs.goumin.com/static/image/smiley/default/dog/86.gif'>", -1)
	return content
}
