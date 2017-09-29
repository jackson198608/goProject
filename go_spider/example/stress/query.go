package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"net/http"
	"strings"
	"time"
)

func qBbsList(p *page.Page) {
	query := p.GetHtmlParser()
	fmt.Println(p.GetRequest().Url)
	query.Find("#threadlist .bm_c .new a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, _ := s.Attr("href")
		start := time.Now()
		realUrl := "http://bbs.goumin.com/" + url
		// req := newRequest("bbsview", realUrl)
		// p.AddTargetRequestWithParams(req)
		resp, err := http.Get(realUrl)
		if err != nil {
			// handle error
		}
		//记录结束时间
		end := time.Since(start)
		//输出执行时间，单位为毫秒。
		fmt.Println(realUrl + ",响应时间：") //
		fmt.Println(end)
		fmt.Println(resp)
		return true
	})

	query.Find("#pgt .pg a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		txt := s.Text()
		if txt == "下一页" {
			url, isExsit := s.Attr("href")
			if isExsit {
				realUrl := "http://bbs.goumin.com/" + url
				logger.Println("[info]find next list page: ", realUrl)
				realUrlTag := "bbslist"
				start := time.Now()
				req := newRequest(realUrlTag, realUrl)
				end := time.Since(start)
				fmt.Println(realUrl + ",响应时间：") //
				fmt.Println(end)
				p.AddTargetRequestWithParams(req)
			} else {
				fmt.Println("dfdfd")
			}
		}
		return true
	})

}

func qMallList(p *page.Page) {
	query := p.GetHtmlParser()
	fmt.Println(p.GetRequest().Url)
	query.Find(".main .list_type .dd a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, _ := s.Attr("href")
		start := time.Now()
		realUrl := "http://mall.goumin.com" + url
		// req := newRequest("bbsview", realUrl)
		// p.AddTargetRequestWithParams(req)
		resp, err := http.Get(realUrl)
		if err != nil {
			// handle error
		}
		//记录结束时间
		end := time.Since(start)
		//输出执行时间，单位为毫秒。
		fmt.Println(resp)
		fmt.Println(realUrl + ",响应时间：") //
		fmt.Println(end)
		return true
	})

	query.Find(".main .list_mc .pic a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, _ := s.Attr("href")
		start := time.Now()
		realUrl := "http://mall.goumin.com" + url
		// req := newRequest("bbsview", realUrl)
		// p.AddTargetRequestWithParams(req)
		resp, err := http.Get(realUrl)
		if err != nil {
			// handle error
		}
		//记录结束时间
		end := time.Since(start)
		//输出执行时间，单位为毫秒。
		fmt.Println(realUrl + ",响应时间：") //
		fmt.Println(end)
		fmt.Println(resp)
		return true
	})

	query.Find(".list_mt .page_foot a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		txt := s.Text()
		if txt == "下一页" {
			url, isExsit := s.Attr("href")
			if isExsit {
				realUrl := "http://mall.goumin.com" + url
				logger.Println("[info]find next list page: ", realUrl)
				realUrlTag := "malllist"
				start := time.Now()
				req := newRequest(realUrlTag, realUrl)
				end := time.Since(start)
				fmt.Println(realUrl + ",响应时间：") //
				fmt.Println(end)
				p.AddTargetRequestWithParams(req)
			}
		}
		return true
	})

}

func qAskList(p *page.Page) {
	query := p.GetHtmlParser()
	fmt.Println(p.GetRequest().Url)
	query.Find(".modle-a-rm .rm-list-page .rm-title a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, _ := s.Attr("href")
		urlstr := strings.Split(url, "http")
		realUrl := ""
		if len(urlstr) > 1 {
			realUrl = url
		} else {
			realUrl = "http://www.goumin.com" + url
		}
		start := time.Now()

		// req := newRequest("bbsview", realUrl)
		// p.AddTargetRequestWithParams(req)
		resp, err := http.Get(realUrl)
		if err != nil {
			// handle error
		}
		//记录结束时间
		end := time.Since(start)
		//输出执行时间，单位为毫秒。
		fmt.Println(realUrl + ",响应时间：") //
		fmt.Println(end)
		fmt.Println(resp)
		return true
	})

	query.Find(".modle-a-rm .page a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		txt := s.Text()
		if txt == "下一页" { //
			url, isExsit := s.Attr("href")
			if isExsit {
				realUrl := "http://www.goumin.com" + url
				logger.Println("[info]find next list page: ", realUrl)
				realUrlTag := "asklist"
				start := time.Now()
				req := newRequest(realUrlTag, realUrl)
				end := time.Since(start)
				fmt.Println(realUrl + ",响应时间：") //
				fmt.Println(end)
				p.AddTargetRequestWithParams(req)
			}
		}
		return true
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
