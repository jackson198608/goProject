package redisClusterPool

import (
	"errors"
	"fmt"
	queue "github.com/oleiade/lane"
	log "github.com/thinkboy/log4go"
	redis "gopkg.in/redis.v4"
	"gouminGitlab/common/tools"
	"sync"
	"time"
)

var loopnum = 3

type RedisPool struct {
	redisInfo *redis.ClusterOptions
	num       int
	poolQueue *queue.Queue
	lock      sync.RWMutex
}

func NewRedisPool(redisInfo *redis.ClusterOptions, num int) (*RedisPool, error) {
	//@todo check params
	redisPool := new(RedisPool)
	if redisPool == nil {
		return redisPool, errors.New("create memory error.can not create it")
	}
	redisPool.redisInfo = redisInfo
	redisPool.num = num
	err := redisPool.initPool()
	if err != nil {
		return redisPool, err
	}

	return redisPool, nil

}

func (r *RedisPool) initPool() error {
	r.poolQueue = queue.NewQueue()
	for i := 0; i < r.num; i++ {
		connection, err := r.createOneConnection()
		// fmt.Println(i, connection, reflect.TypeOf(connection))
		if err != nil {
			//@todo something
			fmt.Println("err create connection", err)
			return err
		} else {
			r.poolQueue.Enqueue(connection)
		}
	}
	return nil
}

func (r *RedisPool) createOneConnection() (*redis.ClusterClient, error) {
	client, err := tools.GetClusterClient(r.redisInfo)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r *RedisPool) Close() error {
	for i := 0; i < r.poolQueue.Size(); i++ {
		connection := r.poolQueue.Dequeue()
		connection.(*redis.ClusterClient).Close()
		//@todo close
	}
	return nil
}

func (r *RedisPool) GetConnection() (*redis.ClusterClient, error) {
	r.lock.Lock()
	fmt.Println("poolQueue size = ", r.poolQueue.Size())
	if r.poolQueue.Size() == 0 {
		r.lock.Unlock()
		return nil, errors.New("there is no more connection can be use ,please wait")
	}
	connection := r.poolQueue.Dequeue()
	r.lock.Unlock()
	return connection.(*redis.ClusterClient), nil

}

func (r *RedisPool) TellMeOneIsBroken() {
	for i := 0; i < loopnum; i++ {
		fmt.Println("tell me times ", i)
		connection, err := r.createOneConnection()
		if err != nil {
			log.Error("create connection to enqueue error", err)
			time.Sleep(10 * time.Millisecond)
			continue
		}
		_, err1 := connection.Ping().Result()
		if err1 != nil {
			log.Error("again create new connection to enqueue error", err)
			time.Sleep(10 * time.Millisecond)
			continue
		}
		err = r.PutConnection(connection)
		if err != nil {
			log.Error("put connection to enqueue error", err)
			continue
		} else {
			break
		}
	}

}

func (r *RedisPool) PutConnection(connection *redis.ClusterClient) error {
	r.lock.Lock()
	r.poolQueue.Enqueue(connection)
	r.lock.Unlock()
	return nil
}
