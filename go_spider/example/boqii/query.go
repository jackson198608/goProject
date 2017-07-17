package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	// "regexp"
	"strconv"
	"strings"
	"math"
)

func qShopCateList(p *page.Page) {
	logger.Println("[info]find shop category list page: ", p.GetRequest().Url)
	query := p.GetHtmlParser()
	// query.Find(".channel_left_menu a").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	url, isExsit := s.Attr("href")
	// 	if isExsit {
	// 		logger.Println("[info]find next list page: ", url)
	// 		realUrlTag := "shopList"
	// 		req := newRequest(realUrlTag, url)
	// 		p.AddTargetRequestWithParams(req)
	// 	}
	// 	return true
	// })
	query.Find(".goodsCateList_li .goodsCateSub a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("href")
		if isExsit {
			logger.Println("[info]find next list page: ", url)
			realUrlTag := "shopList"
			req := newRequest(realUrlTag, url)
			p.AddTargetRequestWithParams(req)
		}
		return true
	})
}

func qShopList(p *page.Page) {
	logger.Println("[info]find shop list page: ", p.GetRequest().Url)
	query := p.GetHtmlParser()

	//find shop list
	query.Find(".product_list .product_list_container .product_name a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			_,isExist := checkShopExist(url)
				if !isExist {
				logger.Println("[info]find detail page: ", url)
				realUrlTag := "shopDetail"
				req := newRequest(realUrlTag, url)
				p.AddTargetRequestWithParams(req)
			}
		}
		return true
	})

	query.Find(".product_container .pagination a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		title, isExsit := s.Attr("title")
        if isExsit {
            if strings.Contains(title,"下一页"){
				url, isExsit := s.Attr("href")
				if isExsit {
					logger.Println("[info]find next list page: ", url)
					realUrlTag := "shopList"
					req := newRequest(realUrlTag, url)
					p.AddTargetRequestWithParams(req)
				}
			}
		}
		return true
	})
}

func qShopDetail(p *page.Page, shopDetailId int64) {
	logger.Println("[info]find shop detail page: ", p.GetRequest().Url)
	query := p.GetHtmlParser()
	//其他规格
	query.Find(".change_no a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("href")
		if isExsit {
			_,isExist := checkShopExist(url)
			if !isExist {
				logger.Println("[info]find other sku: ", url)
				realUrlTag := "shopDetail"
				req := newRequest(realUrlTag, url)
				p.AddTargetRequestWithParams(req)
			}
		}
		return true
	})

	//商品图片
	query.Find(".pro_big_img img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("src")
		if isExsit {
			logger.Println("[info]find goods first image : ", url)
			shopDetailIdStr := strconv.Itoa(int(shopDetailId))
			realUrlTag := "shopImage|" + shopDetailIdStr
			req := newRequest(realUrlTag, url)
			p.AddTargetRequestWithParams(req)
		}
		return true
	})

	//商品详情图片
	// query.Find(".mt40 .mt15 div img").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	url, isExsit := s.Attr("src")
	// 	if isExsit {
	// 		realUrl := url
	// 		logger.Println("[info]find goods detail iamge: ", realUrl)
	// 		shopDetailIdStr := strconv.Itoa(int(shopDetailId))
	// 		realUrlTag := "shopImage|" + shopDetailIdStr
	// 		req := newRequest(realUrlTag, realUrl)
	// 		p.AddTargetRequestWithParams(req)
	// 	}
	// 	return true
	// })

	//商品评论数
	commentNum := 0
	query.Find(".pro_tag_cont a em").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i==0 {
			commentNum,_ =strconv.Atoi(s.Text())
			logger.Println("[info]find goods comment num is :", commentNum)
		}
		return true
	})
	if commentNum==0 {
		logger.Println("[info]find goods comment num is :", 0)
		return 
	}
	maxPage := 49
	count := 13
	page := int(math.Ceil(float64(commentNum) / float64(count)))
	if page<49 {
		maxPage = page
	}
	logger.Println("[info]find goods comment page is :", maxPage)
	sourceUrl := p.GetRequest().Url
	urlArr := strings.Split(sourceUrl, "-")
	urlArr1 := strings.Split(urlArr[1], ".")
	id := urlArr1[0]

	if id=="" {
		logger.Println("[info]find goods id fail ", "")
		return
	}
	// 波奇评价最多可查看49页
	for i := 1; i <= maxPage; i++ {
		url := "http://shop.boqii.com/index.php?app=ajax&ctl=comment&act=commentList&id="+ id +"&cmtype=&action=comment&page="+ strconv.Itoa(i) +"&ordertype=1"	
		shopDetailIdStr := strconv.Itoa(int(shopDetailId))
		realUrlTag := "shopCommentList|" + shopDetailIdStr
		logger.Println("[info]find goods comment next page :", url)
		req := newRequest(realUrlTag, url)
		p.AddTargetRequestWithParams(req)
	}
}