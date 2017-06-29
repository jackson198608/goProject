package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
)

func qBaiduList(p *page.Page) {
	query := p.GetHtmlParser()

	//find shop list
	query.Find(".result.c-container .f13 .c-showurl").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			//result.WriteString(url + "\n")
			realUrlTag := "gouminDetail"
			req := newRequest(realUrlTag, url)
			p.AddTargetRequestWithParams(req)
		}
		return true
	})

	times := 0
	if p.GetUrlTag() == "searchList" {
		query.Find("#page a").EachWithBreak(func(i int, s *goquery.Selection) bool {
			// For each item found, get the band and title
			url, isExsit := s.Attr("href")
			if isExsit {
				if times == 4 {
					return false
				}
				realUrl := "https://www.baidu.com" + url
				realUrlTag := "searchListNextPage"
				req := newRequest(realUrlTag, realUrl)
				p.AddTargetRequestWithParams(req)
				times++
			}
			return true
		})
	}

}
