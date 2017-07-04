package main

import (
	"errors"
	"github.com/jackson198608/gotest/go_spider/core/common/request"
	"net/http"
	"net/url"
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
	h.Add("Referer", "http://www.dianping.com/shop/38076482")

	req := request.NewRequest(url, "text", tag, "GET", "", h, nil, nil, nil)
	return req
}

var redirectCount int = 0

func myRedirect(req *http.Request, via []*http.Request) (e error) {
	redirectCount++
	if redirectCount == 2 {
		redirectCount = 0
		return errors.New(req.URL.String())

	}
	return
}

//redirect
func retRequest(urlstr string) string {
	client := &http.Client{CheckRedirect: myRedirect}
	response, err := client.Get(urlstr)
	if err != nil {
		if e, ok := err.(*url.Error); ok && e.Err != nil {
			remoteUrl := e.URL
			return remoteUrl
		}
	}
	defer response.Body.Close()
	return ""
}
