package main

import (
	// "crypto/md5"
	// "encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"os"
	"path"
	"strconv"
	"strings"
	"encoding/json"
)

type GoodsPrice struct {
    P string `json:"p"`
}

type CommentsCount struct {
    CommentCount int64
    GoodRate   float64
}

type CommentsCountslice struct {
    CommentsCount []CommentsCount
}

type Comments struct {
    Content string `json:"content"`
    CreationTime   string `json:"creationTime"`
	ReferenceId string `json:"referenceId"`
}

type Commentsslice struct {
    Comments []Comments `json:"comments"`
}

func checkCanSave(p *page.Page) bool {
	return true
}

func saveImage(p *page.Page) bool {
	//judge if status is 200
	if !checkCanSave(p) {
		return false
	}

	url := p.GetRequest().Url
	//get fullpath
	abPath := getPathFromUrl(url)
	fullPath := saveDir + abPath
	fullDirPath := path.Dir(fullPath)
	err := os.MkdirAll(fullDirPath, 0664)
	if err != nil {
		logger.Println("[error]create dir error:", err, " ", fullDirPath, " ", url)
		return false
	}

	//save file
	result, err := os.Create(fullPath)
	if err != nil {
		logger.Println("[error]create file error:", err, " ", fullPath, " ", url)
		return false
	}
	logger.Println("[info] save image:", url)
	logger.Println("[info] save in:", fullPath)
	logger.Println("[info] save len:", len(p.GetBodyStr()))
	result.WriteString(p.GetBodyStr())
	result.Close()

	return true
}

func getGoodsName(query *goquery.Document, goodsName *string) {
	query.Find(".sku-name").EachWithBreak(func(i int, s *goquery.Selection) bool {
		*goodsName = s.Text()
		*goodsName = strings.Replace(*goodsName, " ", "", -1)
		*goodsName = strings.Replace(*goodsName, "\n", "", -1)
		return true
	})
}

func getFirstCategory(query *goquery.Document, category *string) {
	categoryStr := ""
	query.Find(".crumb .item a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 1 {
			categoryStr = s.Text()	
		}
		return true
	})
	*category = categoryStr
}

func getSecondCategory(query *goquery.Document, category *string) {
	categoryStr := ""
	query.Find(".crumb .item a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 2 {
			categoryStr = s.Text()
		}
		return true
	})
	*category = categoryStr
}

func getThirdCategory(query *goquery.Document, category *string) {
	categoryStr := ""
	query.Find(".crumb .item a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 3 {
			categoryStr = s.Text()
		}
		return true
	})
	*category = categoryStr
}

//调用接口获取到数据
func getGoodsPrice(query *goquery.Document, goodsPrice *float64) {
	var price GoodsPrice
	jsonStr := ""
	query.Find("body").EachWithBreak(func(i int, s *goquery.Selection) bool {
		jsonStr = s.Text()
		return true
	})
	jsonStr = strings.Replace(jsonStr, "[", "", -1)
	jsonStr = strings.Replace(jsonStr, "]", "", -1)
	json.Unmarshal([]byte(jsonStr), &price)
    priceF,_ := strconv.ParseFloat(price.P,64)
	*goodsPrice = priceF
}

func getGoodsSku(query *goquery.Document, goodsSku *string) {
	sku := ""
	query.Find(".p-choose .dd .item.selected").EachWithBreak(func(i int, s *goquery.Selection) bool {
		skuStr, _ := s.Attr("data-value")
		sku += skuStr + " "
		return true
	})
	*goodsSku = sku
}

func getShape(query *goquery.Document, shape *string) {
	shapeStr := ""
	query.Find(".p-parameter .p-parameter-list li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := s.Text()
        text1 := strings.Split(text, "：")
        if strings.Contains(text1[0],"犬型") {
            shapeStr = text1[1]
        }
        return true
	})
	*shape = shapeStr
}

func getAge(query *goquery.Document, age *string) {
	ageStr := ""
	query.Find(".p-parameter ul li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := s.Text()
        text1 := strings.Split(text, "：")
        if strings.Contains(text1[0],"适用犬龄") {
            ageStr = text1[1]
        }
        return true
	})
	*age = ageStr
}

