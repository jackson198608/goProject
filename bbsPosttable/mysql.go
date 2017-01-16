package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "strconv"
)

type Post struct {

    // The req is Request object that contains the parsed result, which saved in PageItems.
    tid int

    pid int
}

func getNewTaskData(startPid int,endPid int,limit int,offset int) Post {
    
}

//读取post任务数据
func getTaskData(startPid int,endPid int,limit int,offset int) map[int]map[string]string{
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
    }
    defer db.Close()
    tableName := "pre_forum_post"
    rows, err := db.Query("select tid,pid from `"+ tableName +"` where pid>=" + strconv.Itoa(startPid) + " and pid<=" + strconv.Itoa(endPid) +" limit "+ strconv.Itoa(limit) + " offset "+strconv.Itoa(offset)+"")
    if err != nil {
        logger.Println("[error] check business sql prepare error: ", err)
        // return 0, false
    }
    defer rows.Close()
    columns, err := rows.Columns() //读出查询出的列字段名  
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
    // Make a slice for the values
    values := make([][]byte, len(columns))//values是每个列的值，这里获取到byte里  
    // values := make([]sql.RawBytes, len(columns))//values是每个列的值，这里获取到byte里  

    // scanArgs := make([]interface{}, len(columns))//因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度 
    scanArgs := make([]interface{}, len(values))//因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度 
    for i := range values {
        scanArgs[i] = &values[i]
    }

    results := make(map[int]map[string]string) //最后得到的map  
    i := 0 
    for rows.Next() { //循环，让游标往下移动  
        err = rows.Scan(scanArgs...)//query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里  
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
            fmt.Println(err)
            return nil
        }
        row := make(map[string]string)//每行数据  
        var value string
        for k, v := range values {  //每行数据是放在values里面，现在把它挪到row里  
            if v == nil {
                value = "NULL"
            }else{
                value = string(v)
            }
            key := columns[k]
            row[key] = value
        }
        results[i] = row   //装入结果集中  
        i++
    }
    // fmt.Println(results)
    // for _, v := range results {   //查询出来的数组  
    //     // fmt.Println(k)
    //     for kk,vv := range v {
    //         fmt.Println(kk)
    //         fmt.Println(vv)
    //     }
    // }

    if err := rows.Err(); err != nil {
        logger.Println("[error] check business sql get rows error ", err)
    }
    return results
}



//修改帖子表thread的posttableid字段的值
func updateThread(tableid int ,tid int){
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
    }
    defer db.Close()
    posttableid := strconv.Itoa(tableid)
    stmt, err := db.Prepare("update pre_forum_thread SET posttableid=? where tid=?")
    if err != nil {
        logger.Println("[error] insert prepare error: ", err)
    }
    res, err := stmt.Exec(posttableid,strconv.Itoa(tid))
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }
    num, err := res.RowsAffected()
    fmt.Println(num)
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }
}

//如果post分表不存在创建新分表
func createTable(tableid int){
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
    }
    defer db.Close()
    tableName := "pre_forum_post_"+ strconv.Itoa(tableid)
    stmt, err := db.Prepare("create table "+ tableName +" like pre_forum_post")
    if err != nil {
        logger.Println("[error] insert prepare error: ", err)
    }

    res, err := stmt.Exec()
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }

    num, err := res.RowsAffected()
    fmt.Println(num)
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }
}

//把post数据按照tid分表导入新的post_tid%100表中
func insertPost(tableid int,tid int,pid int){
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
    }
    defer db.Close()
    tableName := "pre_forum_post_"+ strconv.Itoa(tableid)
    stmt, err := db.Prepare("INSERT IGNORE INTO "+ tableName +" select * from pre_forum_post where tid="+strconv.Itoa(tid))
    if err != nil {
        logger.Println("[error] insert prepare error: ", err)
    }

    res, err := stmt.Exec()
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }

    num, err := res.RowsAffected()
    fmt.Println(num)
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }
}

//检查数据是否存在
func checkPostExist(tableid int,tid int, pid int ) ( bool) {
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
        //fmt.Println("[error] connect db err")
        return false
    }
    defer db.Close()
    tableName := "pre_forum_post_"+strconv.Itoa(tableid)
    rows, err := db.Query("select pid from `"+ tableName +"` where tid=" + strconv.Itoa(tid) + " and pid=" + strconv.Itoa(pid) +"")
    if err != nil {
        logger.Println("[error] check business sql prepare error: ", err)
        fmt.Println("[error] check sql prepare error: ", err)
        return false
    }

    return rows.Next()

}

//检查分表是否存在
func checkTableExist(tableid int) (bool) {
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
        return false
    }
    defer db.Close()
    tableName := "pre_forum_post_" + strconv.Itoa(tableid)
    rows, err := db.Query("SHOW TABLES LIKE '"+ tableName +"'")
    if err != nil {
        logger.Println("[error] check business sql prepare error: ", err)
        fmt.Println("[error] check sql prepare error: ", err)
        return false
    }
    return rows.Next()

}