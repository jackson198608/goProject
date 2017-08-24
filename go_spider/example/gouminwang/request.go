package main

import (
	"github.com/jackson198608/gotest/go_spider/core/common/request"
	"net/http"
)

func newRequest(tag string, url string) *request.Request {
	h := make(http.Header)
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.94 Safari/537.36")
	// h.Add("User-Agent", "Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)")
	h.Add("Connection", "close")
	h.Add("Accept-Encoding", "gzip")

	req := request.NewRequest(url, "html", tag, "GET", "", h, nil, nil, nil)
	return req
}

func newImageRequest(tag string, url string) *request.Request {
	h := make(http.Header)
	h.Add("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.94 Safari/537.36")
	// h.Add("User-Agent", "MQQBrowser/26 Mozilla/5.0 (Linux; U; Android 2.3.7; zh-cn; MB200 Build/GRJ22; CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1")
	// h.Add("User-Agent", "Sogou Pic Spider/3.0(+http://www.sogou.com/docs/help/webmasters.htm#07)")
	h.Add("Connection", "close")
	h.Add("Accept-Encoding", "gzip")
	h.Add("Referer", "http://weixin.sogou.com/")
	// h.Add("Referer", "http://weixin.sogou.com/weixin?query=狗与爱的世界&_sug_type_=&s_from=input&_sug_=n&type=2&page=1&ie=utf8")

	req := request.NewRequest(url, "text", tag, "GET", "", h, nil, nil, nil)
	return req
}
