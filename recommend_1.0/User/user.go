package user

import (
	"errors"
	"github.com/donnie4w/go-logger/logger"
	"github.com/go-xorm/xorm"
	"gopkg.in/mgo.v2"
	"strconv"
	"strings"
	"gouminGitlab/common/orm/elasticsearch"
	"github.com/olivere/elastic"
	"encoding/json"
)

type User struct {
	mysqlXorm       []*xorm.Engine
	mongoConn       []*mgo.Session
	Uid             int
	esConn          *elastic.Client
	province        string
	species         string
	address         string
	age             string
	myData          elasticsearch.UserData
	notRecommendUid []int
	notRecommendFid []int
}

var recommendClubData []elasticsearch.RecommendClubData
var recommendUserData []elasticsearch.RecommendUserData

func NewUser(mysqlXorm []*xorm.Engine, mongoConn []*mgo.Session, uid string, esConn *elastic.Client) *User {
	if (mysqlXorm == nil) || (mongoConn == nil) || (uid == "") || (esConn == nil) {
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
	u.esConn = esConn
	return u
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
	elkU,err := elasticsearch.NewUser(u.esConn)
	if err != nil {
		return err
	}
	user := elkU.SearchById(u.Uid)

	if user != nil {
		u.myData =  *user
		u.species = u.getSpecies()   //我的宠物品种
		u.age = u.getAge()           //我的宠物年龄
		u.province = u.getProvince() //我的地域
		u.address = u.getAddress()

		if u.myData.FollowUsers != "" {
			follow_users := strings.Split(u.myData.FollowUsers, ",")
			for f, _ := range follow_users {
				follow_uid, _ := strconv.Atoi(follow_users[f])
				u.notRecommendUid = append(u.notRecommendUid, follow_uid)
			}
		}
		if u.myData.FollowClubs != "" {
			follow_clubs := strings.Split(u.myData.FollowClubs, ",")
			for f, _ := range follow_clubs {
				follow_fid, _ := strconv.Atoi(follow_clubs[f])
				u.notRecommendFid = append(u.notRecommendFid, follow_fid)
			}
		}
		return nil
	}
	return errors.New("get my info error by uid is " + strconv.Itoa(u.Uid))
}

func (u *User) getAge() string {
	age := ""
	pets := u.myData.Pets
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
	pets := u.myData.Pets
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
	address := u.myData.Address
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
	address := u.myData.Address
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
	logger.Info("[user] recommend speciesNum is ", strconv.Itoa(speciesNum), " by uid ", u.Uid)
	logger.Info("[user] recommend ageNum is ", strconv.Itoa(ageNum), " by uid ", u.Uid)
	if ageNum+speciesNum < 8 {
		num := 8 - ageNum - speciesNum
		addressNum, _ := u.recommendUserByAddress(num)
		logger.Info("[user] recommend addressNum is ", strconv.Itoa(addressNum), " by uid ", u.Uid)
		if addressNum+speciesNum+ageNum < 8 {
			num := 8 - addressNum - speciesNum - ageNum
			nextSpeciesNum, _ := u.recommendUserBySpecies(1, num) //相同犬种
			logger.Info("[user] recommend nextSpeciesNum is ", strconv.Itoa(nextSpeciesNum), " by uid ", u.Uid)
		}
	}
	u.pushUserRecommend()
	return nil
}

//根据犬种推荐用户
func (u *User) recommendUserBySpecies(isFirst int, num int) (int, error) {
	if u.species != "" {
		speciesItems := strings.Split(u.species, ";")
		elkU,err := elasticsearch.NewUser(u.esConn)
		if err != nil {
			return 0, err
		}
		user := elkU.SearchByPets(speciesItems,u.notRecommendUid, num)
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
		elkU,err := elasticsearch.NewUser(u.esConn)
		if err != nil {
			return 0,err
		}
		user := elkU.SearchByAddress(u.address,u.notRecommendUid, num)
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
		elkU,err := elasticsearch.NewUser(u.esConn)
		if err != nil {
			return 0,err
		}
		user := elkU.SearchByPets(ageItems,u.notRecommendUid, 3)
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
	ur := elasticsearch.NewUserRecommendData(u.esConn)
	err := ur.InsertData(u.Uid,dataStr,2)
	return err
}

// -----   推荐俱乐部  ------

//获取已关注的俱乐部
func (u *User) followClubs() (int, error) {
	c := elasticsearch.NewClub(u.esConn)
	club := c.SearchByFollowIds(u.notRecommendFid, 6)
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
	c := elasticsearch.NewClub(u.esConn)
	club := c.SearchByKeyword(speciesItems, u.notRecommendFid, 6)
	if club != nil {
		u.buildRecommendClubData(club)
		return len(*club), nil
	}
	return 0, nil
}

func (u *User) recommendClubByFid(fid int) (int, error) {
	c := elasticsearch.NewClub(u.esConn)
	club := c.SearchById(fid)
	if club != nil {
		u.buildRecommendClubData(club)
		return len(*club), nil
	}
	return 0, nil
}

func (u *User) recommendClubByAddress() (int, error) {
	if u.province != "" {
		var keyword []string
		keyword = append(keyword,u.province,)
		c := elasticsearch.NewClub(u.esConn)
		club := c.SearchByKeyword(keyword,u.notRecommendFid,6)
		if club != nil {
			u.buildRecommendClubData(club)
			return len(*club), nil
		}
	}
	return 0, nil
}

func (u *User) recommendClubByFup(fup int) (int, error) {
	c := elasticsearch.NewClub(u.esConn)
	club := c.SearchByFup(fup,u.notRecommendFid,4)
	if club != nil {
		u.buildRecommendClubData(club)
		return len(*club), nil
	}
	return 0, nil
}

func (u *User) pushClubRecommend() error {
	jsonData, _ := json.Marshal(recommendClubData)
	dataStr := string(jsonData)
	ur := elasticsearch.NewUserRecommendData(u.esConn)
	err := ur.InsertData(u.Uid,dataStr,1)
	//err := u.insertData(dataStr, 1)
	if err != nil {
		return err
	}
	return nil
}

/**
构造向user_recommend_data 中存储的俱乐部数据
 */
func (u *User) buildRecommendClubData(club *[]elasticsearch.ClubData) {
	//新增数据
	clubItems := *club
	for i, _ := range clubItems {
		clubItem := &clubItems[i]
		membernum, _ := strconv.Atoi(clubItem.Membernum)
		isFollow := u.isFollow(clubItem.Id)
		var data elasticsearch.RecommendClubData
		data = elasticsearch.RecommendClubData{
			clubItem.Id,
			clubItem.Name,
			clubItem.Description,
			clubItem.Icon,
			membernum,
			1,
			isFollow}
		recommendClubData = append(recommendClubData, data)
	}
}

/**
构造向user_recommend_data 中存储的相关用户数据
 */
func (u *User) buildRecommendUserData(user *[]elasticsearch.UserData, dataType int) {
	userItems := *user
	for i, _ := range userItems {
		userItem := &userItems[i]
		//uidInt,_ := strconv.Atoi(userItem.Id)
		u.notRecommendUid = append(u.notRecommendUid, userItem.Id)
		var data elasticsearch.RecommendUserData
		data = elasticsearch.RecommendUserData{
			userItem.Id,
			userItem.Nickname,
			userItem.Avatar,
		0,
			dataType}
		recommendUserData = append(recommendUserData, data)
	}
}

func (u *User) isFollow(fid int) int {
	follows := strings.Split(u.myData.FollowClubs, strconv.Itoa(fid))
	if len(follows) > 1 {
		return 1
	}
	return 0
}
