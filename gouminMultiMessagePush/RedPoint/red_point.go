package RedPoint

import (
	"gopkg.in/redis.v4"
	"strings"
	"github.com/bitly/go-simplejson"
	"strconv"
	"errors"
	log "github.com/thinkboy/log4go"
	"time"
)

var activityRedpointKey = "redpoint_activity_"
var recommendRedpointKey = "redpoint_recommend_"
var serviceRedpointKey = "redpoint_service_"
var totalRedpointKey = "redpoint_totle_"
var activityRedpointType = 1
var recommendRedpointType = 6

type Task struct {
	jobStr    string
	redisConn *redis.ClusterClient
}

func NewTask(jsonString string, redisConn *redis.ClusterClient) (t *Task) {
	var rp Task
	rp.jobStr = jsonString
	rp.redisConn = redisConn
	return &rp
}

/**
获取需要更改小红点的uid及消息类型
 */
func (t *Task) getTask() map[int]int {
	var uidMsgtypes map[int]int
	uidMsgtypes = make(map[int]int)

	jobStrArr := strings.SplitN(t.jobStr, "\n", -1)
	for i, jobInfo := range jobStrArr {
		if i%2 == 1 {
			json, err := simplejson.NewJson([]byte(jobInfo))
			if err != nil {
				errors.New("mcinsert update redpoint fail jobstr:" + jobInfo)
				continue
			}
			msgType, err := json.Get("type").Int()
			if err != nil {
				errors.New("mcinsert update redpoint fail jobstr:" + jobInfo)
				continue
			}

			uid, err := json.Get("uid").Int()
			if err != nil {
				errors.New("mcinsert update redpoint fail jobstr:" + jobInfo)
				continue
			}
			uidMsgtypes[uid] = msgType
		}
	}
	return uidMsgtypes
}

func (t *Task) ChangeRedisKeys() {
	uidMsgtypes := t.getTask()
	for uid, msgType := range uidMsgtypes {
		t.changeRedisKey(uid, msgType)
	}
}
func (t *Task) changeRedisKey(uid int, msgType int) {
	client := t.redisConn
	var key string
	if msgType == activityRedpointType {
		key = activityRedpointKey + strconv.Itoa(uid)
	} else if msgType == recommendRedpointType {
		key = recommendRedpointKey + strconv.Itoa(uid)
	} else {
		key = serviceRedpointKey + strconv.Itoa(uid)
	}
	result := (*client).Incr(key)
	_, err := result.Result()
	if err != nil {
		for i := 0; i < 3; i++ {
			result = (*client).Incr(key)
			_, err = result.Result()
			if err == nil {
				break
			}
			//休眠100毫秒后再次尝试
			time.Sleep(100 * time.Millisecond)
			if i ==3 && err!=nil {
				log.Error("incr key: ", key, " fail")
			}
		}
	}

	totalRedPointKey := totalRedpointKey + strconv.Itoa(uid)
	setResult := (*client).Set(totalRedPointKey, 1, 0)
	_,err = setResult.Result()
	if err != nil {
		for i := 0; i < 3; i++ {
			setResult := (*client).Set(totalRedPointKey, 1, 0)
			_,err = setResult.Result()
			if err == nil {
				break
			}
			//休眠100毫秒后再次尝试
			time.Sleep(100 * time.Millisecond)
			if i ==3 && err!=nil {
				log.Error("set key: ", totalRedPointKey, " fail")
			}
		}
	}
	log.Info("success redpoint uid: ", uid, "message type: ", msgType)
}
