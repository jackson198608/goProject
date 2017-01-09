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

func saveFirstAreaFind(p *page.Page) bool {
	query := p.GetHtmlParser()
	query.Find(".nc-items").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {

		if i == 1 {
			// find FirstBusiness
			s.Find(".nc-items a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
				//find it
				name := s.Text()
				logger.Println("[info]find first bussiness", name)

				//save it
				businessId, isExist := checkBusinessExist(name)
				if !isExist {
					businessId = insertBusiness(name, CityId, 0)
				}

				//when insert firt business fail
				if businessId == 0 {
					logger.Println("[error]insert first bussiness fail", name)
					return true
				}

				url, isExist := s.Attr("href")
				if isExist {
					realUrl := "http://www.dianping.com" + url
					tag := "shopSecondBusinessList|" + strconv.Itoa(int(businessId))
					logger.Println("[info]find second business list url", name, " ", realUrl, " ", tag)

					req := newRequest(tag, realUrl)
					p.AddTargetRequestWithParams(req)

				} else {
					logger.Println("this first bussiness has no second bussiness", name)
				}

				return true
			}, nil)

		} else if i == 2 {
			//find first region
			s.Find(".nc-items a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
				name := s.Text()
				logger.Println("find first region", name)

				//save it
				regionId, isExist := checkRegionExist(name)
				if !isExist {
					regionId = insertRegion(name, CityId, 0)
				}

				//when insert firt business fail
				if regionId == 0 {
					logger.Println("[error]insert first region fail", name)
					return true
				}

				url, isExist := s.Attr("href")
				if isExist {
					realUrl := "http://www.dianping.com" + url
					tag := "shopSecondRegionList|" + strconv.Itoa(int(regionId))
					logger.Println("[info]find second region list url", name, " ", realUrl, " ", tag)

					req := newRequest(tag, realUrl)
					p.AddTargetRequestWithParams(req)

				} else {
					logger.Println("this first region has no second region", name)
				}

				return true
			}, nil)

		} else if i == 3 {
			//find first region
			s.Find(".nc-items a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
				name := s.Text()
				logger.Println("find first metro", name)

				//save it
				metroId, isExist := checkMetroExist(name)
				if !isExist {
					metroId = insertMetro(name, CityId, 0)
				}

				//when insert firt business fail
				if metroId == 0 {
					logger.Println("[error]insert first metro fail", name)
					return true
				}

				url, isExist := s.Attr("href")
				if isExist {
					realUrl := "http://www.dianping.com" + url
					tag := "shopSecondMetroList|" + strconv.Itoa(int(metroId))
					logger.Println("[info]find second metro list url", name, " ", realUrl, " ", tag)

					req := newRequest(tag, realUrl)
					p.AddTargetRequestWithParams(req)

				} else {
					logger.Println("this first metro has no second metro", name)
				}

				return true
			}, nil)

		}

		return true
	}, nil)

	return true

}

func saveSecondBusinessList(p *page.Page, firstBusinessId int64) bool {
	query := p.GetHtmlParser()
	firstBusinessIdStr := strconv.Itoa(int(firstBusinessId))
	query.Find("#bussi-nav-sub a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		secondBusinessName := s.Text()
		secondBusinessUrl, isExist := s.Attr("href")
		logger.Println("[info] find second business", secondBusinessName)
		if i == 0 {
			if !isExist {
				return true
			}

			secondBusinessUrl = "http://www.dianping.com" + secondBusinessUrl
			tag := "shopUpdateOfSecondBusinessList|" + firstBusinessIdStr
			logger.Println("[info] find second business url", secondBusinessName, " ", secondBusinessUrl, " ", tag)
			req := newRequest(tag, secondBusinessUrl)
			p.AddTargetRequestWithParams(req)

		} else {
			secondBusinessId, isExist := checkBusinessExist(secondBusinessName)
			if !isExist {
				secondBusinessId = insertBusiness(secondBusinessName, CityId, firstBusinessId)
			}

			if secondBusinessId == 0 {
				logger.Println("[error]insert second business fail", secondBusinessName)
				return true
			}

			secondBusinessIdStr := strconv.Itoa(int(secondBusinessId))

			if !isExist {
				return true
			}

			secondBusinessUrl = "http://www.dianping.com" + secondBusinessUrl
			tag := "shopUpdateSecondBusinessList|" + firstBusinessIdStr + "|" + secondBusinessIdStr
			logger.Println("[info] find second business url", secondBusinessName, " ", secondBusinessUrl, " ", tag)
			req := newRequest(tag, secondBusinessUrl)
			p.AddTargetRequestWithParams(req)

		}

		return true
	}, nil)
	return true

}

