package inMongo

import (
	"encoding/json"
	"fmt"
	mgo "gopkg.in/mgo.v2"
	"log"
	"strconv"
)

type Worker struct {
	t *Task
}

func NewWorker(t *Task) (w *Worker) {
	//init the worker
	var wR Worker
	wR.t = t
	return &wR
}

func (w Worker) Insert(session *mgo.Session) {
	//convert json string to struct
	var m row
	log.Println("[notice]", w.t.columData)
	if err := json.Unmarshal([]byte(w.t.columData), &m); err != nil {
		fmt.Println("[error]insert mongo error", err, w.t.columData)
	}

	//get the table name
	tableNumber := strconv.Itoa(m.Uid % 1000)
	tableName := "message_push_record_" + tableNumber

	//create mongo session
	c := session.DB("MessageCenter").C(tableName)

	err := c.Insert(&m)
	if err != nil {
		log.Fatal(err)
	}
}
