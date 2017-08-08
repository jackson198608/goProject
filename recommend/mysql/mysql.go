package mysql

import (
	"database/sql"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"fmt"
)

type PostImage struct {
	Aid 			int		"aid"
	Attachment 		string	"attachment"
	MobileThumb 	string	"mobile_thumb"
	MobileMedium 	string	"mobile_medium"
	MobileSmall 	string	"mobile_small"
}

type Post struct {
	Tid 		int		"tid"
	Authorid 	int		"authorid" 
	Subject 	string	"subject"
	Message 	string	"message"
	Name 		string	"name"
	Dateline 	int		"dateline"
}

type Video struct {
	Id 			int		"id"
	Thumb 		string	"thumb"
	Content 	string	"content"
	Created 	string	"created"
}

type Breed struct {
	Bid 	int		"bid"
}

//获取视频信息
func GetVideoData(vid int, db *sql.DB) []*Video {
	//SELECT * FROM `video` WHERE (`id`=48428) AND (`status`=2)
	rows, err := db.Query("SELECT * FROM `video` WHERE (`id`="+ strconv.Itoa(vid) +") AND (`status`=2)")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check video sql prepare error: ", err)
		return nil
	}
	var rowsData []*Video
	for rows.Next() {
		var row = new(Video)
		rows.Scan(&row.Id, &row.Thumb, &row.Content, &row.Created)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

//检查帖子状态
func checkThreadIsExist(tid int, db *sql.DB) int {
	//SELECT `pre_forum_thread`.* FROM `pre_forum_thread` LEFT JOIN `pre_forum_forum` ON `pre_forum_thread`.`fid` = `pre_forum_forum`.`fid` WHERE (((`tid`=4437081) AND (pre_forum_thread.displayorder >= 0 AND closed = 0)) AND (`pre_forum_forum`.`fup` IN (76, 78, 2))) AND (pre_forum_forum.status=1)
	rows, err := db.Query("SELECT `pre_forum_thread`.* FROM `pre_forum_thread` LEFT JOIN `pre_forum_forum` ON `pre_forum_thread`.`fid` = `pre_forum_forum`.`fid` WHERE `tid`="+ strconv.Itoa(tid) +" AND pre_forum_thread.displayorder >= 0 AND closed = 0 AND `pre_forum_forum`.`fup` IN (76, 78, 2) AND pre_forum_forum.status=1")
	defer rows.Close()

	if err != nil {
		logger.Error("[error] check pre_forum_thread sql prepare error: ", err)
		return 0
	}

	for rows.Next() {
		var tid int
		if err := rows.Scan(&tid); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return 0
		}
		logger.Error("[info] check sql find true")
		return tid
	}
	return 0
}

//获取分表ID
func getTableid(tid int, db *sql.DB) int {
	//SELECT `pre_forum_thread`.* FROM `pre_forum_thread` LEFT JOIN `pre_forum_forum` ON `pre_forum_thread`.`fid` = `pre_forum_forum`.`fid` WHERE (((`tid`=4437081) AND (pre_forum_thread.displayorder >= 0 AND closed = 0)) AND (`pre_forum_forum`.`fup` IN (76, 78, 2))) AND (pre_forum_forum.status=1)
	rows, err := db.Query("SELECT posttableid FROM `pre_forum_thread` WHERE `tid`="+ strconv.Itoa(tid) +" AND pre_forum_thread.displayorder >= 0 AND closed = 0")
	defer rows.Close()

	if err != nil {
		logger.Error("[error] check pre_forum_thread sql prepare error: ", err)
		return 0
	}

	for rows.Next() {
		var posttableid int
		if err := rows.Scan(&posttableid); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return 0
		}
		logger.Error("[info] check sql find true")
		return posttableid
	}
	return 0
}

