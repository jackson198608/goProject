package club

import (
	"errors"
	"fmt"
	// "github.com/jackson198608/goProject/common/tools/elkClient"
	"github.com/olivere/elastic"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	Uid       int
	elkDsn    string
}

func NewUser(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, uid int, elkDsn string) *User {
	if (mysqlXorm == nil) || (mongoConn == nil) || (uid == 0) || (elkDsn == "") {
		return nil
	}

	u := new(User)
	if u == nil {
		return nil
	}

	u.mysqlXorm = mysqlXorm
	u.mongoConn = mongoConn
	u.Uid = uid
	u.elkDsn = elkDsn
	return u
}

func (u *User) Do() error {
	client, err := elastic.NewClient(elastic.SetURL("http://192.168.86.5:9200"))
	if err != nil {
		panic(err)
	}
	fmt.Println(client)
	return nil
}

func (u *User) getUser(client *elastic.Client) error {
	termQuery := elastic.NewTermQuery("app", appName)

	res, err := client.Search(indexName).
		Index(indexName).
		Query(termQuery).
		Sort("time", true).
		Do()

	if err != nil {
		return err
	}

	// fmt.Println("Logs found:")
	// var l Log
	// for _, item := range res.Each(reflect.TypeOf(l)) {
	// 	l := item.(Log)
	// 	fmt.Printf("time: %s message: %s\n", l.Time, l.Message)
	// }

	// return nil
}
