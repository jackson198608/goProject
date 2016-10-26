package appPush

import (
	"io/ioutil"
	"testing"
)

const workNum = 10

var p12bytes []byte

func iosPush(c chan int) {
	//redisString format: devicetoken | json
	var redisString string = `0|fb71306452499efc778cc77d1be6614b8e1753e79b7acd8348ce6cb47abd4dc2|{"aps":{"alert":"task","sound":"default","badge":1,"type":6,"mark":""}}`
	t := NewTask(redisString)
	w := NewWorker(t)
	w.Push(p12bytes)

	c <- 1
}

func androidPush(c chan int) {
	var redisString string = `1|fb71306452499efc778cc77d1be6614b8e1753e79b7acd8348ce6cb47abd4dc2|{"action":"pushMessageToSingleAction","clientData":"CgASC3B1c2htZXNzYWdlGgAiFmM2Mk5XOUp4d2w5eFJTRlh6SlAwbjQqFm1RZ3FuY0ZDVXk5aG14aHFLQ3ZHTDkyADoKCgASABoAIgItMUIKCAEQABiuTpAKAEIcCK5OEAMYZOIIAOoIBgoAEgAaAPAIAPgIZJAKAEIHCGQQB5AKAEoJZHVyYXRpb249","transmissionContent":"eyJ0eXBlIjo2LCJtYXJrIjoiIiwic2lsZW50IjowLCJ0aXRsZSI6InRlc3QgcHVzaCB0aXRsZSIsImNvbnRlbnQiOiJ0ZXN0IHB1c2ggY29udGVudCEhISIsInVpZCI6MjAzNTgzMH0=","isOffline":true,"offlineExpireTime":43200000,"pushNetWorkType":1,"appId":"mQgqncFCUy9hmxhqKCvGL9","clientId":"aef16d24053b161bb0b35d3a0f775fcd","alias":null,"type":2,"pushType":"TransmissionMsg","version":"3.0.0.0","appkey":"c62NW9Jxwl9xRSFXzJP0n4"}`
	t := NewTask(redisString)
	w := NewWorker(t)
	w.Push(p12bytes)

	c <- 1
}

func TestFibonacci(t *testing.T) {
	cBytes, err := ioutil.ReadFile("/etc/pro-lingdang.pem")
	if err != nil {
		return
	}

	p12bytes = cBytes

	c := make(chan int, workNum)

	for i := 0; i < workNum; i++ {
		go androidPush(c)
	}
	for i := 0; i < workNum; i++ {
		<-c
	}
}
