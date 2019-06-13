package task

import (
	"github.com/olivere/elastic"
	"github.com/jackson198608/goProject/gouminMultiMessagePush/channels/multiPush"
	"github.com/jackson198608/goProject/gouminMultiMessagePush/channels/singlePush"
	"gopkg.in/redis.v4"
	"github.com/jackson198608/goProject/gouminMultiMessagePush/channels/mcInsert"
	"github.com/bitly/go-simplejson"
	"fmt"
	"strconv"
	"github.com/pkg/errors"
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
	m := multiPush.NewMultipush(t.Jobstr,t.RedisConn,t.P12Bytes)
	err := m.Do()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) channelSingle() error {
	s := singlePush.NewSinglepush(t.Jobstr,t.RedisConn,t.P12Bytes)
	err := s.Do()
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) channelInsert() error {
	mc := mcInsert.NewTask(t.Jobstr)
	if mc.Insert(t.EsConn,t.RedisConn) {
		t.changeRedisKey(t.Jobstr)
	}
	return nil
}

func (t *Task) changeRedisKey(jobStr string) {

	json, err := simplejson.NewJson([]byte(jobStr))
	if err != nil {
		fmt.Println("[error]json format error", jobStr, err)
		return
	}
	Jtype, err := json.Get("type").Int()
	if err != nil {
		fmt.Println("[error]get type from json error", jobStr, err)
		return
	}

	uid, err := json.Get("uid").Int()
	if err != nil {
		fmt.Println("[error]get uid from json error", jobStr, err)
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

	fmt.Println("[info]set key", key)
	(*client).Incr(key)

	totalRedPointKey := totalRedpointKey + strconv.Itoa(uid)
	(*client).Set(totalRedPointKey, 1, 0)
}
