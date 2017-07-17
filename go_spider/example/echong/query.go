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
	// query := p.GetHtmlParser()

	// query.Find(".dogType ul li h3 a").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	url, isExsit := s.Attr("href")
	// 	if isExsit {
	// 		logger.Println("[info]find next list page: ", url)
	// 		realUrlTag := "shopList"
	// 		req := newRequest(realUrlTag, url)
	// 		p.AddTargetRequestWithParams(req)
	// 	}
	// 	return true
	// })
	ids := [12]string {"54","53","55","48","49","5","6","50","3895","3886","2651","2652"}
	for i := 0; i < len(ids); i++ {
        
        for ii := 1; ii <= 46; ii++ {
        	var id string
        	if ii ==1 {	
        		id = ids[i]
        	}else{
        		id =  ids[i] + "b1f" + strconv.Itoa(ii)
        	}
        	url := "http://list.epet.com/"+ id +".html"     
	        realUrlTag := "shopList"
	        req := newRequest(realUrlTag, url)
	        p.AddTargetRequestWithParams(req)
        }
	}
}

func qShopList(p *page.Page) {
	query := p.GetHtmlParser()

	//find shop list
	query.Find(".list-box-con a.gd-photo").EachWithBreak(func(i int, s *goquery.Selection) bool {
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
}

func qNextShopList(p *page.Page) {
	query := p.GetHtmlParser()

	query.Find(".pages a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		title := s.Text()
        if strings.Contains(title,"下一页"){
			url, isExsit := s.Attr("href")
			if isExsit {
				logger.Println("[info]find next list page: ", url)
				realUrlTag := "shopList"
				req := newRequest(realUrlTag, url)
				p.AddTargetRequestWithParams(req)
			}
		}
		return true
	})
}

func qShopDetail(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()
	//其他规格
	query.Find(".norms-con a.norms-a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Find(".goods-select").Length() > 0 {

		} else{
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
        }
		return true
	})

	//商品图片
	query.Find(".goodslogo a.cloud-zoom").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("href")
		if isExsit {
			url = strings.Replace(url, " ", "", -1)
			url = strings.Replace(url, "\n", "", -1)
			logger.Println("[info]find goods first image : ", url)
			shopDetailIdStr := strconv.Itoa(int(shopDetailId))
			realUrlTag := "shopImage|" + shopDetailIdStr
			req := newRequest(realUrlTag, url)
			p.AddTargetRequestWithParams(req)
			return false
		}
		return true
	})

	//商品详情图片
	// query.Find(".gd_details div div img").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	url, isExsit := s.Attr("src0")
	// 	if isExsit {
	// 		realUrl := url
	// 		logger.Println("[info]find goods detail image : ", realUrl)
	// 		shopDetailIdStr := strconv.Itoa(int(shopDetailId))
	// 		realUrlTag := "shopImage|" + shopDetailIdStr
	// 		req := newRequest(realUrlTag, realUrl)
	// 		p.AddTargetRequestWithParams(req)
	// 		return false
	// 	}
	// 	return true
	// })

	//从爬取源URL获取id
	sourceUrl := p.GetRequest().Url
	urlArr := strings.Split(sourceUrl, "com/")
	urlArr1 := strings.Split(urlArr[1], ".html")
	id := urlArr1[0]

	if id=="" {
		logger.Println("[info]find goods id fail ", "")
	}

	// 商品参数
	getParamsUrl := "http://item.epet.com/goods.html?do=GetParamsAndAnnounce&gid="+id
	shopDetailIdStr := strconv.Itoa(int(shopDetailId))
	paramsUrlTag := "shopDetailParams|" + shopDetailIdStr
	logger.Println("[info]find goods params and announce page:", getParamsUrl)
	req := newRequest(paramsUrlTag, getParamsUrl)
	p.AddTargetRequestWithParams(req)

	// 商品评分
	getScoreUrl := "http://item.epet.com/goods.html?do=GetReplys&gid="+ id +"&app=review&page=1&is_img=0"
	scoreUrlTag := "shopDetailScore|" + shopDetailIdStr
	logger.Println("[info]find goods score:", getScoreUrl)
	scorereq := newRequest(scoreUrlTag, getScoreUrl)
	p.AddTargetRequestWithParams(scorereq)

	//商品评论数
	commentNum := 0
	query.Find(".ats-style .c300").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i==0 {
			text := s.Text()
			num := strings.Split(text,"(")
        	num1 := strings.Split(num[1],")")
        	commentNum,_ = strconv.Atoi(num1[0]);
		}
		return true
	})

	if commentNum==0 {
		query.Find(".had-buy .clearfix .fl a span").EachWithBreak(func(i int, s *goquery.Selection) bool {
			text := s.Text() //4.8分 (3344条评论）
			num := strings.Split(text,"(")
	    	num1 := strings.Split(num[1],"条")
	    	commentNum,_ = strconv.Atoi(num1[0]);
			return true
		})
	}

	if commentNum==0 {
		logger.Println("[info]find goods comment num is :", 0)
	}
	if commentNum>0 {
		maxPage := 99
		count := 15
		page := int(math.Ceil(float64(commentNum) / float64(count)))
		if page<99 {
			maxPage = page
		}
		shopDetailIdStr := strconv.Itoa(int(shopDetailId))
		realUrlTag := "shopCommentList|" + shopDetailIdStr
		for i := 1; i <= maxPage; i++ {
			url := "http://item.epet.com/goods.html?do=GetReplys&gid="+ id +"&app=review&page="+ strconv.Itoa(i) +"&is_img=0"
			logger.Println("[info]find goods comment next page :", url)
			req := newRequest(realUrlTag, url)
			p.AddTargetRequestWithParams(req)
		}
	}
}