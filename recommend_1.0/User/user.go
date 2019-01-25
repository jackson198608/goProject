package user

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"gopkg.in/mgo.v2"
	"net/http"
	// "reflect"
	"strconv"
	"strings"
	"time"
	"gouminGitlab/common/orm/elasticsearch"
	"github.com/gin-gonic/gin/json"
)

type elkClubBody struct {
	Id          int
	name        string
	membernum   string
	todayposts  int
	description string
	fup         int
	icon        string
}

type elkUserBody struct {
	Id            int
	nickname      string
	avatar        string
	grouptitle    string
	pets          string
	address       string
	follow_clubs  string
	follow_users  string
	lastlogintime int
	is_welluser   int
	age           string
}

type User struct {
	mysqlXorm       []*xorm.Engine
	mongoConn       []*mgo.Session
	Uid             int
	elkDsn          string
	nodes           []string
	province        string
	species         string
	address         string
	age             string
	myData          elkUserBody
	notRecommendUid []int
	notRecommendFid []int
}

var recommendClubData []elasticsearch.RecommendClubData
var recommendUserData []elasticsearch.RecommendUserData

func NewUser(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, uid string, elkDsn string) *User {
	if (mysqlXorm == nil) || (mongoConn == nil) || (uid == "") || (elkDsn == "") {
		return nil
	}
	logger.Info("start recommend: ", uid)
	u := new(User)
	if u == nil {
		return nil
	}
	u.mysqlXorm = mysqlXorm
	u.mongoConn = mongoConn
	u.Uid, _ = strconv.Atoi(uid)
	u.nodes = strings.SplitN(elkDsn, ",", -1)
	u.elkDsn = u.nodes[0]
	return u
}

func (u *User) setAbuyun() *abuyunHttpClient.AbuyunProxy {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy("", "", "")

	if abuyun == nil {
		fmt.Println("create abuyun error")
		return nil
	}
	return abuyun
}

func (u *User) Do() error {
	if u.Uid <= 0 {
		return errors.New("uid is error; uid is " + strconv.Itoa(u.Uid))
	}
	u.notRecommendUid = append(u.notRecommendUid, u.Uid)
	err := u.getMyData() //获取我的数据
	if err == nil {
		u.getRecommendClub()
		u.getRecommendUser()
	}
	return err
}

func (u *User) getMyData() error {
	uidStr := strconv.Itoa(u.Uid)
	query := u.getUserQueries(uidStr, 1, 1)
	user, _ := u.getUser(query)
	if user != nil {
		data := *user
		if len(data) > 0 {
			u.myData = data[0]
			u.species = u.getSpecies()   //我的宠物品种
			u.age = u.getAge()           //我的宠物年龄
			u.province = u.getProvince() //我的地域
			u.address = u.getAddress()

			if u.myData.follow_users != "" {
				follow_users := strings.Split(u.myData.follow_users, ",")
				for f, _ := range follow_users {
					follow_uid, _ := strconv.Atoi(follow_users[f])
					u.notRecommendUid = append(u.notRecommendUid, follow_uid)
				}
			}
			if u.myData.follow_clubs != "" {
				follow_clubs := strings.Split(u.myData.follow_clubs, ",")
				for f, _ := range follow_clubs {
					follow_fid, _ := strconv.Atoi(follow_clubs[f])
					u.notRecommendFid = append(u.notRecommendFid, follow_fid)
				}
			}
			return nil
		}
	}
	return errors.New("get my info error by uid is " + uidStr)
}

func (u *User) getAge() string {
	age := ""
	pets := u.myData.pets
	if pets != "" {
		petItems := strings.Split(pets, "|")
		for p, _ := range petItems {
			petItem := strings.Split(petItems[p], ",")
			if len(petItem) > 4 {
				age += petItem[4] + ";"
			}
		}
		if age != "" {
			age = string(age[0 : len(age)-1])
		}
	}
	return age
}

func (u *User) getSpecies() string {
	species := ""
	pets := u.myData.pets
	if pets != "" {
		petItems := strings.Split(pets, "|")
		for p, _ := range petItems {
			petItem := strings.Split(petItems[p], ",")
			if len(petItem) > 4 {
				species += petItem[3] + ";"
			}
		}
		if species != "" {
			species = string(species[0 : len(species)-1])
		}
	}
	if species == "" {
		species = "金毛"
	}
	return species
}

func (u *User) getAddress() string {
	formatted_address := ""
	address := u.myData.address
	if address != "" {
		addressItems := strings.Split(address, ";")
		if len(addressItems) > 2 {
			formatted_address = addressItems[2]
		}
	}
	return formatted_address
}

