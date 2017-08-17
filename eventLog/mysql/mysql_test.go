package mysql

import (
	"database/sql"
	"github.com/donnie4w/go-logger/logger"
	"testing"
)

func TestCheckIsFans(t *testing.T) {
	dbName := "test_dz2"
	db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/"+dbName+"?charset=utf8")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	// uid := 1
	uid := 1138687
	follow_id := 1138687
	fans := CheckIsFans(uid, follow_id, db)
	logger.Info(fans)
}

func TestGetFansData(t *testing.T) {
	dbName := "test_dz2"
	db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/"+dbName+"?charset=utf8")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	// uid := 1
	uid := 1138687
	fans := GetFansData(uid, db)
	logger.Info(fans)
}

func TestLoadById(t *testing.T) {
	dbName := "test_dz2"
	db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/"+dbName+"?charset=utf8")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	// uid := 1
	id := 1
	event := LoadById(id, db)
	logger.Info(event.Uid, event.Infoid, event.TypeId, event.Created, event.Status)
}