func getTaste(query *goquery.Document, taste *string) {
	tasteStr := ""
	query.Find(".p-parameter ul li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := s.Text()
        text1 := strings.Split(text, "：")
        if strings.Contains(text1[0],"口味") {
            tasteStr = text1[1]
        }
        return true
	})
	*taste = tasteStr
}

func getComponent(query *goquery.Document, component *string) {
	componentStr := ""
	query.Find(".p-parameter ul li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := s.Text()
        text1 := strings.Split(text, "：")
        if strings.Contains(text1[0],"功效") {
            componentStr = text1[1]
        }
        return true
	})
	*component = componentStr
}

func getBrand(query *goquery.Document, brand *string) {
	query.Find("#parameter-brand li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		*brand,_ = s.Attr("title")
		return true
	})
}

func getGoodsNum(query *goquery.Document, goodsNum *int) {
	goodsNumStr := ""
	query.Find(".p-parameter ul li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := s.Text()
        text1 := strings.Split(text, "：")
        if strings.Contains(text1[0],"商品编号") {
            goodsNumStr = text1[1]
        }
        return true
	})
	*goodsNum,_ = strconv.Atoi(goodsNumStr)
}

//调用接口获取数据
func getCommentNumScore(query *goquery.Document, commentNum *int64, score *float64) {
	jsonStr := ""
	query.Find("body").EachWithBreak(func(i int, s *goquery.Selection) bool {
		jsonStr = s.Text()
		return true
	})

    var s CommentsCountslice
    json.Unmarshal([]byte(jsonStr), &s)
    *commentNum = s.CommentsCount[0].CommentCount
    *score = s.CommentsCount[0].GoodRate
}

func saveShopDetail(p *page.Page) (int64, bool) {
	sourceUrl := p.GetRequest().Url
	query := p.GetHtmlParser()
	// 商品名称
	var goodsName *string = new(string)
	getGoodsName(query, goodsName)
	logger.Println("[info] goodsName: ", *goodsName)

	// 分类
	var firstCategory *string = new(string)
	getFirstCategory(query, firstCategory)
	logger.Println("[info] goods first category: ", *firstCategory)

	var secondCategory *string = new(string)
	getSecondCategory(query, secondCategory)
	logger.Println("[info] goods second category: ", *secondCategory)

	var thirdCategory *string = new(string)
	getThirdCategory(query, thirdCategory)
	logger.Println("[info] goods third category: ", *thirdCategory)

	// 售价
	var goodsPrice *float64 = new(float64)

 	// 规格
	var goodsSku *string = new(string)
	getGoodsSku(query, goodsSku)
	logger.Println("[info] goods sku: ", *goodsSku)

	// 品牌
	var brand *string = new(string)
	getBrand(query, brand)
	logger.Println("[info] goods brand: ", *brand)

	// 商品编号
	var goodsNumber *int = new(int)
	getGoodsNum(query, goodsNumber)
	logger.Println("[info] goodsNumber: ", *goodsNumber)

	// 适用犬型
	var shape *string = new(string)
	getShape(query, shape)
	logger.Println("[info] goods shape: ", *shape)

	// 适用年龄
	var age *string = new(string)
	getAge(query, age)
	logger.Println("[info] use age: ", *age)

 	// 成分
	var goodsComponent *string = new(string)
	getComponent(query, goodsComponent)
	logger.Println("[info] goods component: ", *goodsComponent)

	// 销量
	var salesVolume *int = new(int)

	// 评分
	var score *float64 = new(float64)

	// 评论数
	var commentNum *int = new(int)

	// 成分含量
	var componentPercent *string = new(string)

	// 口味
	var taste *string = new(string)
	getComponent(query, taste)
	logger.Println("[info] goods taste: ", *taste)

	// 谷物成分
	var grain *string = new(string)

	// 颗粒度
	var graininess *string = new(string)
	
	if *goodsName!="" {
		//insert record
		shopDetailId := insertShopDetail(
			*goodsName,
			*goodsNumber,
			*goodsSku,
			*brand,
			*firstCategory +" "+ *secondCategory +" "+ *thirdCategory,
			*goodsPrice,
			*salesVolume,
			*commentNum,
			*score,
			*shape,
			*age,
			*goodsComponent,
			*componentPercent,
			*taste,
			*grain,
			*graininess,1,sourceUrl)
		return shopDetailId, true
	}else{
		logger.Println("[info]again find goods detail: ", sourceUrl)
		realUrlTag := "shopDetail"
		req := newRequest(realUrlTag, sourceUrl)
		p.AddTargetRequestWithParams(req)
	}
	return 0,false
}

