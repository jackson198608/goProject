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
		logger.Println("[Error] not 200:", p.GetRequest().Url)
		fmt.Println("[Error] not 200:", p.GetRequest().Url)
		return
	}

	tag := p.GetUrlTag()
	fmt.Println(tag)
	if tag == "bbslist" {
		qBbsList(p)
	}
	if tag == "malllist" {
		qMallList(p)
	}
	if tag == "asklist" {
		qAskList(p)
	}
	if tag == "masklist" {
		qMaskList(p)
	}
	// if tag == "bbsview" {
	// 	qBbsView(p)
	// }
	// if tag == "mbbslist"{
	//     qMBbsList()
	// }
}

func (this *MyPageProcesser) Finish() {
}
