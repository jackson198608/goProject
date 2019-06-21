package main

import (
	"github.com/jackson198608/goProject/appPush"
	"os"
	"time"
	"strings"
	"github.com/jackson198608/goProject/common/tools"
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
	log "github.com/thinkboy/log4go"
	"gopkg.in/redis.v4"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/jackson198608/goProject/gouminMultiMessagePush/gouminMessagePush/task"
	"io/ioutil"
)

//define the config var
var c Config = Config{5, 100, "192.168.86.80:6380,192.168.86.80:6381,192.168.86.81:6380,192.168.86.81:6380,192.168.86.82:6380,192.168.86.82:6381", "http://192.168.86.230:9200,http://192.168.86.231:9200","/Users/Snow/Work/go/config/gouminMultiMessagePushLog.xml"}
var numForOneLoop int = c.currentNum
var p12Bytes []byte
var timeout time.Duration = c.httpTimeOut

/*
type of job
	multi: multi app push
	single: single app push
	insert: insert data into mongo
*/
var jobType string = "multi"
var redisQueueName string = "mcMulti"

//define the tasks array for each loop
var tasks []redisData
var taskNum int = 0

func Init() {
	getRedisQueueName()
	cBytes, err := ioutil.ReadFile("/etc/pro-lingdang.pem")
	if err != nil {
		return
	}
	p12Bytes = cBytes

	loadConfig()

	appPush.Init(timeout)
}

func getRedisQueueName() {
	switch os.Args[1] {
	case "multi":
		redisQueueName = "mcMulti"
	case "single":
		redisQueueName = "mcSingle"
	case "insert":
		redisQueueName = "mcInsert"

	default:
		redisQueueName = "mcMulti"
	}
}


func main() {

	//init the system process
	Init()
	jobType = os.Args[1]

	//初始化日志配置
	log.LoadConfiguration(c.log)


	var mongoConnInfo []string
	var mysqlInfo []string

	esNodes := strings.SplitN(c.elasticConn, ",", -1)

	redisInfo := tools.FormatRedisOption(c.redisConn)

	r, err := redisEngine.NewRedisEngine(redisQueueName, &redisInfo, mongoConnInfo, mysqlInfo, esNodes, c.currentNum, 1, jobFuc)
	if err != nil {
		log.Error("[NewRedisEngine] ", err)
	}

	err = r.Do()
	if err != nil {
		log.Error("[redisEngine Do] ", err)
	}

}

func jobFuc(job string, redisConn *redis.ClusterClient, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, esConn *elastic.Client, taskarg []string) error {
	if (redisConn == nil) || (esConn==nil){
		return errors.New("redis or elastic conn error")
	}
	esNodes := strings.SplitN(c.elasticConn, ",", -1)
	t, err := task.NewTask(jobType,job,redisConn,esConn,p12Bytes, esNodes[0])
	if err != nil {
		return err
	}
	err = t.Do()
	if err != nil {
		return err
	}
	return err
}
