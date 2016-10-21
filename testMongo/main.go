package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Message struct {
	Uid           int    `json:"uid"`
	Type          int    `json: "type"`
	Mark          int    `json: "mark"`
	Isnew         int    `json: "isnew"`
	From          int    `json: "from"`
	Channel       int    `json: "channel"`
	Channel_types int    `json: "channel_types"`
	Title         string `json: "title"`
	Content       string `json: "content"`
	Image         string `json: "image"`
	Url_type      int    `json: "url_type"`
	Url           string `json: "url"`
	Created       string `json: "created"`
	Modified      string `json: "modified"`
}

func main() {
	session, err := mgo.Dial("210.14.154.198:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")

	jsonStr := `{"uid":1895167,"type":1,"makr":281,"isnew":0,"from":0,"channel":1,"channel_types":2,"title":"狗狗的寂寞都市之殇","content":"小短腿在家漂移，屁股差点没甩掉了~","image":"/messagepush/day_161020/20161020_7a50e50.jpg","url_type":1,"url":"4346101","created":"2016-10-20 14:12:28","modified":"0000-00-00 00:00:00"}`

	var m Message
	if err = json.Unmarshal([]byte(jsonStr), &m); err != nil {
		fmt.Println("error json fotmat")
	}

	err = c.Insert(&m)
	//err = c.Insert(&Message{1895167, 1, 281, 0, 0, 1, 2, "狗狗的寂寞都市之殇", "小短腿在家漂移，屁股差点没甩掉了~", "/messagepush/day_161020/20161020_7a50e50.jpg", 1, "4346101", "2016-10-20 14:12:28", "0000-00-00 00:00:00"})
	if err != nil {
		log.Fatal(err)
	}

	var result Message
	err = c.Find(bson.M{"uid": 1895167}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("title:", result.Title)
}
