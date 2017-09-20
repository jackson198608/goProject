package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
	// "reflect"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strconv"
	"strings"
)

type Question struct {
	Id         int
	Uid        int
	Pid        int
	Typeid     int
	Subject    string
	Content    string
	Images     string
	Created    string
	Varieties  string
	Gender     string
	Age        string
	Browse_num int
	Ans_num    int
}

type Answer struct {
	Id       int
	Uid      int
	Qst_id   int
	Created  string
	Content  string
	Support  int
	Com_num  int
	Source   int
	Is_hide  int
	Audio_id int
}

func LoadQuestionById(id int, db *sql.DB) *Question {
	rows, err := db.Query("select id,uid,pid,type,subject,content,images,created,varieties,gender,age,browse_num,ans_num from `ask`.`ask_question` where is_hide=1 and id=" + strconv.Itoa(int(id)))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check ask_question sql prepare error: ", err)
		return nil
	}
	var varieties sql.NullString
	var gender sql.NullString
	var age sql.NullString
	var images string = ""
	var spe_name_s string = ""
	var row = new(Question)
	for rows.Next() {
		rows.Scan(&row.Id, &row.Uid, &row.Pid, &row.Typeid, &row.Subject, &row.Content, &images, &row.Created, &varieties, &gender, &age, &row.Browse_num, &row.Ans_num)
	}
	if gender.Valid {
		if gender.String == "1" {
			row.Gender = "母"
		} else {
			row.Gender = "公"
		}
	}
	if age.Valid {
		row.Age = age.String
	}
	row.Images = ""
	if images != "" {
		decoder := php_serialize.NewUnSerializer(images)
		result, err := decoder.Decode()
		if err != nil {
			logger.Info(err)
			return nil
		}
		real_result, _ := result.(php_serialize.PhpArray)
		var ids string = ""
		for _, v := range real_result {
			id, _ := v.(int)
			if id == 0 {
				sid, _ := v.(string)
				id, _ = strconv.Atoi(sid)
			}
			ids += strconv.Itoa(id) + ","
		}
		row.Images = getImageUrl(ids, db)
	}
	if varieties.Valid {
		rows1, err := db.Query("select spe_name_s from `new_dog123`.`dog_species` where spe_id=" + varieties.String)
		defer rows1.Close()
		if err != nil {
			logger.Error("[error] check dog_species sql prepare error: ", err)
			return nil
		}
		for rows1.Next() {
			rows1.Scan(&spe_name_s)
		}
	}
	row.Varieties = spe_name_s
	return row
}

func getImageUrl(ids string, db *sql.DB) string {
	ids = strings.Trim(ids, ",")
	rows, err := db.Query("select image from `ask`.`ask_attachment` where id in(" + ids + ")")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check ask_attachment sql prepare error: ", err)
		return ""
	}
	var images string = ""
	for rows.Next() {
		var image string
		rows.Scan(&image)
		images += askDomain + image + ","
	}
	images = strings.Trim(images, ",")
	return images
}

func LoadAnswersById(id int, db *sql.DB) []*Answer {
	rows, err := db.Query("select id,uid,qst_id,created,content,support,com_num,source,is_hide,audio_id from `ask`.`ask_answer` where is_hide=1 and qst_id=" + strconv.Itoa(int(id)) + " order by support desc")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check ask_answer sql prepare error: ", err)
		return nil
	}
	var rowsData []*Answer
	for rows.Next() {
		var row = new(Answer)
		rows.Scan(&row.Id, &row.Uid, &row.Qst_id, &row.Created, &row.Content, &row.Support, &row.Com_num, &row.Source, &row.Is_hide, &row.Audio_id)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

type Comment struct {
	Id        int
	Uid       int
	Ansid     int
	Comid     int
	Typeid    int
	Content   string
	Created   string
	Replyuser string
	Replyuid  int
}

func LoadCommentsById(id int, db *sql.DB) []*Comment {
	rows, err := db.Query("select id,uid,ans_id,com_id,type,content,created from `ask`.`ask_comment` where is_hide=1 and ans_id=" + strconv.Itoa(int(id)) + " order by id asc")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check ask_comment sql prepare error: ", err)
		return nil
	}
	var rowsData []*Comment
	for rows.Next() {
		var row = new(Comment)
		rows.Scan(&row.Id, &row.Uid, &row.Ansid, &row.Comid, &row.Typeid, &row.Content, &row.Created)
		if row.Typeid == 2 {
			oneComment := LoadOneCommentByComid(row.Comid, db)
			row.Replyuid = oneComment.Uid
		}
		rowsData = append(rowsData, row)
	}
	return rowsData
}

