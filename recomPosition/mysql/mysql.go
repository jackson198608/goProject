package mysql

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/go-sql-driver/mysql"
	// "reflect"
	// "github.com/menduo/gobaidumap"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Fuids struct {
	follow_id int "follow_id"
}

type Userinfo struct {
	uid         int    "uid"
	avatar      string "avatar"
	nickname    string "nickname"
	species     string "dog_species"
	pet_age     string "pet_age"
	typeId      int    "type"
	source_desc string
}

type Position struct {
	latitude  string "latitude"
	longitude string "longitude"
}

type Pet struct {
	dog_species int "dog_species"
	dog_birth_y int "dog_birth_y"
	dog_birth_m int "dog_birth_m"
}

type Petuids struct {
	uid int "uid"
}

type Goods struct {
	goods_id   int    "goods_id"
	goods_name int    "goods_name"
	goods_img  string "goods_img"
	price      string "price"
}

type Ad struct {
	aid      int    "aid"
	content  string "content"
	image    string "image"
	title    string "title"
	typeId   int    "type"
	showtime string "showtime"
}

//获取宠物 dog_species
func GetPetBreed(uid int, db *sql.DB) []*Breed {
	rows, err := db.Query("select distinct(dog_species) as bid from dog_doginfo where dog_userid=" + strconv.Itoa(uid))
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
	for i := 0; i < 1; i++ {
		var row = new(Breed)
		rows.Scan(0)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

//-------------------------------

//获取用户信息
func getUsers(uid int, db *sql.DB) []*Userinfo {
	//已经关注的人
	follows := getFollowedUids(uid, db)
	// fmt.Println(follows)
	//获取用户位置信息
	posi := getPositionByUid(uid, db)

	// near := NearbyUser(uid, follows, Position, db)
	return nil
}

//获取用户的认证身份类型
func getRauthinfoByUid(uid int, db *sql.DB) int {
	rows, err := db.Query("SELECT type as typeId FROM rauthentication WHERE (status=1) AND uid=" + strconv.Itoa(uid))
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check follow sql prepare error: ", err)
		return 0
	}
	for rows.Next() {
		var t int
		rows.Scan(&t)
		return t
	}
	return 0
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
func NearbyUser(uid int, followuids []int, Position []string, db *sql.DB) int {
	str := getStrByArr(followuids)
	rows, err := db.Query("SELECT uid FROM `user_location` WHERE latitude > (" + Position[0] + "-1) and latitude < (" + Position[0] + "+1) AND longitude > (" + Position[1] + "-1) and longitude < (" + Position[1] + "+1) AND uid NOT IN (" + str + ") ORDER BY ACOS(SIN(( " + Position[0] + " * 3.1415) / 180 ) *SIN((latitude * 3.1415) / 180 ) +COS((" + Position[0] + " * 3.1415) / 180 ) * COS((latitude * 3.1415) / 180 ) * COS((" + Position[1] + " * 3.1415) / 180 - (longitude * 3.1415) / 180 ) ) * 6380 LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check follow sql prepare error: ", err)
		return 0
	}
	for rows.Next() {
		var uid int
		rows.Scan(&uid)
		return uid
	}
	return 0

}

//获取用户的宠物信息
func getPetInfoByUid(uid int, db *sql.DB) []*Pet {
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
		rows.Scan(&row.dog_species, &row.dog_birth_y, &row.dog_birth_m)
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
			rows.Scan(&row.dog_species, &row.dog_birth_y, &row.dog_birth_m)
			rowsData = append(rowsData, row)
		}
		return rowsData
	}
	return rowsData

}

//相同犬种年龄的人
func getSameAgePetUsers(uid int, followuids []int, Pet []int, db *sql.DB) []*Petuids {
	str := getStrByArr(followuids)
	rows, err := db.Query("SELECT uid FROM pre_ucenter_members LEFT JOIN dog_doginfo ON pre_ucenter_members.uid = dog_doginfo.dog_userid WHERE (uid NOT IN ( " + str + " )) AND (((dog_birth_y=" + strconv.Itoa(Pet[1]) + ") AND (dog_birth_m=" + strconv.Itoa(Pet[2]) + ")) AND (dog_birth_d > 0)) LIMIT 3")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_ucenter_members sql prepare error: ", err)
		return nil
	}
	var rowsData []*Petuids
	for rows.Next() {
		var row = new(Petuids)
		rows.Scan(&row.uid)
		rowsData = append(rowsData, row)
	}
	return rowsData
}

