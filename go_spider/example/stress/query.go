package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"time"
)

func qBbsList(p *page.Page) {
	query := p.GetHtmlParser()
	query.Find("#threadlist .bm_c .common").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, _ := s.Find("a").Attr("href")
		start := time.Now()
		req := newRequest("bbsview", url)
		//记录结束时间
		end := time.Since(start)
		//输出执行时间，单位为毫秒。
		fmt.Println(end / 1000)
		fmt.Println(req)
		// p.AddTargetRequestWithParams(req)
		return false
	})

	query.Find("#pgt .pg .nxt").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			realUrl := "http://bbs.goumin.com/" + url
			logger.Println("[info]find next list page: ", realUrl)
			realUrlTag := "bbslist"
			req := newRequest(realUrlTag, realUrl)
			p.AddTargetRequestWithParams(req)
		}
		return false
	})

}

// func qMBbsList(p *page.Page, num int) {
// 	query := p.GetHtmlParser()
// }

// func qIuserList(p *page.Page, num int) {
// 	query := p.GetHtmlParser()
// }

// func qAskList(p *page.Page, num int) {
// 	query := p.GetHtmlParser()
// }

// func qMaskList(p *page.Page, num int) {
// 	query := p.GetHtmlParser()
// }

// func qMallList(p *page.Page, num int) {
// 	query := p.GetHtmlParser()
// }
