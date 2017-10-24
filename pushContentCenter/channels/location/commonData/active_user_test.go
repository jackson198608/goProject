package commonData

import (
	"fmt"
	// "github.com/jackson198608/goProject/pushContentCenter/channels/location/commonData"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

func testConn() []*mgo.Session {

	mongoConn := "192.168.86.192:27017"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		fmt.Println("[error] connect mongodb err")
		return nil
	}
	var sessionAry []*mgo.Session
	sessionAry = append(sessionAry, session)
	return sessionAry
	// return engine, session
}

func TestLoadDataToHashmap(t *testing.T) {
	mongoConn := testConn()

	fmt.Println(LoadDataToHashmap(mongoConn[0]))
}
