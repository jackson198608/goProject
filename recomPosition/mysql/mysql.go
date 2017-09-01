package mysql

import (
	"database/sql"
	// "fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	// "reflect"
	// "github.com/menduo/gobaidumap"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	// "sort"
)

type Fuids struct {
	follow_id int "follow_id"
}
type Forum struct {
	Name  		string "name"
	Fid         int    "fid"
}

type ForumInfo struct {
	Fid         int    	`json:"fid"`
	Name  		string 	`json:"name"`
	Membernum   int 	`json:"membernum"`
	Icon      	string 	`json:"icon"`
}

type Userinfo struct {
	Uid         int    `json:"uid"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	TypeId      int    `json:"type"`
	SourceDesc  string `json:"source_desc"`
	Species     string `json:"dog_species"`
	PetAge      string `json:"pet_age"`
}

type Position struct {
	latitude  string 	"latitude"
	longitude string 	"longitude"
}

type Near struct {
	uid       int    	"uid"
	latitude  string 	"latitude"
	longitude string 	"longitude"
}

type Pet struct {
	DogSpecies int 	"dog_species"
	DogBirth_y int 	"dog_birth_y"
	DogBirth_m int 	"dog_birth_m"
}

type Petuids struct {
	uid int "uid"
}

type Goods struct {
	GoodsId   		string      `json:"goods_id"`
	GoodsName 		string    	`json:"goods_name"`
	GoodsImg  		string 		`json:"goods_img"`
	Price     		float64 	`json:"price"`
	Stock     	    int 		`json:"stock"`
	SalesCount      int 		`json:"sales_count"`
}


type AdInfo struct {
	Aid      int    `json:"aid"`
	TypeId   int    `json:"type"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Image    string `json:"image"`
}

type Ad struct {
	Aid      int    "id"
	Content  string "content"
	Image    string "image"
	Title    string "title"
	TypeId   int    "type"
	Showtime string "showtime"
}

//-------------------------------
//根据uid获取用户的信息
func GetUserInfoByUids(Fuid int, uids []int, db *sql.DB) []Userinfo {
	var userInfos []Userinfo
	//附近的人
	near := getNearUser(Fuid, db)
	if len(near)>0 {
		uid,_ := strconv.Atoi(near[0])
		m := Userinfo{uid, near[1],near[2],0,near[3],"",""}
		for i := 0; i < 1; i++ {
			userInfos = append(userInfos, m)
		}	
	}

	for key, uid := range uids {
		nickname := GetNickname(uid, db)
		avatar := GetAvatar(uid, db)
		var source_desc string
		if key == 2 {
			source_desc = "同一俱乐部"
		} else {
			source_desc = "可能认识"
		}
		m := Userinfo{uid, nickname,avatar,0,source_desc,"",""}
		userInfos = append(userInfos, m)
	}
	return userInfos
}

//根据犬种和年龄推荐的用户的uid
func GetUids(uid int, db *sql.DB) []int {
	var uids []int
	//已经关注的人
	follows := getFollowedUids(uid, db)
	//不能推荐自己
	follows = append(follows, uid)
	//附近的人
	near := getNearUser(uid, db)
	if len(near)>0 {
		nearuid, _ := strconv.Atoi(near[0])
		follows = append(follows, nearuid)
	}

	//相同犬种的人
	Pet := GetPetInfoByUid(uid, db)
	var pet []int
	for _, v := range Pet {
		pet = append(pet, v.DogSpecies)
		pet = append(pet, v.DogBirth_y)
		pet = append(pet, v.DogBirth_m)
	}

	speciesUids := getSameSpeciesPetUsers(uid, follows, pet, db)
	for _, value := range speciesUids {
		uids = append(uids, value.uid)
		follows = append(follows, value.uid)
	}

	//相同宠物年龄的人
	ageUids := getSameAgePetUsers(uid, follows, pet, db)
	for _, value1 := range ageUids {
		uids = append(uids, value1.uid)
	}
	return uids
}

