package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

)

func testSelect() {

	db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/mall?charset=utf8");
	if err != nil {
		fmt.Printf("connect err");
	}

	rows, err1 := db.Query("select uid,token from user_token");
    if err1 != nil {
		fmt.Println(err1.Error());
		return;
	}

	defer rows.Close();
	fmt.Println("");
	cols, _ := rows.Columns();
	for i := range cols {
		fmt.Print(cols[i]);
		fmt.Print("\t");
	}

	var uid int;
	var token string;

	for rows.Next() {
		if err := rows.Scan(&uid, &token); err == nil {
			fmt.Print(uid);
			fmt.Print("\t");
			fmt.Print(token);
			fmt.Print("\t\r\n");
		}
	}
}

func main() {
	testSelect();
}
