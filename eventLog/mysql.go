package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    "fmt"
    "time"
    "math"
    "github.com/donnie4w/go-logger/logger"
    "github.com/jackson198608/squirrel"
    "github.com/jackson198608/structable"
    // "reflect"
)

type EventLog struct {
    structable.Recorder
    builder squirrel.StatementBuilderType
    // The req is Request object that contains the parsed result, which saved in PageItems.
    id int64
    typeId int
    uid int
    info string
    created string
    infoid int
    status int
    tid int
    isSplit   bool
    logLevel  int
    postTable string
}

//检查分表是否存在
func checkTableExist(tableName string) (bool) {
    db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Error("[error] connect db err")
        return false
    }
    defer db.Close()
    rows, err := db.Query("SHOW TABLES LIKE '"+ tableName +"'")
    defer rows.Close()
    if err != nil {
        logger.Error("[error] check business sql prepare error: ", err)
        fmt.Println("[error] check sql prepare error: ", err)
        return false
    }
    return rows.Next()

}

func insertEventLog(
    typeId int,
    uid int,
    info string,
    created string,
    infoid int,
    status int,
    tid int,
) int64 {

    db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Error("[error] connect db err")
        return 0
    }
    defer db.Close()
    //fmt.Println(time.Now().Format("2006-01-02 15:04:05"))  // 这是个奇葩,必须是这个时间点,记忆方法:6-1-2-3-4-5    year := time.Now().Format("2006")
    t, _ := time.Parse("2006-01-02 15:04:05", created)
    t4,_:= strconv.ParseFloat(t.Format("1"),64)  
    season := math.Ceil(t4 / float64(3)) //根据日期计算所在季度,按照季度进行分表
    year := t.Format("2006")
    // fmt.Println(reflect.TypeOf(season))
    seasons := strconv.FormatFloat(season, 'f',-1, 64)
    tableName := "event_log_"+year+seasons //分表名按照季度分表
    fmt.Println(tableName)
    tableIsExist := checkTableExist(tableName)
    fmt.Println(tableIsExist)
    if tableIsExist == false {
        createTable(tableName)
    }
    eventId, isExist := checkEventExist(tableName, int(uid), info, created, int(infoid))
    if isExist {
        logger.Error("[error] event info is exist: ", eventId)
        return 0
    }
    stmt, err := db.Prepare("INSERT INTO `"+tableName+"` SET type=?,uid=?,info=?,created=?,infoid=?,status=?,tid=?")
    if err != nil {
        logger.Error("[error] insert prepare error: ", err)
        return 0
    }
    defer stmt.Close()
    res, err := stmt.Exec(typeId,uid,info,created,infoid,status,tid)
    if err != nil {
        logger.Error("[error] insert excute error: ", err)
        return 0
    }

    id, err := res.LastInsertId()
    if err != nil {
        logger.Error("[error] get insert id error: ", err)
        return 0
    }
    return id
}

func checkEventExist(tableName string,uid int,info string,created string,infoid int) (int64, bool) {
    db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Error("[error] connect db err")
        //fmt.Println("[error] connect db err")
        return 0, false
    }
    defer db.Close()
    // fmt.Println(reflect.TypeOf(uid))
    rows, err := db.Query("select id from `"+tableName+"` where uid='" +strconv.Itoa(uid)+ "' and info='"+info+"' and created='" + created + "' and infoid="+strconv.Itoa(infoid)+"")
    defer rows.Close()
    if err != nil {
        logger.Error("[error] check checkEventExist sql prepare error: ", err)
        //fmt.Println("[error] check sql prepare error: ", err)
        return 0, false
    }

    for rows.Next() {
        var eventId int64
        if err := rows.Scan(&eventId); err != nil {
            logger.Error("[error] check checkEventExist sql get rows error ", err)
            return 0, false
        }
        logger.Error("[info] check checkEventExist sql find true ", eventId, " ", eventId)
        //fmt.Println("[info] check sql find true", shopName, " ", shopId)
        return eventId, true

    }
    if err := rows.Err(); err != nil {
        logger.Error("[error] check checkEventExist sql get rows error ", err)
        //fmt.Println("[error] check sql get rows error ", err)
        return 0, false
    }
    return 0, false

}

