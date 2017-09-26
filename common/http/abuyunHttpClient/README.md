# introduction
- this package is design for the situation when you want to send the request using abuyun proxy

# Installation

	go get github.com/jackson198608/goProject/common/http/abuyunHttpClient 

# Quick Start

- Create Abuyun

```Go
var abuyun *AbuyunProxy = NewAbuyunProxy(proxyServer,
    proxyUser,
    proxyPasswd)
```

- User Abuyun to send the request 

```Go
targetUrl := "http://m.goumin.com/"

var h http.Header
statusCode, responseHeader, _, err := abuyun.SendRequest(targetUrl, h, true)
```
- targetUrl: the page you want to request
- h: the customHeadr you want to add
- swichip or not
