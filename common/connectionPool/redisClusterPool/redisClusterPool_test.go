package redisClusterPool

import (
	"fmt"
	"gouminGitlab/common/tools"
	"testing"
	"time"
)

const conum = 10

var rp *redisPool
var redisConn = "192.168.86.193:6380,192.168.86.193:6381,192.168.86.193:6382,192.168.86.193:6383,192.168.86.193:6384,192.168.86.193:6385"

func TestIndex(t *testing.T) {
	rp = createPool()
	if rp == nil {
		fmt.Println("create redis pool error")
	}
	// fmt.Println(redisClusterPool)
	c := make(chan int, 10)
	for i := 0; i < 10; i++ {
		go work(i, c)
	}

	for i := 0; i < 10; i++ {
		<-c
	}

}

func work(i int, c chan int) {
	for {
		fmt.Println("work i ", i)
		connection, err := rp.GetConnection()
		if err != nil {
			fmt.Println("work redis connection err", err)
		}
		if connection == nil {
			fmt.Println("get no connection")
			time.Sleep(10 * time.Microsecond)
			continue
		} else {
			// connection.lpush()
			_, err := (*connection).LPush("messageDatas", "messageData").Result()
			if err != nil {
				fmt.Println("lpush err", err)
			}
			rp.PutConnection(connection)
			time.Sleep(5000 * time.Microsecond)
		}
	}
	c <- 1
}

func createPool() *redisPool {
	redisInfo := tools.FormatRedisOption(redisConn)
	rp, err := NewRedisPool(&redisInfo, conum)
	if err != nil {
		fmt.Println("new redis pool err", err)
		return nil
		//to do
	}
	return rp
}
