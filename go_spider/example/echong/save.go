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
	goodsNameStr := ""
	query.Find(".gddes h1").EachWithBreak(func(i int, s *goquery.Selection) bool {
		goodsNameStr = s.Text()
		goodsNameStr = strings.Replace(goodsNameStr, " ", "", -1)
		goodsNameStr = strings.Replace(goodsNameStr, "\n", "", -1)
		return true
	})

	if goodsNameStr=="" {
		//直邮商品样式
		query.Find(".xq-name h1").EachWithBreak(func(i int, s *goquery.Selection) bool {
			goodsNameStr = s.Text()
			goodsNameStr = strings.Replace(goodsNameStr, " ", "", -1)
			goodsNameStr = strings.Replace(goodsNameStr, "\n", "", -1)
			return true
		})
	}
	*goodsName = goodsNameStr
}

func getFirstCategory(query *goquery.Document, category *string) {
	categoryStr := "" 
	query.Find(".pet-onav .onav-cont span").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			categoryStr = s.Text()	
		}
		return true
	})
	*category = categoryStr
}

func getSecondCategory(query *goquery.Document, category *string) {
	categoryStr := "" 
	query.Find(".pet-onav .onav-cont span").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 1 {
			categoryStr = s.Text()	
		}
		return true
	})
	*category = categoryStr
}

func getThirdCategory(query *goquery.Document, category *string) {
	categoryStr := "" 
	query.Find(".pet-onav .onav-cont span").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 2 {
			categoryStr = s.Text()
			categoryStr = strings.Replace(categoryStr, " ", "", -1)
			categoryStr = strings.Replace(categoryStr, "\n", "", -1)
		}
		return true	
	})
	*category = categoryStr
}

func getGoodsPrice(query *goquery.Document, goodsPrice *float64) {
	query.Find(".epet-pprice .mt5 .c93c.goods-font").EachWithBreak(func(i int, s *goquery.Selection) bool {
		price := s.Text()
		price = strings.Replace(price, " ", "", -1)
		price = strings.Replace(price, "\n", "", -1)
		priceF,_ := strconv.ParseFloat(price,64)
		*goodsPrice = priceF
		return true
	})

	// 直邮样式
	query.Find(".this-price #goods-sale-price").EachWithBreak(func(i int, s *goquery.Selection) bool {
		price := s.Text()
		priceF,_ := strconv.ParseFloat(price,64)
		*goodsPrice = priceF
		return true
	})
}

func getGoodsSku(query *goquery.Document, goodsSku *string) {
	sku := ""
	query.Find("a.norms-a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if s.Find(".goods-select").Length() > 0 {
            text := s.Text()
            text = strings.Replace(text, " ", "", -1)
			text = strings.Replace(text, "\n", "", -1)
			sku += text + " "
        }
		return true
	})
	*goodsSku = sku
}

func getAge(query *goquery.Document, age *string) {
	ageStr := ""
	query.Find(".textR").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(),"年龄") {
			ageStr = s.Siblings().Text()
			ageStr = strings.Replace(ageStr, " ", "", -1)
			ageStr = strings.Replace(ageStr, "\n", "", -1)
		}
		return true
	})
	*age = ageStr
}

func getShape(query *goquery.Document, shape *string) {
	shapeStr := ""
	query.Find(".textR").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(),"体型") {
			shapeStr = s.Siblings().Text()
			shapeStr = strings.Replace(shapeStr, " ", "", -1)
			shapeStr = strings.Replace(shapeStr, "\n", "", -1)
		}
		return true
	})
	*shape = shapeStr
}

func getComponent(query *goquery.Document, component *string) {
	componentStr := ""
	query.Find(".textR").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(),"主要成份") {
			componentStr = s.Siblings().Text()
			componentStr = strings.Replace(componentStr, " ", "", -1)
			componentStr = strings.Replace(componentStr, "\n", "", -1)
		}
		return true
	})
	*component = componentStr
}

func getComponentPercent(query *goquery.Document, component *string) {
	componentStr := ""
	query.Find(".textR").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(),"成份含量") {
			componentStr = s.Siblings().Text()
			componentStr = strings.Replace(componentStr, " ", "", -1)
			componentStr = strings.Replace(componentStr, "\n", "", -1)
		}
		return true
	})
	*component = componentStr
}

func getGraininess(query *goquery.Document, graininess *string) {
	graininessStr := ""
	query.Find(".textR").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(),"颗粒大小") {
			graininessStr = s.Siblings().Find("div span").Text()
			graininessStr = strings.Replace(graininessStr, " ", "", -1)
			graininessStr = strings.Replace(graininessStr, "\n", "", -1)
		}
		return true
	})
	*graininess = graininessStr
}