//格式化附近的人
func getNearUser(uid int, db *sql.DB) []string {
	//推荐的用户uids
	//已经关注的人
	follows := getFollowedUids(uid, db)
	//不能推荐自己
	follows = append(follows, uid)
	//获取用户位置信息
	posi := getPositionByUid(uid, db)
	var Position []string
	for _, v := range posi {
		Position = append(Position, v.latitude)
		Position = append(Position, v.longitude)
	}
	// fmt.Println(Position)
	//附近的人信息数据格式化
	near := NearbyUser(uid, follows, Position, db)
	if near != nil {
		var nLatitude string
		var nLongitude string
		var uid int
		for _, v := range near {
			uid = v.uid
			nLatitude = v.latitude
			nLongitude = v.longitude
		}
		var nearInfo []string
		if len(Position)==0 {
			return nearInfo
		}

		lat1, err1 := strconv.ParseFloat(Position[0], 64)
		if err1 != nil {
			logger.Error("[error] can not get latitude error: ", err1)
		}		
		lng1, err2 := strconv.ParseFloat(Position[1], 64)
		if err1 != nil {
			logger.Error("[error] can not get longitude error: ", err2)
		}
		lat2, err3 := strconv.ParseFloat(nLatitude, 64)
		if err1 != nil {
			logger.Error("[error] can not get latitude error: ", err3)
		}
		lng2, err4 := strconv.ParseFloat(nLongitude, 64)
		if err1 != nil {
			logger.Error("[error] can not get longitude error: ", err4)
		}
		//格式化数据
		dis := getDistanceByLatitude(lat1, lng1, lat2, lng2)
		nickname := GetNickname(uid, db)
		avatar := GetAvatar(uid, db)
		source_desc := "相距" + strconv.Itoa(dis) + "米"
		nearInfo = []string{strconv.Itoa(uid), nickname, avatar, source_desc}
		return nearInfo
	}
	return nil
}

//根据经度和纬度得到距离
func getDistanceByLatitude(lat1, lng1, lat2, lng2 float64) int {
	radius := 6378.137
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	dis := dist * radius
	if int(dis) <= 0 {
		return 5
	}
	return int(dis)
}

//获得用户昵称
func GetNickname(uid int, db *sql.DB) string {
	rows, err := db.Query("SELECT mem_nickname FROM `dog_member` WHERE (`uid`=" + strconv.Itoa(uid) + ") LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check dog_member sql prepare error: ", err)
		return ""
	}
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return ""
		}
		return a
	}
	return ""
}

//获得用户头像
func GetAvatar(uid int, db *sql.DB) string {
	rows, err := db.Query("SELECT image FROM `album` WHERE (`uid`=" + strconv.Itoa(uid) + ") AND `type`=5 LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check album sql prepare error: ", err)
		return ""
	}
	if rows == nil {
		rows, err := db.Query("SELECT image FROM `album` WHERE (`uid`=" + strconv.Itoa(uid) + ") AND `type`=25 LIMIT 1")
		defer rows.Close()
		if err != nil {
			logger.Error("[error] check album sql prepare error: ", err)
			return ""
		}
	}
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return ""
		}
		return "http://c1.cdn.goumin.com/diary/" + a
	}
	return "http://c1.cdn.goumin.com/diary/head/cover-s.jpg"
}

