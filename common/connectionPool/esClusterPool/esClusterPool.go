package esClusterPool

import (
	"sync"
	"github.com/olivere/elastic"
	queue "github.com/oleiade/lane"
	log "github.com/thinkboy/log4go"
	"fmt"
	"time"
	"errors"
	"gouminGitlab/common/orm/elasticsearchBase"
)

var loopnum = 3

type EsPool struct {
	esInfo    []string
	num       int
	poolQueue *queue.Queue
	lock      sync.RWMutex
}

func NewEsPool(esInfo []string, num int) (*EsPool, error) {
	//@todo check params
	esPool := new(EsPool)
	if esPool == nil {
		return esPool, errors.New("create memory error.can not create it")
	}
	esPool.esInfo = esInfo
	esPool.num = num
	err := esPool.initPool()
	if err != nil {
		return esPool, err
	}

	return esPool, nil

}

func (r *EsPool) initPool() error {
	r.poolQueue = queue.NewQueue()
	for i := 0; i < r.num; i++ {
		connection, err := r.createOneConnection()
		// fmt.Println(i, connection, reflect.TypeOf(connection))
		if err != nil {
			//@todo something
			fmt.Println("err create connection", err)
			return err
		} else {
			r.poolQueue.Enqueue(connection)
		}
	}
	return nil
}

func (r *EsPool) createOneConnection() (*elastic.Client, error) {
	if r.esInfo == nil {
		return nil, nil
	}
	esR, _ := elasticsearchBase.NewClient(r.esInfo)
	client, err := esR.Run()
	if err != nil {
		return nil, nil
	}
	return client, nil
}

func (r *EsPool) Close() error {
	for i := 0; i < r.poolQueue.Size(); i++ {
		connection := r.poolQueue.Dequeue()
		connection.(*elastic.Client).Stop()
		//@todo close
	}
	return nil
}

func (r *EsPool) GetConnection() (*elastic.Client, error) {
	r.lock.Lock()
	fmt.Println("poolQueue size = ", r.poolQueue.Size())
	if r.poolQueue.Size() == 0 {
		r.lock.Unlock()
		return nil, errors.New("there is no more connection can be use ,please wait")
	}
	connection := r.poolQueue.Dequeue()
	r.lock.Unlock()
	return connection.(*elastic.Client), nil

}

func (r *EsPool) TellMeOneIsBroken() {
	for i := 0; i < loopnum; i++ {
		fmt.Println("tell me times ", i)
		connection, err := r.createOneConnection()
		if err != nil {
			log.Error("create connection to enqueue error", err)
			time.Sleep(10 * time.Millisecond)
			continue
		}
		err = r.PutConnection(connection)
		if err != nil {
			log.Error("put connection to enqueue error", err)
			continue
		} else {
			break
		}
	}

}

func (r *EsPool) PutConnection(connection *elastic.Client) error {
	r.lock.Lock()
	r.poolQueue.Enqueue(connection)
	r.lock.Unlock()
	return nil
}