//创建分表
func createTable(tableName string){
    db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
    if err != nil {
        fmt.Println(err)
        // logger.Println("[error] connect db err")
    }
    defer db.Close()
    stmt, err := db.Prepare("create table "+ tableName +" like event_log")
    if err != nil {
        fmt.Println(err)
        // logger.Println("[error] insert prepare error: ", err)
    }
    defer stmt.Close()
    res, err := stmt.Exec()
    fmt.Println(res)
    if err != nil {
        logger.Error("[error] create table excute error: ", err)
    }

    // num, err := res.RowsAffected()
    // if err != nil {
    //     logger.Error("[error] create event_log_XXX table excute error: ", err,num)
    // }
}

//获取动态数据
func getEventLogData(startId int,endId int,limit int,offset int)[]*EventLog{
    db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Error("[error] connect db err")
    }
    defer db.Close()
    tableName := "event_log"
    
    rows, err := db.Query("select type as typeId,uid,info,created,infoid,status,tid from `"+ tableName +"` where id>=" + strconv.Itoa(startId) + " and id<=" + strconv.Itoa(endId) +" limit "+ strconv.Itoa(limit) + " offset "+strconv.Itoa(offset)+"")
    if err != nil {
        logger.Error("[error] check getEventLogData sql prepare error: ", err)
        // return 0, false
    }
    defer rows.Close()
    var rowsData []*EventLog
    for rows.Next() {
        var row = new(EventLog)
        rows.Scan(&row.typeId,&row.uid,&row.info,&row.created,&row.infoid,&row.status,&row.tid)
        rowsData = append(rowsData,row)
    }

    // for _,ar := range rowsData {
    //     fmt.Println(ar.id,ar.typeId)
    // }
    // // fmt.Println(rowsData)
    return rowsData
}

func NewEvent(logLevel int, db squirrel.DBProxyBeginner, dbFlavor string, id int64, isSplit bool) *EventLog {
    u := new(EventLog)
    logger.SetLevel(logger.LEVEL(logLevel))

    u.isSplit = isSplit
    if (id > 0){
        u.id = id
    }


    if (id > 0){
        u = u.LoadById()
    }

    u.logLevel = logLevel

    return u
}

func (p *EventLog) IdExists() bool {
    isExist, err := p.ExistsWhere("id = ?", p.id)
    if err != nil {
        logger.Error("find exists error", p.id,p.TableName(), err)
    }
    return isExist
}
func (p *EventLog) hasChanged() bool {
    if p.id <= 0{
        logger.Error("have no pid or tid can not continute")
        return false
    }
    if !p.isSplit {
        // p.postTable = p.getTableSplitName()
        p.Recorder.ChangeBindTableName(p.postTable)
        p.isSplit = true
    }
    isExist := p.IdExists()

    // p.backToMain()
    return isExist

}

func (p *EventLog) MoveToSplit() bool {
    id := insertEventLog(p.typeId,p.uid,p.info,p.created,p.infoid,p.status,p.tid)
    if id == 0 {
        logger.Error("insert error",id)
        return false
    }
    return true
}

func (p *EventLog) getTableSplitName() string {
    return ""
}

// LoadByName is a custom loader.
//
// The Load() method on a Recorder loads by ID. This allows us to load by
// a different field -- Name.
func (p *EventLog) LoadByPid() error {
    return p.Recorder.LoadWhere("id = ? limit 0,1", p.id)
}


func (p *EventLog) LoadById() *EventLog{
    db, err := sql.Open("mysql", c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
    if err != nil {
        logger.Error("[error] connect db err")
    }
    defer db.Close()
    tableName := "event_log"
    rows, err := db.Query("select id,type as typeId,uid,info,created,infoid,status,tid from `"+ tableName +"` where id=" + strconv.Itoa(int(p.id))+ "")
    defer rows.Close()
    if err != nil {
        logger.Error("[error] check event_log sql prepare error: ", err)
        return nil
    }
    for rows.Next() {
        var row = new(EventLog)
        rows.Scan(&row.id,&row.typeId,&row.uid,&row.info,&row.created,&row.infoid,&row.status,&row.tid)
        return row
    }
    // for _,ar := range rowsData {
    //     fmt.Println(ar.id,ar.typeId)
    // }
    // // fmt.Println(rowsData)
    return &EventLog{}
}