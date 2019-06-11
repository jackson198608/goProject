package focus

import (
	// "fmt"
	"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/allPersons"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/fansPersons"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"gopkg.in/mgo.v2"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/cardFansPersons"
	"strconv"
	"github.com/olivere/elastic"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/CajFansPersons"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/robots"
	log "github.com/thinkboy/log4go"
)

type Focus struct {
	mysqlXorm []*xorm.Engine
	mongoConn []*mgo.Session
	jobstr    string
	jsonData  *job.FocusJsonColumn
	esConn   *elastic.Client
}

var fansChannel = 0
var fansAndRobotChannel = 1
var robotChannel = 2
func NewFocus(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jobStr string, esConn *elastic.Client) *Focus {
	if (mysqlXorm == nil) || (jobStr == "") ||( esConn == nil) {
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
//TypeId = 5 diary, push fans active persons
//TypeId = 6 video, push fans active persons
//TypeId = 8 ask, push fans and breed active persons
//TypeId = 9 recommend bbs, push all active persons
//TypeId = 15 recommend video, push all active persons
//TypeId = 18 宠家号文章, push fans active persons
//TypeId = 19 宠家号视频, push fans active persons
//TypeId = 30 星球传记(事迹), push fans active persons
func (f *Focus) Do() error {
	//fmt.Println(f.jsonData)
	log.Info("get task: ",f.jsonData)
	if f.jsonData.TypeId == 1 || f.jsonData.TypeId == 18 || f.jsonData.TypeId == 19 {
		//获取原始channel
		originalChannel := f.jsonData.Channel
		if f.jsonData.Channel== fansChannel || f.jsonData.Channel==fansAndRobotChannel {
			f.jsonData.Source = 3
			fp := fansPersons.NewFansPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
			err := fp.Do()
			f.jsonData.Channel = originalChannel //恢复channel
			if err != nil {
				return err
			}
		}

		if f.jsonData.TypeId == 1 && (f.jsonData.Channel == fansAndRobotChannel || f.jsonData.Channel == robotChannel) {
			//推送给机器人
			fr := robots.NewRobots(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
			frErr := fr.Do()
			f.jsonData.Channel = originalChannel //恢复channel
			if frErr != nil {
				return frErr
			}
		}

		// f.jsonData.Source = 2
		// cp := clubPersons.NewClubPersons(f.mysqlXorm, f.mongoConn, f.jsonData)
		// err = cp.Do()
		// if err != nil {
		// 	return err
		// }
	} else if f.jsonData.TypeId == 6 || f.jsonData.TypeId == 5{
		fp := fansPersons.NewFansPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := fp.Do()
		if err != nil {
			return err
		}

		//if f.jsonData.Channel == 1 {
		//	//推送给机器人
		//	fr := robots.NewRobots(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		//	frErr := fr.Do()
		//	if frErr != nil {
		//		return frErr
		//	}
		//}
	} else if f.jsonData.TypeId == 8 {
		f.jsonData.Source = 3
		fp := fansPersons.NewFansPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := fp.Do()
		if err != nil {
			return err
		}

		//f.jsonData.Source = 4
		//bp := breedPersons.NewBreedPersons(f.mysqlXorm, f.mongoConn, f.jsonData,f.esConn)
		//err = bp.Do()
		//if err != nil {
		//	return err
		//}

		//推送给机器人
		//fr := robots.NewRobots(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		//frErr := fr.Do()
		//if frErr != nil {
		//	f.jsonData.Channel = 0
		//	return frErr
		//}

	} else if ((f.jsonData.TypeId == 9) || (f.jsonData.TypeId == 15)) && (f.jsonData.Source) == 1 {
		ap := allPersons.NewAllPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := ap.Do()
		if err != nil {
			return err
		}
	} else if f.jsonData.TypeId == 30 {
		log.Info("card json uid: ", strconv.Itoa(f.jsonData.Uid) ," infoid:" , strconv.Itoa(f.jsonData.Infoid))
		cfp := cardFansPersons.NewCardFansPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := cfp.Do()
		if err != nil {
			return err
		}
	} else if f.jsonData.TypeId == 36{
		caj := CajFansPersons.NewCajFansPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
		err := caj.Do()
		if err != nil {
			return err
		}
	}
	//else {
	//	ap := allPersons.NewAllPersons(f.mysqlXorm, f.mongoConn, f.jsonData, f.esConn)
	//	err := ap.Do()
	//	if err != nil {
	//		return err
	//	}
	//}
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

	jsonC.AdoptId, _ = js.Get("event_info").Get("adopt_id").Int()
	jsonC.PetName,_ = js.Get("event_info").Get("pet_name").String()
	jsonC.PetAge,_ = js.Get("event_info").Get("pet_age").String()
	jsonC.PetBreed,_ = js.Get("event_info").Get("pet_breed").Int()
	jsonC.PetGender,_ = js.Get("event_info").Get("pet_gender").Int()
	jsonC.PetSpecies,_ = js.Get("event_info").Get("pet_species").String()
	jsonC.Province,_ = js.Get("event_info").Get("province").String()
	jsonC.City,_ = js.Get("event_info").Get("city").String()
	jsonC.County,_ = js.Get("event_info").Get("county").String()
	jsonC.Reason,_ = js.Get("event_info").Get("reason").String()
	jsonC.Image,_ = js.Get("event_info").Get("image").String()
	jsonC.PetImmunity, _ = js.Get("event_info").Get("pet_immunity").Int()
	jsonC.PetExpelling, _ = js.Get("event_info").Get("pet_expelling").Int()
	jsonC.PetSterilization, _ = js.Get("event_info").Get("pet_sterilization").Int()
	jsonC.PetStatus, _ = js.Get("event_info").Get("pet_status").Int()
	jsonC.AdoptStatus, _ = js.Get("event_info").Get("adopt_status").Int()
	jsonC.PetIntroduction, _ = js.Get("event_info").Get("pet_introduction").String()
	jsonC.UserIdentity, _ = js.Get("event_info").Get("user_identity").Int()
	jsonC.AdoptTag = js.Get("event_info").Get("adopt_tag").Interface()
	jsonC.PetAgenum, _ = js.Get("event_info").Get("pet_agenum").Int()
	jsonC.RegisterTime,_ = js.Get("event_info").Get("register_time").Int()
	jsonC.Channel,_ = js.Get("channel").Int()  //1推送给粉丝+机器人 0 仅推送给粉丝
	jsonC.Inhome,_ = js.Get("inhome").Int()  //是否是首页推荐内容  1是
	return &jsonC, nil
}
