package redisClusterPool

import (
	"testing"
	"time"
)

const conum = 100

var redisClusterPool *redisClusterPool

func TestIndex(t *testing.T) {
	redisClusterPool = createPool()
	for i := 0; i < 100; i++ {
		go work(i)
	}
}

func work(i int) {
	for {
		connection, error := redisClusterPool.GetConnection()
		if connection == nil {
			time.Sleep(10 * time.Microsecond)
			continue
		} else {
			connection.lpush()
			redisClusterPool.PutConnection(connection)
		}
	}
}

func createPool() *redisClusterPool {
	redisClusterPool := NewRedisPool()
	return redisClusterPool
}
