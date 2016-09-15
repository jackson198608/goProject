package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

)

func testSelect() {

	//db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/mysql?charset=utf8");
	db, err := sql.Open("mysql", "dog123:dog123@tcp(192.168.5.199:3306)/ask?charset=utf8");
	if err != nil {
		fmt.Printf("connect err");
	}

	rows, err1 := db.Query("select uid,subject from ask_question order by id desc limit 0,5");
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
	fmt.Println("here\n");

	var Host int;
	var User string;

	for rows.Next() {
		if err := rows.Scan(&Host, &User); err == nil {
			fmt.Print(Host);
			fmt.Print("\t");
			fmt.Print(User);
			fmt.Print("\t\r\n");
		}
	}
}

func main() {
	testSelect();
}
