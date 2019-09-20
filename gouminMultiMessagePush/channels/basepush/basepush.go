package basepush

import (
	"github.com/jackson198608/goProject/appPush"
	"gopkg.in/redis.v4"
	"github.com/pkg/errors"
)

type Basepush struct {
	jobstr string
	redisConn *redis.ClusterClient
	p12Bytes []byte
}

func Newpush(jobStr string,redisConn *redis.ClusterClient,p12Bytes []byte) *Basepush{
	if (jobStr == "") ||( redisConn == nil) || (p12Bytes == nil){
		return nil
	}

	b := new(Basepush)
	if b == nil {
		return nil
	}
	b.jobstr = jobStr
	b.redisConn = redisConn
	b.p12Bytes = p12Bytes

	return b
}

func (b *Basepush) Do() error {
	//发送push
	t := appPush.NewTask(b.jobstr)
	if t != nil {
		w := appPush.NewWorker(t, b.redisConn)
		result := w.Push(b.p12Bytes)
		if !result {
			return errors.New("multi push fail")
		}
	}
	return nil
}