func getBrand(query *goquery.Document, brand *string) {
	brandStr := ""
	query.Find(".brands-home .ft14 a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i==0 {
			brandStr = s.Text()
		}
		return true
	})
	*brand = brandStr
}

func getGoodsNum(query *goquery.Document, goodsNum *int) {
	goodsNumInt := 0
	query.Find(".ftc .xq-num").EachWithBreak(func(i int, s *goquery.Selection) bool {
		num,_:=strconv.Atoi(s.Text())
		goodsNumInt = num
		return true
	})

	if  goodsNumInt==0 {
		//直邮样式
		query.Find(".clearfix .c999 span").EachWithBreak(func(i int, s *goquery.Selection) bool {
			numStr := strings.Split(s.Text(), "商品编号：")
			num,_:=strconv.Atoi(numStr[1])
			goodsNumInt = num
			return true
		})
	}
	*goodsNum = goodsNumInt
}

func getSalesVolume(query *goquery.Document, salesVolume *int) {
	salesVolumeInt := 0
	query.Find(".ats-style .ce54649").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i==0 {
			num,_:=strconv.Atoi(s.Text())
			salesVolumeInt = num
		}
		return true
	})

	if salesVolumeInt == 0 {
		//直邮样式
		query.Find(".clearfix .fl span").EachWithBreak(func(i int, s *goquery.Selection) bool {
			numStr := strings.Split(s.Text(), "已购买人数")
			num,_:=strconv.Atoi(numStr[1])
			salesVolumeInt = num
			return true
		})
	}
	*salesVolume = salesVolumeInt
}

//没有获取到数据
func getScore(query *goquery.Document, score *float64) {
	// query.Find(".pl_l .pl_score span").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	*result = s.Text()
	// 	return true
	// })
	*score = 0.0

	//直邮样式
	query.Find(".had-buy .clearfix .fl a span").EachWithBreak(func(i int, s *goquery.Selection) bool {
		text := s.Text() //4.8分 (3344条评论）
		num := strings.Split(text,"分")
		scoreF,_ := strconv.ParseFloat(num[0],64)
    	*score = scoreF;
		return true
	})
}

func getCommentNum(query *goquery.Document, commentNum *int) {
	commentNumStr := ""
	query.Find(".ats-style .c300.ce54649").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i==0 {
			text := s.Text()
			num := strings.Split(text,"(")
        	num1 := strings.Split(num[1],")")
        	commentNumStr = num1[0];
		}
		return true
	})

	//直邮样式
	if commentNumStr=="" {
		query.Find(".had-buy .clearfix .fl a span").EachWithBreak(func(i int, s *goquery.Selection) bool {
			text := s.Text() //4.8分 (3344条评论）
			num := strings.Split(text,"(")
	    	num1 := strings.Split(num[1],"条")
	    	commentNumStr = num1[0];
			return true
		})
	}
	*commentNum,_ = strconv.Atoi(commentNumStr);
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

	// 适用犬型
	var shape *string = new(string)
	// 适用年龄
	var age *string = new(string)
 	// 成分
	var goodsComponent *string = new(string)
	// 成分含量
	var componentPercent *string = new(string)
	// 口味
	var taste *string = new(string)
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
	}else{
		logger.Println("[info]again find goods detail: ", sourceUrl)
		realUrlTag := "shopDetail"
		req := newRequest(realUrlTag, sourceUrl)
		p.AddTargetRequestWithParams(req)
	}
	return 0,false
}

func saveShopDetailParams(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()

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

	// 成分含量
	var componentPercent *string = new(string)
	getComponentPercent(query, componentPercent)
	logger.Println("[info] goods component percent: ", *componentPercent)

	// 颗粒度
	var graininess *string = new(string)
	getGraininess(query, graininess)
	logger.Println("[info] goods graininess: ", *graininess)

	updateShopDetail(
		*shape,
		*age,
		*goodsComponent,
		*componentPercent,
		*graininess,
		shopDetailId)

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
	s.Find(".userrg .pingjiatext").EachWithBreak(func(i int, s *goquery.Selection) bool {
        contentStr = s.Text()
        contentStr = strings.Replace(contentStr, " ","",-1)
        contentStr = strings.Replace(contentStr, "\n","",-1)
        return true
	})
	*content = contentStr
}

func getCommentTime(p *page.Page, s *goquery.Selection, commentTime *string) {
	commentTimeStr := ""
	s.Find(".userrg .user-huifu .fr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		commentTimeStr = s.Text()
		return true
	})
	*commentTime = commentTimeStr
}

func saveShopCommentList(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()
	sourceUrl := p.GetRequest().Url
	logger.Println("[info] common source url: ", sourceUrl)

	query.Find(".evaluation").EachWithBreak(func(i int, s *goquery.Selection) bool {
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
			4,
			*commentTime)
		return true
	})
}
