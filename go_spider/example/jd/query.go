package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	// "regexp"
	"strconv"
	"strings"
)

func getNextPageUrl(url string, page string) string {
	paramsIndex := strings.IndexByte(url, '?')
	if paramsIndex == -1 {
		return url + page
	}
	baseUrl := url[0:paramsIndex]
	nextPageUrl := string(baseUrl) + page
	return nextPageUrl
}

func StripQueryString(inputUrl string) string {
	index := strings.Index(inputUrl, "jpg")
	realUrl := inputUrl[0 : index+3]
	return string(realUrl)
}

/*
	type 1 for Update of
	type 2 for Update
	cat 1 for business
	cat 2 for region
	cat 3 for metro
*/

func qShopCateList(p *page.Page) {
    category := [30]string {"幼犬狗粮","成犬狗粮","妙鲜包","罐头","肉类零食","磨牙洁齿","饼干奶酪","狗笼","狗窝","厕所/尿垫","食具水具","清洁消毒","玩具","沐浴露","护毛素","电推剪","美容工具","牵引用品","航空箱","外出包","服饰","拾便器","补钙","美毛","肠胃调理","营养膏","奶粉","驱虫药","皮肤药剂","口/耳/眼护理"}
	for i := 0; i < len(category); i++ {
		keyword := category[i]
		m := 1
	    for ii := 0; ii < 20; ii++ {
	        if ii%2 ==1 {
	            page := strconv.Itoa(ii)
	            if ii==1 {
	                m += 57
	            } else {
	                m += 56
	            }
	            
	            s := strconv.Itoa(m)
	            url := "https://search.jd.com/Search?keyword="+ keyword +"&enc=utf-8&qrst=1&rt=1&stop=1&vt=2&wq="+ keyword +"&stock=1&page="+ page +"&s="+ s +"&click=0"        
		        realUrlTag := "shopList"
		        req := newRequest(realUrlTag, url)
		        p.AddTargetRequestWithParams(req)
	        }
	    }
	    break 
	}
}

func qShopList(p *page.Page) {
	query := p.GetHtmlParser()

	//find shop list
	query.Find(".gl-i-wrap .p-img a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("href")
		if isExsit {
			realUrl := "http:" + url 
			logger.Println("[info]find detail page: ", realUrl)
			realUrlTag := "shopDetail"
			req := newRequest(realUrlTag, realUrl)
			p.AddTargetRequestWithParams(req)
		}
		return true
	})
}

func qShopDetail(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()
	//其他规格
	query.Find(".p-choose .dd .item").EachWithBreak(func(i int, s *goquery.Selection) bool {
		dataSku, isExsit := s.Attr("data-sku")
		if isExsit {
			sourceUrl := p.GetRequest().Url
			urlArr := strings.Split(sourceUrl, "com/")
			urlArr1 := strings.Split(urlArr[1], ".html")
			id := urlArr1[0]
			if !strings.Contains(id,dataSku) {
				url := "https://item.jd.com/"+ dataSku +".html"
				logger.Println("[info]find other sku: ", url)
				realUrlTag := "shopDetail"
				req := newRequest(realUrlTag, url)
				p.AddTargetRequestWithParams(req)
			}
		}
		return true
	})

	//商品图片
	query.Find(".main-img img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("data-origin")
		if isExsit {
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
	query.Find(".mt40 .mt15 div img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("href")
		if isExsit {
			realUrl := url
			logger.Println("[info]find goods detail image : ", realUrl)
			shopDetailIdStr := strconv.Itoa(int(shopDetailId))
			realUrlTag := "shopImage|" + shopDetailIdStr
			req := newRequest(realUrlTag, realUrl)
			p.AddTargetRequestWithParams(req)
			return false
		}
		return true
	})

	// //评论页
	// //商品评论数
	// commentNum := 0
	// query.Find(".pro_tag_cont a em").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	if i==0 {
	// 		commentNum,_ =strconv.Atoi(s.Text())
	// 		logger.Println("[info]find goods comment num is :", commentNum)
	// 	}
	// 	return true
	// })
	// if commentNum==0 {
	// 	logger.Println("[info]find goods comment num is :", 0)
	// }
	// maxPage := 49
	// count := 13
	// page := int(math.Ceil(float64(commentNum) / float64(count)))
	// if page<49 {
	// 	maxPage = page
	// }
	// logger.Println("[info]find goods comment page is :", maxPage)
	// sourceUrl := p.GetRequest().Url
	// urlArr := strings.Split(sourceUrl, "-")
	// urlArr1 := strings.Split(urlArr[1], ".")
	// id := urlArr1[0]

	// if id=="" {
	// 	logger.Println("[info]find goods id fail ", "")
	// }
	// // 波奇评价最多可查看49页
	// for i := 1; i <= maxPage; i++ {
	// 	url := "http://shop.boqii.com/index.php?app=ajax&ctl=comment&act=commentList&id="+ id +"&cmtype=&action=comment&page="+ strconv.Itoa(i) +"&ordertype=1"	
	// 	shopDetailIdStr := strconv.Itoa(int(shopDetailId))
	// 	realUrlTag := "shopCommentList|" + shopDetailIdStr
	// 	logger.Println("[info]find goods comment next page :", url)
	// 	req := newRequest(realUrlTag, url)
	// 	p.AddTargetRequestWithParams(req)
	// }

}

func qShopCommentList(p *page.Page, shopDetailId int64) {
	//get next page
	// query := p.GetHtmlParser()
	// url := p.GetRequest().Url
	// query.Find(".NextPage").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	logger.Println("[info]in the next page")
	// 	url, isExsit := s.Attr("href")
	// 	if isExsit {
	// 		realUrl := getNextPageUrl(*result, url)
	// 		logger.Println("[info]find next page: ", realUrl)
	// 		shopDetailIdStr := strconv.Itoa(int(shopDetailId))
	// 		realUrlTag := "shopCommentList" + "|" + shopDetailIdStr
	// 		req := newRequest(realUrlTag, realUrl)
	// 		p.AddTargetRequestWithParams(req)
	// 	}
	// 	return false
	// })
}