package main

import (
	"github.com/jackson198608/gotest/go_spider/core/pipeline"
	"github.com/jackson198608/gotest/go_spider/core/spider"
	"net/url"
	"os"
)

var result *os.File

func onstart() {
	fileName := os.Args[1]
	resultop, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		panic("can not open file")
	}
	result = resultop

}

func main() {
	onstart()

	keyword := os.Args[2]
	realKeyWord := keyword + "  site:bbs.goumin.com"

	tempUrl := "http://1.com/?wd=" + realKeyWord
	tempUrlP, _ := url.Parse(tempUrl)
	realKeyWordEncode := tempUrlP.Query().Encode()
	result.WriteString(realKeyWordEncode + "\n")

	url1 := "https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&"
	url2 := "&rsv_pq=e04aca4d0000cf55&rsv_t=dabboKw3o5qu1XZAgI43hhbjd2olBB3puS%2Fqgn7abC1zNtc%2BA4jzjem3%2BEI&rqlang=cn&rsv_enter=1&rsv_sug3=28&rsv_sug1=17&rsv_sug7=100&rsv_sug2=0&inputT=145353&rsv_sug4=146135"

	startUrlTag := "searchList"
	//startUrl := "https://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=%E9%87%91%E6%AF%9B%20site%3Abbs.goumin.com&rsv_pq=e04aca4d0000cf55&rsv_t=dabboKw3o5qu1XZAgI43hhbjd2olBB3puS%2Fqgn7abC1zNtc%2BA4jzjem3%2BEI&rqlang=cn&rsv_enter=1&rsv_sug3=28&rsv_sug1=17&rsv_sug7=100&rsv_sug2=0&inputT=145353&rsv_sug4=146135"
	startUrl := url1 + realKeyWordEncode + url2
	req := newRequest(startUrlTag, startUrl)
	spider.NewSpider(NewMyPageProcesser(), "getBaiduResult").
		AddRequest(req).
		AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		SetThreadnum(5).                            // Crawl request by three Coroutines
		Run()
}
