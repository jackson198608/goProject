package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var dsn string = "root:goumintech@tcp(192.168.86.72:3309)/"
var database string = "mall?charset=utf8"

func connect() *sql.DB {
	//db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/mall?charset=utf8")
	db, err := sql.Open("mysql", dsn+database)
	if err != nil {
		fmt.Printf("connect err")
	}
	return db

}

func main() {
	db := connect()
}
