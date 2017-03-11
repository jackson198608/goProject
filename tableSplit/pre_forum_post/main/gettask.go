package main

import (
	"fmt"
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
	rows, err := db.Query("select tid from pre_forum_thread order by tid asc limit " + strconv.Itoa(offset) +
		" offset " + strconv.Itoa(offset))
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