func (u *User) getProvince() string {
	province := ""
	address := u.myData.address
	if address != "" {
		addressItems := strings.Split(address, ";")
		if len(addressItems) > 2 {
			formatted_address := addressItems[2]
			provinceAry := strings.Split(formatted_address, "省")
			if len(provinceAry) > 1 {
				province = provinceAry[0]
			} else {
				provinceAry := strings.Split(formatted_address, "市")
				if len(provinceAry) > 1 {
					province = provinceAry[0]
				}
			}
		}
	}
	return province
}

func (u *User) getRecommendClub() error {
	followNum := 0
	if len(u.notRecommendFid) != 0 {
		followNum, _ = u.followClubs()
		logger.Info("[club] user followNum is ", strconv.Itoa(followNum), " by uid ", strconv.Itoa(u.Uid))
	}
	if followNum < 6 {
		speciesNum, _ := u.recommendClubBySpecies() //犬种
		logger.Info("[club] recommend speciesNum is ", strconv.Itoa(speciesNum), " by uid ", strconv.Itoa(u.Uid))
		if followNum+speciesNum < 6 {
			addressNum, _ := u.recommendClubByAddress() //地域
			logger.Info("[club] recommend addressNum is ", strconv.Itoa(addressNum), " by uid ", strconv.Itoa(u.Uid))
			if followNum+addressNum+speciesNum < 6 {
				fidNum, _ := u.recommendClubByFid(159) //训练
				logger.Info("[club] recommend fidNum is ", strconv.Itoa(fidNum), "by fid is 159 by uid ", strconv.Itoa(u.Uid))
				if followNum+addressNum+speciesNum+fidNum < 6 {
					fid1Num, _ := u.recommendClubByFid(10) //巧手
					logger.Info("[club] recommend fid1Num is ", strconv.Itoa(fid1Num), "by fid is 10 by uid ", strconv.Itoa(u.Uid))
					if followNum+addressNum+speciesNum+fidNum+fid1Num < 6 {
						FupNum, _ := u.recommendClubByFup(2) //综合
						logger.Info("[club] recommend FupNum is ", strconv.Itoa(FupNum), " by uid ", strconv.Itoa(u.Uid))

					}
				}
			}
		}
	}
	u.pushClubRecommend()
	return nil
}

func (u *User) getRecommendUser() error {
	speciesNum, _ := u.recommendUserBySpecies(0, 5) //相同犬种
	ageNum, _ := u.recommendUserByAge()             //相同年龄
	logger.Info("[user] recommend speciesNum is ", strconv.Itoa(speciesNum), " by uid ", strconv.Itoa(u.Uid))
	logger.Info("[user] recommend ageNum is ", strconv.Itoa(ageNum), " by uid ", strconv.Itoa(u.Uid))
	if ageNum+speciesNum < 8 {
		num := 8 - ageNum - speciesNum
		addressNum, _ := u.recommendUserByAddress(num)
		logger.Info("[user] recommend addressNum is ", strconv.Itoa(addressNum), " by uid ", strconv.Itoa(u.Uid))
		if addressNum+speciesNum+ageNum < 8 {
			num := 8 - addressNum - speciesNum - ageNum
			nextSpeciesNum, _ := u.recommendUserBySpecies(1, num) //相同犬种
			logger.Info("[user] recommend nextSpeciesNum is ", strconv.Itoa(nextSpeciesNum), " by uid ", strconv.Itoa(u.Uid))
		}
	}
	u.pushUserRecommend()
	return nil
}

//根据犬种推荐用户
func (u *User) recommendUserBySpecies(isFirst int, num int) (int, error) {
	if u.species != "" {
		speciesItems := strings.Split(u.species, ";")
		speciesKeyword := ""
		for s, _ := range speciesItems {
			speciesKeyword += `\"` + speciesItems[s] + `\"` + ","
		}
		speciesKeyword = string(speciesKeyword[0 : len(speciesKeyword)-1])
		getType := 3
		if isFirst == 0 {
			getType = 4
		}
		query := u.getUserQueries(speciesKeyword, getType, num) //获取根据犬种查询条件
		fmt.Println("species query:")
		fmt.Println(query)
		user, err := u.getUser(query)
		if err != nil {
			fmt.Println("get user error, by " + speciesKeyword)
			return 0, nil
		}
		if user != nil {
			u.buildRecommendUserData(user, 1)
			return len(*user), nil
		}
	}
	return 0, nil
}

