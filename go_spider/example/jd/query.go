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
    category := [7]string {"幼犬狗粮","成犬狗粮","妙鲜包","狗罐头","狗狗沐浴露","狗狗护毛素","狗拾便器"}
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

	category_id := [22]string {"6994,6996,7006","6994,6996,7007","6994,6996,7008","6994,6998,7022","6994,6998,7017","6994,6998,7019","6994,6998,7018","6994,6998,7021","6994,6999","6994,7001,7039","6994,7001,7037","6994,7000,7028","6994,7000,7029","6994,7000,7029","6994,7000,7030","6994,6997,7011","6994,6997,7012","6994,6997,7013","6994,6997,7015","6994,6997,7016","6994,6997,11984","6994,7001,7037"}
	for i := 0; i < len(category_id); i++ {
		keyword := category_id[i]
	    for ii := 1; ii <= 200; ii++ {
	        s := strconv.Itoa(ii)
	 		url := "https://list.jd.com/list.html?cat="+ keyword +"&page="+ s +"&sort=sort_totalsales15_desc&trans=1&JL=6_0_0#J_main"
	        realUrlTag := "shopList"
	        req := newRequest(realUrlTag, url)
	        p.AddTargetRequestWithParams(req)
	    }
        
	}
}

func qShopList(p *page.Page) {
	query := p.GetHtmlParser()

	//find shop list
	query.Find(".gl-i-wrap .p-img a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("href")
		if isExsit {
			urls := strings.Split(url, "https://")
			realUrl := ""

			if len(urls)>1 {
				realUrl = url 
			} else {
				realUrl = "http:" + url 
			}
			_,isExist := checkShopExist(realUrl)
			if !isExist {
				logger.Println("[info]find detail page: ", realUrl)
				realUrlTag := "shopDetail"
				req := newRequest(realUrlTag, realUrl)
				p.AddTargetRequestWithParams(req)
			}
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

	// shopDetailIdStr := strconv.Itoa(int(shopDetailId))

	// //获取商品价格
	// getPriceUrl := "https://p.3.cn/prices/mgets?&type=1&area=1_72_2799_0&pdtk=Tb9taS%252BIexnBRFKsj189v9oGHpVaVXq4WXvMG%252BdvtPyh92O%252BPoSi2ySSqJYKBQOrVtTgJp%252FGekyZ%250A5hrmhI%252FMOQ%253D%253D&pduid=1486720070332751456873&pdpin=&pdbp=0&skuIds=J_"+ strconv.FormatInt(sku_id,10) +"&ext=10000000&source=item-pc" //callback=jQuery9272433

	// logger.Println("[info]find goods price: ", getPriceUrl)
	// priceTag := "goodsPrice|" + shopDetailIdStr
	// priceReq := newRequest(priceTag, getPriceUrl)
	// p.AddTargetRequestWithParams(priceReq)

	// //获取商品评论数、评分
	// getCommentCountUrl := "https://club.jd.com/comment/productCommentSummaries.action?referenceIds="+ strconv.FormatInt(sku_id, 10) +"&_=1497592374339" //&callback=jQuery5551091
	// logger.Println("[info]find comment num and score: ", getCommentCountUrl)
	// CommentCountTag := "goodsCommentNumScore|" + shopDetailIdStr
	// commentCountReq := newRequest(CommentCountTag, getCommentCountUrl)
	// p.AddTargetRequestWithParams(commentCountReq)

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

func qSkuPrice(p *page.Page) {
	logger.Println("[info] in query sku price")
	//获取规格价格
	// 获取价钱为0的商品id
	details := getIdsByPrice()
	for i := 0; i < len(details); i++ {
		id := details[i].id
		sku_id := details[i].sku_id
		getPriceUrl := "https://p.3.cn/prices/mgets?&type=1&area=1_72_2799_0&pdtk=Tb9taS%252BIexnBRFKsj189v9oGHpVaVXq4WXvMG%252BdvtPyh92O%252BPoSi2ySSqJYKBQOrVtTgJp%252FGekyZ%250A5hrmhI%252FMOQ%253D%253D&pduid=1486720070332751456873&pdpin=&pdbp=0&skuIds=J_"+ strconv.Itoa(sku_id) +"&ext=10000000&source=item-pc" //callback=jQuery9272433
		logger.Println("[info] get price by sku_id: ", sku_id)
		shopDetailIdStr := strconv.Itoa(int(id))
		referer := "https://item.jd.com/"+ strconv.Itoa(sku_id) +".html"
		priceTag := "goodsPrice|" + shopDetailIdStr
		priceReq := newJsonRequest(priceTag, getPriceUrl, referer)
		p.AddTargetRequestWithParams(priceReq)
	}
}

func qSkuCommentNum(p *page.Page) {
	logger.Println("[info] in query sku price")
	// 获取评论数的商品id
	details := getIdsByCommentNum()
	for i := 0; i < len(details); i++ {
		id := details[i].id
		sku_id := details[i].sku_id
		//获取商品评论数、评分
		logger.Println("[info] get comment num by sku_id: ", sku_id)
		getCommentCountUrl := "https://club.jd.com/comment/productCommentSummaries.action?referenceIds="+ strconv.Itoa(sku_id) +"&_=1497592374339" //&callback=jQuery5551091
		logger.Println("[info]find comment num and score: ", getCommentCountUrl)
		shopDetailIdStr := strconv.Itoa(int(id))
		CommentCountTag := "goodsCommentNumScore|" + shopDetailIdStr
		referer := "https://item.jd.com/"+ strconv.Itoa(sku_id) +".html"
		commentCountReq := newJsonRequest(CommentCountTag, getCommentCountUrl, referer)
		p.AddTargetRequestWithParams(commentCountReq)
	}
}
