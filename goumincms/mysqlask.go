package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
	// "reflect"
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
	rows, err := db.Query("select image from `ask_attachment` where id in(" + ids + ")")
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
	return images
}

func LoadAnswersById(id int, db *sql.DB) []*Answer {
	rows, err := db.Query("select id,uid,qst_id,created,content,support,com_num,source,is_hide,audio_id from `ask_answer` where is_hide=1 and qst_id=" + strconv.Itoa(int(id)) + " order by support desc")
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
	fmt.Println("...")
	return rowsData
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
	fmt.Println(rowsData)
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

func LoadRelateAsk(id int,db *sql.DB) []{
    
}
