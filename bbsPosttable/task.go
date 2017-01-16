package main

import (
    // "strconv"
    // "strings"
    "fmt"
)

// type Task struct {
//     phoneType   int //o for ios ,1 for android
//     DeviceToken string
//     TaskJson    string
// }

// func NewTask(redisString string) (t *Task) {
//     var tR Task
//     result := strings.Split(redisString, "|")
//     if len(result) != 3 {
//         return nil
//     }
//     tR.phoneType, _ = strconv.Atoi(result[0])
//     tR.DeviceToken = result[1]
//     tR.TaskJson = result[2]
//     return &tR
// }

func NewTask(taskNum int) {
    var startPid int = 30
    var endPid int = 60
    var limit int = 100
    var offset int = 0

    count := getPostCount(startPid,endPid)
    for {
        task := getNewTaskData(startPid,endPid,limit,offset)
        if len(task) == 0 {
            fmt.Println("task data is empty")
            return
        }
        for _,v := range task {
            insertIntoPost(v.tid,v.pid)
        }
        offset += limit
        if offset > count {
            break;
        }
        // fmt.Println(task)
    }
}

func insertIntoPost(tid int, pid int){
    tableid := tid%100
    if tableid == 0 {
        tableid = 100
    }
    tableIsExist := checkTableExist(tableid)
    if tableIsExist == false {
        createTable(tableid)
    }
    updateThread(tableid,tid)
    postIsExist := checkPostExist(tableid,tid,pid)
    if postIsExist == false {
        insertPost(tableid,tid,pid)
    }
}
