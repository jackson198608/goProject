package main

import (
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
	// "reflect"
	"strconv"
)

type Thread struct {
	Tid         int
	Fid         int
	Posttableid int
	Typeid      int
	Author      string
	Authorid    int
	Subject     string
	Dateline    int
	Views       int
	Replies     int
}

type Post struct {
	Pid      int
	Tid      int
	First    int
	Author   string
	Authorid int
	Subject  string
	Dateline int
	Message  string
}

type Forum struct {
	Fid        int
	Threadtype string
	Name       string
}

type Relatelink struct {
	Name string
	Url  string
}

func LoadThreadByTid(tid int, db *sql.DB) *Thread {
	tableName := "pre_forum_thread"
	rows, err := db.Query("select tid,fid,posttableid,typeid,author,authorid,subject,dateline,views,replies from `" + tableName + "` where displayorder in(0,1) and tid=" + strconv.Itoa(int(tid)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_thread sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row = new(Thread)
		rows.Scan(&row.Tid, &row.Fid, &row.Posttableid, &row.Typeid, &row.Author, &row.Authorid, &row.Subject, &row.Dateline, &row.Views, &row.Replies)
		return row
	}
	return &Thread{}
}

func LoadFirstPostByTid(tid int, posttableid int, db *sql.DB) *Post {
	tableName := "pre_forum_post_" + strconv.Itoa(posttableid)
	rows, err := db.Query("select pid,tid,first,author,authorid,subject,dateline,message from `" + tableName + "` where invisible=0 and first=1 and tid=" + strconv.Itoa(int(tid)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_post first sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row = new(Post)
		rows.Scan(&row.Pid, &row.Tid, &row.First, &row.Author, &row.Authorid, &row.Subject, &row.Dateline, &row.Message)
		return row
	}
	return &Post{}
}

func LoadPostsByTid(tid int, posttableid int, db *sql.DB) []*Post {
	tableName := "pre_forum_post_" + strconv.Itoa(posttableid)
	rows, err := db.Query("select pid,tid,first,author,authorid,subject,dateline,message from `" + tableName + "` where invisible=0 and tid=" + strconv.Itoa(int(tid)) + " order by dateline")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_post sql prepare error: ", err)
		return nil
	}
	var rowsData []*Post
	for rows.Next() {
		var row = new(Post)
		rows.Scan(&row.Pid, &row.Tid, &row.First, &row.Author, &row.Authorid, &row.Subject, &row.Dateline, &row.Message)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

func LoadForumByFid(fid int, typeid int, db *sql.DB) *Forum {
	rows, err := db.Query("SELECT ff.fid, f.name,ff.threadtypes FROM pre_forum_forum f LEFT JOIN pre_forum_forumfield ff ON ff.fid=f.fid WHERE f.fid=" + strconv.Itoa(int(fid)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check event_log sql prepare error: ", err)
		return nil
	}
	var row = new(Forum)
	for rows.Next() {
		rows.Scan(&row.Fid, &row.Name, &row.Threadtype)
	}

	decoder := php_serialize.NewUnSerializer(row.Threadtype)
	result, err := decoder.Decode()
	if err != nil {
		logger.Info(err)
		return nil
	}
	real_result, ok := result.(php_serialize.PhpArray)
	var rs interface{}
	if ok {
		real_types_result, ok := real_result["types"].(php_serialize.PhpArray)
		if ok {
			rs = real_types_result[typeid]
			// fmt.Println(reflect.TypeOf(rs))
			value, _ := rs.(string)
			row.Threadtype = value
			return row
		}
	}
	return &Forum{}

}

func LoadRelateLink(db *sql.DB) []*Relatelink {
	tableName := "pre_common_relatedlink"
	rows, err := db.Query("select name,url from `" + tableName + "`")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_common_relatedlink sql prepare error: ", err)
		return nil
	}
	var rowsData []*Relatelink
	for rows.Next() {
		var row = new(Relatelink)
		rows.Scan(&row.Name, &row.Url)
		rowsData = append(rowsData, row)
	}
	return rowsData
}
