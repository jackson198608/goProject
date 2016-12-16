package main

import (
	"testing"
)

func TestInsertShopDetail(t *testing.T) {
	city := "北京"
	shopType := 1
	shopName := "中国店"
	shopAddress := "洛杉矶"
	shopPhone := "18210091845"
	commentNum := 100
	price := 50000
	star := 4
	servicePoint := 2
	envPoint := 3
	weightPoint := 5
	shopTime := "9:00-21:00 周一到周五"
	shopImage := "http://p0.meituan.net/dpmerchantalbum/1505c146620a5f5e24efb948aba2fe2882525.jpg%40240w_180h_1e_1c_1l%7Cwatermark%3D1%26%26r%3D1%26p%3D9%26x%3D5%26y%3D5%26relative%3D1"
	insertShopDetail(city, shopType, shopName, shopAddress, shopPhone, commentNum, price, star, servicePoint, envPoint, weightPoint, shopTime, shopImage)
}

func TestInsertShopPhoto(t *testing.T) {
	shopId := 684
	shopImage := "http://p1.meituan.net/dpmerchantalbum/bb0d36ce78c29da743f03b8362478508200488.jpg%40700w_700h_0e_1l%7Cwatermark%3D1%26%26r%3D1%26p%3D9%26x%3D5%26y%3D5%26relative%3D1"
	insertShopPhoto(shopId, shopImage)
}

func TestInsertShopComment(t *testing.T) {
	shopId := 684
	content := "导航到很方便，罗马家园对面底商最里面，好停车，有收费不贵，大悦城西门附近。因为都是提前预约，到了就做美容，给滴耳朵药，洗澡，去耳毛，挤肛门线，一看就很专业。狗狗寄养和暂时看护都是一箱子一箱子的，很和谐卡通。狗狗游泳不错，促进减肥，骨骼发育，按摩作用，促进胃肠运动，泳池每天换水，主人可以一块和爱宠游泳哦，人直接穿防水衣很方便是不是，价格也不贵哦。新店，服务没得说，很体贴，环境好，有休息区，重点是很干净，工作人员穿工服，很注重卫生，一进门没有宠物味儿，不论是接待区还是洗护区都很干净，大大的赞👍望一直保持下去。一进门有个高级美容工作室，落地窗，直接观看到宠物美毛，和狗狗近距离接触，让爱宠不孤单哈哈。推荐狗狗和猫猫主人来此给狗狗贴心护理。办卡更划算哦。狗狗刚刚剃全毛，有点丑不要介意图片哈"
	username := "小调味"
	avar := "http://i2.s1.dpfile.com/pc/ebd4a8142fc26b29f3fd06991415d5ef(48c48)/thumb.jpg"
	price := 60
	star := 5
	servicePoint := 4
	envPoint := 4
	weightPoint := 4
	created := "2016-11-10"
	insertShopComment(shopId, content, username, avar, price, star, servicePoint, envPoint, weightPoint, created)
}

func TestInsertCommentPhoto(t *testing.T) {
	commentId := 684
	commentImage := "http://qcloud.dpfile.com/pc/pNbzwsrQ1578BWsIuRV9ridwkCZOg4csxT1awQ7nJ_jMyPRMleaDSP-6WmQuJiv6TYGVDmosZWTLal1WbWRW3A.jpg"
	insertCommentPhoto(commentId, commentImage)
}
