package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"regexp"
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

func qGouminList(p *page.Page) {
	query := p.GetHtmlParser()

	//find shop list
	query.Find(".new a.xst").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			realUrl := "http://bbs.goumin.com/" + url
			logger.Println("[info]find detail page: ", realUrl)
			/*
				realUrlTag := "shopDetail"
					req := newRequest(realUrlTag, realUrl)
					p.AddTargetRequestWithParams(req)
			*/
		}
		return true
	}, nil)

	query.Find(".pg .nxt").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			realUrl := "http://bbs.goumin.com/" + url
			logger.Println("[info]find next list page: ", realUrl)
			realUrlTag := "gouminList"
			req := newRequest(realUrlTag, realUrl)
			p.AddTargetRequestWithParams(req)
		}
		return false
	}, nil)

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
func qShopUpdateList(p *page.Page, tag string, first int64, second int64, Type int, cat int) {
	query := p.GetHtmlParser()
	images := make([]string, 20, 20)

	query.Find(".shop-all-list .pic a img").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		// For each item found, get the band and title
		//if i == 0 {
		url, isExsit := s.Attr("data-src")
		if isExsit {
			images[i] = url
			realUrlTag := "shopImage"
			logger.Println("[info] fetch image url: ", url)
			req := newImageRequest(realUrlTag, url)
			p.AddTargetRequestWithParams(req)
		}
		//}
		return true
	}, nil)

	//find shop list
	query.Find(".shop-all-list .txt .tit a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			isMatch, err := regexp.MatchString("^/shop/\\d+$", url)
			if err != nil {
				logger.Println("[error]Match parse error", url)
				return true
			}

			shopImageUrl := images[i]
			shopImage := getPathFromUrl(shopImageUrl)

			if isMatch {
				realUrl := "http://www.dianping.com" + url
				logger.Println("[info]find detail url: ", realUrl)
				logger.Println("[info]find detail url: ", shopImageUrl, shopImage)
				/*
					//find name
					s.Find("h4").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
						shopName := s.Text()
						logger.Println("[info]find detail name: ", shopName, " ", tag, " shop Image:", shopImage)
						_, isExist := checkShopExist(shopName, City)
						if !isExist {
							logger.Println("[info]shop is not exist,fetch it ", shopName, " ", realUrl)
							tag := "shopDetail"
							req := newRequest(tag, realUrl)
							p.AddTargetRequestWithParams(req)

							return false
						}

						updateShopImage(shopName, shopImage)
							if Type == 1 {
								switch cat {
								case 1:
									logger.Println("[info]in the updateShopBusiness", shopName, " ", tag)
									updateShopBusiness(shopName, first)
								case 2:
									logger.Println("[info]in the updateShopRegion", shopName, " ", tag)
									updateShopRegion(shopName, first)
								case 3:
									logger.Println("[info]in the updateShopMetro", shopName, " ", tag)
									updateShopMetro(shopName, first)
								}
							} else if Type == 2 {
								switch cat {
								case 1:
									logger.Println("[info]in the updateShopBusinessWithSub", shopName, " ", tag, " first:", first, " second:", second)
									updateShopBusinessWithSub(shopName, first, second)
								case 2:
									logger.Println("[info]in the updateShopRegionWithSub", shopName, " ", tag)
									updateShopRegionWithSub(shopName, first, second)
								case 3:
									logger.Println("[info]in the updateShopMetroWithSub", shopName, " ", tag)
									updateShopMetroWithSub(shopName, first, second)
								}
							}
						return false
					}, nil)
				*/

			}
		}
		return true
	}, nil)

	query.Find(".page .next").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		// For each item found, get the band and title
		//if i == 0 {
		url, isExsit := s.Attr("href")
		if isExsit {
			realUrl := "http://www.dianping.com" + url
			logger.Println("[info]find next list page: ", realUrl)
			req := newRequest(tag, realUrl)
			p.AddTargetRequestWithParams(req)
		}
		//}
		return true
	}, nil)
}

