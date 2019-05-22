package esClusterPool

import (
	"testing"
	"fmt"
	"time"
	"strings"
	"gouminGitlab/common/orm/elasticsearch"
)

const conum = 10

var ep *EsPool
var esConn = "http://192.168.86.230:9200,http://192.168.86.231:9200"

func TestIndex(t *testing.T) {
	ep = createPool()
	if ep == nil {
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
		connection, err := ep.GetConnection()
		if err != nil {
			fmt.Println("work redis connection err", err)
		}
		if connection == nil {
			fmt.Println("get no connection")
			time.Sleep(10 * time.Microsecond)
			continue
		} else {
			er,err := elasticsearch.NewUserInfo(connection)
			if err!=nil {
				fmt.Println(" NewUserInfo err ", er)
			}
			var uids []int
			uids = append(uids, 2060500)
			rst,err := er.GetActiveUserInfoByUids(uids, 0, 100)
			if err!=nil {
				fmt.Println(" GetActiveUserInfoByUids err ", er)
			}
			fmt.Println(rst)
			ep.PutConnection(connection)
			time.Sleep(5000 * time.Microsecond)
		}
	}
	c <- 1
}

func createPool() *EsPool {
	esNodes := strings.SplitN(esConn, ",", -1)
	ep, err := NewEsPool(esNodes, conum)
	if err != nil {
		fmt.Println("new redis pool err", err)
		return nil
		//to do
	}
	return ep
}