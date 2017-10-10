package main

import (
	"database/sql"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"strconv"
)

func getAskList(page int) []string {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)

		return nil
	}
	defer db.Close()
	rows, err := db.Query("select id,ans_num from `ask`.`ask_question` where is_hide=1 order by id asc limit " + strconv.Itoa(offset) + " offset " + strconv.Itoa(offset*(page-1)))
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

func getThreadTask(page int) []string {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)

		return nil
	}
	defer db.Close()
	tableName := "pre_forum_thread"
	rows, err := db.Query("select tid,posttableid from `" + tableName + "` where displayorder in(0,1) order by tid asc limit " + strconv.Itoa(offset) + " offset " + strconv.Itoa(offset*(page-1)))
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
		totalpages := int(math.Ceil(float64(200) / float64(20))) //page总数
		str := strconv.Itoa(tid) + "|" + strconv.Itoa(totalpages)
		a = append(a, str)
	}
	return a
}
