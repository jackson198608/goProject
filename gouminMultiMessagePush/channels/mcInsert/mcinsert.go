package mcInsert

import (
	"github.com/olivere/elastic"
	//"gouminGitlab/common/orm/elasticsearch"
	//"github.com/bitly/go-simplejson"
	//"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/redis.v4"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"net/http"
	"encoding/json"
	log "github.com/thinkboy/log4go"
	"github.com/pkg/errors"
	"time"
)

type Task struct {
	columData string
}

type Worker struct {
	t *Task
}

func NewTask(jsonString string) (t *Task) {
	var tR Task
	tR.columData = jsonString
	return &tR
}

func (w *Task) Insert(es *elastic.Client, redisConn *redis.ClusterClient, esInfo string) error {
	var target = esInfo+"/_bulk"
	var h http.Header = make(http.Header)
	h.Set("Content-Type","application/x-ndjson")
	h.Set("timeout", "30")
	abuyun := w.getAbuyun()
	statusCode, _, body, err := abuyun.SendPostRequest(target, h, w.columData, true)
	if err != nil {
		return err
	}
	if statusCode == 200 {
		var result map[string]interface{}
		if err:=json.Unmarshal([]byte(body),&result);err==nil{
			if(result["errors"] == true){
				//curl请求es一次休眠100毫秒
				time.Sleep(100 * time.Millisecond)
				log.Error("bulk sync to es fail, error: ",body, " task: " ,w.columData)
				return errors.New("bulk sync to es fail, error message: " + body)
			}
		}
		log.Info("bulk to es success ", w.columData)
	}
	return nil
}

func (p *Task) getAbuyun() *abuyunHttpClient.AbuyunProxy {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy("", "", "")
	if abuyun == nil {
		//fmt.Println("create abuyun error")
		log.Error("create abuyun error")
		return nil
	}
	return abuyun
}
