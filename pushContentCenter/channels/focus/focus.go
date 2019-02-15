package focus

import (
	// "fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/allPersons"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/breedPersons"

	"github.com/jackson198608/goProject/pushContentCenter/channels/location/fansPersons"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/cardFansPersons"
	"fmt"
	"strconv"
	"github.com/olivere/elastic"
)

const (
	mongoConn = "192.168.5.22:27017,192.168.5.26:27017,192.168.5.200:27017"
	//mongoConn = "192.168.86.80:27017,192.168.86.81:27017,192.168.86.82:27017" //@todo change online dsn
)

var m map[int]bool

type Focus struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jobstr    string
	jsonData  *job.FocusJsonColumn
	esConn   *elastic.Client
}

//func init() {
//	m = make(map[int]bool)
//	m = loadActiveUserToMap()
//}

func NewFocus(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jobStr string, esConn *elastic.Client) *Focus {
	if (mysqlXorm == nil) || (mongoConn == nil) || (jobStr == "") {
		return nil
	}

	f := new(Focus)
	if f == nil {
		return nil
	}

	f.mysqlXorm = mysqlXorm
	f.mongoConn = mongoConn
	f.jobstr = jobStr
	f.esConn = esConn

	//@todo pass params
	jsonColumn, err := f.parseJson()
	if err != nil {
		return nil
	}
	f.jsonData = jsonColumn
	return f
}

//TypeId = 1 bbs, push fans active persons
//TypeId = 6 video, push fans active persons
//TypeId = 8 ask, push fans and breed active persons
//TypeId = 9 recommend bbs, push all active persons
//TypeId = 15 recommend video, push all active persons
//TypeId = 18 宠家号文章, push fans active persons
//TypeId = 19 宠家号视频, push fans active persons
//TypeId = 30 星球传记(事迹), push fans active persons
func (f *Focus) Do() error {
	//fmt.Println(f.jsonData)
	if f.jsonData.TypeId == 1 || f.jsonData.TypeId == 18 || f.jsonData.TypeId == 19 {
		//fmt.Println(f.jsonData.TypeId)
		f.jsonData.Source = 3
		fp := fansPersons.NewFansPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := fp.Do()
		if err != nil {
			return err
		}

		// f.jsonData.Source = 2
		// cp := clubPersons.NewClubPersons(f.mysqlXorm, f.mongoConn, f.jsonData)
		// err = cp.Do()
		// if err != nil {
		// 	return err
		// }
	} else if f.jsonData.TypeId == 6 {
		fp := fansPersons.NewFansPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := fp.Do()
		if err != nil {
			return err
		}
	} else if f.jsonData.TypeId == 8 {
		f.jsonData.Source = 3
		fp := fansPersons.NewFansPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := fp.Do()
		if err != nil {
			return err
		}

		f.jsonData.Source = 4
		bp := breedPersons.NewBreedPersons(f.mysqlXorm, f.mongoConn, f.jsonData,f.esConn)
		err = bp.Do()
		if err != nil {
			return err
		}
	} else if ((f.jsonData.TypeId == 9) || (f.jsonData.TypeId == 15)) && (f.jsonData.Source) == 1 {
		ap := allPersons.NewAllPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := ap.Do()
		if err != nil {
			return err
		}
	} else if f.jsonData.TypeId == 30 {
		fmt.Println("card json uid:" + strconv.Itoa(f.jsonData.Uid) + " infoid:" + strconv.Itoa(f.jsonData.Infoid))
		cfp := cardFansPersons.NewCardFansPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := cfp.Do()
		if err != nil {
			return err
		}
	} else {
		ap := allPersons.NewAllPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := ap.Do()
		if err != nil {
			return err
		}
	}
	return nil
}

//change json colum to object private member
func (f *Focus) parseJson() (*job.FocusJsonColumn, error) {
	var jsonC job.FocusJsonColumn
	js, err := simplejson.NewJson([]byte(f.jobstr))
	if err != nil {
		return &jsonC, err
	}

	jsonC.Uid, _ = js.Get("uid").Int()
	jsonC.TypeId, _ = js.Get("event_type").Int()
	jsonC.Created, _ = js.Get("time").String()
	jsonC.Tid, _ = js.Get("tid").Int()
	jsonC.Infoid, _ = js.Get("infoid").Int()
	jsonC.Bid, _ = js.Get("event_info").Get("bid").Int()
	jsonC.Title, _ = js.Get("event_info").Get("title").String()
	jsonC.Content, _ = js.Get("event_info").Get("content").String()
	jsonC.Forum, _ = js.Get("event_info").Get("forum").String()
	jsonC.Imagenums, _ = js.Get("event_info").Get("image_num").Int()
	jsonC.ImageInfo, _ = js.Get("event_info").Get("images").String()
	jsonC.VideoUrl, _ = js.Get("event_info").Get("video_url").String()
	jsonC.IsVideo, _ = js.Get("event_info").Get("is_video").Int()
	jsonC.Tag, _ = js.Get("event_info").Get("tag").Int()
	jsonC.Qsttype, _ = js.Get("event_info").Get("qst_type").Int()
	jsonC.Fid, _ = js.Get("event_info").Get("fid").Int()
	jsonC.Source, _ = js.Get("event_info").Get("source").Int()
	jsonC.Status, _ = js.Get("status").Int()
	jsonC.PetId, _ = js.Get("pet_id").Int()     //星球卡片id
	jsonC.PetType, _ = js.Get("pet_type").Int() // 宠物类型 1猫 2狗
	jsonC.Action, _ = js.Get("action").Int()    //行为 -1 删除 0 插入 1 修改

	return &jsonC, nil
}

//func loadActiveUserToMap() map[int]bool   {
//	var nodes []string
//	//@todo change to online
//	nodes = append(nodes, "http://192.168.86.230:9200")
//	nodes = append(nodes, "http://192.168.86.231:9200")
//
//	// is test config
//	//nodes = append(nodes, "http://192.168.6.50:9200")
//
//	//is online config
//	//nodes = append(nodes, "192.168.5.87:9500")
//	//nodes = append(nodes, "192.168.5.30:9500")
//	//nodes = append(nodes, "192.168.5.71:9500")
//	r,_ := elasticsearchBase.NewClient(nodes)
//	esConn,_:=r.Run()
//	var m map[int]bool
//	m = make(map[int]bool)
//	er := elasticsearch.NewUser(esConn)
//	from := 0
//	size := 1000
//	i :=1
//	for {
//		rst := er.SearchAllActiveUser(from, size)
//		total := rst.Hits.TotalHits
//		if total> 0 {
//			for _, hit := range rst.Hits.Hits {
//				//var userinfo elasticsearch.UserData
//				//err := json.Unmarshal(*hit.Source, &userinfo) //另外一种取数据的方法
//				//if err != nil {
//				//	fmt.Println("Deserialization failed", err)
//				//}
//				uid,_ := strconv.Atoi(hit.Id)
//				m[uid] = true
//			}
//		}
//		if int(total) < from {
//			break
//		}
//		i++
//		from = (i-1)*size
//	}
//	fmt.Println(m)
//	return m
//}
