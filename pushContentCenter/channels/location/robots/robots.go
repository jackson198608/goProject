package robots

import (
	"errors"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
	"github.com/jackson198608/goProject/pushContentCenter/channels/location/job"
	"github.com/olivere/elastic"
	"strconv"
	//"gouminGitlab/common/orm/elasticsearch"
	"gouminGitlab/common/orm/mysql/member"
	log "github.com/thinkboy/log4go"
	"gouminGitlab/common/orm/mysql/robot"
	"time"
	"gouminGitlab/common/orm/mysql/recommend_data"
)

type Robots struct {
	mysqlXorm      []*xorm.Engine //@todo to be []
	mongoConn      []*mgo.Session //@todo to be []
	jsonData       *job.FocusJsonColumn
	esConn  *elastic.Client
}

//问题type = 8
const type_question = 8

const count = 100
func NewRobots(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, jsonData *job.FocusJsonColumn, esConn *elastic.Client) *Robots {
	if (mysqlXorm == nil) || (jsonData == nil) || (esConn ==nil){
		return nil
	}

	r := new(Robots)
	if r == nil {
		return nil
	}

	r.mysqlXorm = mysqlXorm
	r.mongoConn = mongoConn
	r.jsonData = jsonData
	r.esConn = esConn

	return r
}

func (r *Robots) Do() error {
	startId := 0
	for {
		//获取机器人数据
		currentPersionList,err := r.getPersons(startId)
		if err != nil {
			return err
		}
		endId, err := r.pushPersons(currentPersionList)
		startId = endId
		if err != nil {
			return err
		}
		if len(*currentPersionList) < count {
			break
		}
	}
	return nil
}

func (r *Robots) pushPersons(robots *[]member.PublishUser) (int, error) {
	if robots == nil {
		return 0, errors.New("push to robot : you have no person to push " + strconv.Itoa(r.jsonData.Infoid))
	}
	persons := *robots
	r.jsonData.Channel=1 //机器人数据，channel=1

	var endId int
	//elx,err := elasticsearch.NewEventLogX(r.esConn, r.jsonData)
	//
	//if err !=nil {
	//	return 0, err
	//}

	for _, person := range persons {
		//err := elx.PushPerson(person.RealUid)
		//改成存储mysql：如果是问答，只推有审核员身份的水军
		if r.jsonData.TypeId == type_question {
			//检查水军是否是审核员
			is_checker := r.isChecker(person.RealUid)
			if(is_checker == 0){
				continue
			}
		}
		_,err := r.saveContent(person.RealUid,r.jsonData)

		if err != nil {
			for i := 0; i < 5; i++ {
				log.Info("push robot ", person.RealUid, " try ", i, " by ",r.jsonData)
				//err := elx.PushPerson(person.RealUid)
				_,err := r.saveContent(person.RealUid,r.jsonData)
				if err == nil {
					break
				}
			}
		}
		endId = person.Id
	}

	return endId, nil
}

/**
	获取机器人uid
 */
//get fans persons by uid
func (r *Robots) getPersons(startId int) (*[]member.PublishUser, error){
	// var persons []int
	var uids []member.PublishUser
	err := r.mysqlXorm[3].Where("robot_uid =? and robot_nums>? and id>?", 0, 0, startId).Asc("id").Limit(count).Find(&uids)
	if err != nil {
		return &uids,err
	}

	return &uids,nil
}


func (r *Robots) formatCreatedToInt() int {
	date := r.jsonData.Created
	timeLayout := "2006-01-02 15:04:05"                       //转化所需模板
	loc,_ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, date, loc) //使用模板在对应时区转化为time.time类型
	datetime_str_to_timestamp := theTime.Unix()               //转化为时间戳 类型是int64
	createdInt := int(datetime_str_to_timestamp)
	return createdInt
}
/**
存储用户发送数据
 */
 func (r *Robots) saveContent(fuid int,jsonData *job.FocusJsonColumn) (int,error) {
 	 createdInt := r.formatCreatedToInt()
 	 exec := robot.UserPublishContent{Type:jsonData.TypeId,Fuid:fuid,Uid:jsonData.Uid,Title:jsonData.Title,Content:jsonData.Content,Images:jsonData.ImageInfo,Image:jsonData.Image,ImageNum:jsonData.Imagenums,ImageWidth:jsonData.ImageWidth,ImageHeight:jsonData.ImageHeight,VideoUrl:jsonData.VideoUrl,Duration:jsonData.Duration,Infoid:jsonData.Infoid,Pid:jsonData.Pid,Created:createdInt,TagInfo:jsonData.TagInfo,City:jsonData.City,IsDigest:jsonData.IsDigest,Source:jsonData.Source,RegisterTime:jsonData.RegisterTime,Inhome:jsonData.Inhome,Channel:jsonData.Channel,ThreadStatus:jsonData.ThreadStatus,IsPotentialKol:jsonData.IsPotentialKol}
 	 num, err := r.mysqlXorm[4].Insert(&exec)
	 if err != nil {
		 return 0, err
	 }
	 return int(num), nil

 }
 /**
 检查水军是否是审核员
  */
  func (r *Robots) isChecker(userid int) int {
	  var checker []recommend_data.WaterChecker
	  err := r.mysqlXorm[5].Cols("uid").Where("uid=?", userid).Find(&checker)
	  if err != nil {
		  log.Error("get table watre_checker error", err," uid:",userid)
		  return 0
	  }
	  if(checker != nil){
	  	return 1
	  }
  	  return 0
  }

