package main

import (
	"fmt"
	mgo "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"os"
	"strconv"
)

var fileName = "/tmp/event.log"

func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

//创建索引
func main() {
	jobType := os.Args[1]
	// err := appendToFile(fileName, "dfdf")
	// fmt.Println(err)
	// session, err := mgo.Dial("192.168.86.68:10001")
	session, err := mgo.Dial("192.168.5.22:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	if jobType == "event" {
		tableName := "event_log" //动态表
		fmt.Println(tableName)
		c := session.DB("EventLog").C(tableName)
		c.EnsureIndexKey("type")
		c.EnsureIndexKey("uid")
		// c.EnsureIndexKey("info")
		c.EnsureIndexKey("created")
		c.EnsureIndexKey("infoid")
		c.EnsureIndexKey("status")
		c.EnsureIndexKey("tid")
		c.EnsureIndexKey("bid")
		c.EnsureIndexKey("tag")
		c.EnsureIndexKey("source")
		c.EnsureIndexKey("is_read")
		c.EnsureIndexKey("qst_type")

		// x := session.DB("EventLog").C("ids")
		// x.Insert(bson.M{"_id": 0, "id": 0})
	}

	if jobType == "fans" {
		for i := 1; i <= 100; i++ {
			tableName1 := "event_log_" + strconv.Itoa(i) //粉丝表
			fmt.Println(tableName1)
			c := session.DB("FansData").C(tableName1)
			c.EnsureIndexKey("type")
			c.EnsureIndexKey("uid")
			c.EnsureIndexKey("fuid")
			c.EnsureIndexKey("created")
			c.EnsureIndexKey("infoid")
			c.EnsureIndexKey("status")
			c.EnsureIndexKey("tid")
			c.EnsureIndexKey("bid")
			c.EnsureIndexKey("tag")
			c.EnsureIndexKey("source")
			c.EnsureIndexKey("is_read")
			c.EnsureIndexKey("qst_type")
			// x := session.DB("FansData").C("ids" + strconv.Itoa(i))
			// x.Insert(bson.M{"_id": 0, "id": 0})
		}
	}

}
