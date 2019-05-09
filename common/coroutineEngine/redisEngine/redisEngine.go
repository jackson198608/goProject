package redisEngine

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/tools"
	"github.com/olivere/elastic"
	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v4"
	"gouminGitlab/common/orm/elasticsearchBase"
	"strconv"
	"strings"
	"time"
)

const tryTimeLimit = 5
const sleepTime = 5

type coroutineResult struct {
	err error
}

type RedisEngine struct {
	queueName     string
	redisInfo     *redis.ClusterOptions //require
	mongoConnInfo []string              //custom @todo need to multi
	mysqlInfo     []string              //the result format like tools.GetMysqlDsn return value,pass to task
	esInfo        []string              //the result format like tools.GetElasticsearchNodes return value,pass to task
	coroutinNum   int
	daemon        int
	taskArgs      []string //somethin you want to give task
	workFun       func(job string, redisConn *redis.ClusterClient, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, esConn *elastic.Client, taskarg []string) error
}

func NewRedisEngine(queueName string,
	redisInfo *redis.ClusterOptions,
	mongoConnInfo []string,
	mysqlInfo []string,
	esInfo []string,
	coroutinNum int,
	daemon int,
	workFun func(job string, redisConn *redis.ClusterClient, tmysqlConns []*xorm.Engine, mgoConns []*mgo.Session, esConn *elastic.Client, taskarg []string) error,
	taskArgs ...string,
) (*RedisEngine, error) {
	//check param
	if (queueName == "") || (redisInfo == nil) || (coroutinNum <= 0) || (workFun == nil) {
		return nil, errors.New("params can not be null")
	}

	//create struct
	r := new(RedisEngine)
	if r == nil {
		return nil, errors.New("there is no more space for create new struct")
	}

	r.queueName = queueName
	r.redisInfo = redisInfo
	r.mysqlInfo = mysqlInfo
	r.mongoConnInfo = mongoConnInfo
	r.coroutinNum = coroutinNum
	r.workFun = workFun
	r.taskArgs = taskArgs
	r.daemon = daemon
	r.esInfo = esInfo

	return r, nil

}