//根据地域推荐用户
func (u *User) recommendUserByAddress(num int) (int, error) {
	if u.address != "" {
		query := u.getUserQueries(u.address, 0, num) //获取根据地址查询条件
		fmt.Println("address query:")
		fmt.Println(query)
		user, err := u.getUser(query)
		if err != nil {
			fmt.Println("get user error, by " + u.address)
			return 0, nil
		}
		if user != nil {
			u.buildRecommendUserData(user, 2)
			return len(*user), nil
		}
	}
	return 0, nil
}

//根据年龄推荐用户
func (u *User) recommendUserByAge() (int, error) {
	if u.age != "" {
		ageItems := strings.Split(u.age, ";")
		ageKeyword := ""
		for s, _ := range ageItems {
			ageKeyword += `\"` + ageItems[s] + `\"` + ","
		}
		ageKeyword = string(ageKeyword[0 : len(ageKeyword)-1])
		query := u.getUserQueries(ageKeyword, 3, 3) //获取根据年龄查询条件
		fmt.Println("age query:")
		fmt.Println(query)
		user, err := u.getUser(query)
		if err != nil {
			fmt.Println("get user error, by " + u.age)
			return 0, nil
		}
		if user != nil {
			u.buildRecommendUserData(user, 0)
			return len(*user), nil
		}
	}
	return 0, nil
}

// 存储用户数据
func (u *User) pushUserRecommend() error {
	jsonData, _ := json.Marshal(recommendUserData)
	dataStr := string(jsonData)
	err := u.insertData(dataStr, 2)
	return err
}

// -----   推荐俱乐部  ------

//获取已关注的俱乐部
func (u *User) followClubs() (int, error) {
	followIds := ""
	for i, _ := range u.notRecommendFid {
		followIds += strconv.Itoa(u.notRecommendFid[i]) + ","
	}
	followIds = string(followIds[0 : len(followIds)-1])
	query := u.getClubQueries(followIds, 0, 0, followIds) //获取根据已关注的俱乐部查询条件
	club, err := u.getClub(query)
	if err != nil {
		fmt.Println("get club error, by " + u.species)
		return 0, nil
	}
	if club != nil {
		u.buildRecommendClubData(club)
		return len(*club), nil
	}
	return 0, nil
}

//根据犬种推荐俱乐部
func (u *User) recommendClubBySpecies() (int, error) {
	if u.species == "" {
		return 0, nil
	}
	speciesItems := strings.Split(u.species, ";")
	speciesKeyword := ""
	for s, _ := range speciesItems {
		speciesKeyword += `\"` + speciesItems[s] + `\"` + ","
	}
	speciesKeyword = string(speciesKeyword[0 : len(speciesKeyword)-1])
	query := u.getClubQueries(speciesKeyword, 0, 0, "") //获取根据犬种查询条件
	club, err := u.getClub(query)
	if err != nil {
		fmt.Println("get club error, by " + u.species)
		return 0, nil
	}
	if club != nil {
		u.buildRecommendClubData(club)
		return len(*club), nil
	}
	return 0, nil
}

func (u *User) recommendClubByFid(fid int) (int, error) {
	query := u.getClubQueries("", 0, fid, "") //获取根据犬种查询条件
	club, err := u.getClub(query)
	if err != nil {
		fmt.Println("get club error, by " + strconv.Itoa(fid))
		return 0, nil
	}
	if club != nil {
		u.buildRecommendClubData(club)
		return len(*club), nil
	}
	return 0, nil
}

func (u *User) recommendClubByAddress() (int, error) {
	if u.province != "" {
		provinceKeyword := `\"` + u.province + `\"`
		query := u.getClubQueries(provinceKeyword, 0, 0, "") //获取根据地址查询条件
		club, err := u.getClub(query)
		if err != nil {
			fmt.Println("get club error, by " + u.province)
			return 0, nil
		}
		if club != nil {
			u.buildRecommendClubData(club)
			return len(*club), nil
		}
	}
	return 0, nil
}

func (u *User) recommendClubByFup(fup int) (int, error) {
	query := u.getClubQueries("", fup, 0, "") //获取根据地址查询条件
	club, err := u.getClub(query)
	if err != nil {
		fmt.Println("get club error, by " + strconv.Itoa(fup))
		return 0, nil
	}
	if club != nil {
		u.buildRecommendClubData(club)
		return len(*club), nil
	}
	return 0, nil
}

func (u *User) pushClubRecommend() error {
	jsonData, _ := json.Marshal(recommendClubData)
	dataStr := string(jsonData)
	err := u.insertData(dataStr, 1)
	if err != nil {
		return err
	}
	return nil
}