func saveSecondRegionList(p *page.Page, firstRegionId int64) bool {
	query := p.GetHtmlParser()
	firstRegionIdStr := strconv.Itoa(int(firstRegionId))
	query.Find("#region-nav-sub a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		secondRegionName := s.Text()
		secondRegionUrl, isExist := s.Attr("href")
		logger.Println("[info] find second region", secondRegionName)
		if i == 0 {
			if !isExist {
				return true
			}

			secondRegionUrl = "http://www.dianping.com" + secondRegionUrl
			tag := "shopUpdateOfSecondRegionList|" + firstRegionIdStr
			logger.Println("[info] find second region url", secondRegionName, " ", secondRegionUrl, " ", tag)
			req := newRequest(tag, secondRegionUrl)
			p.AddTargetRequestWithParams(req)

		} else {
			secondRegionId, isExist := checkRegionExist(secondRegionName)
			if !isExist {
				secondRegionId = insertRegion(secondRegionName, CityId, firstRegionId)
			}

			if secondRegionId == 0 {
				logger.Println("[error]insert second region fail", secondRegionName)
				return true
			}

			secondRegionIdStr := strconv.Itoa(int(secondRegionId))

			if !isExist {
				return true
			}

			secondRegionUrl = "http://www.dianping.com" + secondRegionUrl
			tag := "shopUpdateSecondRegionList|" + firstRegionIdStr + "|" + secondRegionIdStr
			logger.Println("[info] find second region url", secondRegionName, " ", secondRegionUrl, " ", tag)
			req := newRequest(tag, secondRegionUrl)
			p.AddTargetRequestWithParams(req)

		}

		return true
	}, nil)
	return true

}

func saveSecondMetroList(p *page.Page, firstMetroId int64) bool {
	query := p.GetHtmlParser()
	firstMetroIdStr := strconv.Itoa(int(firstMetroId))
	query.Find("#metro-nav-sub a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		secondMetroName := s.Text()
		secondMetroUrl, isExist := s.Attr("href")
		logger.Println("[info] find second metro", secondMetroName)
		if i == 0 {
			if !isExist {
				return true
			}

			secondMetroUrl = "http://www.dianping.com" + secondMetroUrl
			tag := "shopUpdateOfSecondMetroList|" + firstMetroIdStr
			logger.Println("[info] find second metro url", secondMetroName, " ", secondMetroUrl, " ", tag)
			req := newRequest(tag, secondMetroUrl)
			p.AddTargetRequestWithParams(req)

		} else {
			secondMetroId, isExist := checkMetroExist(secondMetroName)
			if !isExist {
				secondMetroId = insertMetro(secondMetroName, CityId, firstMetroId)
			}

			if secondMetroId == 0 {
				logger.Println("[error]insert second metro fail", secondMetroName)
				return true
			}

			secondMetroIdStr := strconv.Itoa(int(secondMetroId))

			if !isExist {
				return true
			}

			secondMetroUrl = "http://www.dianping.com" + secondMetroUrl
			tag := "shopUpdateSecondMetroList|" + firstMetroIdStr + "|" + secondMetroIdStr
			logger.Println("[info] find second metro url", secondMetroName, " ", secondMetroUrl, " ", tag)
			req := newRequest(tag, secondMetroUrl)
			p.AddTargetRequestWithParams(req)

		}

		return true
	}, nil)
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

func getFirtBusiness(query *goquery.Document) {
	query.Find(".nc-items").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 1 {
			s.Find(".nc-items a span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
				name := s.Text()
				logger.Println("find first bussiness", name)
				return false
			}, nil)

		}
		return false
	}, nil)
}

