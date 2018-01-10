package abuyunHttpClient

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type AbuyunProxy struct {
	proxyServer string
	appID       string
	appSecret   string
	client      *http.Client
}

func NewAbuyunProxy(proxyServer string, appID string, appSecret string) *AbuyunProxy {
	abuyunProxy := new(AbuyunProxy)
	if abuyunProxy == nil {
		fmt.Println("[error]create fail")
		return nil
	}

	abuyunProxy.appID = appID
	abuyunProxy.appSecret = appSecret
	abuyunProxy.proxyServer = proxyServer
	if (proxyServer == "") || (appID == "") || (appSecret == "") {
		abuyunProxy.makeClientWithoutProxy()
	} else {

		abuyunProxy.makeClient()
	}

	return abuyunProxy

}

func abuyunHttpClient(req *http.Request, via []*http.Request) error {
	return errors.New("this url has redirect" + req.URL.String())
}

func (p *AbuyunProxy) makeClient() {
	proxyUrl, err := url.Parse("http://" + p.appID + ":" + p.appSecret + "@" + p.proxyServer)
	if err != nil {
		fmt.Println("[error]url parse error", p.appID, p.appSecret, p.proxyServer)
	}

	p.client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, Timeout: 5 * time.Second}
	p.client.CheckRedirect = abuyunHttpClient
}

func (p *AbuyunProxy) makeClientWithoutProxy() {
	p.client = &http.Client{Timeout: 5 * time.Second}
	p.client.CheckRedirect = abuyunHttpClient
}

func (p *AbuyunProxy) Close() {
}

// todo: make customHeader to be pointer
func (p *AbuyunProxy) SendRequest(targetUrl string, customHeader http.Header, requestBody string, switchip bool) (int, *http.Header, string, error) {
	request, err := http.NewRequest("GET", targetUrl, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return 0, nil, "", err
	}

	if switchip {
		request.Header.Set("Proxy-Switch-Ip", "yes")
	}

	if len(customHeader) != 0 {
		for k, v := range customHeader {
			//@todo v can be a list to add on header
			request.Header.Set(k, v[0])
			if k == "Host" {
				request.Host = v[0]
			}
		}
	}

	response, err := p.client.Do(request)

	if err != nil {
		return 0, nil, "", err
	} else {
		bodyByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return 0, nil, "", err
		}
		response.Body.Close()

		body := string(bodyByte)

		return response.StatusCode, &response.Header, body, nil

	}

}
