package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const ProxyServer = "proxy.abuyun.com:9010"

type AbuyunProxy struct {
	AppID     string
	AppSecret string
}

func (p AbuyunProxy) ProxyClient() http.Client {
	proxyURL, _ := url.Parse("http://" + p.AppID + ":" + p.AppSecret + "@" + ProxyServer)
	return http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
}

func main() {
	targetURI := "http://www.dianping.com/search/category/1/95/g25147"
	//targetURI := "https://www.abuyun.com/switch-ip"
	//targetURI := "https://www.abuyun.com/current-ip"

	// 初始化 proxy http client
	client := AbuyunProxy{AppID: "HK71T41EZ21304GP", AppSecret: "75FE0C4E23EEA0E7"}.ProxyClient()

	request, _ := http.NewRequest("GET", targetURI, bytes.NewBuffer([]byte(``)))

	// 设置IP切换头 (只支持 HTTP)
	request.Header.Set("Proxy-Switch-Ip", "yes")

	response, err := client.Do(request)

	if err != nil {
		panic("failed to connect: " + err.Error())
	} else {
		bodyByte, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("读取 Body 时出错", err)
			return
		}
		response.Body.Close()

		body := string(bodyByte)

		fmt.Println("Response Status:", response.Status)
		fmt.Println("Response Header:", response.Header)
		fmt.Println("Response Body:\n", body)
	}
}
