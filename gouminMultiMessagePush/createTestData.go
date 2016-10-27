package main

import (
	"fmt"
	redis "gopkg.in/redis.v4"
)

func conn(conn string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}

func testRpush(client *redis.Client) {
	var iosRedisString string = `0|fb71306452499efc778cc77d1be6614b8e1753e79b7acd8348ce6cb47abd4dc2|{"aps":{"alert":"task","sound":"default","badge":1,"type":6,"mark":""}}`

	var mongoString string = `{"uid":1895167,"type":1,"makr":281,"isnew":0,"from":0,"channel":1,"channel_types":2,"title":"狗狗的寂寞都市之殇","content":"小短腿在家漂移，屁股差点没甩掉了~","image":"/messagepush/day_161020/20161020_7a50e50.jpg","url_type":1,"url":"4346101","created":"2016-10-20 14:12:28","modified":"0000-00-00 00:00:00"}`
	var androidRedisString string = `1|fb71306452499efc778cc77d1be6614b8e1753e79b7acd8348ce6cb47abd4dc2|{"action":"pushMessageToSingleAction","clientData":"CgASC3B1c2htZXNzYWdlGgAiFmM2Mk5XOUp4d2w5eFJTRlh6SlAwbjQqFm1RZ3FuY0ZDVXk5aG14aHFLQ3ZHTDkyADoKCgASABoAIgItMUIKCAEQABiuTpAKAEIcCK5OEAMYZOIIAOoIBgoAEgAaAPAIAPgIZJAKAEIHCGQQB5AKAEoJZHVyYXRpb249","transmissionContent":"eyJ0eXBlIjo2LCJtYXJrIjoiIiwic2lsZW50IjowLCJ0aXRsZSI6InRlc3QgcHVzaCB0aXRsZSIsImNvbnRlbnQiOiJ0ZXN0IHB1c2ggY29udGVudCEhISIsInVpZCI6MjAzNTgzMH0=","isOffline":true,"offlineExpireTime":43200000,"pushNetWorkType":1,"appId":"mQgqncFCUy9hmxhqKCvGL9","clientId":"aef16d24053b161bb0b35d3a0f775fcd","alias":null,"type":2,"pushType":"TransmissionMsg","version":"3.0.0.0","appkey":"c62NW9Jxwl9xRSFXzJP0n4"}`

	//insert mulit-push data
	for i := 0; i < numForOneLoop; i++ {
		if i%2 == 0 {
			err := (*client).RPush("mcMulti", iosRedisString).Err()
			if err != nil {
				//insert error
			}
		} else {
			err := (*client).RPush("mcMulti", androidRedisString).Err()
			if err != nil {
				//insert error
			}
		}
	}

	//insert single-push data
	for i := 0; i < numForOneLoop; i++ {
		if i%2 == 0 {
			err := (*client).RPush("mcSingle", iosRedisString).Err()
			if err != nil {
				//insert error
			}
		} else {
			err := (*client).RPush("mcSingle", androidRedisString).Err()
			if err != nil {
				//insert error
			}
		}
	}

	//insert mongo
	for i := 0; i < numForOneLoop; i++ {
		err := (*client).RPush("mcInsert", mongoString).Err()
		if err != nil {
			//insert error
		}
	}

}

func testCreateTestData() {
	client := conn(c.redisConn)
	testRpush(client)
}

func testPrintTasks() {
	for i := 0; i < taskNum; i++ {
		fmt.Println(tasks[i])
	}
}
