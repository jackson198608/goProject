package abuyunHttpClient

import (
	"fmt"
	"net/http"
	"testing"
)

const proxyServer = ""
const proxyUser = ""
const proxyPasswd = ""

func TestAbuyun(t *testing.T) {
	var abuyun *AbuyunProxy = NewAbuyunProxy(proxyServer,
		proxyUser,
		proxyPasswd)

	t.Log("begin the test")

	if abuyun == nil {
		t.Error("create abuyun error")
		return
	}

	targetUrl := "http://m.goumin.com/ask/47693.html"

	var h http.Header = make(http.Header)
	h.Set("a", "1")
	statusCode, responseHeader, _, err := abuyun.SendRequest(targetUrl, h, `{"query": {"query_string":{"query":"假証"}}}`, true)
	if err != nil {
		t.Log("http request error", err)
		t.Fail()
	}
	fmt.Println(statusCode)
	fmt.Println(responseHeader)
	//fmt.Println(body)
	fmt.Println(err)

}

func TestAbuyunPost(t *testing.T) {
	var abuyun *AbuyunProxy = NewAbuyunProxy(proxyServer,
		proxyUser,
		proxyPasswd)

	t.Log("begin the test")

	if abuyun == nil {
		t.Error("create abuyun error")
		return
	}

	targetUrl := "http://lingdang.goumin.com/v4/travels"

	var h http.Header = make(http.Header)
	h.Set("a", "1")
	statusCode, responseHeader, body, err := abuyun.SendPostRequest(targetUrl, h, `{"token":"waety80tw344gzyk7skd63ewheu2q7od5svtpnbde47ksuuqd3uawyx30ivheai8","uid":"881050","ver":"2.0","seqnum":"GMCAMERAARD13891521453876401605642603","data":{}}`, true)
	if err != nil {
		t.Log("http request error", err)
		t.Fail()
	}
	fmt.Println(statusCode)
	fmt.Println(responseHeader)
	fmt.Println(body)
	fmt.Println(err)

}

func TestAbuyunNoProxy(t *testing.T) {
	var abuyun *AbuyunProxy = NewAbuyunProxy("", "", "")

	t.Log("begin the test")

	if abuyun == nil {
		t.Error("create abuyun error")
		return
	}

	targetUrl := "http://m.goumin.com/ask/47693.html"

	var h http.Header = make(http.Header)
	h.Set("a", "1")
	statusCode, responseHeader, body, err := abuyun.SendRequest(targetUrl, h, `{"query": {"query_string":{"query":"假証"}}}`, true)
	if err != nil {
		t.Log("http request error", err)
		t.Fail()
	}
	fmt.Println(statusCode)
	fmt.Println(responseHeader)
	fmt.Println(body)
	fmt.Println(err)
}
