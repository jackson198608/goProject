package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"log/syslog"
	"time"
)

type PreUcenterMembers struct {
	Uid      int    `xorm:"int(11)  NOT NULL"`
	Username string `xorm:"char(20) NOT NULL DEFAULT ''"`
	Salt     string `xorm:"char(6) NOT NULL"`
	Regdate  int    `xorm:"int(10) NOT NULL DEFAULT '0'"`
	Password string `xorm:"char(32) NOT NULL DEFAULT ''"`
	Email    string `xorm:"char(32) NOT NULL DEFAULT ''"`
	// Myid          string `xorm:"char(30) NOT NULL DEFAULT ''"`
	// Myidkey       string `xorm:"char(16) NOT NULL DEFAULT ''"`
	// Regip         string `xorm:"char(15) NOT NULL DEFAULT ''"`
	// Lastloginip   int    `xorm:"int(10) NOT NULL DEFAULT '0'"`
	// Lastlogintime int    `xorm:"int(10) unsigned NOT NULL DEFAULT '0'"`
	// Secques       string `xorm:"char(8) NOT NULL DEFAULT ''"`
	// EmailAuth     string `xorm:"varchar(20) DEFAULT NULL"`
	// AuthStatus    int    `xorm:"int(11) DEFAULT '0'"`
	// Phone         string `xorm:"varchar(20) NOT NULL DEFAULT "`
	// LoginKey      string `xorm:"varchar(128) NOT NULL DEFAULT ''"`
	// Passwd        string `xorm:"varchar(200) DEFAULT NULL"`
	// Created  time.Time `xorm:"created"`
	// Updated  time.Time `xorm:"updated"`
}

type Report struct {
	Id       int
	Cid      int
	Type     int
	Status   int
	Created  time.Time `xorm:"created"`
	Modified time.Time `xorm:"created"`
}

func getUser() {
	dataSourceName := dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println(err)
	}
	logWriter, err := syslog.New(syslog.LOG_DEBUG, "rest-xorm-example")
	if err != nil {
		log.Fatalf("Fail to create xorm system logger: %v\n", err)
	}

	logger := xorm.NewSimpleLogger(logWriter)
	logger.ShowSQL(true)
	engine.SetLogger(logger)
	// engine.ShowSQL(true) //则会在控制台打印出生成的SQL语句；
	// err0 := engine.Sync2(new(PreUcenterMembers))
	// if err0 != nil {
	// 	fmt.Println(err0)
	// }
	// 返回的结果类型为 []map[string][]byte
	// results, err1 := engine.Query("select * from pre_ucenter_members where uid=881050")
	// QueryString 返回 []map[string]string
	// results, err1 := engine.QueryString("select * from pre_ucenter_members where uid=881050")
	// if err1 != nil {
	// 	fmt.Println(err1)
	// }
	// var uid string = "0"
	// for _, v := range results {
	// 	uid = v["uid"]
	// 	fmt.Println(uid)
	// }
	// fmt.Println(engine)
	// err = engine.CreateTables(&PreUcenterMembers{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	/**var email string = "lidings@goumin.com"
	affected, err2 := engine.Exec("update pre_ucenter_members set email = ? where uid = ?", email, uid)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(affected)
	**/
	// user := PreUcenterMembers{}
	// user := &PreUcenterMembers{}
	// has, err3 := engine.Get(user)
	// if err3 != nil {
	// 	fmt.Println(err3)
	// }
	// fmt.Println(has)

	var valuesMap = &PreUcenterMembers{}
	var id string = "881050"
	has1, err4 := engine.Where("uid = ?", id).Get(valuesMap)
	if err4 != nil {
		fmt.Println(err4)
	}
	fmt.Println(valuesMap.Uid)
	fmt.Println(has1)

	// var valuesMap = make(map[string]string)
	// has3, err6 := engine.Where("id = ?", id).Get(&valuesMap)
	// fmt.Println(has3)
	// fmt.Println(err6)
	// fmt.Println(valuesMap)

	// var user PreUcenterMembers
	// has2, err5 := engine.Where("uid = ?", id).Cols("username", "salt").Get(&user)
	// if err5 != nil {
	// 	fmt.Println(err5)
	// }
	// fmt.Println(user.Username)
	// fmt.Println(has2)

	user := new(PreUcenterMembers)
	rows, err := engine.Where("uid <?", 6).Cols("username", "uid").Desc("uid").Rows(user)
	if err != nil {
	}
	defer rows.Close()
	var rowsData []*PreUcenterMembers
	for rows.Next() {
		row := new(PreUcenterMembers)
		rows.Scan(row)
		rowsData = append(rowsData, row)
	}
	fmt.Println(rowsData)
	for _, v := range rowsData {
		fmt.Println(v.Uid)
	}

	user1 := &PreUcenterMembers{Uid: 3}
	has6, err6 := engine.Get(user1)
	if err6 != nil {
		fmt.Println(err6, has6)
	}
	fmt.Println(user1)

	// user2 := new(PreUcenterMembers)
	// user2.Username = "myname1"
	// user2.Email = "goumins1@goumin.com"
	user2 := PreUcenterMembers{Username: "myname3", Email: "sss@goumin.com"}
	has7, err7 := engine.Where("uid=?", 1).Update(user2)
	fmt.Println(err7, has7)

	user4 := PreUcenterMembers{Username: "sdjdjdj", Email: "sds@goumin.com"}
	has9, err9 := engine.Insert(&user4)
	fmt.Println(err9, has9)

	// user3 := new(PreUcenterMembers)
	// var user3 PreUcenterMembers
	// has8, err8 := engine.Where("uid=?", 1).Cols("email", "username").Get(&user3)

	var user3 []PreUcenterMembers
	err8 := engine.Where("uid=?", 1).Cols("email", "username").Find(&user3)
	fmt.Println(user3, err8)

	report0 := Report{Cid: 2, Type: 3}
	has11, err11 := engine.Insert(&report0)
	fmt.Println(err11, has11)

	var report []Report
	// var uids string = "1,2,3"
	has10 := engine.In("id", 1, 2, 3).Find(&report)
	// has10 := engine.In("id", []int{1, 2, 3}).Find(&report)
	fmt.Println(has10, report)

	var user5 []PreUcenterMembers
	err12 := engine.In("uid", 1).Find(&user5)
	fmt.Println(user4, err12)

	var report6 []Report
	err13 := engine.Limit(10, 0).Find(&report6)
	fmt.Println(report6, err13)

	var report7 []Report
	err14 := engine.GroupBy("type").Find(&report7)
	fmt.Println(report7, err14)

	// var report8 []Report
	// affected, err15 := engine.Where("id=?", 1).Delete(&report8)
	// fmt.Println(err15, affected)

	report9 := new(Report)
	affected16, err16 := engine.Where("id=?", 2).Delete(report9)
	fmt.Println(affected16, err16)
}
