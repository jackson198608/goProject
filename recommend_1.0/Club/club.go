package club

import (
	"fmt"
	"github.com/jackson198608/goProject/common/tools/elkClient"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Club struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	Uid       int
}

func NewClub(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, Uid int) *Club {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jobStr == "") {
		return nil
	}

	c := new(Club)
	if c == nil {
		return nil
	}

	c.mysqlXorm = mysqlXorm
	c.mongoConn = mongoConn
	c.Uid = Uid
	return c
}

func (c *Club) Do() error {

}

func (c *Club) getUser() error {
	client, err = elkClient.NewClient("192.168.86.5:9200")
	if !err {
		fmt.Println(err)
	}
}
