package abuyunHttpClient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AbuyunProxy struct {
	proxyServer string
	appID       string
	appSecret   string
	client      http.Client
}

func NewAbuyunProxy(proxyServer string, appID string, appSecret string) *AbuyunProxy {
	if (proxyServer == "") || (appID == "") || (appSecret == "") {
		return nil
	}
	var abuyunProxy AbuyunProxy = AbuyunProxy{proxyServer: proxyServer,
		appID:     appID,
		appSecret: appSecret,
		client:    http.Client{},
	}
	abuyunProxy.makeClient()
	return &abuyunProxy

}

func (p AbuyunProxy) makeClient() {
	proxyUrl, err := url.Parse("http://" + p.appID + ":" + p.appSecret + "@" + p.proxyServer)
	if err != nil {
		fmt.Println("[error]url parse error", p.appID, p.appSecret, p.proxyServer)
	}

	p.client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
}

func (p AbuyunProxy) SendRequest(targetUrl string, customHeader http.Header, switchip bool) (int, *http.Header, string, error) {
	request, err := http.NewRequest("GET", targetUrl, bytes.NewBuffer([]byte(``)))
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
