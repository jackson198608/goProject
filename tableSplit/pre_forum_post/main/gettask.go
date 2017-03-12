package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func getTask(page int) []int64 {
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		fmt.Println("[error] connect db err")
		return nil
	}
	defer db.Close()

	offset := page * numloops
	sql := "select tid from pre_forum_thread where tid < " + strconv.Itoa(lastTid) + " and tid >= " + strconv.Itoa(firstTid) + " order by tid asc limit " + strconv.Itoa(numloops) + " offset " + strconv.Itoa(offset)
	fmt.Println(sql)
	rows, err := db.Query(sql)

	if err != nil {
		fmt.Println("[error] query error")
		return nil
	}

	tids := make([]int64, 0, numloops)

	for rows.Next() {
		var tid int64
		if err := rows.Scan(&tid); err != nil {
			fmt.Println("[error] check sql get rows error ", err)
			return nil
		}
		tids = append(tids, tid)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("[error] check sql get rows error ", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return nil

	}
	return tids
}
