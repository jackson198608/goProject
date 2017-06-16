package main

import (
	"github.com/jackson198608/gotest/go_spider/core/common/request"
	"net/http"
)

func newRequest(tag string, url string) *request.Request {
	h := make(http.Header)
	h.Add("User-Agent", "Baiduspider")
	h.Add("Connection", "close")
	h.Add("Accept-Encoding", "gzip")

	req := request.NewRequest(url, "html", tag, "GET", "", h, nil, nil, nil)
	return req
}

func newImageRequest(tag string, url string) *request.Request {
	h := make(http.Header)
	h.Add("User-Agent", "Baiduspider")
	h.Add("Connection", "close")
	h.Add("Accept-Encoding", "gzip")
	h.Add("Referer", "http://www.jd.com")

	req := request.NewRequest(url, "text", tag, "GET", "", h, nil, nil, nil)
	return req
}
