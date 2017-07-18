package mysql

import (
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	// "reflect"
)

type EventLog struct {
	Id        int64
	TypeId    int
	Uid       int
	Info      string
	Created   string
	Infoid    int
	Status    int
	Tid       int
	isSplit   bool
	logLevel  int
	postTable string
}

type Follow struct {
	Follow_id int
	// user_id   int
}

func GetFansData(uid int, db *sql.DB) []*Follow {
	tableName := "follow"
	rows, err := db.Query("select distinct(follow_id) from `" + tableName + "` where user_id=" + strconv.Itoa(int(uid)) + " and fans_active=1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check event_log sql prepare error: ", err)
		return nil
	}
	var rowsData []*Follow
	for rows.Next() {
		var row = new(Follow)
		rows.Scan(&row.Follow_id)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

func LoadById(id int, db *sql.DB) *EventLog {
	tableName := "event_log"
	rows, err := db.Query("select id,type as typeId,uid,info,created,infoid,status,tid from `" + tableName + "` where id=" + strconv.Itoa(int(id)) + "")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check event_log sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row = new(EventLog)
		rows.Scan(&row.Id, &row.TypeId, &row.Uid, &row.Info, &row.Created, &row.Infoid, &row.Status, &row.Tid)
		return row
	}
	return &EventLog{}
}

func GetMysqlData(fans int, uid int, count int, page int, db *sql.DB, loopNum int, fansLimit int, eventLimit int, pushLimit int, dateLimit string) []*EventLog {
	offset := page * loopNum
	var sql string
	if fans >= fansLimit {
		sql = "select type,uid,created,infoid,status,tid from `event_log` where status=1 and uid=" + strconv.Itoa(uid) + " order by id desc limit " + strconv.Itoa(loopNum) + " offset " + strconv.Itoa(offset)
		// sql = "select type,uid,created,infoid,status,tid from `event_log` where uid=" + strconv.Itoa(uid) + " and created >='" + c.dateLimit + "' order by id desc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	} else if fans < fansLimit && count > pushLimit {
		sql = "select type,uid,created,infoid,status,tid from `event_log` where status=1 and uid=" + strconv.Itoa(uid) + " and created >='" + dateLimit + "' order by id desc limit " + strconv.Itoa(loopNum) + " offset " + strconv.Itoa(offset)
		// sql = "select type,uid,created,infoid,status,tid from `event_log` where uid=" + strconv.Itoa(uid) + " order by id desc limit " + strconv.Itoa(c.numloops) + " offset " + strconv.Itoa(offset)
	} else {
		sql = "select type,uid,created,infoid,status,tid from `event_log` where status=1 and uid=" + strconv.Itoa(uid) + " and created >='" + dateLimit + "' order by id desc limit " + strconv.Itoa(loopNum) + " offset " + strconv.Itoa(offset)
	}
	rows, err := db.Query(sql)
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check event_log sql prepare error: ", err)
		return nil
	}
	var rowsData []*EventLog
	for rows.Next() {
		var row = new(EventLog)
		rows.Scan(&row.TypeId, &row.Uid, &row.Created, &row.Infoid, &row.Status, &row.Tid)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

func GetFollowForumIds(uid int, db *sql.DB) []int {
	tableName := "forumfollow"
	rows, err := db.Query("select forum_id from `" + tableName + "` where user_id=" + strconv.Itoa(int(uid)))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check forum follow user error: ", err)
		return nil
	}
	var a []int
	for rows.Next() {
		var forum_id int
		rows.Scan(&forum_id)
		a = append(a, forum_id)
	}
	return a
}

func GetFollowForumId(tid int, db *sql.DB) int {
	tableName := "pre_forum_thread"
	rows, err := db.Query("select fid from `" + tableName + "` where tid=" + strconv.Itoa(tid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check forum follow user error: ", err)
		return 0
	}
	for rows.Next() {
		var fid int
		rows.Scan(&fid)
		return fid
	}
	return 0
}

func GetGoodsOrderUids(goodsId int, price string, db *sql.DB) []int {
	// _, priceInt := strconv.Atoi(price)
	priceFloat, _ := strconv.ParseFloat(price, 32)
	var rows *sql.Rows
	var err error
	if priceFloat > 0 {
		rows, err = db.Query("select uid from mall.`order` as o left join mall.order_info as oi on o.id=oi.order_id where oi.goods_sku_id=" + strconv.Itoa(goodsId) + " and oi.goods_price>" + price + " group by uid")
		defer rows.Close()
	} else {
		rows, err = db.Query("select uid from mall.`order` as o left join mall.order_info as oi on o.id=oi.order_id where oi.goods_sku_id=" + strconv.Itoa(goodsId) + " group by uid ")
		defer rows.Close()
	}
	if err != nil {
		logger.Error("[error] check forum follow user error: ", err)
		return nil
	}
	var a []int
	for rows.Next() {
		var uid int
		rows.Scan(&uid)
		a = append(a, uid)
	}
	return a
}
