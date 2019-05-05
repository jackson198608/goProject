package main

import (
	"errors"
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/jackson198608/goProject/common/tools"
	"github.com/jackson198608/goProject/pushContentCenter/task"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v4"
	"os"
	// "strconv"
	"fmt"
	"strings"
	"github.com/olivere/elastic"
)

var c Config = Config{
	"192.168.86.193:3307", //mysql dsn
	"new_dog123",          //mysql dbName
	"card",
	"adoption",//mysql dbName
	"member",
	"dog123:dog123",       //mysqldbAuth
	"127.0.0.1:6379",      //redis info
	1,                     //thread num
	"pushContentCenter",   //queuename
	//"192.168.86.192:27017",
	"http://192.168.86.230:9200,http://192.168.86.231:9200"} // mongo

func init() {
	loadConfig()
}

func main() {
	//var clubId string

	params := os.Args[1]
	jobType := params
	//jobTypeclubId := strings.Split(params, "_")
	//if len(jobTypeclubId) == 2 {
	//	jobType = jobTypeclubId[0]
	//	clubId = jobTypeclubId[1]
	//}
	switch jobType {
	case "push": //push content conter
		var mongoConnInfo []string
		//mongoConnInfo = append(mongoConnInfo, c.mongoConn)
		var mysqlInfo []string
		mysqlInfo = append(mysqlInfo, c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")
		if c.dbName1 != "" {
			mysqlInfo = append(mysqlInfo, c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName1+"?charset=utf8mb4")
		}
		if c.dbName2 != "" {
			mysqlInfo = append(mysqlInfo, c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName2+"?charset=utf8mb4")
		}
		if c.dbName3 != "" {
			mysqlInfo = append(mysqlInfo, c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName3+"?charset=utf8mb4")
		}
		esNodes := strings.SplitN(c.elkNodes, ",", -1)

		redisInfo := tools.FormatRedisOption(c.redisConn)
		logger.Info("start work")
		r, err := redisEngine.NewRedisEngine(c.queueName, &redisInfo, mongoConnInfo, mysqlInfo, esNodes, c.coroutinNum, 1, jobFuc)
		if err != nil {
			logger.Error("[NewRedisEngine] ", err)
		}

		err = r.Do()
		if err != nil {
			logger.Error("[redisEngine Do] ", err)
		}
	//case "allindex": // all collection create index
	//	err := ClubData.AllCollectionsCreateIndex(c.mongoConn)
	//	if err != nil {
	//		logger.Error("all collections create index error! ", err)
	//	}
	//case "singleindex": // single collection create index
	//	err := ClubData.SingleCollectionCreateIndex(nil, "forum_content_"+clubId, c.mongoConn)
	//	if err != nil {
	//		logger.Error("single collection create index error by clubId is ", clubId, err)
	//	}
	case "--help":
		help()
	default:
		fmt.Println("unsupported params")
	}
}

func jobFuc(job string, redisConn *redis.ClusterClient, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, esConn *elastic.Client, taskarg []string) error {
	if (mysqlConns == nil) || (esConn==nil){
		return errors.New("mysql or mongo or elastic conn error")
	}
	t, err := task.NewTask(job, mysqlConns, mgoConns, esConn)
	if err != nil {
		return err
	}
	err = t.Do()
	if err != nil {
		return err
	}
	return err
}

func help() {
	fmt.Println("usage: pushApp [options]")
	fmt.Println("Options:")
	fmt.Println("  allindex\t\t\t\tThe index of all collections is deleted, and then a new index is created")
	fmt.Println("  singleindex_clubid\t\t\tSpecify a collection to create an index")
	fmt.Println("  push\t\t\t\t\tpush data")
	fmt.Println("  --help\t\t\t\tshow this usage information")
}
