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
    counts int
}



func getNewTaskData(startPid int,endPid int,limit int,offset int) []*Post{
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
    var rowsData []*Post
    for rows.Next() {
        var row = new(Post)
        rows.Scan(&row.tid,&row.pid)
        rowsData = append(rowsData,row)
    }

    // for _,ar := range rowsData {
    //     fmt.Println(ar.tid,ar.pid)
    // }
    // fmt.Println(rowsData)
    return rowsData
}

func getPostCount(startPid int,endPid int) (int){
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
    }
    defer db.Close()
    tableName := "pre_forum_post"
    user := new(Post)
    row  :=db.QueryRow("SELECT count(*) as counts FROM `"+ tableName +"` where pid>=" + strconv.Itoa(startPid) + " and pid<=" + strconv.Itoa(endPid) +"")
    row.Scan(&user.counts)
    return user.counts
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
        logger.Println("[error] cloumns err",err) // proper error handling instead of panic in your app
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
            logger.Println("[error] rows scan",err)// proper error handling instead of panic in your app
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
        fmt.Println(err)
        // logger.Println("[error] insert prepare error: ", err)
    }
    res, err := stmt.Exec(posttableid,strconv.Itoa(tid))
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }
    num, err := res.RowsAffected()
    if err != nil {
        logger.Println("[error] insert excute error: ", err,num)
    }
}

//如果post分表不存在创建新分表
func createTable(tableid int){
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        fmt.Println(err)
        // logger.Println("[error] connect db err")
    }
    defer db.Close()
    tableName := "pre_forum_post_"+ strconv.Itoa(tableid)
    stmt, err := db.Prepare("create table "+ tableName +" like pre_forum_post")
    if err != nil {
        fmt.Println(err)
        // logger.Println("[error] insert prepare error: ", err)
    }

    res, err := stmt.Exec()
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }

    num, err := res.RowsAffected()
    if err != nil {
        logger.Println("[error] insert excute error: ", err,num)
    }
}

//把post数据按照tid分表导入新的post_tid%100表中
func insertPost(tableid int,tid int,pid int){
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        fmt.Println(err)
        // logger.Println("[error] connect db err")
    }
    defer db.Close()
    tableName := "pre_forum_post_"+ strconv.Itoa(tableid)
    stmt, err := db.Prepare("INSERT IGNORE INTO "+ tableName +" select * from pre_forum_post where tid="+strconv.Itoa(tid))
    if err != nil {
        fmt.Println(err)
        // logger.Println("[error] insert prepare error: ", err)
    }

    res, err := stmt.Exec()
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }

    num, err := res.RowsAffected()
    if err != nil {
        logger.Println("[error] insert excute error: ", err,num)
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

type EventLog struct {
    counts int
    id int
    infoid int
}

// event log

func getEventLogCount() (int){
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
    }
    defer db.Close()
    tableName := "event_log"
    event := new(EventLog)
    row  :=db.QueryRow("SELECT count(*) as counts FROM `"+ tableName +"` where type=2 and tid=0")
    row.Scan(&event.counts)
    return event.counts
}

func getEventLogTask(limit int,offset int) []*EventLog{
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
    }
    defer db.Close()
    tableName := "event_log"
    
    rows, err := db.Query("select id,infoid from `"+ tableName +"` where type=2 and tid=0 limit "+ strconv.Itoa(limit) + " offset "+strconv.Itoa(offset)+"")

    if err != nil {
        logger.Println("[error] check business sql prepare error: ", err)
        // return 0, false
    }
    defer rows.Close()
    var rowsData []*EventLog
    for rows.Next() {
        var row = new(EventLog)
        rows.Scan(&row.id,&row.infoid)
        rowsData = append(rowsData,row)
    }

    return rowsData
}

func updateEventLogTid(id int, tid int) {
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
    }
    defer db.Close()
    stmt, err := db.Prepare("update event_log SET tid=? where id=?")
    if err != nil {
        fmt.Println(err)
        // logger.Println("[error] insert prepare error: ", err)
    }
    res, err := stmt.Exec(strconv.Itoa(tid),strconv.Itoa(id))
    if err != nil {
        logger.Println("[error] insert excute error: ", err)
    }
    num, err := res.RowsAffected()
    if err != nil {
        logger.Println("[error] insert excute error: ", err,num)
    } 
}


//检查数据是否存在
func checkEventPostExist(pid int ) (int64, bool) {
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
        //fmt.Println("[error] connect db err")
        return 0,false
    }
    defer db.Close()
    tableName := "pre_forum_post"

    rows, err := db.Query("select tid from `"+ tableName +"` where pid=" + strconv.Itoa(pid) +"")
    if err != nil {
        logger.Println("[error] check business sql prepare error: ", err)
        fmt.Println("[error] check sql prepare error: ", err)
        return 0,false
    }

    for rows.Next() {
        var tid int64
        if err := rows.Scan(&tid); err != nil {
            logger.Println("[error] check post sql get rows error ", err)
            return 0, false
        }
        return tid, true

    }

    return 0,false

}

func checkEventLogExist(id int ) ( bool) {
    db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Println("[error] connect db err")
        //fmt.Println("[error] connect db err")
        return false
    }
    defer db.Close()
    tableName := "event_log"
    rows, err := db.Query("select id from `"+ tableName +"` where id=" + strconv.Itoa(id) +"")
    if err != nil {
        logger.Println("[error] check business sql prepare error: ", err)
        fmt.Println("[error] check sql prepare error: ", err)
        return false
    }

    return rows.Next()

}