package redisEngine

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v4"
)

const tryTimeLimit = 5

type coroutineResult struct {
	err error
}

type RedisEngine struct {
	queueName     string
	redisInfo     *redis.Options //require
	mongoConnInfo []string       //custom @todo need to multi
	mysqlInfo     []string       //the result format like tools.GetMysqlDsn return value,pass to task
	coroutinNum   int
	taskArgs      []string //somethin you want to give task
	workFun       func(job string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error
}

func NewRedisEngine(queueName string,
	redisInfo *redis.Options,
	mongoConnInfo []string,
	mysqlInfo []string,
	coroutinNum int,
	workFun func(job string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error,
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
	r.coroutinNum = coroutinNum
	r.workFun = workFun
	r.taskArgs = taskArgs

	return r, nil

}

// create redis connection and return it to the caller
func redisConnect(redisInfo *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(redisInfo)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// create xorms engines base on mysqlInfo and return it to the caller
func (r *RedisEngine) mysqlConnect() ([]*xorm.Engine, error) {
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
	if err != nil {
		return nil, err
	}
	return engine, nil
}

// create mongo session base on mongoConnInfo and return it to the caller
func (r *RedisEngine) mgoConnect() ([]*mgo.Session, error) {
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
	session, err := mgo.Dial(mgoInfo)
	if err != nil {
		return nil, err
	}

	return session, nil
}

//create several coroutin to do the job and controll the error is job fail
func (r *RedisEngine) Do() error {
	c := make(chan coroutineResult, r.coroutinNum)

	for i := 0; i < r.coroutinNum; i++ {

	}

	return nil
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

	//init redis client
	redisConn, err := redisConnect(r.redisInfo)
	if r.checkError(&result, c, err) {
		return
	}

	defer redisConn.Close()

	//get task data from redis,and invoke the callback fun
	for {
		//get task
		raw, err := redisConn.LPop(r.queueName).Result()
		if r.checkError(&result, c, err) {
			return
		}

		//get the raw parese result to decide whethere going to the next step or not
		realraw, trytimes, err := r.parseRaw(raw)
		if r.checkError(&result, c, err) {
			return
		}

		if trytimes > tryTimeLimit {
			result.err = errors.New("task over trytimes limit")
			c <- result
			return
		}

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

		//if goint to here ,call the invoke
		err = r.workFun(realraw, mysqlConns, mgoConns, r.taskArgs)
		if err != nil {
			if trytimes == tryTimeLimit {
				result.err = err
				c <- result
				return
			} else {
				err = r.pushFails(realraw, trytimes)
				if r.checkError(&result, c, err) {
					return
				}
			}
		}

	} //end of for loop

}

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
func (r *RedisEngine) pushFails(realraw string, trytimes int) error {
	return nil
}

func (r *RedisEngine) parseRaw(raw string) (string, int, error) {
	//maybe realraw may have the sep string,so we can not use strings.split
	rawSlice := []byte(raw)
	rawLen := len(rawSlice)
	for i := rawLen - 1; i >= 0; i-- {
		fmt.Println(rawSlice[i])
	}
	return "", 0, nil
}
