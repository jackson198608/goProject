package main

import (
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func getTask(page int) []int64 {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)

		return nil
	}
	defer db.Close()

	offset := page * c.numloops
	//获取正常显示和隐藏的数据
	sql := "select id from event_log where status in (-1,1) and id <= " + strconv.Itoa(c.lastId) + " and id >= " + strconv.Itoa(c.firstId) + " order by id asc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	logger.Info(sql)
	rows, err := db.Query(sql)

	if err != nil {
		logger.Error("sql query error", sql)
		return nil
	}

	ids := make([]int64, 0, c.numloops)

	for rows.Next() {
		var id int64
		// var uid string
		// var created string
		if err := rows.Scan(&id); err != nil {
			logger.Error("check sql get rows error", err)
			return nil
		}
		// ids = append(ids, id)
		// var fans int
		// fan := db.QueryRow("SELECT count(*) as counts FROM `follow` where user_id=" + uid)
		// fan.Scan(&fans)
		// fansLimit, _ := strconv.Atoi(c.fansLimit)
		// if fansLimit > 0 && fans > fansLimit {
		// 	if created > c.dateLimit {
		// 		ids = append(ids, id)
		// 	}
		// } else {
		ids = append(ids, id)
		// }
	}
	if err := rows.Err(); err != nil {
		logger.Error("check sql get rows error", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return nil

	}
	return ids
}

func getFollowTask(page int) []string {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)

		return nil
	}
	defer db.Close()

	offset := page * c.numloops
	//获取正常显示和隐藏的数据
	sql := "select distinct(user_id) from follow order by id asc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	// sql := "select distinct(user_id) from follow where id <= " + strconv.Itoa(c.followLastId) + " and id >= " + strconv.Itoa(c.followFirstId) + " order by id asc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	// sql := "select distinct(user_id) from follow where user_id in(1,881050,881052,1138687,49567,1138689,1140002,1140013,1140001,1140009,1139968,1139934,1139976) and id < " + strconv.Itoa(c.followLastId) + " and id >= " + strconv.Itoa(c.followFirstId) + " order by id asc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	logger.Info(sql)
	rows, err := db.Query(sql)

	if err != nil {
		logger.Error("sql query error", sql)
		return nil
	}

	uids := make([]string, 0, c.numloops)

	for rows.Next() {
		var user_id int
		// var created string
		if err := rows.Scan(&user_id); err != nil {
			logger.Error("check sql get rows error", err)
			return nil
		}
		// ids = append(ids, id)
		var fans int
		fan := db.QueryRow("SELECT count(*) as counts FROM `follow` where user_id=" + strconv.Itoa(user_id))
		fan.Scan(&fans)
		var eventNums = 0
		// if fans < c.fansLimit {
		eventNum := db.QueryRow("SELECT count(*) as counts FROM `event_log` where status=1 and uid=" + strconv.Itoa(user_id) + " and created >='" + c.dateLimit + "'")
		eventNum.Scan(&eventNums)
		// }
		follow := strconv.Itoa(user_id) + "|" + strconv.Itoa(fans) + "|" + strconv.Itoa(eventNums)
		logger.Info("follow.data", strconv.Itoa(user_id))
		uids = append(uids, follow)
	}
	if err := rows.Err(); err != nil {
		logger.Error("check sql get rows error", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return nil

	}
	return uids
}

func getFansTask(page int) []int64 {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)

		return nil
	}
	defer db.Close()

	offset := page * c.numloops
	//获取正常显示和隐藏的数据
	sql := "select distinct(follow_id) from follow order by id asc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	// sql := "select distinct(user_id) from follow where id <= " + strconv.Itoa(c.followLastId) + " and id >= " + strconv.Itoa(c.followFirstId) + " order by id asc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	// sql := "select distinct(user_id) from follow where user_id in(1,881050,881052,1138687,49567,1138689,1140002,1140013,1140001,1140009,1139968,1139934,1139976) and id < " + strconv.Itoa(c.followLastId) + " and id >= " + strconv.Itoa(c.followFirstId) + " order by id asc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	logger.Info(sql)
	rows, err := db.Query(sql)

	if err != nil {
		logger.Error("sql query error", sql)
		return nil
	}

	uids := make([]int64, 0, c.numloops)

	for rows.Next() {
		var follow_id int64
		// var created string
		if err := rows.Scan(&follow_id); err != nil {
			logger.Error("check sql get rows error", err)
			return nil
		}
		uids = append(uids, follow_id)
	}
	if err := rows.Err(); err != nil {
		logger.Error("check sql get rows error", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return nil

	}
	return uids
}
