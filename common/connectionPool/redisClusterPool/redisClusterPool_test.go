package redisClusterPool

import (
	"fmt"
	"gouminGitlab/common/tools"
	"testing"
	"time"
)

const conum = 100

var redisClusterPool *redisPool
var redisConn = "192.168.86.193:6380,192.168.86.193:6381,192.168.86.193:6382,192.168.86.193:6383,192.168.86.193:6384,192.168.86.193:6385"

func TestIndex(t *testing.T) {
	redisClusterPool = createPool()
	// fmt.Println(redisClusterPool)
	for i := 0; i < 100; i++ {
		go work(i)
	}
}

func work(i int) {
	for {
		fmt.Println("i:", i)
		connection, err := redisClusterPool.GetConnection()
		fmt.Println("connection:", connection)
		if err != nil {
			fmt.Println("work redis connection err", err)
		}
		if connection == nil {
			time.Sleep(10 * time.Microsecond)
			continue
		} else {
			// connection.lpush()
			_, err := (*connection).LPush("messageDatas", "messageData").Result()
			if err != nil {
				fmt.Println("lpush err", err)
			}
			redisClusterPool.PutConnection(connection)
		}
	}
}

func createPool() *redisPool {
	redisInfo := tools.FormatRedisOption(redisConn)
	redisClusterPool, err := NewRedisPool(&redisInfo, conum)
	if err != nil {
		fmt.Println("new redis pool err", err)
		return redisClusterPool
		//to do
	}
	return redisClusterPool
}
