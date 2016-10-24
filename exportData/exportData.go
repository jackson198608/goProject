package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func testSelect(searchSql string, dbName string) {

	db, err := sql.Open("mysql", "dog123:dog123@tcp(210.14.154.198:3306)/"+dbName+"?charset=utf8")
	//db, err := sql.Open("mysql", "dog123:dog123@tcp(192.168.5.199:3306)/ask?charset=utf8");
	if err != nil {
		fmt.Printf("connect err")
	}

	rows, err1 := db.Query(searchSql)
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}

	defer rows.Close()
	fmt.Println("")
	cols, _ := rows.Columns()
	for i := range cols {
		fmt.Print(cols[i])
		fmt.Print("\t")
	}

	var Host string
	var User string
	fmt.Println(rows)

	for rows.Next() {
		if err := rows.Scan(&Host, &User); err == nil {
			fmt.Print(Host)
			fmt.Print("\t")
			fmt.Print(User)
			fmt.Print("\t\r\n")
		}
	}
}

func main() {
	dbName := os.Args[1]
	searchSql := os.Args[2]
	fmt.Println("sql:", searchSql, " db:", dbName)
	testSelect(searchSql, dbName)
}
