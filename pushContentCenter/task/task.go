package main

import (
	"fmt"
	mgo "gopkg.in/mgo.v2"
)

type Task struct {
	job string
}

//job: redisQueue pop string
//taskarg: mongoHost,mongoDatabase,mongoReplicaSetName
func NewTask(job string, taskarg ...string) *task {
	//@todo check prams

	t = new(Task)
	if t == nil {
		return nil
	}

	//@todo pass param
	//@todo

	return t

}

func (t *Task) Do() error {
	return nil
}
