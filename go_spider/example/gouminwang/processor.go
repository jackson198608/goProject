package main

import (
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	//"time"
	"fmt"
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
		logger.Println("[Error]not 200: ", p.GetRequest().Url)
		return
	}

	tag := p.GetUrlTag()
	fmt.Println(tag)
	strArr := strings.Split(tag, "|")
	if strArr[0] == "articleDetail" {
		logger.Println("[info]article detail page:", p.GetRequest().Url)
		//save shop into mysql
		success := saveArticleDetail(p)
		if success {
			logger.Println("[info]save article detail success")
		}
	} else if strArr[0] == "articleList" {
		logger.Println("[info]find article list by tag : ", tag, p.GetRequest().Url)
		num, _ := strconv.Atoi(strArr[1])
		fmt.Println(num)
		if num < 2 {
			qArticleList(p, num)
		}

	} else if strArr[0] == "shopImage" {
		logger.Println("[info]find article list by tag : ", tag, p.GetRequest().Url)
		saveImage(p)

	}
}

func (this *MyPageProcesser) Finish() {
}