//-----  base function  ---

func (u *User) getUserQueries(keyword string, getType int, size int) string {
	query := ""
	mustNotQuery := ""
	filterQuery := ""
	filterQuery += "\"post_filter\":{\"bool\":{\"must_not\":["
	for m, _ := range u.notRecommendUid {
		mustNotQuery += "{\"term\":{\"id\":" + strconv.Itoa(u.notRecommendUid[m]) + "}},"
	}
	filterQuery += string(mustNotQuery[0 : len(mustNotQuery)-1])
	filterQuery += "]}},"
	sizeStr := strconv.Itoa(size)
	// 达人数据
	if getType == 2 {
		query = "{\"query\": {\"query_string\":{\"query\":\"1\",\"fields\":[\"is_welluser\"]}},\"sort\": { \"lastlogintime\": { \"order\": \"desc\" }}}"
	} else if getType == 3 {
		//相同品种或年龄
		query = "{\"size\" : " + sizeStr + ",\"query\": {\"query_string\":{\"query\":\"" + keyword + "\",\"fields\":[\"pets\"]}}," + filterQuery + "\"sort\": { \"lastlogintime\": { \"order\": \"desc\" }}}"
	} else if getType == 1 {
		//我的数据
		query = "{\"query\": {\"query_string\":{\"query\":\"" + keyword + "\",\"fields\":[\"id\"]}}}"
	} else if getType == 4 {
		//补偿机制 相同品种
		query = "{\"size\" : " + sizeStr + ",\"query\": {\"query_string\":{\"query\":\"" + keyword + "\",\"fields\":[\"pets\"]}}," + filterQuery + "\"sort\": { \"lastlogintime\": { \"order\": \"desc\" }}}"
	} else {
		//地域相近
		query = "{\"size\" : " + sizeStr + ",\"query\": {\"query_string\":{\"query\":\"" + keyword + "\",\"fields\":[\"address\"]}}," + filterQuery + "\"sort\": { \"lastlogintime\": { \"order\": \"desc\" }}}"
	}
	return query
}

//获取用户数据
func (u *User) getUser(query string) (*[]elkUserBody, error) {
	abuyun := u.setAbuyun()
	targetUrl := u.elkDsn + "/user/user_info/_search?pretty"
	var h http.Header = make(http.Header)

	h.Set("Content-Type", "application/json")
	statusCode, _, body, err := abuyun.SendRequest(targetUrl, h, query, true)
	if err != nil {
		fmt.Println("http request error", err)
		return nil, err
	}
	if statusCode == 200 {
		user, err := u.formatUser(body)
		if err != nil {
			fmt.Println("format user data error", err)
			return nil, err
		}
		return user, nil
	}
	return nil, err
}

func (u *User) formatUser(body string) (*[]elkUserBody, error) {
	var user []elkUserBody
	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return &user, err
	}
	hits, _ := js.Get("hits").Get("hits").Array()
	for i, _ := range hits {
		var userBody elkUserBody
		source := js.Get("hits").Get("hits").GetIndex(i).Get("_source")
		userBody.Id, _ = source.Get("id").Int()
		userBody.nickname, _ = source.Get("nickname").String()
		userBody.grouptitle, _ = source.Get("grouptitle").String()
		userBody.follow_clubs, _ = source.Get("follow_clubs").String()
		userBody.follow_users, _ = source.Get("follow_users").String()
		userBody.is_welluser, _ = source.Get("is_welluser").Int()
		userBody.lastlogintime, _ = source.Get("lastlogintime").Int()
		userBody.pets, _ = source.Get("pets").String()
		userBody.address, _ = source.Get("address").String()
		user = append(user, userBody)
	}
	return &user, nil
}

//根据关键词搜索相关俱乐部
//fup=76 各地俱乐部
//fup=78 犬种俱乐部
//fup=2 综合论坛
//fup=0时, 不限

