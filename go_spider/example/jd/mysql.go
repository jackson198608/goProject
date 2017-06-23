package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"strconv"
)

func checkShopExist(sourceUrl string) (int64, bool) {
	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		//fmt.Println("[error] connect db err")
		return 0, false
	}
	defer db.Close()

	rows, err := db.Query("select id from goods_source_info where source_url='" + sourceUrl + "'")
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
		logger.Println("[info] check sql find true", sourceUrl)
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

func findSkuId(
	shopDetailId int64,
) (int64,bool) {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0, false
	}
	defer db.Close()

	rows, err := db.Query("select sku_id from goods_source_info where id=" + strconv.FormatInt(shopDetailId, 10))
	if err != nil {
		logger.Println("[error] check sql prepare error: ", err)
		//fmt.Println("[error] check sql prepare error: ", err)
		return 0, false
	}

	for rows.Next() {
		var sku_id int64
		if err := rows.Scan(&sku_id); err != nil {
			logger.Println("[error] check sql get rows error ", err)
			return 0, false
		}
		logger.Println("[info] check sql find true", shopDetailId)
		return sku_id, true

	}
	if err := rows.Err(); err != nil {
		logger.Println("[error] check sql get rows error ", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return 0, false

	}
	return 0, false
}

func findCommentNum(
	shopDetailId int64,
) (int64,bool) {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0, false
	}
	defer db.Close()

	rows, err := db.Query("select comment_num from goods_source_info where id=" + strconv.FormatInt(shopDetailId, 10))
	if err != nil {
		logger.Println("[error] check sql prepare error: ", err)
		//fmt.Println("[error] check sql prepare error: ", err)
		return 0, false
	}

	for rows.Next() {
		var comment_num int64
		if err := rows.Scan(&comment_num); err != nil {
			logger.Println("[error] check sql get rows error ", err)
			return 0, false
		}
		logger.Println("[info] check sql find true", shopDetailId)
		return comment_num, true

	}
	if err := rows.Err(); err != nil {
		logger.Println("[error] check sql get rows error ", err)
		//fmt.Println("[error] check sql get rows error ", err)
		return 0, false

	}
	return 0, false
}

func updateShopDetail(
	shape string,
	age string,
	component string,
	componentPercent string,
	graininess string,
	shopDetailId int64,
) bool {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE goods_source_info SET shape=?,age=?,component=?, component_percent=?,graininess=? where id=?")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(shape, age, component, componentPercent, graininess, shopDetailId)
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

func updateGoodsPrice(
	price float64,
	shopDetailId int64,
) bool {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE goods_source_info SET price=? where id=?")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(price*100, shopDetailId)
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
func updateCommentNumAndScore(
	score float64,
	commentNum int64,
	shopDetailId int64,
) bool {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE goods_source_info SET score=?,comment_num=? where id=?")
	if err != nil {
		logger.Println("[error] update prepare error: ", err)
		return false
	}

	res, err := stmt.Exec(score*100, commentNum, shopDetailId)
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
	goodsName string,
	goodsNumber int,
	goodsSku string,
	brand string,
	category string,
	goodsPrice float64,
	salesVolume int,
	commentNum int,
	score float64,
	shape string,
	age string,
	component string,
	componentPercent string,
	taste string,
	grain string,
	graininess string,
	source int,
	sourceUrl string,
) int64 {

	id,isExist := checkShopExist(sourceUrl)
	if isExist {
		logger.Println("[info] find goods")
		return id
	} else {
		db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
		if err != nil {
			logger.Println("[error] connect db err")
			return 0
		}
		defer db.Close()

		stmt, err := db.Prepare("INSERT goods_source_info SET goods_name=?, sku_id=?,sku_name=?, brand=?,category=?,price=?,sales_volume=?,comment_num=?,score=?,shape=?,age=?,component=?,component_percent=?,taste=?,grain=?,graininess=?,source=?,created=?,source_url=?")
		if err != nil {
			logger.Println("[error] insert prepare error: ", err)
			return 0
		}

		tm := time.Unix(time.Now().Unix(), 0)

	    createTime := tm.Format("2006-01-02 03:04:05")

		res, err := stmt.Exec(goodsName, goodsNumber, goodsSku, brand, category, goodsPrice, salesVolume, commentNum, score, shape, age, component, componentPercent, taste, grain, graininess, source, createTime, sourceUrl)
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
}

func insertShopPhoto(
	shopId int64,
	shopImage string,
	imageType int64,
) int64 {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT shop_image SET goods_source_id=?,path=?,type=?")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(shopId, shopImage, imageType)
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
	goodsSourceId int64,
	content string,
	source int,
	created string,
) int64 {

	db, err := sql.Open("mysql", dbAuth+"@tcp("+dbDsn+")/"+dbName+"?charset=utf8mb4")
	if err != nil {
		logger.Println("[error] connect db err")
		return 0
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT goods_source_comment SET goods_source_id=?,content=?, source=?,created=?")
	if err != nil {
		logger.Println("[error] insert prepare error: ", err)
		return 0
	}

	res, err := stmt.Exec(goodsSourceId, content, source, created)
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
