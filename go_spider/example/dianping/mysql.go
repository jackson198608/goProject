package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func checkShopExist(shopName string, shopCity string) (int64, bool) {
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		//fmt.Println("[error] connect db err")
		return 0, false
	}
	defer db.Close()

	rows, err := db.Query("select id from shop where name='" + shopName + "' and city='" + shopCity + "'")
	if err != nil {
		logger.Println("[error] check sql prepare error: ", err)
		//fmt.Println("[error] check sql prepare error: ", err)
		return 0, false
	}

	for rows.Next() {
		var shopId int64
		if err := rows.Scan(&shopId); err != nil {
			logger.Println("[error] check sql get rows error ", err)
			return 0, false
		}
		logger.Println("[info] check sql find true", shopName, " ", shopId)
		//fmt.Println("[info] check sql find true", shopName, " ", shopId)
		return shopId, true

	}
	if err := rows.Err(); err != nil {
		logger.Println("[error] check sql get rows error ", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return 0, false

	}
	return 0, false

}

func checkBusinessExist(name string) (int64, bool) {
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		//fmt.Println("[error] connect db err")
		return 0, false
	}
	defer db.Close()

	rows, err := db.Query("select id from business where name='" + name + "' and city_id=" + CityIdStr + "")
	if err != nil {
		logger.Println("[error] check business sql prepare error: ", err)
		//fmt.Println("[error] check sql prepare error: ", err)
		return 0, false
	}

	for rows.Next() {
		var businessId int64
		if err := rows.Scan(&businessId); err != nil {
			logger.Println("[error] check business sql get rows error ", err)
			return 0, false
		}
		logger.Println("[info] check business sql find true ", name, " ", CityId)
		//fmt.Println("[info] check sql find true", shopName, " ", shopId)
		return businessId, true

	}
	if err := rows.Err(); err != nil {
		logger.Println("[error] check business sql get rows error ", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return 0, false
	}
	return 0, false

}

func checkRegionExist(name string) (int64, bool) {
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		//fmt.Println("[error] connect db err")
		return 0, false
	}
	defer db.Close()

	rows, err := db.Query("select id from region where name='" + name + "' and city_id=" + CityIdStr + "")
	if err != nil {
		logger.Println("[error] check region sql prepare error: ", err)
		//fmt.Println("[error] check sql prepare error: ", err)
		return 0, false
	}

	for rows.Next() {
		var regionId int64
		if err := rows.Scan(&regionId); err != nil {
			logger.Println("[error] check region sql get rows error ", err)
			return 0, false
		}
		logger.Println("[info] check region sql find true ", name, " ", CityId)
		//fmt.Println("[info] check sql find true", shopName, " ", shopId)
		return regionId, true

	}
	if err := rows.Err(); err != nil {
		logger.Println("[error] check region sql get rows error ", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return 0, false
	}
	return 0, false

}

func checkMetroExist(name string) (int64, bool) {
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		//fmt.Println("[error] connect db err")
		return 0, false
	}
	defer db.Close()

	rows, err := db.Query("select id from metro where name='" + name + "' and city_id=" + CityIdStr + "")
	if err != nil {
		logger.Println("[error] check metro sql prepare error: ", err)
		//fmt.Println("[error] check sql prepare error: ", err)
		return 0, false
	}

	for rows.Next() {
		var metroId int64
		if err := rows.Scan(&metroId); err != nil {
			logger.Println("[error] check metro sql get rows error ", err)
			return 0, false
		}
		logger.Println("[info] check metro sql find true ", name, " ", CityId)
		//fmt.Println("[info] check sql find true", shopName, " ", shopId)
		return metroId, true

	}
	if err := rows.Err(); err != nil {
		logger.Println("[error] check metro sql get rows error ", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return 0, false
	}
	return 0, false

}

func updateShopBusiness(
	shopName string,
	business int64,
) bool {
	shopId, isExist := checkShopExist(shopName, City)
	if !isExist {
		return false
	}

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("update shop set business = ? where id = ? ")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(business, shopId)
	if err != nil {
		logger.Println("[error] update excute error: ", err)
		return false
	}

	num, err := res.RowsAffected()
	if err != nil {
		logger.Println("[error] get insert id error: ", err, " num:", num)
		return false
	}

	return true
}

func updateShopImage(
	shopName string,
	image string,
) bool {
	shopId, isExist := checkShopExist(shopName, City)
	if !isExist {
		return false
	}

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("update shop set image = ? where id = ? ")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(image, shopId)
	if err != nil {
		logger.Println("[error] update excute error: ", err)
		return false
	}

	num, err := res.RowsAffected()
	if err != nil {
		logger.Println("[error] get insert id error: ", err, " num:", num)
		return false
	}

	return true
}

func updateShopBusinessWithSub(
	shopName string,
	business int64,
	business_sub int64,
) bool {
	shopId, isExist := checkShopExist(shopName, City)
	if !isExist {
		logger.Println("[error] the shop is not exist for update", shopName, " ", business, " ", business_sub)
		return false
	}

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("update shop set business = ? , business_sub = ? where id = ? ")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(business, business_sub, shopId)
	if err != nil {
		logger.Println("[error] update excute error: ", err)
		return false
	}

	num, err := res.RowsAffected()
	if err != nil {
		logger.Println("[error] get insert id error: ", err, " num:", num)
		return false
	}

	return true
}

