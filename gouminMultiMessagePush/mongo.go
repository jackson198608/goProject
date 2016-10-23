package main

import (
	"github.com/jackson198608/gotest/inMongo"
	mgo "gopkg.in/mgo.v2"
	"log"
)

func doInMongo(c chan int, session *mgo.Session, i int) {
	mongoStr := tasks[i].insertStr
	if mongoStr == "0" {
		log.Println("[Warn]no need to insert mongo")
		c <- 1
		return
	}
	t := inMongo.NewTask(mongoStr)
	w := inMongo.NewWorker(t)
	w.Insert(session)

	c <- 1
}

func insertMongo() {
	//init session
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		log.Println("[Error] mongo connect error")
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//create channel
	c := make(chan int, taskNum)

	for i := 0; i < taskNum; i++ {
		go doInMongo(c, session, i)
	}
	for i := 0; i < taskNum; i++ {
		<-c
	}
}
