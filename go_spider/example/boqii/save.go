package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"os"
	"path"
	"strconv"
	"strings"
)

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
	query.Find(".shop_name").EachWithBreak(func(i int, s *goquery.Selection) bool {
		*goodsName = s.Text()
		*goodsName = strings.Replace(*goodsName, " ", "", -1)
		*goodsName = strings.Replace(*goodsName, "\t", "", -1)
		*goodsName = strings.Replace(*goodsName, "\n", "", -1)
		return true
	})
}

func getFirstCategory(query *goquery.Document, category *string) {
	categoryStr := "" 
	query.Find(".crumbs_content .crumbs a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 1 {
			categoryStr = s.Text()
		}
		return true
	})
	*category = categoryStr
}

func getSecondCategory(query *goquery.Document, category *string) {
	categoryStr := "" 
	query.Find(".crumbs_content .crumbs a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 2 {
			categoryStr = s.Text()
		}
		return true
	})
	*category = categoryStr
}

func getThirdCategory(query *goquery.Document, category *string) {
	categoryStr := "" 
	query.Find(".crumbs_content .crumbs a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 3 {
			categoryStr = s.Text()
		}
		return true
	})
	*category = categoryStr
}

func getGoodsPrice(query *goquery.Document, goodsPrice *float64) {
	query.Find(".bq_price #yhhcast").EachWithBreak(func(i int, s *goquery.Selection) bool {
		price,_ := s.Attr("value")
		priceF,_ := strconv.ParseFloat(price,64)
		*goodsPrice = priceF
		return true
	})
}

func getGoodsSku(query *goquery.Document, goodsSku *string) {
	query.Find(".change.current").EachWithBreak(func(i int, s *goquery.Selection) bool {
		*goodsSku = s.Text()
		return true
	})
}

func getShape(query *goquery.Document, shape *string) {
	shapeStr := ""
	query.Find(".left p span.cfe7247").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "适用犬类:") {
            shapeStr = s.Siblings().Text()
		}
		return true
	})
	*shape = shapeStr
}

func getAge(query *goquery.Document, age *string) {
	ageStr := ""
	query.Find(".left p span.cfe7247").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "适用年龄:") {
            ageStr = s.Siblings().Text()
		}
		return true
	})
	*age = ageStr
}

func getComponent(query *goquery.Document, component *string) {
	componentStr := ""
	query.Find(".left p span.cfe7247").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "主要成分:") {
            componentStr = s.Siblings().Text()
		}
		return true
	})
	*component = componentStr
}

func getBrand(query *goquery.Document, brand *string) {
	query.Find(".brand dl dd").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i==0 {
			*brand = s.Text()
		}
		return true
	})
}

func getGoodsNum(query *goquery.Document, goodsNum *int) {
	query.Find(".brand dl dd").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i==1 {
			num,_:=strconv.Atoi(s.Text())
			*goodsNum = num
		}
		return true
	})
}

func getSalesVolume(query *goquery.Document, salesVolume *int) {
	query.Find(".brand dl dd").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i==2 {
			numStr := strings.Split(s.Text(), "件")
            num,_:=strconv.Atoi(numStr[0])
			*salesVolume = num
		}
		return true
	})
}

func getScore(query *goquery.Document, score *float64) {
	query.Find(".pl_l .pl_score span").EachWithBreak(func(i int, s *goquery.Selection) bool {
		scoreStr := s.Text()
		scoreF,_ := strconv.ParseFloat(scoreStr, 64)
		*score = scoreF
		return true
	})
}

