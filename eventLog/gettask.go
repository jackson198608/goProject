package main

import (
    "database/sql"
    "github.com/donnie4w/go-logger/logger"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    // "fmt"
)

func getTask(page int) []int64 {
    db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Error("can not connect to mysql", c.dbDsn, c.dbName, c.dbAuth)

        return nil
    }
    defer db.Close()

    offset := page * c.numloops
    sql := "select id from event_log where id < " + strconv.Itoa(c.lastId) + " and id >= " + strconv.Itoa(c.firstId) + " order by id asc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
    logger.Info(sql)
    rows, err := db.Query(sql)

    if err != nil {
        logger.Error("sql query error", sql)
        return nil
    }

    ids := make([]int64, 0, c.numloops)

    for rows.Next() {
        var id int64
        if err := rows.Scan(&id); err != nil {
            logger.Error("check sql get rows error", err)
            return nil
        }
        ids = append(ids, id)
    }
    if err := rows.Err(); err != nil {
        logger.Error("check sql get rows error", err)
        //fmt.Println("[error] check sql get rows error ", err)
        return nil

    }
    return ids
}
