package mysql

import (
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	// "reflect"
)

type EventLog struct {
	Id        int64
	TypeId    int
	Uid       int
	Info      string
	Created   string
	Infoid    int
	Status    int
	Tid       int
	isSplit   bool
	logLevel  int
	postTable string
}

type Follow struct {
	follow_id int
	// user_id   int
}

func GetFansData(uid int, db *sql.DB) []int {
	// db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	// if err != nil {
	//  logger.Error("[error] connect db err")
	// }
	// defer db.Close()
	tableName := "follow"
	rows, err := db.Query("select distinct(follow_id) from `" + tableName + "` where user_id=" + strconv.Itoa(int(uid)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check event_log sql prepare error: ", err)
		return nil
	}
	// var rowsData []*Follow
	// for rows.Next() {
	// 	var row = new(Follow)
	// 	rows.Scan(&row.follow_id)
	// 	rowsData = append(rowsData, row)
	// }
	pids := make([]int, 0, 100)
	for rows.Next() {
		var pid int
		if err := rows.Scan(&pid); err != nil {
			logger.Error("get pid error after row.next", err)
		}
		pids = append(pids, pid)
	}
	return pids
}

func LoadById(id int, db *sql.DB) *EventLog {
	tableName := "event_log"
	rows, err := db.Query("select id,type as typeId,uid,info,created,infoid,status,tid from `" + tableName + "` where id=" + strconv.Itoa(int(id)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check event_log sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row = new(EventLog)
		rows.Scan(&row.Id, &row.TypeId, &row.Uid, &row.Info, &row.Created, &row.Infoid, &row.Status, &row.Tid)
		return row
	}
	return &EventLog{}
}
