package main

import (
	"fmt"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"strconv"
	"strings"
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

	if tag == "searchList" {
		qBaiduList(p, 0)
	}
	strArr := strings.Split(tag, "|")
	if strArr[0] == "searchListNextPage" {
		if len(strArr) == 2 {
			num, _ := strconv.Atoi(strArr[1])
			fmt.Println(num)
			if num < 10 {
				qBaiduList(p, num)
			}
		}
	}

	if strArr[0] == "domainUrl" {
		rank, _ := strconv.Atoi(strArr[1])
		saveKeyWordRank(p, rank, strArr[2])
	}
}

func (this *MyPageProcesser) Finish() {
}