func (u *User) getClubQueries(keyword string, fup int, fid int, followIds string) string {
	query := ""
	mustNotQuery := ""
	filterQuery := ""
	if len(u.notRecommendFid) > 0 {
		filterQuery += "\"post_filter\":{\"bool\":{\"must_not\":["
		for m, _ := range u.notRecommendFid {
			mustNotQuery += "{\"term\":{\"id\":" + strconv.Itoa(u.notRecommendFid[m]) + "}},"
		}
		filterQuery += string(mustNotQuery[0 : len(mustNotQuery)-1])
		filterQuery += "]}},"
	}
	//综合版区
	if followIds != "" {
		query = "{\"size\" : 6,\"query\": {\"bool\": {\"must\": {\"terms\": {\"id\": [" + followIds + "]}}}}}"
	} else if fup == 2 {
		fupStr := strconv.Itoa(fup)
		query = "{\"size\" : 4,\"query\": {\"query_string\":{\"query\":\"" + fupStr + "\",\"fields\":[\"fup\"]}}," + filterQuery + "\"sort\": { \"todayposts\": { \"order\": \"desc\" }}}"
	} else if fid != 0 {
		fidStr := strconv.Itoa(fid)
		query = "{\"query\": {\"query_string\":{\"query\":\"" + fidStr + "\",\"fields\":[\"id\"]}}}"
	} else {
		query = "{\"size\" : 6,\"query\": {\"query_string\":{\"query\":\"" + keyword + "\",\"fields\":[\"name\",\"description\"]}}," + filterQuery + "\"sort\": { \"todayposts\": { \"order\": \"desc\" }}}"
	}
	return query
}

func (u *User) getClub(query string) (*[]elkClubBody, error) {
	abuyun := u.setAbuyun()
	targetUrl := u.elkDsn + "/club/club_info/_search"
	var h http.Header = make(http.Header)
	h.Set("Content-Type", "application/json")
	statusCode, _, body, err := abuyun.SendRequest(targetUrl, h, query, true)
	if err != nil {
		fmt.Println("http request error", err)
		return nil, err
	}
	if statusCode == 200 {
		club, err := u.formatClub(body)
		if err != nil {
			fmt.Println("format club data error", err)
		}
		return club, nil
	}
	return nil, err
}

func (u *User) formatClub(body string) (*[]elkClubBody, error) {
	var club []elkClubBody
	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return &club, err
	}
	hits, _ := js.Get("hits").Get("hits").Array()
	for i, _ := range hits {
		var clubBody elkClubBody
		source := js.Get("hits").Get("hits").GetIndex(i).Get("_source")
		clubBody.Id, _ = source.Get("id").Int()
		clubBody.name, _ = source.Get("name").String()
		clubBody.icon, _ = source.Get("icon").String()
		clubBody.description, _ = source.Get("description").String()
		clubBody.membernum, _ = source.Get("membernum").String()
		clubBody.todayposts, _ = source.Get("todayposts").Int()
		clubBody.fup, _ = source.Get("fup").Int()
		club = append(club, clubBody)
	}
	return &club, nil
}

/**
构造向user_recommend_data 中存储的俱乐部数据
 */
func (u *User) buildRecommendClubData(club *[]elkClubBody) {
	//新增数据
	clubItems := *club
	for i, _ := range clubItems {
		clubItem := &clubItems[i]
		membernum, _ := strconv.Atoi(clubItem.membernum)
		isFollow := u.isFollow(clubItem.Id)
		var data elasticsearch.RecommendClubData
		data = elasticsearch.RecommendClubData{
			clubItem.Id,
			clubItem.name,
			clubItem.description,
			clubItem.icon,
			membernum,
			1,
			isFollow}
		recommendClubData = append(recommendClubData, data)
	}
}

/**
构造向user_recommend_data 中存储的相关用户数据
 */
func (u *User) buildRecommendUserData(user *[]elkUserBody, dataType int) {
	userItems := *user
	for i, _ := range userItems {
		userItem := &userItems[i]
		u.notRecommendUid = append(u.notRecommendUid, userItem.Id)
		var data elasticsearch.RecommendUserData
		data = elasticsearch.RecommendUserData{
			userItem.Id,
			userItem.nickname,
			userItem.avatar,
		0,
			dataType}
		recommendUserData = append(recommendUserData, data)
	}
}

func (u *User) insertData(recommendData string, recommendType int) error {
	//新增数据
	created := time.Now().Format("2006-01-02")
	var _id = strconv.Itoa(recommendType)+"_"+ created+"_" + strconv.Itoa(u.Uid)
	var data elasticsearch.UserRecommendListData
	data = elasticsearch.UserRecommendListData{_id,
			u.Uid,
recommendData,
		recommendType,
			created}
	ru := elasticsearch.NewUserRecommendData(u.nodes)
	err := ru.Create(&data)
	if err != nil {
		fmt.Println("insert user error")
		return err
	}
	return nil
}

func (u *User) isFollow(fid int) int {
	follows := strings.Split(u.myData.follow_clubs, strconv.Itoa(fid))
	if len(follows) > 1 {
		return 1
	}
	return 0
}