func getShopType(query *goquery.Document, shopType *string) {
	query.Find(".breadcrumb a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 2 {
			*result = s.Text()
			*result = strings.Replace(*result, " ", "", -1)
			*result = strings.Replace(*result, "\n", "", -1)
			return false
		}

		return true
	}, shopType)
}

func getShopName(query *goquery.Document, shopName *string) {
	query.Find(".breadcrumb span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		*result = s.Text()
		return false
	}, shopName)
}

func getShopAddress(query *goquery.Document, shopAddress *string) {
	query.Find(".address span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 2 {
			*result, _ = s.Attr("title")
			return false
		}
		return true
	}, shopAddress)
}

func getShopCity(query *goquery.Document, shopCity *string) {
	query.Find(".J-city").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		*result = s.Text()
		return false
	}, shopCity)
}

func getShopPhone(query *goquery.Document, shopPhone *string) {
	query.Find(".tel span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i >= 1 {
			*result = *result + " " + s.Text()
		}
		return true

	}, shopPhone)
}

func getShopStar(query *goquery.Document, shopStar *string) {
	query.Find(".mid-rank-stars").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		Class, isExist := s.Attr("class")
		if isExist {
			starClass := strings.Split(Class, " ")[1]
			lenStarClass := len(starClass)
			realStar := starClass[lenStarClass-2 : lenStarClass]
			*result = string(realStar)
		}
		return false

	}, shopStar)
}

func deleteHeadBlank(former string) string {
	//find first not null byte position
	var i int = 0
	for i = 0; i < len(former); i++ {
		if former[i] != ' ' {
			break
		}
	}
	realTime := former[i:]

	return string(realTime)
}

func getShopTime(query *goquery.Document, shopTime *string) {
	//get 营业时间 position
	var dataI string = "-1"
	query.Find(".J-other .info-name").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		infoName := s.Text()
		logger.Println("infoName: ", infoName)
		if infoName == "营业时间：" {
			*result = strconv.Itoa(i)
			return false
		}
		return true
	}, &dataI)

	//get 营业时间 from position
	query.Find(".J-other .item").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		j, err := strconv.Atoi(*result)
		if err != nil {
			return false
		}
		if i == j {
			*result = s.Text()
			return false
		}
		return true
	}, &dataI)
	logger.Println("[info]dataI: ", dataI)

	//remove addional \n
	results := strings.Split(dataI, "\n")
	*shopTime = deleteHeadBlank(results[1])
}

func getShopPrice(query *goquery.Document, shopPrice *string) {
	query.Find(".brief-info span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 2 {
			//人均: 390元
			priceText := s.Text()
			*result = getRealPriceFromText(priceText)
			return false
		}
		return true
	}, shopPrice)
}

func getShopServicePoint(query *goquery.Document, shopServicePoint *string) {
	query.Find(".brief-info span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 3 {
			//服务：9.3
			pointsText := s.Text()
			*result = getRealPriceFromText(pointsText)
			return false
		}
		return true
	}, shopServicePoint)
}

func getShopEnvPoint(query *goquery.Document, shopEnvPoint *string) {
	query.Find(".brief-info span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 4 {
			//环境：9.3
			pointsText := s.Text()
			*result = getRealPriceFromText(pointsText)
			return false
		}
		return true
	}, shopEnvPoint)
}

func getShopWeightPoint(query *goquery.Document, shopWeightPoint *string) {
	query.Find(".brief-info span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 5 {
			//划算：9.4
			pointsText := s.Text()
			*result = getRealPriceFromText(pointsText)
			return false
		}
		return true
	}, shopWeightPoint)
}

func getImageFromPage(query *goquery.Document, shopImage *string) {
	query.Find(".pic-list-b img").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		shopUrl, isExist := s.Attr("src")
		if !isExist {
			return false
		}
		*result = shopUrl
		return false
	}, shopImage)
}