//获取帖子信息
func GetPostData(tid int, db *sql.DB) []*Post {
	tableName := "pre_forum_post"
	// tableid := tid%100;
	tableid := getTableid(tid, db)
	if  tableid != 0 {
		tableName = "pre_forum_post_" + strconv.Itoa(tableid)		
	}
	//select authorid,subject,message,f.name from pre_forum_post as p left join pre_forum_forum as f on p.fid=f.fid where invisible=0 and first=1 and tid=2073071;
	rows, err := db.Query("select tid,authorid,subject,message,name,Dateline from `"+ tableName +"` as p left join pre_forum_forum as f on p.fid=f.fid where f.status=1 and invisible=0 and first=1 and tid="+ strconv.Itoa(tid))
	fmt.Println("select tid,authorid,subject,message,name,Dateline from `"+ tableName +"` as p left join pre_forum_forum as f on p.fid=f.fid where f.status=1 and invisible=0 and first=1 and tid="+ strconv.Itoa(tid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check "+ tableName +" sql prepare error: ", err)
		return nil
	}
	var rowsData []*Post
	for rows.Next() {
		var row = new(Post)
		rows.Scan(&row.Tid, &row.Authorid, &row.Subject, &row.Message, &row.Name, &row.Dateline)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

//获取帖子标签信息
func GetTagData(tid int, db *sql.DB) int {
	//SELECT * FROM `tag_content` WHERE (`ctype`=3) AND (`cid`=4437081) LIMIT 1
	rows, err := db.Query("SELECT tid FROM `tag_content` WHERE (`ctype`=3) AND (`cid`="+ strconv.Itoa(tid) +") LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check tag_content sql prepare error: ", err)
		return 0
	}
	for rows.Next() {
		var tid int
		if err := rows.Scan(&tid); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return 0
		}
		logger.Error("[info] check sql find true")
		return tid
	}
	return 0
}

// //获取3张帖子图片
// func GetImageData(tid int, db *sql.DB) []*PostImage {
// 	tableid := tid % 10;
// 	tableName := "pre_forum_attachment_" + strconv.Itoa(tableid)

// 	//SELECT * FROM `pre_forum_attachment_1` WHERE (`tid`=4437081) AND (`isimage` IN (-1, 1)) LIMIT 3
// 	rows, err := db.Query("select aid,attachment,mobile_thumb,mobile_medium,mobile_small from `"+ tableName +"` where isimage in(-1,1) and tid="+ strconv.Itoa(tid) +" order by aid asc limit 3")
// 	defer rows.Close()
// 	if err != nil {
// 		logger.Error("[error] check "+ tableName +" sql prepare error: ", err)
// 		return nil
// 	}
// 	var rowsData []*PostImage
// 	for rows.Next() {
// 		var row = new(PostImage)
// 		rows.Scan(&row.Aid, &row.Attachment, &row.MobileThumb, &row.MobileMedium, &row.MobileSmall)
// 		rowsData = append(rowsData, row)
// 	}
// 	return rowsData
// }

// //获取帖子图片数
// func GetImageNum(tid int, db *sql.DB) int {
// 	tableid := tid % 10;
// 	tableName := "pre_forum_attachment_" + strconv.Itoa(tableid)

// 	//SELECT COUNT(*) FROM `pre_forum_attachment_1` WHERE `tid`=4437081
// 	rows, err := db.Query("select count(*) as num from `"+ tableName +"` where isimage in(-1,1) and tid="+ strconv.Itoa(tid))
// 	defer rows.Close()
// 	if err != nil {
// 		logger.Error("[error] check "+ tableName +" sql prepare error: ", err)
// 		return 0
// 	}
// 	// var rowsData []*PostImageNum
// 	for rows.Next() {
// 		var num int
// 		rows.Scan(&num)
// 		return num
// 	}
// 	return 0
// }

//获取宠物 dog_species
func GetPetBreed(uid int, db *sql.DB) []*Breed {
	rows, err := db.Query("select distinct(dog_species) as bid from member_pets where dog_userid="+ strconv.Itoa(uid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check member_pets sql prepare error: ", err)
		return nil
	}
	var rowsData []*Breed
	for rows.Next() {
		var row = new(Breed)
		rows.Scan(&row.Bid)
		rowsData = append(rowsData, row)
	}
	return rowsData
}