//相同犬种的人
func getSameSpeciesPetUsers(uid int, followuids []int, Pet []int, db *sql.DB) []*Petuids {
	str := getStrByArr(followuids)
	rows, err := db.Query("SELECT uid FROM pre_ucenter_members LEFT JOIN dog_doginfo ON pre_ucenter_members.uid = dog_doginfo.dog_userid WHERE (uid NOT IN (" + str + ")) AND (dog_species=" + strconv.Itoa(Pet[0]) + ") LIMIT 4")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check pre_ucenter_members sql prepare error: ", err)
		return nil
	}
	var rowsData []*Petuids
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
func getSpeciesnameBySpeciesid(species int, db *sql.DB) string {
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
	str := getStrByArr(followfids)
	rows, err := db.Query("SELECT fid FROM pre_forum_forum WHERE ((name like '" + species + "%') AND (fup IN (76, 78, 2))) AND fid NOT IN (" + str + ") AND (status=1) LIMIT 1")
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
	}
	resp, err := client.Do(req)
	if err != nil {
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	return string(body)
}

//根据经纬度获得城市
func getCity(latitude string, longitude string) string {
	url := "http://api.map.baidu.com/geocoder?output=json&location=" + latitude + "," + longitude
	jsonStr := getUrl(url)
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
	str := getStrByArr(followfids)
	rows, err := db.Query("SELECT fid FROM pre_forum_forum WHERE ((name like '" + province + "%') AND (fup IN (76, 78, 2))) AND fid NOT IN (" + str + ") AND (status=1) LIMIT 1")
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
	str := getStrByArr(followfids)
	rows, err := db.Query("SELECT fid FROM pre_forum_forum WHERE ((status=1) AND (fid NOT IN (" + str + "))) AND (fup IN (76, 78, 2)) ORDER BY todayposts DESC LIMIT " + strconv.Itoa(count))
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
func getGoods(tag string) [][]string {
	tag = url.QueryEscape(tag)
	// age = url.QueryEscape(age)
	common := url.QueryEscape("通用")

	solr_url := "http://210.14.154.199:8983/solr/mall_goods/select?q=tags%3A(" + tag + "+OR+" + common + ")&fq=cat_id%3A*&fq=-stock%3A0&sort=sum_sales_count+desc&wt=json&indent=true"
	jsonStr := getUrl(solr_url)
	js, _ := simplejson.NewJson([]byte(jsonStr))
	status, _ := js.Get("responseHeader").Get("status").Int()
	var goods_infos [][]string
	if status == 0 {
		numFound, _ := js.Get("response").Get("numFound").Int()
		var docsLen int
		if numFound > 10 {
			docsLen = 10
		} else {
			docsLen = numFound
		}
		docs := js.Get("response").Get("docs")
		for i := 0; i < docsLen; i++ {
			name, _ := docs.GetIndex(i).Get("name").String()
			goods_id, _ := docs.GetIndex(i).Get("id").String()
			price, _ := docs.GetIndex(i).Get("lowest_price").String()
			img, _ := docs.GetIndex(i).Get("img").String()
			sales_count, _ := docs.GetIndex(i).Get("sum_sales_count").Int()
			stock, _ := docs.GetIndex(i).Get("stock").Int()
			salesCountStr := strconv.Itoa(sales_count)
			stockStr := strconv.Itoa(stock)
			goods_info := []string{goods_id, name, img, price, stockStr, salesCountStr}
			goods_infos = append(goods_infos, goods_info)
		}
	}
	return goods_infos
}

//------------------广告
func getAdInfo(db *sql.DB) []*Ad {
	t := time.Now().Format("2006-01-02")
	s := t + " 00:00:00"
	e := t + " 23:59:59"

	rows, err := db.Query("SELECT id as aid,type as typeId,title,content,image FROM `ads_recommend` WHERE showtime > '" + s + "' AND showtime < '" + e + "' order by weight desc LIMIT 1")
	defer rows.Close()
	if err != nil {
		logger.Error("[error] check ads_recommend sql prepare error: ", err)
		return nil
	}
	var rowsData []*Ad
	for rows.Next() {
		var row = new(Ad)
		rows.Scan(&row.aid, &row.typeId, &row.title, &row.content, &row.image)
		fmt.Println(row)
		rowsData = append(rowsData, row)
	}
	return rowsData
}