func qShopList(p *page.Page) {
	query := p.GetHtmlParser()

	//find shop list
	query.Find(".shop-all-list .txt .tit a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			isMatch, err := regexp.MatchString("^/shop/\\d+$", url)
			if err != nil {
				logger.Println("[error]Match parse error", url)
				return true
			}

			if isMatch {
				realUrl := "http://www.dianping.com" + url
				logger.Println("[info]find detail page: ", realUrl)
				realUrlTag := "shopDetail"
				req := newRequest(realUrlTag, realUrl)
				p.AddTargetRequestWithParams(req)
			}
		}
		return true
	}, nil)

	query.Find(".page .next").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		// For each item found, get the band and title
		//if i == 0 {
		url, isExsit := s.Attr("href")
		if isExsit {
			realUrl := "http://www.dianping.com" + url
			logger.Println("[info]find next list page: ", realUrl)
			realUrlTag := "shopList"
			req := newRequest(realUrlTag, realUrl)
			p.AddTargetRequestWithParams(req)
		}
		//}
		return true
	}, nil)

}

func qShopDetail(p *page.Page, shopDetailId int64) {
	query := p.GetHtmlParser()

	//find shop image1
	query.Find(".photo-header a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		// For each item found, get the band and title
		if i == 0 {
			url, isExsit := s.Attr("href")
			if isExsit {
				realUrl := "http://www.dianping.com" + url
				logger.Println("[info]find image start page: ", realUrl)
				shopDetailIdStr := strconv.Itoa(int(shopDetailId))
				realUrlTag := "shopImage1|" + shopDetailIdStr
				req := newRequest(realUrlTag, realUrl)
				p.AddTargetRequestWithParams(req)
				return false
			}
		}
		return true
	}, nil)

	//find commentlist
	query.Find(".J-tab a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		// For each item found, get the band and title
		if i == 1 {
			url, isExsit := s.Attr("href")
			if isExsit {
				realUrl := "http://www.dianping.com" + url
				logger.Println("[info]find comment list page: ", realUrl)
				shopDetailIdStr := strconv.Itoa(int(shopDetailId))
				realUrlTag := "shopCommentList|" + shopDetailIdStr
				req := newRequest(realUrlTag, realUrl)
				p.AddTargetRequestWithParams(req)
				return false
			}
		}
		return true
	}, nil)

}

func qShopCommentList(p *page.Page, shopDetailId int64) {
	//get next page
	query := p.GetHtmlParser()
	url := p.GetRequest().Url
	query.Find(".NextPage").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		logger.Println("[info]in the next page")
		url, isExsit := s.Attr("href")
		if isExsit {
			realUrl := getNextPageUrl(*result, url)
			logger.Println("[info]find next page: ", realUrl)
			shopDetailIdStr := strconv.Itoa(int(shopDetailId))
			realUrlTag := "shopCommentList" + "|" + shopDetailIdStr
			req := newRequest(realUrlTag, realUrl)
			p.AddTargetRequestWithParams(req)
		}
		return false
	}, &url)

}

func qShopImage(p *page.Page, shopDetailId int64, imageNum int) {
	query := p.GetHtmlParser()
	imageNum = imageNum + 1
	imageNumStr := strconv.Itoa(imageNum)

	//find all the list detail page
	query.Find(".pic-nav-wrap a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		if i == 1 {
			url, isExsit := s.Attr("href")
			if isExsit {
				realUrl := "http://www.dianping.com" + url
				logger.Println("[info]find image start page: ", realUrl)
				shopDetailIdStr := strconv.Itoa(int(shopDetailId))
				realUrlTag := "shopImage" + *result + "|" + shopDetailIdStr
				req := newRequest(realUrlTag, realUrl)
				p.AddTargetRequestWithParams(req)
				return false
			}
		}
		return true
	}, &imageNumStr)
}

func qCommentPhotoPage(p *page.Page, s *goquery.Selection, shopCommentId int64) bool {
	s.Find(".shop-photo a").EachWithBreak(func(i int, s *goquery.Selection, result *string) bool {
		photoPageUrl, isExist := s.Attr("href")
		if isExist {
			realPhotoPageUrl := "http://www.dianping.com" + photoPageUrl
			tag := "shopCommentImage|" + strconv.Itoa(int(shopCommentId))
			logger.Println("[info]photoPageUrl:", realPhotoPageUrl, " tag:", tag)
			req := newRequest(tag, realPhotoPageUrl)
			p.AddTargetRequestWithParams(req)
		}
		return true
	}, nil)

	return true
}
