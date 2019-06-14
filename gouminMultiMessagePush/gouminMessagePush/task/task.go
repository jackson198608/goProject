package task

import (
	"github.com/olivere/elastic"
	"gopkg.in/redis.v4"
	"github.com/jackson198608/goProject/gouminMultiMessagePush/channels/mcInsert"
	"github.com/bitly/go-simplejson"
	"strconv"
	"github.com/pkg/errors"
	"github.com/jackson198608/gotest/gouminMultiMessagePush/channels/basepush"
	log "github.com/thinkboy/log4go"
)

var activityRedpointKey = "redpoint_activity_"
var recommendRedpointKey = "redpoint_recommend_"
var serviceRedpointKey = "redpoint_service_"
var totalRedpointKey = "redpoint_totle_"

type Task struct {
	Jobstr    string
	JobType   string
	EsConn  *elastic.Client
	RedisConn *redis.ClusterClient
	P12Bytes []byte
}

func NewTask(jobType string,jobStr string,redisConn *redis.ClusterClient,esConn *elastic.Client,p12Bytes []byte) (*Task,error){
	t := new(Task)
	if t == nil {
		return nil, errors.New("there is no space to create struct")
	}
	t.Jobstr = jobStr
	t.JobType = jobType
	t.EsConn = esConn
	t.RedisConn = redisConn
	t.P12Bytes = p12Bytes

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
	err := mc.Insert(t.EsConn,t.RedisConn)
	if err == nil {
		//更新小红点
		t.changeRedisKey(t.Jobstr)
		log.Info("insert success jobstr:",t.Jobstr)
		return nil
	}
	log.Info("insert fail err:",err," jobstr:",t.Jobstr)
	return err
}

func (t *Task) changeRedisKey(jobStr string) {

	json, err := simplejson.NewJson([]byte(jobStr))
	if err != nil {
		errors.New("mcinsert update redpoint fail jobstr:"+jobStr)
		return
	}
	Jtype, err := json.Get("type").Int()
	if err != nil {
		errors.New("mcinsert update redpoint fail jobstr:"+jobStr)
		return
	}

	uid, err := json.Get("uid").Int()
	if err != nil {
		errors.New("mcinsert update redpoint fail jobstr:"+jobStr)
		return
	}

	client := t.RedisConn

	var key string
	if Jtype == 1 {
		key = activityRedpointKey + strconv.Itoa(uid)

	} else if Jtype == 6 {
		key = recommendRedpointKey + strconv.Itoa(uid)

	} else {
		key = serviceRedpointKey + strconv.Itoa(uid)
	}
	(*client).Incr(key)

	totalRedPointKey := totalRedpointKey + strconv.Itoa(uid)
	(*client).Set(totalRedPointKey, 1, 0)
	log.Info("success redpoint jobstr:",jobStr)
}
