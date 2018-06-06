package task

import (
	"github.com/go-xorm/xorm"
	mgo "gopkg.in/mgo.v2"
	"errors"
	"strings"
	//"gouminGitlab/common/weixin/accessToken/accesstokenManager"
	"github.com/jackson198608/goProject/gouminMultiMessagePush/channels/wxProgram"
	"github.com/donnie4w/go-logger/logger"
	//"fmt"
	"gopkg.in/redis.v4"
)

type Task struct {
	Raw       string         //the data get from redis queue
	Appid     string
	Token     string
	Jobstr    string         //private member parse from raw
}

func NewTask(raw string, mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session) (*Task, error) {
	//check prams
	//if (raw == "") || (mysqlXorm == nil) || (mongoConn == nil) {
	//	return nil, errors.New("params can not be null")
	//}

	t := new(Task)
	if t == nil {
		return nil, errors.New("there is no space to create struct")
	}

	//pass params
	t.Raw = raw

	//create private member
	err := t.parseRaw()
	if err != nil {
		return nil, errors.New("raw format error ,can not find jobstr and jobtype " + raw)
	}

	return t, nil

}
func (t *Task) parseRaw() error {
	rawSlice := []byte(t.Raw)
	rawLen := len(rawSlice)
	lastIndex := strings.LastIndex(t.Raw, "|")
	t.Appid = string(rawSlice[0:lastIndex])
	t.Jobstr = string(rawSlice[lastIndex+1 : rawLen])
	return nil
}

func (t *Task) Do(redisConn *redis.ClusterClient,token string) error {

	//tokens := wxTokens.AccessTokens
	//tokenValue := tokens[t.Appid].Value
	//if tokenValue == "" {
	//	return errors.New("empty token value, appid:"+t.Appid)
	//}
	//请求微信，发送服务通知
	program := new(wxProgram.Task)
	program.AppId = t.Appid
	program.TaskJson = t.Jobstr
	program.AccessToken = token

	err := program.SendRequest()
	if err != nil {
		logger.Info("request weixin fail",err)
	}else{
		logger.Info("request weixn success")
	}
	return nil
}