func getRealShopImageUrl(query *goquery.Document, realShopUrl *string) {
	query.Find(".picture-list .J_list .img a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		shopUrl, isExist := s.Attr("href")
		if !isExist {
			return false

		}
		*result = "http://www.dianping.com" + shopUrl
		return false

	}, realShopUrl)

}

func getRealPriceFromText(priceText string) string {
	realPriceText := priceText[9:len(priceText)]
	return string(realPriceText)
}

func getFloatPrice(realPriceText string, deleteYuan bool) int {
	if realPriceText == "-" {
		return 0
	} else {
		if deleteYuan {
			realPriceTextWithoutYuan := realPriceText[0 : len(realPriceText)-3]
			realPriceText = string(realPriceTextWithoutYuan)
		}
		realFloatPrice, err := strconv.ParseFloat(realPriceText, 10)
		if err == nil {
			return int(realFloatPrice * 100)
		} else {
			logger.Println("[error] get int price error,make it to be 0", err, "realPriceText:", realPriceText)
			return 0
		}

	}
}

func saveShopDetail(p *page.Page) (int64, bool) {
	query := p.GetHtmlParser()
	//find name
	var shopType *string = new(string)
	var shopIntType int = 0
	getShopType(query, shopType)
	logger.Println("[info] shopType: ", *shopType)

	switch *shopType {
	case "宠物店":
		shopIntType = 1
	case "宠物医院":
		shopIntType = 2
	}
	logger.Println("[info] shopIntType: ", shopIntType)

	var shopName *string = new(string)
	getShopName(query, shopName)
	logger.Println("[info] shopName: ", *shopName)

	var shopCity *string = new(string)
	getShopCity(query, shopCity)
	logger.Println("[info] shopCity: ", *shopCity)

	shopFindId, isExist := checkShopExist(*shopName, *shopCity)
	if isExist {
		logger.Println("[info]find shop exist ", *shopName, " ", shopFindId)
		return shopFindId, false
	}

	//find address
	var shopAddress *string = new(string)
	getShopAddress(query, shopAddress)
	logger.Println("[info] shopAddress: ", *shopAddress)

	//find phone
	var shopPhone *string = new(string)
	getShopPhone(query, shopPhone)
	logger.Println("[info] shopPhone: ", *shopPhone)

	//find price
	var shopPrice *string = new(string)
	getShopPrice(query, shopPrice)
	shopIntPrice := getFloatPrice(*shopPrice, true)
	logger.Println("[info] shopPrice: ", shopIntPrice)

	//find star
	var shopStar *string = new(string)
	getShopStar(query, shopStar)
	shopIntStar, err := strconv.Atoi(*shopStar)
	if err != nil {
		logger.Println("[error] parsInt error for shop star", err, "star:", *shopStar)
		if (*shopStar)[len(*shopStar)-1] == '0' {
			logger.Println("[info] fix error for shop star make it to be 0")
			shopIntStar = 0
		} else {
			logger.Println("[error] still error for shop star return")
			return 0, false
		}

	}
	logger.Println("[info] shopIntStar: ", shopIntStar)

	//find service point
	var shopServicePoint *string = new(string)
	getShopServicePoint(query, shopServicePoint)
	shopIntServicePoint := getFloatPrice(*shopServicePoint, false)
	logger.Println("[info] shopServicePoint: ", shopIntServicePoint)

	//find env point
	var shopEnvPoint *string = new(string)
	getShopEnvPoint(query, shopEnvPoint)
	shopIntEnvPoint := getFloatPrice(*shopEnvPoint, false)
	logger.Println("[info] shopEnvPoint: ", shopIntEnvPoint)

	//find weight point
	var shopWeightPoint *string = new(string)
	getShopWeightPoint(query, shopWeightPoint)
	shopIntWeightPoint := getFloatPrice(*shopWeightPoint, false)
	logger.Println("[info] shopWeightPoint: ", shopIntWeightPoint)

	//find shop time
	var shopTime *string = new(string)
	getShopTime(query, shopTime)
	logger.Println("[info] shopTime: ", *shopTime)

	//insert record
	shopDetailId := insertShopDetail(
		*shopCity,
		//shopIntType,
		Type,
		*shopName,
		*shopAddress,
		*shopPhone,
		0,
		shopIntPrice,
		shopIntStar,
		shopIntServicePoint,
		shopIntEnvPoint,
		shopIntWeightPoint,
		*shopTime, "http://www.baidu.com", 0, 0, 0, 0, 0, 0)

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

func saveShopImagePath(p *page.Page, shopDetailId int64, isFirst bool) (int64, bool) {
	query := p.GetHtmlParser()
	var shopImage *string = new(string)
	getImageFromPage(query, shopImage)
	if *shopImage == "" {
		logger.Println("[info] no Image found:")
		var realShopUrl *string = new(string)
		getRealShopImageUrl(query, realShopUrl)
		logger.Println("[info] no Image found,try to find real:", *realShopUrl)
		req := newRequest(p.GetUrlTag(), *realShopUrl)
		p.AddTargetRequestWithParams(req)

		return 0, false
	}
	logger.Println("[info]Image url:", *shopImage)
	imageRealPath := getPathFromUrl(*shopImage)
	logger.Println("[info]Image path:", imageRealPath)

	//save image into shop_image table
	insertShopPhoto(shopDetailId, imageRealPath)

	if isFirst {
		saveSucess := updateImageForFirstPhoto(shopDetailId, imageRealPath)
		if !saveSucess {
			logger.Println("[error] update first image into shop detail fail", shopDetailId, imageRealPath)
		}
	}

	//生成对图片的抓取任务
	req := newImageRequest("shopImage", *shopImage)
	p.AddTargetRequestWithParams(req)
	return 0, true
}

func saveShopCommentImagePath(p *page.Page, shopCommentId int64) (int64, bool) {
	query := p.GetHtmlParser()
	var shopImage *string = new(string)
	getImageFromPage(query, shopImage)
	if *shopImage == "" {
		return 0, false
	}
	logger.Println("[info]Image url:", *shopImage)
	imageRealPath := getPathFromUrl(*shopImage)
	logger.Println("[info]Image path:", imageRealPath)

	//save image into comment_photo table
	commentPhotoId := insertCommentPhoto(shopCommentId, imageRealPath)

	//生成对图片的抓取任务
	req := newImageRequest("shopImage", *shopImage)
	p.AddTargetRequestWithParams(req)
	return commentPhotoId, true
}

func getPathFromUrl(url string) string {
	//find first ? position
	var i int = 0
	len := len(url)
	for i = 0; i < len; i++ {
		if url[i] == '%' {
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

func getCommentAvar(p *page.Page, s *goquery.Selection) string {
	var avar string
	s.Find(".pic img").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		avarUrl, isExist := s.Attr("src")
		if isExist {
			avar = getPathFromUrl(avarUrl)
			//生成对图片的抓取任务
			req := newImageRequest("shopImage", avarUrl)
			p.AddTargetRequestWithParams(req)
		}
		return false
	}, nil)
	return avar

}

func getCommentUsername(p *page.Page, s *goquery.Selection) string {
	var username string
	s.Find(".pic .name a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		username = s.Text()
		return false
	}, nil)
	return username
}

func getCommentContent(p *page.Page, s *goquery.Selection) string {
	var content string
	s.Find(".content .comment-txt .J_brief-cont").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		content = s.Text()
		return false
	}, nil)
	return content
}

