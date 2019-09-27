package main

import (
	"fmt"
	"github.com/jackson198608/goProject/common/tools"
	log "github.com/thinkboy/log4go"
	"strconv"
	"time"
)

var c Config = Config{5, 100, "192.168.86.80:6380,192.168.86.80:6381,192.168.86.81:6380,192.168.86.81:6380,192.168.86.85:6380,192.168.86.85:6381", "http://192.168.86.230:9200,http://192.168.86.231:9200","/etc/gouminMultiMessagePushLog.xml"}


func main()  {
	redisInfo := tools.FormatRedisOption(c.redisConn)
	log.Info("start work")

	redisConn, err := tools.GetClusterClient(&redisInfo)
	_, err = redisConn.Ping().Result()
	if err != nil {
		for t:=1; t<5;t++  {
			redisConn, err = tools.GetClusterClient(&redisInfo)
			_, err = redisConn.Ping().Result()
			if err == nil {
				break
			}else{
				fmt.Println("[error] redis connection crashed ,going to create a new one, try times ",t)
			}
		}

		fmt.Println("[error] redis connection crashed ,going to create a new one")
	}

	defer redisConn.Close()
	i:=0
	for i=0; i<10000000;i++  {
		//redisConn.RPush("test_list_value", i)
		//rst := redisConn.RPop("test_list_value")
		rst:=redisConn.Get("redpoint_totle_2264057")
		fmt.Println("i--"+strconv.Itoa(i),rst)
		time.Sleep(10)
	}

	//_, err := redisEngine.NewRedisEngine(c.queueName, &redisInfo, mongoConnInfo, mysqlInfo, esNodes, c.coroutinNum, 1, jobFuc)
	if err != nil {
		log.Error("[NewRedisEngine] ", err)
		fmt.Println("err:",err)
	}
}
