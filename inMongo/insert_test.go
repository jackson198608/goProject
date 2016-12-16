package inMongo

import (
	mgo "gopkg.in/mgo.v2"
	"testing"
)

const workNum = 10

func do(c chan int, session *mgo.Session) {
	var redisString string = `{"uid":1895167,"type":1,"makr":281,"isnew":0,"from":0,"channel":1,"channel_types":2,"title":"狗狗的寂寞都市之殇","content":"小短腿在家漂移，屁股差点没甩掉了~","image":"/messagepush/day_161020/20161020_7a50e50.jpg","url_type":1,"url":"4346101","created":"2016-10-20 14:12:28","modified":"0000-00-00 00:00:00"}`
	t := NewTask(redisString)
	w := NewWorker(t)
	w.Insert(session)

	c <- 1
}

func TestFibonacci(t *testing.T) {
	//init session
	session, err := mgo.Dial("210.14.154.198:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	//create channel
	c := make(chan int, workNum)

	for i := 0; i < workNum; i++ {
		go do(c, session)
	}
	for i := 0; i < workNum; i++ {
		<-c
	}
}
