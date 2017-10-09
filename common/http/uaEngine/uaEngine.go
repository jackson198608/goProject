package uaEngine

import (
	"math/rand"
	"time"
)

var mobileUserUa [5]string = [5]string{"Mozilla/5.0 (Linux; U; Android 6.0.1; zh-cn; OPPO R9sk Build/MMB29M) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/38.0.0.0 Mobile Safari/537.36 OppoBrowser/4.2.8",
	"Mozilla/5.0 (Linux; U; Android 6.0.1; zh-cn; OPPO R9sk Build/MMB29M) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/38.0.0.0 Mobile Safari/537.36 OppoBrowser/4.2.8",
	"Mozilla/5.0 (Linux; Android 6.0; vivo Y67A Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/38.0.0.0 Mobile Safari/537.36 VivoBrowser/5.0.10",
	"Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/8.0 /6.1 Mobile/15A372 Safari/600.1.4",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_2 like Mac OS X; zh-CN) AppleWebKit/537.51.1 (KHTML, like Gecko) Mobile/14A456 UCBrowser/11.3.0.895 Mobile  AliApp(TUnionSDK/0.1.6)"}

var pcUserUa [5]string = [5]string{"Mozilla/4.0 (compatible; MSIE 9.0; Windows NT 6.1)",
	"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.132 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:10.0.3) Gecko/20100101 Firefox/10.0.3",
	"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.115",
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 1.1.4322)"}

type UaEngine struct {
	currentUa string
}

func NewUaEngine(currentUa string) *UaEngine {
	uaEngine := new(UaEngine)
	uaEngine.currentUa = currentUa
	return uaEngine
}

func (u *UaEngine) GetPcRandomeEngine() string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(5)
	u.currentUa = pcUserUa[index]
	return u.currentUa
}

func (u *UaEngine) GetMobileRandomeEngine() string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(5)
	u.currentUa = mobileUserUa[index]
	return u.currentUa
}
