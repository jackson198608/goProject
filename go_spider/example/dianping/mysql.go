package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

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
) int64 {

	// db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/travel?charset=utf8mb4")
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(127.0.0.1:3306)/activitydb?charset=utf8mb4")
	//db, err := sql.Open("mysql", "dog123:dog123@tcp(192.168.5.199:3306)/shop?charset=utf8mb4")
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

func insertShopPhoto(
	shopId int64,
	shopImage string,
) int64 {
	// db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/travel?charset=utf8mb4")
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(127.0.0.1:3306)/activitydb?charset=utf8mb4")
	//db, err := sql.Open("mysql", "dog123:dog123@tcp(192.168.5.199:3306)/shop?charset=utf8mb4")
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
	// db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/travel?charset=utf8mb4")
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(127.0.0.1:3306)/activitydb?charset=utf8mb4")
	//	db, err := sql.Open("mysql", "dog123:dog123@tcp(192.168.5.199:3306)/shop?charset=utf8mb4")
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
	// db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/travel?charset=utf8mb4")
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(127.0.0.1:3306)/activitydb?charset=utf8mb4")
	//db, err := sql.Open("mysql", "dog123:dog123@tcp(192.168.5.199:3306)/shop?charset=utf8mb4")
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
	// db, err := sql.Open("mysql", "root:goumintech@tcp(192.168.86.72:3309)/travel?charset=utf8mb4")
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(127.0.0.1:3306)/activitydb?charset=utf8mb4")
	//db, err := sql.Open("mysql", "dog123:dog123@tcp(192.168.5.199:3306)/shop?charset=utf8mb4")
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
