package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"strconv"
	"strings"
)

func qArticleList(p *page.Page, num int) {
	logger.Println("[info]find article list page: ", p.GetRequest().Url)
	query := p.GetHtmlParser()
	//find article list
	query.Find(".news-list .img-box a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			// url = "http://mp.weixin.qq.com" + url
			logger.Println("[info]find detail page: ", url)
			realUrlTag := "articleDetail"
			req := newRequest(realUrlTag, url)
			p.AddTargetRequestWithParams(req)
		}
		return true
	})

	query.Find(".news-box .p-fy a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		title := s.Text()
		if strings.Contains(title, "下一页") {
			url, isExsit := s.Attr("href")
			if isExsit {
				realUrl := "http://weixin.sogou.com/weixin" + url
				logger.Println("[info]find next article list page: ", realUrl)
				num++
				numstr := strconv.Itoa(num)
				realUrlTag := "articleList|" + numstr
				req := newRequest(realUrlTag, realUrl)
				p.AddTargetRequestWithParams(req)
			}
		}
		return true
	})
}
