package main

import (
	"github.com/jackson198608/gotest/go_spider/core/common/page"
)

type MyPageProcesser struct {
}

func NewMyPageProcesser() *MyPageProcesser {
	return &MyPageProcesser{}
}

// Parse html dom here and record the parse result that we want to Page.
// Package goquery (http://godoc.org/github.com/PuerkitoBio/goquery) is used to parse html.
func (this *MyPageProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		logger.Println("[Error]not 200: ", p.GetRequest().Url)
		return
	}

	tag := p.GetUrlTag()
	switch tag {
	case "threadForumQuanzhong":
		logger.Println("[info]forum quanzhong page:", p.GetRequest().Url)
		qTheadForumQuanzhong(p)
	case "threadForumDifang":
		logger.Println("[info]forum difang page:", p.GetRequest().Url)
		qTheadForumDifang(p)
	case "threadForumZonghe":
		logger.Println("[info]forum zonghe page:", p.GetRequest().Url)
		qTheadForumZonghe(p)
	case "threadList":
		logger.Println("[info]list page:", p.GetRequest().Url)
		qTheadList(p)
	case "threadDetail":
		logger.Println("[info]detail page:", p.GetRequest().Url)
		save(p)
		qTheadDetail(p)
	case "image":
		logger.Println("[info]image for detail page:", p.GetRequest().Url)
		saveImage(p)
	case "askList":
		logger.Println("[info]ask list page:", p.GetRequest().Url)
		qAskList(p)
	case "askDetail":
		logger.Println("[info]ask Detail page:", p.GetRequest().Url)
		save(p)
	}
}

func (this *MyPageProcesser) Finish() {
}
