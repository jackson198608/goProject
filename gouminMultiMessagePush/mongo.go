package main

import (
	"github.com/jackson198608/gotest/inMongo"
	mgo "gopkg.in/mgo.v2"
)

func do(c chan int, session *mgo.Session, i int) {
	t := inMongo.NewTask(tasks[i].insertStr)
	w := inMongo.NewWorker(t)
	w.Insert(session)

	c <- 1
}

func insertMongo() {
	//init session
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//create channel
	c := make(chan int, taskNum)

	for i := 0; i < taskNum; i++ {
		go do(c, session, i)
	}
	for i := 0; i < taskNum; i++ {
		<-c
	}
}
