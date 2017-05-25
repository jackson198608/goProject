package main

import (
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	// "reflect"
)

type EventLog struct {
	id        int64
	typeId    int
	uid       int
	info      string
	created   string
	infoid    int
	status    int
	tid       int
	isSplit   bool
	logLevel  int
	postTable string
}

func LoadById(id int, db *sql.DB) *EventLog {
	// db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	// if err != nil {
	// 	logger.Error("[error] connect db err")
	// }
	// defer db.Close()
	tableName := "event_log"
	rows, err := db.Query("select id,type as typeId,uid,info,created,infoid,status,tid from `" + tableName + "` where id=" + strconv.Itoa(int(id)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check event_log sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row = new(EventLog)
		rows.Scan(&row.id, &row.typeId, &row.uid, &row.info, &row.created, &row.infoid, &row.status, &row.tid)
		return row
	}
	// for _,ar := range rowsData {
	//     fmt.Println(ar.id,ar.typeId)
	// }
	// // fmt.Println(rowsData)
	return &EventLog{}
}

type Follow struct {
	follow_id int
	// user_id   int
}

func GetFansData(uid int, db *sql.DB) []*Follow {
	// db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	// if err != nil {
	// 	logger.Error("[error] connect db err")
	// }
	// defer db.Close()
	tableName := "follow"
	rows, err := db.Query("select follow_id from `" + tableName + "` where user_id=" + strconv.Itoa(int(uid)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check event_log sql prepare error: ", err)
		return nil
	}
	var rowsData []*Follow
	for rows.Next() {
		var row = new(Follow)
		rows.Scan(&row.follow_id)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

func getMysqlData(fans int, uid int, count int, page int, db *sql.DB) []*EventLog {
	offset := page * c.numloops
	var sql string
	if fans >= c.fansLimit {
		sql = "select type,uid,created,infoid,status,tid from `event_log` where uid=" + strconv.Itoa(uid) + " order by id desc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
		// sql = "select type,uid,created,infoid,status,tid from `event_log` where uid=" + strconv.Itoa(uid) + " and created >='" + c.dateLimit + "' order by id desc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	} else if fans < c.fansLimit && count > c.pushLimit {
		sql = "select type,uid,created,infoid,status,tid from `event_log` where uid=" + strconv.Itoa(uid) + " and created >='" + c.dateLimit + "' order by id desc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
		// sql = "select type,uid,created,infoid,status,tid from `event_log` where uid=" + strconv.Itoa(uid) + " order by id desc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	} else {
		sql = "select type,uid,created,infoid,status,tid from `event_log` where uid=" + strconv.Itoa(uid) + " and created >='" + c.dateLimit + "' order by id desc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	}
	rows, err := db.Query(sql)
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check event_log sql prepare error: ", err)
		return nil
	}
	var rowsData []*EventLog
	for rows.Next() {
		var row = new(EventLog)
		rows.Scan(&row.typeId, &row.uid, &row.created, &row.infoid, &row.status, &row.tid)
		rowsData = append(rowsData, row)
	}
	return rowsData
}
