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
	Status   int
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

type UserInfo struct {
	Uid        int
	Username   string
	Nickname   string
	Grouptitle string
	Avatar     string
}

type AttachmentX struct {
	Aid        int
	Attachment string
	Thumb      string
	Medium     string
	Small      string
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
	rows, err := db.Query("select pid,tid,first,author,authorid,subject,dateline,message,status from `" + tableName + "` where invisible=0 and first=1 and tid=" + strconv.Itoa(int(tid)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_post first sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row = new(Post)
		rows.Scan(&row.Pid, &row.Tid, &row.First, &row.Author, &row.Authorid, &row.Subject, &row.Dateline, &row.Message, &row.Status)
		return row
	}
	return &Post{}
}

func LoadPostsByTid(tid int, posttableid int, db *sql.DB) []*Post {
	tableName := "pre_forum_post_" + strconv.Itoa(posttableid)
	rows, err := db.Query("select pid,tid,first,author,authorid,subject,dateline,message,status from `" + tableName + "` where invisible=0 and tid=" + strconv.Itoa(int(tid)) + " order by dateline")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_post sql prepare error: ", err)
		return nil
	}
	var rowsData []*Post
	for rows.Next() {
		var row = new(Post)
		rows.Scan(&row.Pid, &row.Tid, &row.First, &row.Author, &row.Authorid, &row.Subject, &row.Dateline, &row.Message, &row.Status)
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

func LoadUserinfoByUid(uid int, db *sql.DB) *UserInfo {
	tableName := "pre_common_member"
	rows, err := db.Query("select username,grouptitle from `" + tableName + "` as a left join pre_common_usergroup as b on a.groupid=b.groupid where uid=" + strconv.Itoa(int(uid)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_common_usergroup sql prepare error: ", err)
		return nil
	}
	row1, err := db.Query("SELECT mem_nickname FROM `dog_member` WHERE `uid`=" + strconv.Itoa(int(uid)))
	defer row1.Close()
	if err != nil {
		logger.Error("[error] check dog_member sql prepare error: ", err)
		return nil
	}
	var avatar string
	row2, err := db.Query("SELECT image FROM `album` WHERE type=5 and `uid`=" + strconv.Itoa(int(uid)))
	defer row2.Close()
	if err != nil {
		logger.Error("[error] check album type=5 sql prepare error: ", err)
		return nil
	}
	for row2.Next() {
		row2.Scan(&avatar)
	}
	if avatar == "" {
		row3, err := db.Query("SELECT image FROM `album` WHERE type=25 and `uid`=" + strconv.Itoa(int(uid)))
		defer row3.Close()
		if err != nil {
			logger.Error("[error] check album type=25 sql prepare error: ", err)
			return nil
		}
		for row3.Next() {
			row3.Scan(&avatar)
		}
	}
	if avatar == "" {
		avatar = "head/cover-s.jpg"
	}

	var username string
	var grouptitle string
	var nickname string
	for rows.Next() {
		rows.Scan(&username, &grouptitle)
	}
	for row1.Next() {
		row1.Scan(&nickname)
	}
	var row = new(UserInfo)
	row.Uid = uid
	row.Username = username
	row.Nickname = nickname
	row.Grouptitle = grouptitle
	row.Avatar = diaryDomain + avatar
	return row
}

func LoadAttachmentByPid(tid int, pid int, db *sql.DB) []*AttachmentX {
	i := tid % 10
	tableName := "pre_forum_attachment_" + strconv.Itoa(i)
	rows, err := db.Query("select aid,attachment,mobile_thumb,mobile_medium,mobile_small from `" + tableName + "` where isimage in(-1,1) and pid=" + strconv.Itoa(pid) + " and tid=" + strconv.Itoa(int(tid)))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_post sql prepare error: ", err)
		return nil
	}
	var rowsData []*AttachmentX
	for rows.Next() {
		var row = new(AttachmentX)
		rows.Scan(&row.Aid, &row.Attachment, &row.Thumb, &row.Medium, &row.Small)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

//SELECT * FROM `pre_common_member` WHERE `uid` IN (57172, 1)
//SELECT * FROM `pre_common_usergroup` WHERE `groupid`=1
//SELECT * FROM `dog_member` WHERE `uid`=57172
//SELECT * FROM `dog_head_image` WHERE head_userid=57172
//select * from album where type=5 and uid=57172\G;
//select * from album where type=25 and uid=57172\G;
