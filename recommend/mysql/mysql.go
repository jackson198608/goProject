package mysql

import (
	"database/sql"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	// "fmt"
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
	Uid 		int		"uid"
	Thumb 		string	"thumb"
	Content 	string	"content"
	Created 	string	"created"
}

type Breed struct {
	Bid 	int		"bid"
}

//获取视频信息
func GetVideoData(vid int, db *sql.DB) []*Video {
	rows, err := db.Query("SELECT id,uid,thumb,content,created FROM `video` WHERE (`id`="+ strconv.Itoa(vid) +") AND (`status`=2)")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check video sql prepare error: ", err)
		return nil
	}
	var rowsData []*Video
	for rows.Next() {
		var row = new(Video)
		rows.Scan(&row.Id, &row.Uid, &row.Thumb, &row.Content, &row.Created)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

//检查帖子状态
func CheckThreadIsExist(tid int, db *sql.DB) int {
	//SELECT `pre_forum_thread`.* FROM `pre_forum_thread` LEFT JOIN `pre_forum_forum` ON `pre_forum_thread`.`fid` = `pre_forum_forum`.`fid` WHERE (((`tid`=4437081) AND (pre_forum_thread.displayorder >= 0 AND closed = 0)) AND (`pre_forum_forum`.`fup` IN (76, 78, 2))) AND (pre_forum_forum.status=1)
	rows, err := db.Query("SELECT `pre_forum_thread`.tid FROM `pre_forum_thread` LEFT JOIN `pre_forum_forum` ON `pre_forum_thread`.`fid` = `pre_forum_forum`.`fid` WHERE `tid`="+ strconv.Itoa(tid) +" AND pre_forum_thread.displayorder >= 0 AND closed = 0 AND `pre_forum_forum`.`fup` IN (76, 78, 2) AND pre_forum_forum.status=1")
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
		logger.Info("[info] check sql find true")
		return tid
	}
	return 0
}

//获取分表ID
func getTableid(tid int, db *sql.DB) int {
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
		logger.Info("[info] check sql find true")
		return posttableid
	}
	return 0
}

//获取帖子信息
func GetPostData(tid int, db *sql.DB) []*Post {
	tableName := "pre_forum_post"
	tableid := getTableid(tid, db)
	if  tableid != 0 {
		tableName = "pre_forum_post_" + strconv.Itoa(tableid)		
	}
	rows, err := db.Query("select tid,authorid,subject,message,name,Dateline from `"+ tableName +"` as p left join pre_forum_forum as f on p.fid=f.fid where f.status=1 and invisible=0 and first=1 and tid="+ strconv.Itoa(tid))
	logger.Info("select tid,authorid,subject,message,name,Dateline from `"+ tableName +"` as p left join pre_forum_forum as f on p.fid=f.fid where f.status=1 and invisible=0 and first=1 and tid="+ strconv.Itoa(tid))
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
		logger.Info("[info] check sql find true")
		return tid
	}
	return 0
}

//获取宠物 dog_species
func GetPetBreed(uid int, db *sql.DB) []*Breed {
	rows, err := db.Query("select distinct(dog_species) as bid from dog_doginfo where dog_userid="+ strconv.Itoa(uid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check dog_doginfo sql prepare error: ", err)
		return nil
	}
	var rowsData []*Breed
	for rows.Next() {
		var row = new(Breed)
		rows.Scan(&row.Bid)
		rowsData = append(rowsData, row)
	}
	for i:=0; i<1; i++ {
		var row = new(Breed)
		rows.Scan(0)
		rowsData = append(rowsData, row)
	}
	return rowsData
}
