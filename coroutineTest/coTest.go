package main

import (
	"fmt"
)

/*
func printVimGo(ch chan int) {
	ch <- 1
	fmt.Println("vim-go")
}

func main() {
	var chs [5]chan int
	for i := 0; i < 5; i++ {
		chs[i] = make(chan int)
		go printVimGo(chs[i])
	}

	for _, ch := range chs {
		<-ch
	}
}
*/
func Count(ch chan int, i int) {
	fmt.Println("Counting", i)
	ch <- 10 + i
}

func main() {

	chs := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		go Count(chs[i], i)
	}

	for _, ch := range chs {
		<-ch
	}
}
