package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	// "ioutil"
	"os"
	"path"
	"strconv"
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
	pcount := len(posts)
	fmt.Println(forum.Threadtype)
	fmt.Println(pcount)
	var subject = thread.Subject
	// for k, v := range posts {

	// }
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
	n, err1 := f.WriteString(subject) //写入文件(字符串)
	check(err1)
	fmt.Printf("写入 %d 个字节n", n)
	return nil
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
