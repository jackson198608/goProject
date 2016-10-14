package iosPush

import (
	"testing"
)

const workNum = 10

func work(c chan int) {
	//redisString format: devicetoken | json
	var redisString string = `fb71306452499efc778cc77d1be6614b8e1753e79b7acd8348ce6cb47abd4dc2|{"aps":{"alert":"task","sound":"default","badge":1,"type":6,"mark":""}}`
	t := NewTask(redisString)
	w := NewWorker(t)
	w.Push()
	c <- 1
}

func TestFibonacci(t *testing.T) {
	c := make(chan int, workNum)

	for i := 0; i < workNum; i++ {
		go work(c)
	}
	for i := 0; i < workNum; i++ {
		<-c
	}
}
