package mysql

import (
	"database/sql"

	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"testing"
)

func TestGetFansData(t *testing.T) {
	dbName := "new_dog123"
	// dbName := "backend"
	db, err := sql.Open("mysql", "dog123:dog123@tcp(192.168.86.193:3307)/"+dbName+"?charset=utf8")
	if err != nil {
		logger.Error("[error] connect db err")
	}
	// uid := 1
	uid := 68296
	// followuids := []int{68937, 187638, 68237}
	// Position := []string{"39.9", "118.9"}
	// fmt.Println(Position[0])
	// Pet := []int{60, 2017, 7}
	// species := "贵宾"
	// province := "北京"
	// followfids := []int{36, 10, 159}

	// fans := getRauthinfoByUid(uid, db)
	fans := getUsers(uid, db)
	// fans := getHotClubs(3, followfids, db)
	// fans := UnicodeIndex("山东省")
	// fans := getCity("34.60411", "119.2164")
	// fans := getPetClubByUid(species, followfids, db)
	// fans := getSpeciesnameBySpeciesid(species, db)
	// fans := getSameSpeciesPetUsers(uid, followuids, Pet, db)
	// fans := getPetInfoByUid(uid, db)
	// fans := NearbyUser(uid, followuids, Position, db)
	// for _, v := range fans {

	// 	fmt.Println(v.latitude)
	// }
	fmt.Println(fans)
	logger.Info(fans)
}

// func TestLoadById(t *testing.T) {
// 	dbName := "test_dz2"
// 	db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/"+dbName+"?charset=utf8")
// 	if err != nil {
// 		logger.Error("[error] connect db err")
// 	}
// 	// uid := 1
// 	id := 1
// 	event := LoadById(id, db)
// 	logger.Info(event.Uid, event.Infoid, event.TypeId, event.Created, event.Status)
// }