//已经关注的人
func getFollowedUids(uid int, db *sql.DB) []int {
	tableName := "follow"
	rows, err := db.Query("select follow_id from " + tableName + " where user_id=" + strconv.Itoa(uid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check follow sql prepare error: ", err)
		return nil
	}
	var a []int
	for rows.Next() {
		var follow_id int
		rows.Scan(&follow_id)
		a = append(a, follow_id)
	}
	return a
}

//根据uid获取经度和纬度
func getPositionByUid(uid int, db *sql.DB) []*Position {
	rows, err := db.Query("SELECT latitude,longitude FROM `user_location` WHERE `uid`=" + strconv.Itoa(uid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check follow sql prepare error: ", err)
		return nil
	}

	var rowsData []*Position
	for rows.Next() {
		var row = new(Position)
		rows.Scan(&row.latitude, &row.longitude)
		rowsData = append(rowsData, row)
	}
	if rowsData == nil {
		return []*Position{{"39.9", "118.9"}}
	}
	return rowsData
}

//一个附近的人
func NearbyUser(uid int, followuids []int, Position []string, db *sql.DB) []*Near {
	var rowsData []*Near
	// str := getStrByArr(followuids)
	var uidsql string
	if len(Position)==0 {
		return rowsData
	}
	if len(followuids)!=0 {
		str := getStrByArr(followuids)
		uidsql = " AND uid NOT IN (" + str + ")"
	}
	rows, err := db.Query("SELECT uid,latitude,longitude FROM `user_location` WHERE latitude > (" + Position[0] + "-1) and latitude < (" + Position[0] + "+1) AND longitude > (" + Position[1] + "-1) and longitude < (" + Position[1] + "+1)"+ uidsql +" ORDER BY ACOS(SIN(( " + Position[0] + " * 3.1415) / 180 ) *SIN((latitude * 3.1415) / 180 ) +COS((" + Position[0] + " * 3.1415) / 180 ) * COS((latitude * 3.1415) / 180 ) * COS((" + Position[1] + " * 3.1415) / 180 - (longitude * 3.1415) / 180 ) ) * 6380 LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check follow sql prepare error: ", err)
		return nil
	}

	for rows.Next() {
		var row = new(Near)
		rows.Scan(&row.uid, &row.latitude, &row.longitude)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

//获取用户的宠物信息
func GetPetInfoByUid(uid int, db *sql.DB) []*Pet {
	tableName := "dog_doginfo"
	rows, err := db.Query("select dog_species,dog_birth_y,dog_birth_m from " + tableName + " where dog_userid= " + strconv.Itoa(uid) + " and 'default' = 1 limit 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check dog_doginfo sql prepare error: ", err)
		return nil
	}

	var rowsData []*Pet
	for rows.Next() {
		var row = new(Pet)
		rows.Scan(&row.DogSpecies, &row.DogBirth_y, &row.DogBirth_m)
		rowsData = append(rowsData, row)
	}
	if len(rowsData) == 0 {
		rows, err := db.Query("select dog_species,dog_birth_y,dog_birth_m from " + tableName + " where dog_userid= " + strconv.Itoa(uid) + " limit 1")
		defer rows.Close()
		if err != nil {
			logger.Error("[error] check dog_doginfo sql prepare error: ", err)
			return nil
		}
		for rows.Next() {
			var row = new(Pet)
			rows.Scan(&row.DogSpecies, &row.DogBirth_y, &row.DogBirth_m)
			rowsData = append(rowsData, row)
		}
		return rowsData
	}
	return rowsData
}

//相同犬种年龄的人
func getSameAgePetUsers(uid int, followuids []int, Pet []int, db *sql.DB) []*Petuids {
	// str := getStrByArr(followuids)
	var rowsData []*Petuids
	if len(Pet)==0 {
		return rowsData
	}
	var uidsql string
	if len(followuids)!=0 {
		str := getStrByArr(followuids)
		uidsql = " AND uid NOT IN (" + str + ")"
	}
	rows, err := db.Query("SELECT uid FROM pre_ucenter_members LEFT JOIN dog_doginfo ON pre_ucenter_members.uid = dog_doginfo.dog_userid WHERE dog_birth_y=" + strconv.Itoa(Pet[1]) + " AND dog_birth_m=" + strconv.Itoa(Pet[2]) + " AND dog_birth_d > 0 "+ uidsql +" LIMIT 3")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_ucenter_members sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row = new(Petuids)
		rows.Scan(&row.uid)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

//相同犬种的人
func getSameSpeciesPetUsers(uid int, followuids []int, Pet []int, db *sql.DB) []*Petuids {
	// str := getStrByArr(followuids)
	var rowsData []*Petuids
	if len(Pet)==0 {
		return rowsData
	}
	var uidsql string
	if len(followuids)!=0 {
		str := getStrByArr(followuids)
		uidsql = " AND uid NOT IN (" + str + ")"
	}
	rows, err := db.Query("SELECT uid FROM pre_ucenter_members LEFT JOIN dog_doginfo ON pre_ucenter_members.uid = dog_doginfo.dog_userid WHERE dog_species=" + strconv.Itoa(Pet[0]) + uidsql +" LIMIT 4")
	defer rows.Close()

	if err != nil {
		logger.Error("[error] check pre_ucenter_members sql prepare error: ", err)
		return nil
	}
	for rows.Next() {
		var row = new(Petuids)
		rows.Scan(&row.uid)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

//int型数组转为字符串
func getStrByArr(arr []int) string {
	str := ""
	for i := 0; i < len(arr); i++ {
		if i < len(arr)-1 {
			str += strconv.Itoa(arr[i]) + ","
		} else {
			str += strconv.Itoa(arr[i])
		}
	}
	return str
}

//-------------------俱乐部
//俱乐部数据格式化
func GetClubsInfo(fids []int, db *sql.DB) []ForumInfo {
	var clubsInfo []ForumInfo
	if len(fids)==0 {
		return clubsInfo
	}
	str := getStrByArr(fids)
	rows, err := db.Query("SELECT name,fid FROM pre_forum_forum WHERE (fid IN (" + str + ")) AND (status=1)")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_forum sql prepare error: ", err)
		return clubsInfo
	}
	for rows.Next() {
		var row = new(Forum)
		var name string
		var fid int
		rows.Scan(&row.Name,&row.Fid)
		name = row.Name
		fid = row.Fid
		icon := getClubIcon(fid, db)
		membernum := getClubMembers(fid, db)
		club := ForumInfo{fid, name, membernum, icon}
		clubsInfo = append(clubsInfo, club)
	}
	return clubsInfo
}

//俱乐部图标
func getClubIcon(fid int, db *sql.DB) string {
	fidStr := strconv.Itoa(fid)
	rows, err := db.Query("SELECT mobile_icon_thumb as icon FROM pre_forum_forumfield WHERE fid=" + fidStr)
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_forumfield sql prepare error: ", err)
		return ""
	}
	for rows.Next() {
		var icon string
		if err := rows.Scan(&icon); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return ""
		}
		return "http://f1.cdn.goumin.com/attachments/" + icon
	}
	return "http://c1.cdn.goumin.com/cms/picture/day_150814/20150814_4a95a1f.jpg"
}

//俱乐部总人数
func getClubMembers(fid int, db *sql.DB) int {
	rows, err := db.Query("SELECT COUNT(*) as numbers FROM forumfollow WHERE forum_id=" + strconv.Itoa(fid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check forumfollow sql prepare error: ", err)
		return 0
	}
	for rows.Next() {
		var numbers int
		if err := rows.Scan(&numbers); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return 0
		}
		return numbers
	}
	return 0
}

//获得推荐的俱乐部的fids
func GetFids(uid int, db *sql.DB) []int {
	var num int = 0
	var fids []int
	//已经加入的俱乐部
	followfids := getFollowedClubs(uid, db)
	//用户的宠物
	Pet := GetPetInfoByUid(uid, db)
	var species int
	for _, v := range Pet {
		species = v.DogSpecies
	}
	species_name := GetSpeciesnameBySpeciesid(species, db)
	//犬种俱乐部id
	fid1 := getPetClubByUid(species_name, followfids, db)
	if fid1 > 0 {
		followfids = append(followfids, fid1)
		num = num + 1
		fids = append(fids, fid1)
	}

	//获取用户位置信息
	posi := getPositionByUid(uid, db)
	var latitude string
	var longitude string
	for _, v := range posi {
		latitude = v.latitude
		longitude = v.longitude
	}
	province := getCity(latitude, longitude)
	//地域俱乐部id
	if province!="" {
		fid2 := getAreaClubByUid(province, followfids, db)
		if fid2 > 0 {
			followfids = append(followfids, fid2)
			num = num + 1
			fids = append(fids, fid2)
		}
	}

	count := 4 - num
	//热门俱乐部
	hotfids := getHotClubs(count, followfids, db)
	//数组合并
	fids = append(fids, hotfids...)
	return fids
}

//已经加入的俱乐部
func getFollowedClubs(uid int, db *sql.DB) []int {
	tableName := "forumfollow"
	rows, err := db.Query("select forum_id from " + tableName + " where user_id=" + strconv.Itoa(uid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check forumfollow sql prepare error: ", err)
		return nil
	}
	var a []int
	for rows.Next() {
		var forum_id int
		rows.Scan(&forum_id)
		a = append(a, forum_id)
	}
	return a
}

//根据犬种id获得犬种名称
func GetSpeciesnameBySpeciesid(species int, db *sql.DB) string {
	rows, err := db.Query("SELECT spe_name_s FROM dog_species WHERE spe_id=" + strconv.Itoa(species) + " LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check dog_species sql prepare error: ", err)
		return ""
	}
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return ""
		}
		logger.Info("[info] check sql find true")
		s := strings.Split(a, "/")
		return s[0]
	}
	return ""
}

