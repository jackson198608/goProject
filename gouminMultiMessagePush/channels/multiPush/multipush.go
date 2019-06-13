package multiPush

import (
	"github.com/jackson198608/goProject/appPush"
	"gopkg.in/redis.v4"
	"github.com/pkg/errors"
)

var queuename = "mcMulti"

type Multipush struct {
	jobstr string
	redisConn *redis.ClusterClient
	p12Bytes []byte
}

func NewMultipush(jobStr string,redisConn *redis.ClusterClient,p12Bytes []byte) *Multipush{
	if (jobStr == "") ||( redisConn == nil) || (p12Bytes == nil){
		return nil
	}

	m := new(Multipush)
	if m == nil {
		return nil
	}
	m.jobstr = jobStr
	m.redisConn = redisConn
	m.p12Bytes = p12Bytes

	return m
}

func (m *Multipush) Do() error {
	//发送push
	t := appPush.NewTask(m.jobstr)
	if t != nil {
		w := appPush.NewWorker(t)
		result := w.Push(m.p12Bytes)
		if !result {
			return errors.New("multi push fail")
		}
	}
	return nil
}

