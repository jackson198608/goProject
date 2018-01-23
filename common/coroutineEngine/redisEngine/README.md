# 1. introduction
- this package is design for the situation when you want to  parallel your code with redis queue 

# 2. Installation

	go get github.com/jackson198608/goProject/common/coroutineEngine/redisEngine	

# 3. Quick Start

## 3.1 Create redisEngine 
- engine needs to Know: queueName ,redisInfo, xormEngines,mongoConnections,threadNum,JobFun
- JobFun will be invoked in coroutine which created inside the engine

```Go
func newtask() (*RedisEngine, error) {

	redisInfo := redis.Options{
		Addr: "127.0.0.1:6379",
	}
	//getXormEngine
	connStr := tools.GetMysqlDsn(dbAuth, dbDsn, dbName)
	conns := []string{connStr}

	//get mongo session
	mgos := []string{mongoConn}

	r, err := NewRedisEngine("test", &redisInfo, conns, mgos, 3, 1, jobFunc)
	if err != nil {
		return nil, err
	}
	return r, nil

}

func jobFunc(job string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {

	fmt.Println("this is jobFunc", job)
	return errors.New("job func fail")
}

```

### 3.1.1 new params 
#### 3.1.1.1 queueName
- tell redisEngine which queueName you can got the job 
- the queue should be list type in redis

#### 3.1.1.2 redisInfo
- tell redisEngine the connect info of redis 
- note that: just info not the real connection	

#### 3.1.1.3 xormEngines
- type : the slice  []string 
- each element represent each conn info for mysql's db
- each element shoud be the format like : dbAuth + "@tcp(" + dbDsn + ")/" + dbName + "?charset=utf8mb4"

#### 3.1.1.4 mongos
- type: the slice []string
- each element represent each conn info for mongo server
- each element shoud be the format like "127.0.0.1:27017" or "192.168.1.1:27017,192.168.1.2:27017,192.168.1.3:27017" for replica set

#### 3.1.1.5 thread Num
- just tell engine how many thread you want to create to call the jobFunc

#### 3.1.1.6 taskArgs
- type: ...string
- user define string params which will by pass to jobFunc directly

#### 3.1.1.7 jobFunc
- the format like: func(job string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error 
- job string:  pop one from queueName of redis Server 
- []*xorm.Engine:   engine just convert xormEngines params to []*xorm.Engine ,change every string element to xorm.Engine element and pass to jobFunc with the same order
- []*mgo.Session:   engine just convert mongos params to []*mgo.Session ,change every string element to mgo.Session element and pass to jobFunc with the same order
- taskarg []string: which passed by taskArgs... string

## 3.2 do the job

```Go
	err:=r.Do()
```

### 3.2.1 job logic
- you can see it to be graph in some other markdown editor

```sequence
caller->redisEngine: new and pass params
redisEngine->redisEngine: new Object and pass params to  private member
redisEngine->caller: object Engine
caller->redisEngine: invoke Do()
redisEngine->redisEngineCoroutine:  create ThreadNum's coroutine and wait for result ,get the coroutint return error
redisEngineCoroutine-> redisEngineCoroutine:  create redis connection by redisinfo,if can not connect return error
redisEngineCoroutine-> redisEngineCoroutine:  for begin
redisEngineCoroutine-> redisEngineCoroutine:  pop raw from redis queue,if no job ,sleep for 5 seconds
redisEngineCoroutine-> redisEngineCoroutine:  parse if find realRaw and trytimes
redisEngineCoroutine-> redisEngineCoroutine:  if trytimes > limit  {break , return error}
redisEngineCoroutine-> redisEngineCoroutine:  change mysqlConninfo to []*xorm.Engine,if connect error {break ,return error}
redisEngineCoroutine-> redisEngineCoroutine:  change mongoConnInfo to []*mgo.Session,if connect error {break ,return error}
redisEngineCoroutine-> jobFunc: call jobFunc with realRaw,[]*xorm.Engine,[]*mgo.Session ,taskArgs
jobFunc-> redisEngineCoroutine: do job return status ,if error return error
redisEngineCoroutine-> redisEngineCoroutine: if job error and not reach limit ,push it back to queue with trytimes++
redisEngineCoroutine-> redisEngineCoroutine: else break return error
redisEngineCoroutine-> redisEngineCoroutine: for end
```
