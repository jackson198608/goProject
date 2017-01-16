package main 

import (
    "fmt"
    // "math/rand"
    // "time"
    "log"
)

var dbAuth string = "root:goumintech"

//var dbAuth string = "dog123:dog123"

var dbDsn string = "192.168.86.72:3309"

//var dbDsn string = "192.168.5.199:3306"

var dbName string = "test_dz2"

var logger *log.Logger
var logPath string = "/var/spider.log"

func main() {
    // fmt.Println("hello,world!")
    // r := rand.New(rand.NewSource(time.Now().UnixNano()))
    // fmt.Println(r.Intn(100))
    // insertPost(9,2730909,47991496)
    // updateThread(9,2730909)
    // insertPost(2731136217,47993364)
    // postIsExist := checkPostExist(17,2731136217,47993364)
    // // tableExist := checkTableExist(107)
    // arr := getTaskData(1,10,10,0)
    // NewTask(2)
    // EventTask(1)
    // data := getNewTaskData(1,10,10,0)
    // fmt.Println(data)
    // count := getPostCount(1,20)
    fmt.Println("")
    // test()
    task()
}