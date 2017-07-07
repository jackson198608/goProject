package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/jackson198608/gotest/inMongo"
	mgo "gopkg.in/mgo.v2"
	"strconv"
)

func doInMongo(c chan int, session *mgo.Session, i int) {
	mongoStr := tasks[i].insertStr
	if mongoStr == "0" {
		c <- 1
		return
	}
	t := inMongo.NewTask(mongoStr)
	w := inMongo.NewWorker(t)
	if w.Insert(session) {
		changeRedisKey(mongoStr)
	}

	c <- 1
}

func changeRedisKey(mongoStr string) {

	json, err := simplejson.NewJson([]byte(mongoStr))
	if err != nil {
		fmt.Println("[error]json format error", mongoStr, err)
		return
	}
	Jtype, err := json.Get("type").Int()
	if err != nil {
		fmt.Println("[error]get type from json error", mongoStr, err)
		return
	}

	uid, err := json.Get("uid").Int()
	if err != nil {
		fmt.Println("[error]get uid from json error", mongoStr, err)
		return
	}

	client := connect(c.redisConn)
	defer client.Close()

	var key string
	if Jtype == 1 {
		key = "redpoint_activity_" + strconv.Itoa(uid)

	} else if Jtype == 6 {
		key = "redpoint_recommend_" + strconv.Itoa(uid)

	} else {
		key = "redpoint_service_" + strconv.Itoa(uid)
	}

	fmt.Println("[info]set key", key)
	(*client).Set(key, 1, 0)
}

func insertMongo() {
	//init session
	session, err := mgo.Dial(c.mongoConn)
	if err != nil {
		return
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
