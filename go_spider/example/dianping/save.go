package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/page"
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
			logger.Println("[error] get int price error", err, realPriceText)
			return 0
		}

	}
}

func saveShopDetail(p *page.Page) (int64, bool) {
	query := p.GetHtmlParser()
	//find name
	var shopName *string = new(string)
	getShopName(query, shopName)
	logger.Println("[info] shopName: ", *shopName)

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
		logger.Println("[error] parsInt error for shop int start", err)
		return 0, false
	}
	logger.Println("[info] shopStar: ", *shopStar)

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
		City,
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
		*shopTime, "http://www.baidu.com")

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
	return "2016-" + commentTime
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
	commentTime := getCommentTime(p, s)
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