func getCommentNum(query *goquery.Document, commentNum *int) {
	query.Find(".pro_tag_cont a em").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i==0 {
			num,_:=strconv.Atoi(s.Text())
			*commentNum = num
		}
		return true
	})
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
	getGoodsPrice(query, goodsPrice)
	logger.Println("[info] goods price: ", *goodsPrice)

 	// 规格
	var goodsSku *string = new(string)
	getGoodsSku(query, goodsSku)
	logger.Println("[info] goods sku: ", *goodsSku)

	// 适用犬型
	var shape *string = new(string)
	getShape(query, shape)
	logger.Println("[info] use shape: ", *shape)

	// 适用年龄
	var age *string = new(string)
	getAge(query, age)
	logger.Println("[info] use age: ", *age)

 	// 成分
	var goodsComponent *string = new(string)
	getComponent(query, goodsComponent)
	logger.Println("[info] goods component: ", *goodsComponent)

	// 品牌
	var brand *string = new(string)
	getBrand(query, brand)
	logger.Println("[info] goods brand: ", *brand)

	// 商品编号
	var goodsNumber *int = new(int)
	getGoodsNum(query, goodsNumber)
	logger.Println("[info] goodsNumber: ", *goodsNumber)

	// 销量
	var salesVolume *int = new(int)
	getSalesVolume(query, salesVolume)
	logger.Println("[info] salesVolume: ", *salesVolume)

	// 评分
	var score *float64 = new(float64)
	getScore(query, score)
	logger.Println("[info] score: ", *score)

	// 评论数
	var commonNum *int = new(int)
	getCommentNum(query, commonNum)
	logger.Println("[info] common Num: ", *commonNum)

	// 成分含量
	var componentPercent *string = new(string)

	// 口味
	var taste *string = new(string)

	// 谷物成分
	var grain *string = new(string)

	// 颗粒度
	var graininess *string = new(string)
	
	//insert record
	shopDetailId := insertShopDetail(
		*goodsName,
		*goodsNumber,
		*goodsSku,
		*brand,
		*firstCategory +" "+ *secondCategory +" "+ *thirdCategory,
		*goodsPrice,
		*salesVolume,
		*commonNum,
		*score,
		*shape,
		*age,
		*goodsComponent,
		*componentPercent,
		*taste,
		*grain,
		*graininess,3,sourceUrl)

	return shopDetailId, true
}

func save(p *page.Page) bool {
	//judge if status is 200
	if !checkCanSave(p) {
		return false
	}

	//get md5
	h := md5.New()
	url := p.GetRequest().Url
	h.Write([]byte(url))
	md5Str := hex.EncodeToString(h.Sum(nil))

	//get fullpath
	abPath := getPath(md5Str)
	fullDirPath := saveDir + abPath
	err := os.MkdirAll(fullDirPath, 0664)
	if err != nil {
		logger.Println("[error]create dir error:", err, " ", fullDirPath, " ", url)
		return false
	}

	//save file
	fileName := fullDirPath + "/" + path.Base(url)
	result, err := os.Create(fileName)
	if err != nil {
		logger.Println("[error]create file error:", err, " ", fileName, " ", url)
		return false
	}
	logger.Println("[info] save page:", url)
	logger.Println("[info] save in:", fileName)
	result.WriteString(url + "\n")
	result.WriteString(p.GetBodyStr())
	result.Close()

	return true
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

func getCommentContent(p *page.Page, s *goquery.Selection, content *string) {
	contentStr := ""
	s.Find("dl dt").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "内容：") {
            contentStr = s.Parent().Find("dd").Text()
        }
        return true
	})
	*content = contentStr
}

func getCommentTime(p *page.Page, s *goquery.Selection, commentTime *string) {
	commentTimeStr := ""
	s.Find("dl dd span i").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i%2!=0 {
			commentTimeStr = s.Text()
		}
		return true
	})
	*commentTime = commentTimeStr
}

func saveShopCommentList(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()
	sourceUrl := p.GetRequest().Url
	logger.Println("[info] common source url: ", sourceUrl)

	query.Find(".pl_list .pl_right").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// 评论内容
		var content *string = new(string)
		getCommentContent(p, s, content)
		logger.Println("[info] common content: ", *content, i)

		// 评论时间
		var commentTime *string = new(string)
		getCommentTime(p, s, commentTime)
		logger.Println("[info] common time: ", *commentTime,)

		//insert record
		insertShopComment(
			shopDetailId,
			*content,
			3,
			*commentTime)
		return true
	})
}