//犬种俱乐部
func getPetClubByUid(species string, followfids []int, db *sql.DB) int {
	var fidsql string
	if len(followfids)!=0 {
		str := getStrByArr(followfids)
		fidsql = " AND fid NOT IN (" + str + ")"
	}
	rows, err := db.Query("SELECT fid FROM pre_forum_forum WHERE (name like '" + species + "%') AND fup IN (76, 78, 2)"+ fidsql +" AND status=1 LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_forum sql prepare error: ", err)
		return 0
	}
	for rows.Next() {
		var fid int
		if err := rows.Scan(&fid); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return 0
		}
		logger.Info("[info] check sql find true")
		return fid
	}
	return 0
}

func getUrl(url string) string {
	client := &http.Client{}
	logger.Info("get url address", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("getUrl error ", err)
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("getUrl error ", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("getUrl error ", err)
	}
	return string(body)
}

//根据经纬度获得城市
func getCity(latitude string, longitude string) string {
	url := "http://api.map.baidu.com/geocoder?output=json&location=" + latitude + "," + longitude
	jsonStr := getUrl(url)
	if jsonStr=="" {
		return ""
	}
	js, _ := simplejson.NewJson([]byte(jsonStr))
	status, _ := js.Get("status").String()
	if status == "OK" {
		city, _ := js.Get("result").Get("addressComponent").Get("province").String()
		c := UnicodeIndex(city)
		return c
	}
	return ""
}

//字符串截取
func UnicodeIndex(str string) string {
	p := strings.Trim(str, "省")
	c := strings.Trim(p, "市")
	// fmt.Println(c)
	return c
}

//地域俱乐部
func getAreaClubByUid(province string, followfids []int, db *sql.DB) int {
	var fidsql string
	if len(followfids)!=0 {
		str := getStrByArr(followfids)
		fidsql = " AND fid NOT IN (" + str + ")"
	}
	rows, err := db.Query("SELECT fid FROM pre_forum_forum WHERE (name like '" + province + "%') AND fup IN (76, 78, 2)"+ fidsql +" AND status=1 LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_forum sql prepare error: ", err)
		return 0
	}
	for rows.Next() {
		var fid int
		if err := rows.Scan(&fid); err != nil {
			logger.Error("[error] check sql get rows error ", err)
			return 0
		}
		return fid
	}
	return 0
}

