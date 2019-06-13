package mcInsert

import (
	"github.com/olivere/elastic"
	"gouminGitlab/common/orm/elasticsearch"
	"github.com/bitly/go-simplejson"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/redis.v4"
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


func (w *Task) Insert(es *elastic.Client,redisConn *redis.ClusterClient) error {
	//convert json string to struct
	job,jobErr := w.parseJson()
	if jobErr != nil {
		return jobErr
	}

	//create elastic
	m := elasticsearch.NewMessagePushRecord(es,job)
	err := m.CreateRecord(redisConn)
	if err != nil {
		return err
	}
	return nil
}
/**
json任务串转成struct
 */
func (w *Task) parseJson() (*job.MsgPushRecordJsonColumn, error) {
	var jsonC job.MsgPushRecordJsonColumn
	js, err := simplejson.NewJson([]byte(w.columData))
	if err != nil {
		return &jsonC, err
	}

	jsonC.Uid, _ = js.Get("uid").Int()
	jsonC.Type, _ = js.Get("type").Int()
	jsonC.Created, _ = js.Get("created").String()
	jsonC.Mark, _ = js.Get("mark").Int()
	jsonC.Isnew, _ = js.Get("isnew").Int()
	jsonC.From, _ = js.Get("from").Int()
	jsonC.Title, _ = js.Get("title").String()
	jsonC.Content, _ = js.Get("content").String()
	jsonC.Channel, _ = js.Get("channel").Int()
	jsonC.ChannelTypes, _ = js.Get("channel_types").Int()
	jsonC.Image, _ = js.Get("image").String()
	jsonC.UrlType, _ = js.Get("url_type").Int()
	jsonC.Url, _ = js.Get("url").String()
	jsonC.Modified, _ = js.Get("modified").String()

	return &jsonC, nil
}
