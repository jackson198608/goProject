package task

import (
	"errors"
	"github.com/jackson198608/goProject/image/compress"
	// "strings"
)

type Task struct {
	Raw string //the data get from redis queue
}

//job: redisQueue pop string
//taskarg: mongoHost,mongoDatabase,mongoReplicaSetName
func NewTask(raw string) (*Task, error) {
	//check prams
	if raw == "" {
		return nil, errors.New("params can not be null")
	}

	t := new(Task)
	if t == nil {
		return nil, errors.New("there is no space to create struct")
	}

	//pass params
	t.Raw = raw

	return t, nil

}

// public interface for task
// if you have New channles you need to add logic here
func (t *Task) Do() error {
	c := compress.NewCompress(t.Raw)
	err := c.Do()
	if err != nil {
		return err
	}
	return nil
}
