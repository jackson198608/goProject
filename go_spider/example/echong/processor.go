package main

import (
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"strconv"
	"strings"
	//"time"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

func getDetailId(tag string) (int, bool) {
	tags := strings.Split(tag, "|")
	shopDetailIdStr := tags[1]
	shopDetailId, err := strconv.Atoi(shopDetailIdStr)
	if err != nil {
		logger.Println("[error]invaild shop id ", tag)
		return 0, false
	}
	return shopDetailId, true

}

// Parse html dom here and record the parse result that we want to Page.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
	//time.Sleep(1 * time.Second)
	if !p.IsSucc() {
		logger.Println("[Error]not 200: ", p.GetRequest().Url)
		return
	}

	tag := p.GetUrlTag()
	if tag == "shopDetail" {
		logger.Println("[info]shop detail page:", p.GetRequest().Url)
		//save shop into mysql
		shopDetailId, success := saveShopDetail(p)
		if success {
			logger.Println("[info]Get detail id", shopDetailId)
			//get all query
			qShopDetail(p, shopDetailId)
		}
	} else if tag == "DogCateList" {
		qShopCateList(p)

	}  else if tag == "shopList" {
		qShopList(p)

	} else if tag == "shopImage" {
		saveImage(p)

	} else if strings.Contains(tag, "shopImage") {

		shopDetailId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		
		saveShopImagePath(p, int64(shopDetailId))

	} else if strings.Contains(tag, "shopCommentList") {
		shopDetailId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		saveShopCommentList(p, int64(shopDetailId))

	}
}

func (this *MyPageProcesser) Finish() {
}
