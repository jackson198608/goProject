package main

import(
    "github.com/donnie4w/go-logger/logger"
    // "github.com/jackson198608/goProject/eventLog/task"
    "os"
    "fmt"
    // "log"
)

// var dbAuth string = "root:goumintech"
// var dbDsn string = "192.168.86.72:3309"
// var dbName string = "test_dz2"
// var logger *log.Logger
// var logPath string = "/tmp/spider.log"
var c Config = Config{
    "192.168.86.72:3309",
    "test_dz2",
    "root:goumintech",
    1,
    10,//2545,
    1,
    "127.0.0.1:6379",
    "moveEvent",
    "/tmp/moveEvent.log", 0}

func pushALLEventIdFromStartToEnd() {
    r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName)
    page := 0
    for {
        ids := getTask(page)
        if len(ids) == 0 {
            break
        }
        if ids == nil {
            break
        }
        r.PushTaskData(ids)
        page++
    }
}

func do() {
    r := NewRedisEngine(c.logLevel, c.queueName, c.redisConn, "", 0, c.numloops, c.dbAuth, c.dbDsn, c.dbName)
    r.Loop()
}

func Init() {

    loadConfig()
    logger.SetConsole(true)
    logger.SetLevel(logger.DEBUG)
    logger.Error(logger.DEBUG)

}
func main() {
    Init()
    // data := getEventLogData(1,10,10,0)
    // fmt.Println(data)
    // NewTask(1)
    jobType := os.Args[1]
    fmt.Println(jobType)
    switch jobType {
    case "create":
        logger.Info("in the create", 10)
        pushALLEventIdFromStartToEnd()

    case "do":
        logger.Info("in the do")
        do()
    default:

    }
}