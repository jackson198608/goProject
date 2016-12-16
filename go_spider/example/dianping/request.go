package main

import (
	"github.com/jackson198608/gotest/go_spider/core/common/request"
	"net/http"
	// "fmt"
)

func newRequest(tag string, url string) *request.Request {
	h := make(http.Header)
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.94 Safari/537.36")
	h.Add("Connection", "close")
	h.Add("Accept-Encoding", "gzip")

	ProxyIp := LRange()
	// req := request.NewRequest(url, "html", tag, "GET", "", h, nil, nil, nil)
	req := request.NewRequestWithProxy(url, "html", tag, "GET", "", h, nil, ProxyIp, nil, nil) //proxy
	return req
}

func newImageRequest(tag string, url string) *request.Request {
	h := make(http.Header)
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.94 Safari/537.36")
	h.Add("Connection", "close")
	h.Add("Accept-Encoding", "gzip")
	h.Add("Referer", "http://www.dianping.com/shop/38076482")

	req := request.NewRequest(url, "text", tag, "GET", "", h, nil, nil, nil)
	return req
}