func getCommentStar(p *page.Page, s *goquery.Selection) int {
	var star string
	var intStar int
	s.Find(".content .user-info span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		star, _ = s.Attr("class")
		return false
	}, nil)
	starArr := strings.Split(star, " ")
	if len(starArr) != 2 {
		return 0
	}
	realStarStr := starArr[1]

	if len(realStarStr) < 9 {
		return 0
	}

	realStar := realStarStr[8:]
	intStar, err := strconv.Atoi(realStar)
	if err != nil {
		return 0
	}
	return intStar
}

func getCommentPrice(p *page.Page, s *goquery.Selection) int {
	var price string
	var intPrice int
	s.Find(".content .comm-per").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		price = s.Text()
		return false
	}, nil)
	priceArr := strings.Split(price, " ")
	if len(priceArr) != 2 {
		return 0
	}
	realPriceStr := priceArr[1]

	if len(realPriceStr) < 4 {
		return 0
	}

	realPrice := realPriceStr[3:]
	intPrice, err := strconv.Atoi(realPrice)
	if err != nil {
		return 0
	}
	return intPrice
}

func getCommentServicePoint(p *page.Page, s *goquery.Selection) int {
	var servicePoint string
	var intServicePoint int
	s.Find(".content .user-info .comment-rst span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 0 {
			servicePoint = s.Text()
			return false
		}
		return true
	}, nil)

	if len(servicePoint) < 8 {
		return 0
	}

	realServicePoint := servicePoint[6]
	intServicePoint, err := strconv.Atoi(string(realServicePoint))
	if err != nil {
		return 0
	}
	return intServicePoint * 100
}

