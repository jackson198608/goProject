package redisClusterPool

import (
	"errors"
	queue "github.com/oleiade/lane"
	redis "gopkg.in/redis.v4"
	"sync"
)

type redisPool struct {
	redisInfo redis.ClusterOptions
	num       int
	poolQueue queue.Queue
	lock      sync.RWMutex
}

func NewRedisPool(redisInfo redis.ClusterOptions, num int) (redisPool, error) {
	//@todo check params
	redisPool := new(redisPool)
	if redisPool == nil {
		return nil, errors.New("create memory error.can not create it")
	}

	err := redisPool.initPool()
	if err != nil {
		return nil, err
	}

	return redisPool, nil

}

func (r *redisPool) initPool() error {
	for i := 0; i < r.num; i++ {
		connection, err := r.createOneConnection()
		if err != nil {
			//@todo something
		} else {
			r.poolQueue.Enqueue()
		}
	}
}

func (r *redisPool) createOneConnection() (*redis.ClusterClient, error) {

}

func (r *redisPool) Close() error {
	for i := 0; i < r.poolQueue.Size(); i++ {
		connection := r.poolQueue.Dequeue()
		//@todo close
	}
}

func (r *redisPool) GetConnection() (*redis.ClusterClient, error) {
	r.lock.Lock()
	defer r.lock.RLock()
	if r.poolQueue.Size() == 0 {
		return nil, errors.New("there is no more connection can be use ,please wait")
	} else {
		connection := r.poolQueue.Dequeue()
		return connection, nil
	}

}

func (r *redisPool) PutConnection(connection *redis.ClusterClient) error {

	r.lock.Lock()
	r.poolQueue.Enqueue(connection)
	r.lock.RLock()
	return nil

}
