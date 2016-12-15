package main

import (
	"github.com/hu17889/go_spider/core/common/request"
	"net/http"
)

func newRequest(tag string, url string) *request.Request {
	h := make(http.Header)
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.94 Safari/537.36")
	h.Add("Connection", "close")
	h.Add("Accept-Encoding", "gzip")

	req := request.NewRequest(url, "html", tag, "GET", "", h, nil, nil, nil)
	return req
}

func newImageRequest(tag string, url string) *request.Request {
	h := make(http.Header)
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.94 Safari/537.36")
	h.Add("Connection", "close")
	h.Add("Accept-Encoding", "gzip")
	h.Add("Referer", "http://bbs.aigou.com/bbs/post/view/556_124176059_1__1_0.html")

	req := request.NewRequest(url, "text", tag, "GET", "", h, nil, nil, nil)
	return req
}