// create redis connection and return it to the caller
func redisConnect(redisInfo *redis.ClusterOptions) (*redis.ClusterClient, error) {
	client, err := tools.GetClusterClient(redisInfo)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// create xorms engines base on mysqlInfo and return it to the caller
func (r *RedisEngine) mysqlConnect() ([]*xorm.Engine, error) {

	//if you do not need mysql for job func
	if r.mysqlInfo == nil {
		return nil, nil
	}

	// if you need make it for you, and info must be correct
	mysqls := []*xorm.Engine{}
	for _, mysqlInfo := range r.mysqlInfo {
		x, err := r.mysqlSingleConnect(mysqlInfo)
		if err != nil {
			//close former connection
			r.closeMysqlConn(mysqls)
			return nil, err
		}
		mysqls = append(mysqls, x)
	}
	return mysqls, nil

}
func (r *RedisEngine) mysqlSingleConnect(mysqlInfo string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", mysqlInfo)
	err = engine.Ping()
	if err != nil {
		return nil, err
	}
	return engine, nil
}

// create mongo session base on mongoConnInfo and return it to the caller
func (r *RedisEngine) mgoConnect() ([]*mgo.Session, error) {
	//if you do not need mongo for job func ,just return
	if r.mongoConnInfo == nil {
		return nil, nil
	}

	// if you need mongo connection for job func ,make sure info you have is correct
	mgos := []*mgo.Session{}
	for _, mgoInfo := range r.mongoConnInfo {
		m, err := r.mgoSingleConnect(mgoInfo)
		if err != nil {
			//close former connection
			r.closeMgoConn(mgos)
			return nil, err
		}
		mgos = append(mgos, m)
	}
	return mgos, nil

}

func (r *RedisEngine) mgoSingleConnect(mgoInfo string) (*mgo.Session, error) {
	var session *mgo.Session
	var err error
	mgoInfos := strings.Split(mgoInfo, ",")
	if len(mgoInfos) == 1 {
		session, err = tools.GetStandAloneConnecting(mgoInfo)
	} else {
		session, err = tools.GetReplicaConnecting(mgoInfos)
	}
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *RedisEngine) esConnect() (*elastic.Client, error) {
	if r.esInfo == nil {
		return nil, nil
	}
	esR, _ := elasticsearchBase.NewClient(r.esInfo)
	client, err := esR.Run()
	if err != nil {
		return nil, nil
	}
	return client, nil
}

//create several coroutin to do the job and controll the error is job fail
//@todo make error to be []error
func (r *RedisEngine) Do() error {
	c := make(chan coroutineResult, r.coroutinNum)

	tempResult := coroutineResult{err: nil}
	lastResult := coroutineResult{err: nil}

	for i := 0; i < r.coroutinNum; i++ {
		go r.coroutinFunc(c, i)
	}

	for i := 0; i < r.coroutinNum; i++ {
		tempResult = <-c
		if tempResult.err != nil {
			lastResult.err = tempResult.err
		}
	}

	return lastResult.err
}

func (r *RedisEngine) checkError(result *coroutineResult, c chan coroutineResult, err error) bool {
	if err != nil {
		result.err = err
		c <- *result
		return true
	}

	return false
}

func (r *RedisEngine) coroutinFunc(c chan coroutineResult, i int) {
	//create result struct
	result := coroutineResult{
		err: nil,
	}

	//@todo  connection fail need to be retry
	//init redis client
	redisConn, err := redisConnect(r.redisInfo)
	if r.checkError(&result, c, err) {
		return
	}

	defer redisConn.Close()

	//prepare and check the connections for mysql
	mysqlConns, err := r.mysqlConnect()
	if r.checkError(&result, c, err) {
		return
	}

	defer r.closeMysqlConn(mysqlConns)

	//prepare and check the connections for mgo
	mgoConns, err := r.mgoConnect()
	if r.checkError(&result, c, err) {
		return
	}

	defer r.closeMgoConn(mgoConns)

	esConn, err := r.esConnect()
	if r.checkError(&result, c, err) {
		return
	}

	//get task data from redis,and invoke the callback fun
	for {
		//get task
		raw, err := redisConn.LPop(r.queueName).Result()
		if err != nil {
			if r.daemon == 0 {
				c <- result
				return
			}
			//there is no more job,sleep a while
			time.Sleep(sleepTime * time.Second)
			continue
		}

		//get the raw parese result to decide whethere going to the next step or not
		realraw, trytimes, err := r.parseRaw(raw)
		if err != nil {
			fmt.Println("[error] parseRaw error ,goint to next and drop the current one: ", err, raw)
			continue
		}

		if trytimes >= tryTimeLimit {
			fmt.Println("[error] retry overseed ,drop it and continue", raw)
			continue
		}

		redisConn, err = r.checkAndRepairRedisConnectins(redisConn)

		mysqlConns, err = r.checkAndRepairMysqlConnectins(mysqlConns)

		//esConn, esErr :=r.checkAndRepairEsConnectins(esConn)

		//only if repair connection sucess  going to finish the workfun
		if err == nil {
			err = r.workFun(realraw, redisConn, mysqlConns, mgoConns, esConn, r.taskArgs)
		}

		if err != nil {
			fmt.Println("[error]jobFunc get error ,but still can be retry", err)
			err = r.pushFails(redisConn, realraw, trytimes)
			if err != nil {
				fmt.Println("[error]pushFails error ,drop it", err)
				err = r.pushFails(redisConn, realraw, trytimes)
				continue
			}
		}

	} //end of for loop

	if result.err != nil {
		fmt.Println("[error] coroutine exit with error: ", result.err)
	}
	c <- result
	return

}

//检查redis链接, 链接异常时重连
func (r *RedisEngine) checkAndRepairRedisConnectins(redisConn *redis.ClusterClient) (*redis.ClusterClient, error) {
	_, err := redisConn.Ping().Result()
	if err != nil {
		redisConn, err := redisConnect(r.redisInfo)
		if err != nil {
			return nil, err
		}
		defer redisConn.Close()
	}
	return redisConn, nil
}

//检查mysql链接, 链接异常时重连
func (r *RedisEngine) checkAndRepairMysqlConnectins(mysqlConns []*xorm.Engine) ([]*xorm.Engine, error) {
	//for _, mysqlInfo := range r.mysqlInfo {
	for i, mysqlConn := range mysqlConns {
		err := mysqlConn.Ping()
		if err != nil {
			x, err := r.mysqlSingleConnect(r.mysqlInfo[i])
			if err != nil {
				//close former connection
				r.closeMysqlConn(mysqlConns)
				return nil, err
			}
			mysqlConns[i] = x
			defer r.closeMysqlConn(mysqlConns)
		}
	}
	return mysqlConns, nil
}

//检查es链接, 链接异常时重连
//func (r *RedisEngine) checkAndRepairEsConnectins(esConn *elastic.Client) (*elastic.Client,error) {
//	return esConn,nil
//}

//@todo
func (r *RedisEngine) closeMysqlConn(mysqlConns []*xorm.Engine) {
	if mysqlConns == nil {
		return
	}
	for _, conn := range mysqlConns {
		conn.Close()
	}
	return
}

// @todo
func (r *RedisEngine) closeMgoConn(mgoConns []*mgo.Session) {
	if mgoConns == nil {
		return
	}

	for _, conn := range mgoConns {
		conn.Close()
	}

	return
}

//if trytimes < tryTimeLimit ,just push it back to redis
//if push fail ,it will return error
func (r *RedisEngine) pushFails(redisConn *redis.ClusterClient, realraw string, tryTimes int) error {
	//@todo check params
	backRaw := realraw + "_" + strconv.Itoa(tryTimes+1)
	redisConn.RPush(r.queueName, backRaw)
	return nil
}

// tryTimes only suport 0-9 ,if >9 ,the function should be overwritted
func (r *RedisEngine) parseRaw(raw string) (string, int, error) {
	//maybe realraw may have the sep string,so we can not use strings.split
	rawSlice := []byte(raw)
	rawLen := len(rawSlice)
	//当是用户推荐时，raw是uid,长度可能会<2
	if rawLen < 2 {
		return raw, 0, nil
	}
	if (rawSlice[rawLen-2] == '_') && (rawLen > 2) {
		tryTimesStr := string(rawSlice[rawLen-1])
		tryTimesInt, err := strconv.Atoi(tryTimesStr)
		if err != nil {
			return raw, 0, nil
		} else {
			return string(rawSlice[0 : rawLen-2]), tryTimesInt, nil
		}
	} else {
		return raw, 0, nil
	}

}
