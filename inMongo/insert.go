package inMongo

import (
	//"encoding/json"
	//mgo "gopkg.in/mgo.v2"
	//"strconv"
	"github.com/olivere/elastic"
	"gouminGitlab/common/orm/elasticsearch"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"github.com/bitly/go-simplejson"
	"gopkg.in/redis.v4"
)

type Worker struct {
	t *Task
}

func NewWorker(t *Task) (w *Worker) {
	//init the worker
	var wR Worker
	wR.t = t
	return &wR
}

//func (w Worker) Insert(session *mgo.Session) bool {
//	//convert json string to struct
//	var m row
//	if err := json.Unmarshal([]byte(w.t.columData), &m); err != nil {
//		//fmt.Println("[error] mongo json error", err, w.t.columData)
//		return false
//	}
//
//	//get the table name
//	tableNumber := strconv.Itoa(m.Uid % 1000)
//	tableName := "message_push_record_" + tableNumber
//
//	//create mongo session
//	c := session.DB("MessageCenter").C(tableName)
//
//	err := c.Insert(&m)
//	if err != nil {
//		//fmt.Println("[Error]insert into mongo error", err)
//		return false
//	}
//	return true
//}

func (w Worker) Insert(es *elastic.Client,redisConn *redis.ClusterClient) bool {
	//convert json string to struct
	job,jobErr := w.parseJson()
	if jobErr != nil {
		return false
	}

	//create elastic
	m := elasticsearch.NewMessagePushRecord(es,job)
	err := m.CreateRecord(redisConn)
	if err != nil {
		//fmt.Println("[Error]insert into mongo error", err)
		return false
	}
	return true
}
/**
json任务串转成struct
 */
func (w *Worker) parseJson() (*job.MsgPushRecordJsonColumn, error) {
	var jsonC job.MsgPushRecordJsonColumn
	js, err := simplejson.NewJson([]byte(w.t.columData))
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