func getPath(md5Str string) string {
	abPath := make([]byte, 48, 48)
	j := 0
	p := 0
	for i := 0; i <= 16; i++ {
		if p >= 50 || j >= 32 {
			break
		}
		abPath[p] = byte('/')
		abPath[p+1] = md5Str[j]
		abPath[p+2] = md5Str[j+1]
		p = p + 3
		j = j + 2
	}
	return string(abPath)
}

func getImageFromPage(query *goquery.Document, shopImage *string) {
	query.Find(".pro_big_img img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		shopUrl, isExist := s.Attr("src")
		if !isExist {
			return false
		}
		*shopImage = shopUrl
		return true
	})
}

func saveShopImagePath(p *page.Page, shopDetailId int64) (int64, bool) {
	var shopImage *string = new(string)
	*shopImage = p.GetRequest().Url

	logger.Println("[info]Image url:", *shopImage)
	imageRealPath := getPathFromUrl(*shopImage)
	logger.Println("[info]Image path:", imageRealPath)

	//save image into shop_image table
	insertShopPhoto(shopDetailId, imageRealPath,0)

	//生成对图片的抓取任务
	req := newImageRequest("shopImage", *shopImage)
	p.AddTargetRequestWithParams(req)
	return 0, true
}

func getPathFromUrl(url string) string {
	//find first ? position
	var i int = 0
	len := len(url)
	for i = 0; i < len; i++ {
		if url[i] == '?' {
			break
		}
	}
	if i == (len - 1) {
		//have no ?
		return url
	} else {
		path := url[7:i]
		return string(path)
	}
}

func saveShopCommentList(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()

	sku_id,_ := findSkuId(shopDetailId)

	var s Commentsslice
	jsonStr := ""
	query.Find("body").EachWithBreak(func(i int, s *goquery.Selection) bool {
		jsonStr = s.Text()
		return true
	})
    json.Unmarshal([]byte(jsonStr), &s)
    
    for i := 0; i < len(s.Comments); i++ {
    	logger.Println("[info] goods comment content: ",i, s.Comments[i].Content)
    	logger.Println("[info] goods comment CreationTime: ",i, s.Comments[i].CreationTime)
    	comment_sku_id,_ := strconv.ParseInt(s.Comments[i].ReferenceId, 10, 64)
    	if sku_id == comment_sku_id {
	    	insertShopComment(
	    		sku_id,
	    		s.Comments[i].Content,
	    		1,
	    		s.Comments[i].CreationTime,
	    		)
	    }
    }
}

func saveGoodsPrice(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()

	// 售价
	var goodsPrice *float64 = new(float64)
	getGoodsPrice(query, goodsPrice)
	logger.Println("[info] goods price: ", *goodsPrice)
	if *goodsPrice!=0 {
		updateGoodsPrice(*goodsPrice, shopDetailId)
	}
}

func saveGoodsCommentNumAndScore(p *page.Page, shopDetailId int64) (bool){
	query := p.GetHtmlParser()

	// 评分
	var score *float64 = new(float64)

	// 评论数
	var commentNum *int64 = new(int64)

	getCommentNumScore(query, commentNum, score)

	logger.Println("[info] score: ", *score)
	logger.Println("[info] common Num: ", *commentNum)

	status := updateCommentNumAndScore(*score, *commentNum, shopDetailId)
	return status
}
