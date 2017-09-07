package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
	// "reflect"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
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

type RelateThread struct {
	Tid      int
	Subject  string
	Views    int
	Dateline int
}

type RelateAsk struct {
	Id      int
	Subject string
	Views   int
}

type RelateDog struct {
	Speid   int
	Spename string
	Img     string
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
	if posttableid == 0 {
		tableName = "pre_forum_post"
	}
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
	if posttableid == 0 {
		tableName = "pre_forum_post"
	}
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
	var dir string = diaryDomain

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
		row3, err := db.Query("SELECT head_id,head_fileext,head_cdate FROM `dog_head_image` WHERE head_userid=" + strconv.Itoa(int(uid)) + " limit 1")
		defer row3.Close()
		if err != nil {
			logger.Error("[error] check album type=25 sql prepare error: ", err)
			return nil
		}
		var head_id int = 0
		var head_fileext string = ""
		var head_cdate int = 0
		for row3.Next() {
			row3.Scan(&head_id, &head_fileext, &head_cdate)
		}
		if head_id > 0 {
			sub0 := head_id
			sub1 := sub0 >> 8
			sub2 := sub1 >> 8
			sub3 := sub2 >> 8
			sub4 := sub3 >> 8
			time := time.Now().Unix()
			if head_cdate < int(time)-3600 {
				dir = "http://hd2.goumin.com"
			} else {
				dir = "http://www.goumin.com"
			}
			avatar = "/attachments/head/" + strconv.Itoa(sub4) + "/" + strconv.Itoa(sub3) + "/" + strconv.Itoa(sub2) + "/" + strconv.Itoa(sub1) + "/" + strconv.Itoa(head_id) + "s." + head_fileext
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
	row.Avatar = dir + avatar
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

type ThreadsRecommend struct {
	Id         int   "_id"
	Related    []int "related"
	RelatedAsk []int "related_ask"
}

func LoadRelateThread(tid int, fid int, db *sql.DB, session *mgo.Session) []*RelateThread {
	var rowsData []*RelateThread
	ms := new(ThreadsRecommend)
	c := session.DB("BigData").C("threads_recommend")
	err := c.Find(&bson.M{"_id": tid}).One(&ms)
	if err != nil {
		logger.Error("BigData threads_recommend relate thread: ", err)
	}
	if len(ms.Related) == 0 {
		rows, err := db.Query("select tid,subject,views,dateline from `pre_forum_thread` where displayorder>=0 and fid=" + strconv.Itoa(fid) + " order by dateline desc limit 5")
		defer rows.Close()
		if err != nil {
			logger.Error("[error] check pre_forum_post sql prepare error: ", err)
			return nil
		}
		for rows.Next() {
			var row = new(RelateThread)
			rows.Scan(&row.Tid, &row.Subject, &row.Views, &row.Dateline)
			rowsData = append(rowsData, row)
		}
	} else {
		tidstring := ""
		for k, v := range ms.Related {
			if k <= 4 {
				tidstring += strconv.Itoa(v) + ","
			}
		}
		rows, err := db.Query("select tid,subject,views,dateline from `pre_forum_thread` where displayorder>=0 and tid in (" + strings.Trim(tidstring, ",") + ") ")
		defer rows.Close()
		if err != nil {
			logger.Error("check pre_forum_post sql prepare error: ", err)
			return nil
		}
		for rows.Next() {
			var row = new(RelateThread)
			rows.Scan(&row.Tid, &row.Subject, &row.Views, &row.Dateline)
			rowsData = append(rowsData, row)
		}
	}
	return rowsData
}

func LoadRelateAsk(tid int, db *sql.DB, session *mgo.Session) []*RelateAsk {
	var rowsData []*RelateAsk
	ms := new(ThreadsRecommend)
	c := session.DB("BigData").C("threads_recommend")
	err := c.Find(&bson.M{"_id": tid}).One(&ms)
	if err != nil {
		logger.Error("BigData threads_recommend relate ask: ", err)
	}
	if len(ms.RelatedAsk) == 0 {
		rows, err := db.Query("select id,subject,browse_num from `ask`.`ask_question` where is_hide=1 order by ans_num desc limit 5")
		defer rows.Close()
		if err != nil {
			logger.Error("[error] check ask_question sql prepare error: ", err)
			return nil
		}
		for rows.Next() {
			var row = new(RelateAsk)
			rows.Scan(&row.Id, &row.Subject, &row.Views)
			rowsData = append(rowsData, row)
		}
	} else {
		idstring := ""
		for k, v := range ms.RelatedAsk {
			if k <= 4 {
				idstring += strconv.Itoa(v) + ","
			}
		}
		rows, err := db.Query("select id,subject,browse_num from `ask`.`ask_question` where id in (" + strings.Trim(idstring, ",") + ")")
		defer rows.Close()
		if err != nil {
			logger.Error("[error] check ask_question sql prepare error: ", err)
			return nil
		}
		for rows.Next() {
			var row = new(RelateAsk)
			rows.Scan(&row.Id, &row.Subject, &row.Views)
			rowsData = append(rowsData, row)
		}
	}
	return rowsData
}

type Catedog struct {
	Id    int
	Type  int
	Value int
}

func LoadRelateDog(tid int, db *sql.DB, session *mgo.Session) []*RelateDog {
	var doginfo string
	ms := new(Catedog)
	c := session.DB("BigData").C("cate_dog")
	err := c.Find(&bson.M{"id": tid, "type": 1}).One(&ms)
	if err != nil {
		logger.Error("BigData cate_dog relate dog: ", err)
		// return nil
		doginfo = "60,62,16,22,70,35"
	}
	rows, err := db.Query("select related from `related_pet` where spe_id=" + strconv.Itoa(ms.Value))
	defer rows.Close()
	if err != nil {
		logger.Error("check relate_pet sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row string
		rows.Scan(&row)
		doginfo += row + ","
	}
	doginfo = strings.Trim(doginfo, ",")
	if doginfo == "" {
		doginfo = "60,62,16,22,70,35"
	}
	fmt.Println("doginfo" + doginfo)
	rows1, err := db.Query("select dn.spe_id,spe_name_f,image_list from `dog_new_species` as dn left join dog_species as ds on dn.spe_id=ds.spe_id where dn.spe_id in(" + doginfo + ")")
	defer rows1.Close()
	if err != nil {
		logger.Error("check dog_new_species sql prepare error: ", err)
		return nil
	}
	var rowsData []*RelateDog
	for rows1.Next() {
		var row = new(RelateDog)
		rows1.Scan(&row.Speid, &row.Spename, &row.Img)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

func getThreadTask(page int) []int {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)

		return nil
	}
	defer db.Close()
	tableName := "pre_forum_thread"
	rows, err := db.Query("select tid from `" + tableName + "` where displayorder in(0,1) order by tid asc limit " + strconv.Itoa(offset) + " offset " + strconv.Itoa(offset*(page-1)))
	defer rows.Close()
	if err != nil {
		logger.Error("check pre_forum_thread sql prepare error: ", err)
		return nil
	}
	var a []int
	for rows.Next() {
		var tid int
		rows.Scan(&tid)
		a = append(a, tid)
	}
	return a
}

//SELECT * FROM `pre_common_member` WHERE `uid` IN (57172, 1)
//SELECT * FROM `pre_common_usergroup` WHERE `groupid`=1
//SELECT * FROM `dog_member` WHERE `uid`=57172
//SELECT * FROM `dog_head_image` WHERE head_userid=57172
//select * from album where type=5 and uid=57172\G;
//select * from album where type=25 and uid=57172\G;
