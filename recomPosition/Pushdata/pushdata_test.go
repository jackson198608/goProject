package Pushdata

import (
	"database/sql"
	"github.com/donnie4w/go-logger/logger"
	mgo "gopkg.in/mgo.v2"
	"testing"
)

func TestSaveMongoEventLog(t *testing.T) {
	dbName := "test_dz2"
	db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/"+dbName+"?charset=utf8")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	mongoConn := "192.168.86.68:27017"
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		logger.Error("[error] connect mongodb err")
		return
	}
	oid := 0
	id := 0
	status := "0"
	fuid := 1138687
	uid := 1

	m := NewEventLogNew(1, oid, id, db, session)
	if m != nil {
		if oid > 0 {
			logger.Info("export event to mongo")
			m.SaveMongoEventLog(oid)
		}
		if id > 0 {
			logger.Info("update mongo event status")
			m.UpdateMongoEventLogStatus(id, status)
		}
		if fuid > 0 && uid > 0 {
			logger.Info("remove fans event")
			m.RemoveFansEventLog(fuid, uid)
		}
	}
}
