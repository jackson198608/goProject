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
	} else if tag == "shopList" {
		qShopList(p)

	} else if tag == "gouminList" {
		qGouminList(p)

	} else if tag == "shopImage" {
		saveImage(p)

	} else if strings.Contains(tag, "shopImage") {

		shopDetailId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		//judge if this is the first image for shop
		isFirst := false
		tags := strings.Split(tag, "|")
		ImageWithNum := tags[0]
		ImageNum := ImageWithNum[len(ImageWithNum)-1:]
		ImageIntNum, err := strconv.Atoi(ImageNum)
		if err != nil {
			logger.Println("[error] parse int imageNume error")
			return
		}
		if ImageIntNum == 1 {
			isFirst = true
		}
		saveShopImagePath(p, int64(shopDetailId), isFirst)

		if ImageIntNum != 5 {
			logger.Println("[info] find next image")
			qShopImage(p, int64(shopDetailId), ImageIntNum)
		}

	} else if strings.Contains(tag, "shopCommentList") {
		shopDetailId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		qShopCommentList(p, int64(shopDetailId))
		saveShopCommentList(p, int64(shopDetailId))

	} else if strings.Contains(tag, "shopCommentImage") {
		shopCommentId, success := getDetailId(tag)
		if !success {
			logger.Println("[error] get detail id error", tag, p.GetRequest().Url)
			return
		}
		saveShopCommentImagePath(p, int64(shopCommentId))

	} else if tag == "shopFirstAreaFind" {
		saveFirstAreaFind(p)
	} else if strings.Contains(tag, "shopSecondBusinessList") {
		tags := strings.Split(tag, "|")
		firstBusinessIdStr := tags[1]
		firstBusinessId, err := strconv.Atoi(firstBusinessIdStr)
		if err != nil {
			logger.Println("[error]get first business id error ", tag)
			return
		}
		logger.Println("[info]process second business list ", tag, " ", firstBusinessId)
		saveSecondBusinessList(p, int64(firstBusinessId))
	} else if strings.Contains(tag, "shopSecondRegionList") {
		tags := strings.Split(tag, "|")
		firstRegionIdStr := tags[1]
		firstRegionId, err := strconv.Atoi(firstRegionIdStr)
		if err != nil {
			logger.Println("[error]get first region id error ", tag)
			return
		}

		logger.Println("[info]process second region list ", tag, " ", firstRegionId)
		saveSecondRegionList(p, int64(firstRegionId))
	} else if strings.Contains(tag, "shopSecondMetroList") {
		tags := strings.Split(tag, "|")
		firstMetroIdStr := tags[1]
		firstMetroId, err := strconv.Atoi(firstMetroIdStr)
		if err != nil {
			logger.Println("[error]get first metro id error ", tag)
			return
		}

		logger.Println("[info]process second metro list ", tag, " ", firstMetroId)
		saveSecondMetroList(p, int64(firstMetroId))
	} else if strings.Contains(tag, "shopUpdateOfSecondBusinessList") {
		tags := strings.Split(tag, "|")
		firstBusinessIdStr := tags[1]
		firstBusinessId, err := strconv.Atoi(firstBusinessIdStr)
		if err != nil {
			logger.Println("[error]get first business id error ", tag)
			return
		}
		logger.Println("[info]process update of second business", tag)
		qShopUpdateList(p, tag, int64(firstBusinessId), 0, 1, 1)
	} else if strings.Contains(tag, "shopUpdateSecondBusinessList") {
		tags := strings.Split(tag, "|")
		firstBusinessIdStr := tags[1]
		secondBusinessIdStr := tags[2]
		firstBusinessId, err := strconv.Atoi(firstBusinessIdStr)
		if err != nil {
			logger.Println("[error]get first business id error ", tag)
			return
		}
		secondBusinessId, err := strconv.Atoi(secondBusinessIdStr)
		if err != nil {
			logger.Println("[error]get first business id error ", tag)
			return
		}

		logger.Println("[info]process update of second business", tag)
		qShopUpdateList(p, tag, int64(firstBusinessId), int64(secondBusinessId), 2, 1)
	} else if strings.Contains(tag, "shopUpdateOfSecondRegionList") {
		tags := strings.Split(tag, "|")
		firstRegionIdStr := tags[1]
		firstRegionId, err := strconv.Atoi(firstRegionIdStr)
		if err != nil {
			logger.Println("[error]get first region id error ", tag)
			return
		}
		qShopUpdateList(p, tag, int64(firstRegionId), 0, 1, 2)

	} else if strings.Contains(tag, "shopUpdateSecondRegionList") {
		tags := strings.Split(tag, "|")
		firstRegionIdStr := tags[1]
		secondRegionIdStr := tags[2]
		firstRegionId, err := strconv.Atoi(firstRegionIdStr)
		if err != nil {
			logger.Println("[error]get first region id error ", tag)
			return
		}
		secondRegionId, err := strconv.Atoi(secondRegionIdStr)
		if err != nil {
			logger.Println("[error]get first region id error ", tag)
			return
		}

		qShopUpdateList(p, tag, int64(firstRegionId), int64(secondRegionId), 2, 2)

	} else if strings.Contains(tag, "shopUpdateOfSecondMetroList") {
		tags := strings.Split(tag, "|")
		firstMetroIdStr := tags[1]
		firstMetroId, err := strconv.Atoi(firstMetroIdStr)
		if err != nil {
			logger.Println("[error]get first metro id error ", tag)
			return
		}
		qShopUpdateList(p, tag, int64(firstMetroId), 0, 1, 3)

	} else if strings.Contains(tag, "shopUpdateSecondMetroList") {
		tags := strings.Split(tag, "|")
		firstMetroIdStr := tags[1]
		secondMetroIdStr := tags[2]
		firstMetroId, err := strconv.Atoi(firstMetroIdStr)
		if err != nil {
			logger.Println("[error]get first region id error ", tag)
			return
		}
		secondMetroId, err := strconv.Atoi(secondMetroIdStr)
		if err != nil {
			logger.Println("[error]get first region id error ", tag)
			return
		}

		qShopUpdateList(p, tag, int64(firstMetroId), int64(secondMetroId), 2, 3)
	} else if strings.Contains(tag, "shopUpdateDetail") {
	}

}

func (this *MyPageProcesser) Finish() {
}
