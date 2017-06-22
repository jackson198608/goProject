package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	// "regexp"
	"strconv"
	"strings"
	"math"
	"encoding/json"
)

type GoodsDesc struct {
    Date string `json:"date"`
    Content string `json:"content"`
}

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
	    for ii := 0; ii < 200; ii++ {
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
func qGoodsCommentList(p *page.Page, shopDetailId int64) {
	sku_id,_ := findSkuId(shopDetailId)

	//商品评论数
	commentNum,_ := findCommentNum(shopDetailId)

	if commentNum==0 {
		logger.Println("[info]find goods comment num is :", 0)
	}
	maxPage := 99
	count := 10
	page := int(math.Ceil(float64(commentNum) / float64(count)))
	if page<99 {
		maxPage = page
	}
	logger.Println("[info]find goods comment page is :", maxPage)

	// jd评价最多可查看99页
	for i := 1; i <= maxPage; i++ {
		url := "https://sclub.jd.com/comment/productPageComments.action?productId="+ strconv.FormatInt(sku_id,10) +"&score=0&sortType=5&page="+ strconv.Itoa(i) +"&pageSize=10&isShadowSku=0&fold=1" //&callback=jQuery5551091
		shopDetailIdStr := strconv.Itoa(int(shopDetailId))
		realUrlTag := "shopCommentList|" + shopDetailIdStr
		logger.Println("[info]find goods comment next page :", url)
		req := newRequest(realUrlTag, url)
		p.AddTargetRequestWithParams(req)
	}
}

func qGoodsDescImage(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()

	var s GoodsDesc
	jsonStr := ""
	query.Find("body").EachWithBreak(func(i int, s *goquery.Selection) bool {
		jsonStr = s.Text()
		return true
	})
    json.Unmarshal([]byte(jsonStr), &s)
    // logger.Println(jsonStr)
    // logger.Println("json ", s)
	
	// url := "https://sclub.jd.com/comment/productPageComments.action?productId="+ strconv.FormatInt(sku_id,10) +"&score=0&sortType=5&page="+ strconv.Itoa(i) +"&pageSize=10&isShadowSku=0&fold=1" //&callback=jQuery5551091
	// shopDetailIdStr := strconv.Itoa(int(shopDetailId))
	// realUrlTag := "shopImage|" + shopDetailIdStr
	// logger.Println("[info]find goods comment next page :", url)
	// req := newRequest(realUrlTag, url)
	// p.AddTargetRequestWithParams(req)
}

func qShopDetail(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()
	//获取当前skuid
	sku_id,_ := findSkuId(shopDetailId)

	shopDetailIdStr := strconv.Itoa(int(shopDetailId))

	//获取商品价格
	getPriceUrl := "https://p.3.cn/prices/mgets?&type=1&area=1_72_2799_0&pdtk=Tb9taS%252BIexnBRFKsj189v9oGHpVaVXq4WXvMG%252BdvtPyh92O%252BPoSi2ySSqJYKBQOrVtTgJp%252FGekyZ%250A5hrmhI%252FMOQ%253D%253D&pduid=1486720070332751456873&pdpin=&pdbp=0&skuIds=J_"+ strconv.FormatInt(sku_id,10) +"&ext=10000000&source=item-pc" //callback=jQuery9272433

	logger.Println("[info]find goods price: ", getPriceUrl)
	priceTag := "goodsPrice|" + shopDetailIdStr
	priceReq := newRequest(priceTag, getPriceUrl)
	p.AddTargetRequestWithParams(priceReq)

	//获取商品评论数、评分
	getCommentCountUrl := "https://club.jd.com/comment/productCommentSummaries.action?referenceIds="+ strconv.FormatInt(sku_id, 10) +"&_=1497592374339" //&callback=jQuery5551091
	logger.Println("[info]find comment num and score: ", getCommentCountUrl)
	CommentCountTag := "goodsCommentNumScore|" + shopDetailIdStr
	commentCountReq := newRequest(CommentCountTag, getCommentCountUrl)
	p.AddTargetRequestWithParams(commentCountReq)

	//其他规格
	query.Find(".p-choose .dd .item").EachWithBreak(func(i int, s *goquery.Selection) bool {
		dataSku, isExsit := s.Attr("data-sku")
		if isExsit {
			id := strconv.FormatInt(sku_id, 10)
			if !strings.Contains(id,dataSku) {
				url := "https://item.jd.com/"+ dataSku +".html"
				logger.Println("get goods sku url:", url)
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
	query.Find(".main-img img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("data-origin")
		if !isExsit {
			//全球购
			url, isExsit = s.Attr("src")
		}
		logger.Println("[info]find goods first image : ", "http:" + url)
		shopDetailIdStr := strconv.Itoa(int(shopDetailId))
		realUrlTag := "shopImage|" + shopDetailIdStr
		req := newRequest(realUrlTag, "http:" + url)
		p.AddTargetRequestWithParams(req)
		return true
	})

	//商品详情图片
	// getImageUrl := "https://cd.jd.com/description/channel?skuId="+ strconv.FormatInt(sku_id, 10) +"&mainSkuId="+ strconv.FormatInt(sku_id, 10) +"&cdn=2" //&callback=showdesc

	// logger.Println("[info]find goods detail image: ", getImageUrl)
	// imageTag := "goodsDescImage|" + shopDetailIdStr
	// imageReq := newRequest(imageTag, getImageUrl)
	// p.AddTargetRequestWithParams(imageReq)
}