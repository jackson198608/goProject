package singlePush


import (
	"github.com/jackson198608/goProject/appPush"
	"gopkg.in/redis.v4"
	"github.com/pkg/errors"
)

var queuename = "mcSingle"

type Singlepush struct {
	jobstr string
	redisConn *redis.ClusterClient
	p12Bytes []byte
}

func NewSinglepush(jobStr string,redisConn *redis.ClusterClient,p12Bytes []byte) *Singlepush{
	if (jobStr == "") ||( redisConn == nil) || (p12Bytes == nil) {
		return nil
	}

	s := new(Singlepush)
	if s == nil {
		return nil
	}
	s.jobstr = jobStr
	s.redisConn = redisConn
	s.p12Bytes = p12Bytes

	return s
}

func (s *Singlepush) Do() error {
	//发送push
	t := appPush.NewTask(s.jobstr)
	if t != nil {
		w := appPush.NewWorker(t)
		result := w.Push(s.p12Bytes)
		if !result {
			return errors.New("multi push fail")
		}
	}
	return nil
}