func LoadOneCommentByComid(com_id int, db *sql.DB) *Comment {
	rows, err := db.Query("select id,uid from `ask`.`ask_comment` where is_hide=1 and id=" + strconv.Itoa(int(com_id)))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check ask_comment sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row = new(Comment)
		rows.Scan(&row.Id, &row.Uid)
		return row
	}
	return &Comment{}
}

type Doctors struct {
	Uid        int
	Name       string
	Avatar     string
	Departname string
	Hospital   string
}

func LoadHealthDoctor(db *sql.DB) []*Doctors {
	rows, err := db.Query("select user_id,name,120_img,depart,belong_url from `new_dog123`.`dog_health` order by lasttime desc limit 3")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check dog_health sql prepare error: ", err)
		return nil
	}
	var rowsData []*Doctors
	for rows.Next() {
		var row = new(Doctors)
		var avatar string
		var departid int
		var belongurl int
		rows.Scan(&row.Uid, &row.Name, &avatar, &departid, &belongurl)
		row.Avatar = doctorUrl + avatar
		row.Departname = getDepartName(departid, db)
		row.Hospital = getShopName(belongurl, db)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

func getDepartName(departid int, db *sql.DB) string {
	rows, err := db.Query("select cat_name from `new_dog123`.`dog_ask_catalog` where cat_id=" + strconv.Itoa(departid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check dog_ask_catalog sql prepare error: ", err)
		return ""
	}
	var catname string = ""
	for rows.Next() {
		rows.Scan(&catname)
	}
	return catname
}

func getShopName(urlid int, db *sql.DB) string {
	rows, err := db.Query("select shop_name from `new_dog123`.`dog_shopinfo` where shop_id=" + strconv.Itoa(urlid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check dog_shopinfo sql prepare error: ", err)
		return ""
	}
	var shopname string = ""
	for rows.Next() {
		rows.Scan(&shopname)
	}
	return shopname
}

type AskRecommend struct {
	Id         int   "_id"
	Related    []int "related"
	RelatedAsk []int "related_ask"
}

func LoadRelateAskByAsk(id int, pid int, db *sql.DB, session *mgo.Session) []*RelateAsk {
	var rowsData []*RelateAsk
	ms := new(AskRecommend)
	c := session.DB("BigData").C("ask_recommend")
	err := c.Find(&bson.M{"_id": id}).One(&ms)
	if err != nil {
		logger.Error("BigData ask_recommend relate ask: ", err)
	}
	if len(ms.RelatedAsk) > 0 {
		idstring := ""
		for k, v := range ms.Related {
			if k <= 4 {
				idstring += strconv.Itoa(v) + ","
			}
		}
		rows, err := db.Query("select id,subject,browse_num from `ask`.`ask_question` where id in (" + strings.Trim(idstring, ",") + ")")
		defer rows.Close()
		if err != nil {
			logger.Error("[error] check ask_question sql prepare error: ", err)
			return nil
		}
		for rows.Next() {
			var row = new(RelateAsk)
			rows.Scan(&row.Id, &row.Subject, &row.Views)
			rowsData = append(rowsData, row)
		}
	}
	if len(ms.RelatedAsk) == 0 || len(rowsData) == 0 {
		if pid == 0 {
			pid = 10
		}
		rows, err := db.Query("select id,subject,browse_num from `ask`.`ask_question` where is_hide=1 and pid in (" + strconv.Itoa(pid) + ") order by ans_num desc limit 5")
		defer rows.Close()
		if err != nil {
			logger.Error("[error] check ask_question sql prepare error: ", err)
			return nil
		}
		for rows.Next() {
			var row = new(RelateAsk)
			rows.Scan(&row.Id, &row.Subject, &row.Views)
			rowsData = append(rowsData, row)
		}
	}
	return rowsData
}

func LoadRelateThreadByAsk(id int, db *sql.DB, session *mgo.Session) []*RelateThread {
	var rowsData []*RelateThread
	ms := new(AskRecommend)
	c := session.DB("BigData").C("ask_recommend")
	err := c.Find(&bson.M{"_id": id}).One(&ms)
	if err != nil {
		logger.Error("BigData ask_recommend relate thread: ", err)
	}
	if len(ms.Related) > 0 {
		idstring := ""
		for k, v := range ms.Related {
			if k <= 4 {
				idstring += strconv.Itoa(v) + ","
			}
		}
		rows, err := db.Query("select tid,subject,views,dateline from `new_dog123`.`pre_forum_thread` where displayorder>=0 and tid in (" + strings.Trim(idstring, ",") + ") ")
		defer rows.Close()
		if err != nil {
			logger.Error("check pre_forum_post sql prepare error: ", err)
			return nil
		}
		for rows.Next() {
			var row = new(RelateThread)
			rows.Scan(&row.Tid, &row.Subject, &row.Views, &row.Dateline)
			rowsData = append(rowsData, row)
		}
	}
	return rowsData
}

func LoadRelateDogByAsk(id int, pid int, db *sql.DB) []*RelateDog {
	if pid == 0 {
		pid = 10
	}
	rows, err := db.Query("select related from `new_dog123`.`related_pet` where spe_id=" + strconv.Itoa(pid))
	defer rows.Close()
	if err != nil {
		logger.Error("check relate_pet sql prepare error: ", err)
		return nil
	}
	doginfo := ""
	for rows.Next() {
		var row string
		rows.Scan(&row)
		doginfo += row + ","
	}
	doginfo = strings.Trim(doginfo, ",")
	if doginfo == "" {
		doginfo = "60,62,16,22,70,35"
	}
	fmt.Println("doginfo" + doginfo)
	rows1, err := db.Query("select dn.spe_id,spe_name_f,image_list from `new_dog123`.`dog_new_species` as dn left join `new_dog123`.`dog_species` as ds on dn.spe_id=ds.spe_id where dn.spe_id in(" + doginfo + ")")
	defer rows1.Close()
	if err != nil {
		logger.Error("check dog_new_species sql prepare error: ", err)
		return nil
	}
	var rowsData []*RelateDog
	for rows1.Next() {
		var row = new(RelateDog)
		rows1.Scan(&row.Speid, &row.Spename, &row.Img)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

func LoadDefaultRelateThreadByAsk(db *sql.DB) string {
	rows1, err := db.Query("select tid from `new_dog123`.`latestpost` limit 5")
	defer rows1.Close()
	if err != nil {
		logger.Error("check pre_forum_post sql prepare error: ", err)
		return ""
	}
	var ids string = ""
	for rows1.Next() {
		var tid int
		rows1.Scan(&tid)
		ids += strconv.Itoa(tid) + ","
	}
	rows, err := db.Query("select tid,subject,views,dateline from `new_dog123`.`pre_forum_thread` where displayorder>=0 and tid in (" + strings.Trim(ids, ",") + ") ")
	defer rows.Close()
	if err != nil {
		logger.Error("check pre_forum_post sql prepare error: ", err)
		return ""
	}
	var s string = ""
	for rows.Next() {
		var row = new(RelateThread)
		rows.Scan(&row.Tid, &row.Subject, &row.Views, &row.Dateline)
		if row.Views < 3000 {
			row.Views = rand.Intn(5000)
		}
		s += "<li><a href=\"" + mipBbsUrl + "thread-" + strconv.Itoa(row.Tid) + "-1-1.html\"><h3>" + row.Subject + "</h3><span>" + strconv.Itoa(row.Views) + "次浏览</span></a></li>"
	}
	return s
}
