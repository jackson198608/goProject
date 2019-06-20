package task

import (
	"github.com/olivere/elastic"
	"gopkg.in/redis.v4"
	"github.com/jackson198608/goProject/gouminMultiMessagePush/channels/mcInsert"
	//"github.com/bitly/go-simplejson"
	//"strconv"
	"github.com/pkg/errors"
	"github.com/jackson198608/gotest/gouminMultiMessagePush/channels/basepush"
	log "github.com/thinkboy/log4go"
	"github.com/jackson198608/goProject/gouminMultiMessagePush/RedPoint"
)

type Task struct {
	Jobstr    string
	JobType   string
	EsConn  *elastic.Client
	RedisConn *redis.ClusterClient
	P12Bytes []byte
	esInfo string
}

func NewTask(jobType string,jobStr string,redisConn *redis.ClusterClient,esConn *elastic.Client,p12Bytes []byte, esInfo string) (*Task,error){
	t := new(Task)
	if t == nil {
		return nil, errors.New("there is no space to create struct")
	}
	t.Jobstr = jobStr
	t.JobType = jobType
	t.EsConn = esConn
	t.RedisConn = redisConn
	t.P12Bytes = p12Bytes
	t.esInfo = esInfo

	return t,nil
}

func (t *Task) Do() error{
	log.Info("task jobType:",t.JobType," jobstr:",t.Jobstr)
	switch t.JobType {
	case "multi":
		err := t.channelMulti()
		if err != nil {
			return err
		} else {
			return nil
		}
	case "single":
		err := t.channelSingle()
		if err != nil {
			return err
		} else {
			return nil
		}
	case "insert":
		err := t.channelInsert()
		if err != nil {
			return err
		} else {
			return nil
		}
	}

	return nil
}

func (t *Task) channelMulti() error {
	log.Info("multi jobstr:",t.Jobstr)
	m := basepush.Newpush(t.Jobstr,t.RedisConn,t.P12Bytes)
	err := m.Do()
	if err != nil {
		log.Info("multi fail err:",err," jobstr:",t.Jobstr)
		return err
	}
	log.Info("multi success jobstr:",t.Jobstr)
	return nil
}

func (t *Task) channelSingle() error {
	log.Info("single jobstr:",t.Jobstr)
	s := basepush.Newpush(t.Jobstr,t.RedisConn,t.P12Bytes)
	err := s.Do()
	if err != nil {
		log.Info("single fail err:",err," jobstr:",t.Jobstr)
		return err
	}
	log.Info("single success jobstr:",t.Jobstr)
	return nil
}

func (t *Task) channelInsert() error {
	log.Info("insert jobstr:",t.Jobstr)
	mc := mcInsert.NewTask(t.Jobstr)
	err := mc.Insert(t.EsConn,t.RedisConn, t.esInfo)
	if err == nil {
		//更新小红点
		rp := RedPoint.NewTask(t.Jobstr, t.RedisConn)
		rp.ChangeRedisKeys()
		log.Info("insert success jobstr:",t.Jobstr)
		return nil
	}
	log.Info("insert fail err:",err," jobstr:",t.Jobstr)
	return err
}