//热门俱乐部
func getHotClubs(count int, followfids []int, db *sql.DB) []int {
	var fidsql string
	if len(followfids)!=0 {
		str := getStrByArr(followfids)
		fidsql = " AND fid NOT IN (" + str + ")"
	}
	rows, err := db.Query("SELECT fid FROM pre_forum_forum WHERE status=1 AND fup IN (76, 78, 2)"+ fidsql +" ORDER BY todayposts DESC LIMIT " + strconv.Itoa(count))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_forum_forum sql prepare error: ", err)
		return nil
	}
	var a []int
	for rows.Next() {
		var forum_id int
		rows.Scan(&forum_id)
		a = append(a, forum_id)
	}
	return a
}

//-------------商品

//根据宠物犬种搜索商品
func GetGoods(tag string) []Goods {
	tag = url.QueryEscape(tag)
	// age = url.QueryEscape(age)
	common := url.QueryEscape("通用")

	// solr_url := "http://210.14.154.199:8983/solr/mall_goods/select?q=tags%3A(" + tag + "+OR+" + common + ")&fq=cat_id%3A*&fq=-stock%3A0&sort=sum_sales_count+desc&wt=json&indent=true"
	solr_url := "http://192.168.5.75:8983/solr/mall_goods/select?q=tags%3A(" + tag + "+OR+" + common + ")&fq=cat_id%3A*&fq=-stock%3A0&sort=sum_sales_count+desc&wt=json&indent=true"
	jsonStr := getUrl(solr_url)
	js, _ := simplejson.NewJson([]byte(jsonStr))
	status, _ := js.Get("responseHeader").Get("status").Int()
	var goods_infos []Goods
	if status == 0 {
		numFound, _ := js.Get("response").Get("numFound").Int()
		var docsLen int
		if numFound > 8 {
			docsLen = 8
		} else {
			docsLen = numFound
		}
		docs := js.Get("response").Get("docs")
		for i := 0; i < docsLen; i++ {
			goods_name, _ := docs.GetIndex(i).Get("name").String()
			goods_id, _ := docs.GetIndex(i).Get("id").String()
			price, _ := docs.GetIndex(i).Get("lowest_price").Float64()
			goods_img, _ := docs.GetIndex(i).Get("img").String()
			sales_count, _ := docs.GetIndex(i).Get("sum_sales_count").Int()
			stock, _ := docs.GetIndex(i).Get("stock").Int()
			goods_info := Goods{goods_id, goods_name, "http://c1.cdn.goumin.com/cms"+goods_img, price, stock, sales_count}
			goods_infos = append(goods_infos, goods_info)
		}
	}
	return goods_infos
}

//------------------广告
func GetAdInfo(db *sql.DB) []AdInfo {
	var rowsData []AdInfo
	t := time.Now().Format("2006-01-02")
	s := t + " 00:00:00"
	e := t + " 23:59:59"
	rows, err := db.Query("SELECT id,type,title,content,image FROM backend.ads_recommend WHERE showtime > '" + s + "' AND showtime < '" + e + "' order by weight desc LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check ads_recommend sql prepare error: ", err)
		return rowsData
	}
	for rows.Next() {
		var row = new(Ad)
		rows.Scan(&row.Aid, &row.TypeId, &row.Title, &row.Content, &row.Image)
		m := AdInfo{row.Aid, row.TypeId, row.Title, row.Content, "http://c1.cdn.goumin.com/cms"+row.Image}
		rowsData = append(rowsData, m)
	}
	return rowsData
}
