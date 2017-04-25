package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func insertEventLog(
    id int64,
    type int64,
    uid int64,
    info string,
    created string,
    infoid int64,
    status int64,
    tid int64,
) int64 {

    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
        return 0
    }
    defer db.Close()

    stmt, err := db.Prepare("INSERT business SET city_id=?,pid=?,name=?,lat=0,lng=0,sort=0")
    if err != nil {
        logger.Println("[error] insert prepare error: ", err)
        return 0
    }

    res, err := stmt.Exec(cityId, pid, name)
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