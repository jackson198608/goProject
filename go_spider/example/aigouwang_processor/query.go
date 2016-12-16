package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
)

func qTheadForumQuanzhong(p *page.Page) {
	query := p.GetHtmlParser()

	//find all the list detail page
	query.Find(".tab_con_item ").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		logger.Println("[info]each num: ", i)

		if i == 0 {
			s.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
				url, isExsit := s.Attr("href")
				if isExsit {
					logger.Println("[info]find brand: ", url)
					req := newRequest("threadList", url)
					p.AddTargetRequestWithParams(req)
				}
				return true
			})
		}
		return true
	})

}

func qTheadForumDifang(p *page.Page) {
	query := p.GetHtmlParser()

	//find all the list detail page
	query.Find(".tab_con_item ").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		logger.Println("[info]each num: ", i)

		if i == 1 {
			s.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
				url, isExsit := s.Attr("href")
				if isExsit {
					logger.Println("[info]find brand: ", url)
					req := newRequest("threadList", url)
					p.AddTargetRequestWithParams(req)
				}
				return true
			})
		}
		return true
	})

}

func qTheadForumZonghe(p *page.Page) {
	query := p.GetHtmlParser()

	//find all the list detail page
	query.Find(".tab_con_item ").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		logger.Println("[info]each num: ", i)

		if i == 2 {
			s.Find("a").EachWithBreak(func(i int, s *goquery.Selection) bool {
				url, isExsit := s.Attr("href")
				if isExsit {
					logger.Println("[info]find brand: ", url)
					req := newRequest("threadList", url)
					p.AddTargetRequestWithParams(req)
				}
				return true
			})
		}
		return true
	})

}

func qTheadList(p *page.Page) {

	query := p.GetHtmlParser()

	//find all the list detail page
	query.Find(".all .one .oneRt a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			req := newRequest("threadDetail", url)
			p.AddTargetRequestWithParams(req)
		}
		return true
	})

	//find the next list page
	query.Find(".pagebox a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			text := s.Text()
			if text == "下一页" {
				req := newRequest("threadList", url)
				p.AddTargetRequestWithParams(req)
				return false
			}
		}
		return true
	})

}

func qTheadDetail(p *page.Page) {
	query := p.GetHtmlParser()
	query.Find(".pagebox a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			text := s.Text()
			if text == "下一页" {
				req := newRequest("threadDetail", url)
				p.AddTargetRequestWithParams(req)
				return false
			}
		}
		return true
	})

	//find all the image arrear in post
	query.Find(".forum-text img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("src")
		if isExsit {
			req := newImageRequest("image", url)
			p.AddTargetRequestWithParams(req)
			logger.Println("[info]image: ", url)
		}
		return true
	})

}

func qAskList(p *page.Page) {

	query := p.GetHtmlParser()

	//find all the list detail page
	query.Find(".text-info a").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		if isExsit {
			logger.Println("[info]find detail page: ", url)
			req := newRequest("askDetail", url)
			p.AddTargetRequestWithParams(req)
		}
		return true
	})

	//find the next list page
	query.Find(".next-right").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		url, isExsit := s.Attr("href")
		url = "http://zhidao.aigou.com/" + url
		if isExsit {
			logger.Println("[info]find next page: ", url)
			req := newRequest("askList", url)
			p.AddTargetRequestWithParams(req)
			return false
		}
		return true
	})

}

func qAskDetail(p *page.Page) {
}
