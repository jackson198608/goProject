package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/hu17889/go_spider/core/common/page"
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
				/*
					realUrlTag := "shopDetail"
						req := newRequest(realUrlTag, realUrl)
						p.AddTargetRequestWithParams(req)
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
