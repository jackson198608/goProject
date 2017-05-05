package main

// import (
// 	"github.com/jackson198608/goProject/appPush"
// )

// func doPush(c chan int, i int) {
// 	t := appPush.NewTask(tasks[i].pushStr)
// 	if t != nil {
// 		w := appPush.NewWorker(t)
// 		result := w.Push(p12Bytes)
// 		if !result {
// 			//fmt.Println("[Warn] push fail,put back to redis:  ", tasks[i].pushStr)
// 			putFailOneBack(i)
// 		}
// 	}
// 	c <- 1
// }

// func push() {
// 	c := make(chan int, taskNum)

// 	for i := 0; i < taskNum; i++ {
// 		go doPush(c, i)
// 	}
// 	for i := 0; i < taskNum; i++ {
// 		<-c
// 	}
// }
