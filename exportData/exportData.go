package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func doSelect(searchSql string, dbName string) {

	//db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/"+dbName+"?charset=utf8");
	db, err := sql.Open("mysql", "dog123:dog123@tcp(192.168.5.199:3306)/"+dbName+"?charset=utf8")
	if err != nil {
		fmt.Printf("connect err")
	}

	rows, err1 := db.Query(searchSql)
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}

	defer rows.Close()
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, _ := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			fmt.Print(v)
			fmt.Print(",")
		}
		fmt.Print("\n")
	}
}

func main() {
	dbName := os.Args[1]
	searchSql := os.Args[2]
	fmt.Println("sql:", searchSql, " db:", dbName)
	doSelect(searchSql, dbName)
}
