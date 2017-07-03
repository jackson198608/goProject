package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

func checkKeywordExist(keyword string) (int, int, bool) {
	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		//fmt.Println("[error] connect db err")
		return 0, 0, false
	}
	defer db.Close()
	date := time.Now().Format("2006-01-02")
	rows, err := db.Query("select id,rank from rank_list where keyword='" + strings.Trim(keyword, " ") + "' and date='" + date + "'")
	if err != nil {
		logger.Println("[error] check sql prepare error: ", err)
		//fmt.Println("[error] check sql prepare error: ", err)
		return 0, 0, false
	}
	for rows.Next() {
		var Id int
		var Rank int
		if err := rows.Scan(&Id, &Rank); err != nil {
			logger.Println("[error] check sql get rows error ", err)
			return 0, 101, false
		}
		logger.Println("[info] check sql find true", keyword, " ", Id)
		return Id, Rank, true

	}
	return 0, 0, false
}

func saveKeywordRankData(
	keyword string,
	rank int,
	url string,
	domain string,
) int64 {

	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()
	date := time.Now().Format("2006-01-02")
	stmt, err := db.Prepare("INSERT rank_list SET keyword=?, rank=?,date=?, url=?,domain=?")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(keyword, rank, date, url, domain)
	if err != nil {
		logger.Println("[error] insert excute error: ", err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Println("[error] get insert id error: ", err)
		return 0
	}
	return id

}

func updateKeywordRank(
	id int,
	keyword string,
	rank int,
	url string,
	domain string,
) bool {

	db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("update rank_list set rank = ?,url = ?,domain = ? where id=?")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(rank, url, domain, id)
	if err != nil {
		logger.Println("[error] update excute error: ", err)
		return false
	}

	num, err := res.RowsAffected()
	if err != nil {
		logger.Println("[error] get insert id error: ", err, " num:", num)
		return false
	}

	return true
}
