package main

import (
	"errors"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/recommend_1.0/task"
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/jackson198608/goProject/common/tools"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v4"
	"os"
	"strconv"
	"gouminGitlab/common/orm/elasticsearch"
	"strings"
	"github.com/olivere/elastic"
	"gouminGitlab/common/orm/elasticsearchBase"
)

var c Config = Config{
	"192.168.86.193:3307", //mysql dsn
	"new_dog123",          //mysql dbName
	"dog123:dog123",       //mysqldbAuth
	"127.0.0.1:6379",      //redis info
	1,                     //thread num
	"pushContentCenter",   //queuename
	//"192.168.86.192:27017",
	"192.168.86.230:9200,192.168.5.231:9200"} // mongo

func init() {
	loadConfig()
}

func pushAllActiveUserToRedis(queueName string, channel string) bool {
	redisInfo := tools.FormatRedisOption(c.redisConn)
	rc, _ := tools.GetClusterClient(&redisInfo)
	if rc == nil {
		fmt.Println("error")
		return false
	}
	nodes := strings.SplitN(c.elkDsn, ",", -1)
	r,_ := elasticsearchBase.NewClient(nodes)
	esConn ,_ := r.Run()

	er,err := elasticsearch.NewUser(esConn)
	if err !=nil {
		return false
	}
	from := 0
	size := 100
	i :=1
	for {
		rst := er.SearchAllActiveUser(from, size)
		total := rst.Hits.TotalHits
		if total> 0 {
			for _, hit := range rst.Hits.Hits {
				uid,_ := strconv.Atoi(hit.Id)
				err := rc.RPush(queueName, strconv.Itoa(uid)+"|"+channel).Err()
				if err != nil {
					logger.Error("insert redis error", err)
					return false
				}
			}
		}
		if int(total) < from {
			break
		}
		i++
		from = (i-1)*size
	}
	return true
}

func main() {
	jobType := os.Args[1]
	switch jobType {
	case "recommend": //push content conter
		var mongoConnInfo []string
		//mongoConnInfo = append(mongoConnInfo, c.mongoConn)
		var mysqlInfo []string
		mysqlInfo = append(mysqlInfo, c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")

		redisInfo := tools.FormatRedisOption(c.redisConn)

		//生产任务
		tastStatus := pushAllActiveUserToRedis(c.queueName, "follow")
		if !tastStatus {
			logger.Error("[NewRedisEngine] ", errors.New("create task fail"))
		}
		nodes := strings.SplitN(c.elkDsn, ",", -1)

		logger.Info("start work")
		r, err := redisEngine.NewRedisEngine(c.queueName, &redisInfo, mongoConnInfo, mysqlInfo,nodes, c.coroutinNum, 0, jobFuc, c.elkDsn)
		if err != nil {
			logger.Error("[NewRedisEngine] ", err)
		}

		err = r.Do()
		if err != nil {
			logger.Error("[redisEngine Do] ", err)
		}
	case "content": //push content conter
		var mongoConnInfo []string
		//mongoConnInfo = append(mongoConnInfo, c.mongoConn)
		var mysqlInfo []string
		mysqlInfo = append(mysqlInfo, c.dbAuth+"@tcp("+c.dbDsn+")/"+c.dbName+"?charset=utf8mb4")

		redisInfo := tools.FormatRedisOption(c.redisConn)
		//生产任务
		pushAllActiveUserToRedis(c.queueName, "content")
		nodes := strings.SplitN(c.elkDsn, ",", -1)

		logger.Info("start work")
		r, err := redisEngine.NewRedisEngine(c.queueName, &redisInfo, mongoConnInfo, mysqlInfo,nodes, c.coroutinNum, 0, jobFuc, c.elkDsn)
		if err != nil {
			logger.Error("[NewRedisEngine] ", err)
		}

		err = r.Do()
		if err != nil {
			logger.Error("[redisEngine Do] ", err)
		}
	//case "allindex": // all collection create index
	//	err := RecommendData.AllCollectionsCreateIndex(c.mongoConn)
	//	if err != nil {
	//		logger.Error("all collections create index error! ", err)
	//	}
	case "--help":
		help()
	default:
		fmt.Println("unsupported params")
	}
}

func jobFuc(job string,redisConn *redis.ClusterClient, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, esConn *elastic.Client,taskarg []string) error {
	if (mysqlConns == nil) || (esConn == nil){
		return errors.New("mysql or mongo or esConn conn error")
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
	// fmt.Println("  singleindex_userId\t\t\tSpecify a collection to create an index")
	fmt.Println("  recommend\t\t\t\t\tRecommending related clubs and unconcerned users")
	fmt.Println("  content\t\t\t\t\tRecommending a selection of content to the user")
	fmt.Println("  --help\t\t\t\tshow this usage information")
}
