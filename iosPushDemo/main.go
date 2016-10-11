package main

import (
	"fmt"
)

func work(c chan int) {
	var redisString string = `fb71306452499efc778cc77d1be6614b8e1753e79b7acd8348ce6cb47abd4dc2|{"aps":{"alert":"task","sound":"default","badge":1,"type":6,"mark":""}}`
	t := NewTask(redisString)
	w := NewWorker(t)
	w.Push()
	c <- 1
}

func main() {
	c := make(chan int)
	go work(c)
	result := <-c
	fmt.Println(result)
}
