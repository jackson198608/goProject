package main

import (
	"fmt"
	mgo "gopkg.in/mgo.v2"
	"strconv"
)

func main() {
	session, err := mgo.Dial("192.168.5.22:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	for i := 10; i < 1000; i++ {
		tableName := "message_push_record_" + strconv.Itoa(i)
		fmt.Println(tableName)
		c := session.DB("MessageCenter").C(tableName)
		c.EnsureIndexKey("uid")
		c.EnsureIndexKey("created")
		c.EnsureIndexKey("mark")
		c.EnsureIndexKey("isnew")
		c.EnsureIndexKey("from")
	}
}