func getCommentEnvPoint(p *page.Page, s *goquery.Selection) int {
	var envPoint string
	var intEnvPoint int
	s.Find(".content .user-info .comment-rst span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 1 {
			envPoint = s.Text()
			return false
		}
		return true
	}, nil)

	if len(envPoint) < 8 {
		return 0
	}

	realEnvPoint := envPoint[6]
	intEnvPoint, err := strconv.Atoi(string(realEnvPoint))
	if err != nil {
		return 0
	}
	return intEnvPoint * 100
}

func getCommentTime(p *page.Page, s *goquery.Selection) string {
	var commentTime string
	s.Find(".content .misc-info .time").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		commentTimeStr := s.Text()
		if len(commentTimeStr) >= 5 {
			realCommentTime := commentTimeStr[0:5]
			commentTime = string(realCommentTime)
		}
		return false
	}, nil)
	return "2016-" + commentTime + " 00:00:00"
}

func getCommentWeightPoint(p *page.Page, s *goquery.Selection) int {
	var weightPoint string
	var intWeightPoint int
	s.Find(".content .user-info .comment-rst span").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 2 {
			weightPoint = s.Text()
			return false
		}
		return true
	}, nil)

	if len(weightPoint) < 8 {
		return 0
	}

	realWeightPoint := weightPoint[6]
	intWeightPoint, err := strconv.Atoi(string(realWeightPoint))
	if err != nil {
		return 0
	}
	return intWeightPoint * 100
}

func processEachComment(p *page.Page, s *goquery.Selection, shopDetailId int64) int64 {

	//find comment content
	content := getCommentContent(p, s)
	logger.Println("[info]comment content:", content)

	//find user name
	username := getCommentUsername(p, s)
	logger.Println("[info]comment username:", username)

	//find user avar
	avar := getCommentAvar(p, s)
	logger.Println("[info]comment avar:", avar)

	//find price
	price := getCommentPrice(p, s)
	logger.Println("[info]comment price:", price)

	//find star
	star := getCommentStar(p, s)
	logger.Println("[info]comment star:", star)

	//find servicePoint
	servicePoint := getCommentServicePoint(p, s)
	logger.Println("[info]comment servicePoint:", servicePoint)

	//find envPoint
	envPoint := getCommentEnvPoint(p, s)
	logger.Println("[info]comment envPoint:", envPoint)

	//find weightPoint
	weightPoint := getCommentWeightPoint(p, s)
	logger.Println("[info]comment weightPoint:", weightPoint)

	//find commentTime
	//commentTime := getCommentTime(p, s)
	commentTime := "2017-01-01 00:00:00"
	logger.Println("[info]comment commentTime:", commentTime)

	shopCommentId := insertShopComment(shopDetailId, content, username, avar, price, star, servicePoint, envPoint, weightPoint, commentTime)

	//find comment photo

	return shopCommentId

}

func saveShopCommentList(p *page.Page, shopDetailId int64) bool {
	query := p.GetHtmlParser()

	query.Find(".comment-list li").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		liId, isExist := s.Attr("id")
		if isExist && strings.Contains(liId, "rev_") {
			commentId := processEachComment(p, s, shopDetailId)
			logger.Println("[info] insert comment id:", commentId)
			qCommentPhotoPage(p, s, commentId)
		}
		return true
	}, nil)

	return true
}