func updateShopRegion(
	shopName string,
	region int64,
) bool {
	shopId, isExist := checkShopExist(shopName, City)
	if !isExist {
		return false
	}

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("update shop set region = ? where id = ? ")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(region, shopId)
	if err != nil {
		logger.Println("[error] update excute error: ", err)
		return false
	}

	num, err := res.RowsAffected()
	if err != nil {
		logger.Println("[error] get insert id error: ", err, " num:", num)
		return false
	}

	return true
}

func updateShopRegionWithSub(
	shopName string,
	region int64,
	region_sub int64,
) bool {
	shopId, isExist := checkShopExist(shopName, City)
	if !isExist {
		return false
	}

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("update shop set region = ? , region_sub = ? where id = ? ")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(region, region_sub, shopId)
	if err != nil {
		logger.Println("[error] update excute error: ", err)
		return false
	}

	num, err := res.RowsAffected()
	if err != nil {
		logger.Println("[error] get insert id error: ", err, " num:", num)
		return false
	}

	return true
}

func updateShopMetro(
	shopName string,
	metro int64,
) bool {
	shopId, isExist := checkShopExist(shopName, City)
	if !isExist {
		return false
	}

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("update shop set metro = ? where id = ? ")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(metro, shopId)
	if err != nil {
		logger.Println("[error] update excute error: ", err)
		return false
	}

	num, err := res.RowsAffected()
	if err != nil {
		logger.Println("[error] get insert id error: ", err, " num:", num)
		return false
	}

	return true
}

func updateShopMetroWithSub(
	shopName string,
	metro int64,
	metro_sub int64,
) bool {
	shopId, isExist := checkShopExist(shopName, City)
	if !isExist {
		return false
	}

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("update shop set metro = ? , metro_sub= ? where id = ? ")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(metro, metro_sub, shopId)
	if err != nil {
		logger.Println("[error] update excute error: ", err)
		return false
	}

	num, err := res.RowsAffected()
	if err != nil {
		logger.Println("[error] get insert id error: ", err, " num:", num)
		return false
	}

	return true
}

func insertShopDetail(
	city string,
	shopType int,
	shopName string,
	shopAddress string,
	shopPhone string,
	commentNum int,
	price int,
	star int,
	servicePoint int,
	envPoint int,
	weightPoint int,
	shopTime string,
	shopImage string,
	business int64,
	business_sub int64,
	region int64,
	region_sub int64,
	metro int64,
	metro_sub int64,
) int64 {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT shop SET city=?, type=?,name=?, address=?,phone=?,comment_num=?,price=?,star=?,service_point=?,env_point=?,weight_point=?,shop_time=?,image=?")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(city, shopType, shopName, shopAddress, shopPhone, commentNum, price, star, servicePoint, envPoint, weightPoint, shopTime, shopImage)
	if err != nil {
		logger.Println("[error] insert excute error: ", err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Println("[error] get insert id error: ", err)
		return 0
	}
	return id

}

func insertBusiness(
	name string,
	cityId int64,
	pid int64,
) int64 {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT business SET city_id=?,pid=?,name=?,lat=0,lng=0,sort=0")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(cityId, pid, name)
	if err != nil {
		logger.Println("[error] insert excute error: ", err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Println("[error] get insert id error: ", err)
		return 0
	}
	return id
}

func insertRegion(
	name string,
	cityId int64,
	pid int64,
) int64 {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT region SET city_id=?,pid=?,name=?,lat=0,lng=0,sort=0")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(cityId, pid, name)
	if err != nil {
		logger.Println("[error] insert excute error: ", err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Println("[error] get insert id error: ", err)
		return 0
	}
	return id
}

func insertMetro(
	name string,
	cityId int64,
	pid int64,
) int64 {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT metro SET city_id=?,pid=?,name=?,lat=0,lng=0,sort=0")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(cityId, pid, name)
	if err != nil {
		logger.Println("[error] insert excute error: ", err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Println("[error] get insert id error: ", err)
		return 0
	}
	return id
}

func insertShopPhoto(
	shopId int64,
	shopImage string,
) int64 {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT shop_photo SET shop_id=?,image=?")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(shopId, shopImage)
	if err != nil {
		logger.Println("[error] insert excute error: ", err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Println("[error] get insert id error: ", err)
		return 0
	}
	return id

}

func insertShopComment(
	shopId int64,
	content string,
	username string,
	avar string,
	price int,
	star int,
	servicePoint int,
	envPoint int,
	weightPoint int,
	created string,
) int64 {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT shop_comment SET shop_id=?,content=?, username=?,avar=?,price=?,star=?,service_point=?,env_point=?,weight_point=?,created=?")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(shopId, content, username, avar, price, star, servicePoint, envPoint, weightPoint, created)
	if err != nil {
		logger.Println("[error] insert excute error: ", err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Println("[error] get insert id error: ", err)
		return 0
	}
	return id

}

func insertCommentPhoto(
	commentId int64,
	commentImage string,
) int64 {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT comment_photo SET comment_id=?,image=?")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(commentId, commentImage)
	if err != nil {
		logger.Println("[error] insert excute error: ", err)
		return 0
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Println("[error] get insert id error: ", err)
		return 0
	}
	return id
}

func updateImageForFirstPhoto(
	shopDetailId int64,
	shopImage string,
) bool {
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("Update shop SET image=? where id=?")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return false
	}

	_, err = stmt.Exec(shopImage, shopDetailId)
	if err != nil {
		logger.Println("[error] insert excute error: ", err)
		return false
	}
	return true
}
