package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type InfoNew struct {
	db     *sql.DB
	id     int
	typeid int
}

func NewInfo(logLevel int, id int, typeid int, db *sql.DB) *InfoNew {
	logger.SetLevel(logger.LEVEL(logLevel))
	e := new(InfoNew)
	e.db = db
	e.id = id
	e.typeid = typeid
	return e
}

func (e *InfoNew) CreateThreadHtmlContent(tid int) error {
	thread := LoadThreadByTid(tid, e.db)
	posts := LoadPostsByTid(tid, thread.Posttableid, e.db)
	forum := LoadForumByFid(thread.Fid, thread.Typeid, e.db)
	firstpost := LoadFirstPostByTid(tid, thread.Posttableid, e.db)
	relatelink := LoadRelateLink(e.db)

	content := groupContent(tid, thread, posts, forum, firstpost, relatelink)
	saveContentToHtml(content, tid)
	return nil
}

func saveContentToHtml(content string, tid int) bool {
	var filename = saveDir + "thread-" + strconv.Itoa(tid) + "-1-1.html"
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

func groupContent(tid int, thread *Thread, posts []*Post, forum *Forum, firstpost *Post, relatelink []*Relatelink) string {
	var subject = thread.Subject
	var url = "http://bbs.goumin.com/thread-" + strconv.Itoa(tid) + "-1-1.html"
	var title = thread.Subject + " - " + forum.Name + " -  狗民网｜铃铛宠物App"
	var subMessage = show_substr(firstpost.Message, 160)
	var keyword = forum.Name + "，" + subject + subMessage + " ..."
	var description = subMessage + " ... " + subject + " ,狗民网｜铃铛宠物App"
	content := url + title + description + keyword
	html := getTemplateHtml()
	fmt.Println(html)

	return content
}

func getTemplateHtml() string {
	html := ""
	if checkFileIsExist(templatefile) {
		fi, err := os.Open(templatefile)
		if err != nil {
			check(err)
		}
		defer fi.Close()
		fd, err := ioutil.ReadAll(fi)
		html = string(fd)
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

func regexp_string(src string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "")

	return strings.TrimSpace(src)
}
