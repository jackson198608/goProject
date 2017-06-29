package main

import (
	"fmt"
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
	//time.Sleep(1 * time.Second)
	if !p.IsSucc() {
		return
	}

	tag := p.GetUrlTag()
	if tag == "searchList" || tag == "searchListNextPage" {
		qBaiduList(p)
	}

	if tag == "gouminDetail" {
		fmt.Println(p.GetHeader())
	}
}

func (this *MyPageProcesser) Finish() {
}
