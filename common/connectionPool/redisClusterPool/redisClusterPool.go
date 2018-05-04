package redisClusterPool

import (
	"errors"
	"fmt"
	queue "github.com/oleiade/lane"
	redis "gopkg.in/redis.v4"
	"gouminGitlab/common/tools"
	"reflect"
	"sync"
)

type redisPool struct {
	redisInfo *redis.ClusterOptions
	num       int
	poolQueue queue.Queue
	lock      sync.RWMutex
}

func NewRedisPool(redisInfo *redis.ClusterOptions, num int) (*redisPool, error) {
	//@todo check params
	redisPool := new(redisPool)
	if redisPool == nil {
		return redisPool, errors.New("create memory error.can not create it")
	}
	redisPool.redisInfo = redisInfo
	err := redisPool.initPool()
	if err != nil {
		return redisPool, err
	}

	return redisPool, nil

}

func (r *redisPool) initPool() error {
	for i := 0; i < r.num; i++ {
		connection, err := r.createOneConnection()
		fmt.Println(i, connection, reflect.TypeOf(connection))
		if err != nil {
			//@todo something
		} else {
			r.poolQueue.Enqueue(connection)
		}
	}
	return nil
}

func (r *redisPool) createOneConnection() (*redis.ClusterClient, error) {
	client, err := tools.GetClusterClient(r.redisInfo)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r *redisPool) Close() error {
	for i := 0; i < r.poolQueue.Size(); i++ {
		connection := r.poolQueue.Dequeue()
		connection.(*redis.ClusterClient).Close()
		//@todo close
	}
	return nil
}

func (r *redisPool) GetConnection() (*redis.ClusterClient, error) {
	r.lock.Lock()
	defer r.lock.RLock()
	fmt.Println("get connection ", r)
	if r.poolQueue.Size() == 0 {
		return nil, errors.New("there is no more connection can be use ,please wait")
	}
	connection := r.poolQueue.Dequeue()
	fmt.Println("get connection ", connection)
	return connection.(*redis.ClusterClient), nil

}

func (r *redisPool) PutConnection(connection *redis.ClusterClient) error {

	r.lock.Lock()

	r.poolQueue.Enqueue(connection)
	r.lock.RLock()
	return nil

}
