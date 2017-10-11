package main

import (
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"math"
	// "reflect"
	"strconv"
)

func getMaxId(jobType string) int {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)
		return 0
	}
	defer db.Close()
	var rows *sql.Rows
	if jobType == "asksave" {
		rows, err = db.Query("select id from `ask`.`ask_question` where is_hide=1 order by id desc limit 1")
	}
	if jobType == "threadsave" {
		rows, err = db.Query("select tid from `pre_forum_thread` where displayorder in(0,1) order by tid desc limit 1")
	}
	defer rows.Close()
	if err != nil {
		logger.Error("check ask_question or pre_forum_thread sql prepare error: ", err)
		return 0
	}
	var maxid int = 0
	for rows.Next() {
		rows.Scan(&maxid)
	}
	return maxid
}

func getAskList(startId int, endId int) []string {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)

		return nil
	}
	defer db.Close()
	// rows, err := db.Query("select id,ans_num from `ask`.`ask_question` where is_hide=1 order by id asc limit " + strconv.Itoa(offset) + " offset " + strconv.Itoa(offset*(page-1)))
	rows, err := db.Query("select id,ans_num from `ask`.`ask_question` where is_hide=1 and id>" + strconv.Itoa(startId) + " and id<=" + strconv.Itoa(endId) + " order by id asc")
	defer rows.Close()
	if err != nil {
		logger.Error("check ask_question sql prepare error: ", err)
		return nil
	}
	var a []string
	for rows.Next() {
		var id int
		var ans_num int
		rows.Scan(&id, &ans_num)
		totalpages := int(math.Ceil(float64(ans_num) / float64(5))) //page总数
		str := strconv.Itoa(id) + "|" + strconv.Itoa(totalpages)
		a = append(a, str)
	}
	return a
}

func getThreadTask(startId int, endId int) []string {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)

		return nil
	}
	defer db.Close()
	tableName := "pre_forum_thread"
	// rows, err := db.Query("select tid,posttableid from `" + tableName + "` where displayorder in(0,1) order by tid asc limit " + strconv.Itoa(offset) + " offset " + strconv.Itoa(offset*(page-1)))
	rows, err := db.Query("select tid,posttableid from `" + tableName + "` where displayorder in(0,1) and tid>" + strconv.Itoa(startId) + " and tid<=" + strconv.Itoa(endId) + " order by tid asc")
	defer rows.Close()
	if err != nil {
		logger.Error("check pre_forum_thread sql prepare error: ", err)
		return nil
	}
	var a []string
	for rows.Next() {
		var tid int
		var posttableid int
		rows.Scan(&tid, &posttableid)
		postCount := getPostCount(tid, posttableid, db)
		totalpages := int(math.Ceil(float64(postCount) / float64(20))) //page总数
		str := strconv.Itoa(tid) + "|" + strconv.Itoa(totalpages)
		a = append(a, str)
	}
	return a
}

func getPostCount(tid int, posttableid int, db *sql.DB) int {
	tableName := "pre_forum_post_" + strconv.Itoa(posttableid)
	if posttableid == 0 {
		tableName = "pre_forum_post"
	}
	//h5 bbs 获取用户回复post数据sql语句
	// rows, err := db.Query("select count(*) from `" + tableName + "` as p inner join `pre_common_member` as m on p.authorid=m.uid where m.groupid!=4 and invisible=0 and tid=" + strconv.Itoa(int(tid)))
	//pc bbs 获取用户回复post数据sql语句
	rows, err := db.Query("select count(*) from `" + tableName + "` where invisible=0 and tid=" + strconv.Itoa(int(tid)))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_post sql prepare error: ", err)
		return 0
	}
	for rows.Next() {
		var count int
		rows.Scan(&count)
		return count
	}
	return 0
}
