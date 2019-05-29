package appPush

import (
	"fmt"
	"io/ioutil"
	"testing"
)

const workNum = 10

var p12bytes []byte

func iosPush(c chan int) {
	//redisString format: devicetoken | json
	//var redisString string = `0|ed4c6d809ec99ac445e0f5d23262a8aab86745c2b23a0c8b0749553791efc970|{"aps":{"alert":"task","sound":"default","badge":1,"type":6,"mark":""}}`
	var redisString string = `0|85e0fba08d3d7173051adb20c0f1a3cd4f8ad65a98eefceca4892145b37be4c6|{"aps":{"alert":"task","sound":"default","badge":1,"type":6,"mark":""}}`
	//var redisString string = `0|11153d0425dd2457a0fbcc4f5f61c3f0a3d7f95499cc5276693cf6eb9311747c|{"aps":{"alert":"task","sound":"default","badge":1,"type":20,"mark":"G3535-216#216"}}`
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
		fmt.Println("hi!")
		return
	}

	p12bytes = cBytes

	c := make(chan int, workNum)

	for i := 0; i < workNum; i++ {
		go iosPush(c)
	}
	for i := 0; i < workNum; i++ {
		<-c
	}
}
