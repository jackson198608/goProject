package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"net/url"
	"strconv"
	"strings"
)

type UrlJson struct {
	Fm      string
	Ensrcid string
	Order   string
	Mu      string
}

func qBaiduList(p *page.Page, num int) {
	query := p.GetHtmlParser()

	fmt.Println(p.GetUrlTag())
	var keyword string = ""
	d, _ := url.Parse(p.GetRequest().Url)
	m, _ := url.ParseQuery(d.RawQuery)
	for k, v := range m {
		if k == "wd" || k == "word" {
			keyword = v[0]
		}
	}

	fmt.Println("[info] maybe  find next page realUrl:", p.GetRequest().Url, p.GetUrlTag(), keyword)
	query.Find(".result").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		logger.Println("get keyword baidu rank : ", keyword)
		datalog, _ := s.Attr("data-log")
		str := strings.Replace(datalog, "'", "\"", -1)
		fmt.Println(str)
		jsonArr := new(UrlJson)
		err := json.Unmarshal([]byte(str), &jsonArr)
		if err != nil {
			logger.Println("json decode error : ", err)
			return false
		}
		realurl := jsonArr.Mu                         //url
		rankReal, err1 := strconv.Atoi(jsonArr.Order) //排序
		if err1 != nil {
			logger.Println("strconv Atoi error : ", err1)
			return false
		}
		var orderRank int = 0
		if rankReal != 11 {
			orderRank = rankReal + (10 * num)
			logger.Println("realurl rank : ", orderRank) //真实rank
		}
		if realurl != "" {
			saveRealUrl(realurl, keyword, orderRank)
			logger.Println("search data list realurl : ", realurl) //真实地址
		}
		fmt.Println("??????????", keyword)
		if realurl == "" && rankReal != 11 {
			fmt.Println("真实地址为空,获取百度地址")
			url, _ := s.Find(".c-container a").Attr("href")
			// url := "http://www.baidu.com/link?url=avo8Hf21hX49dNiUng8uxFsRHA0ZV3tPSl6Y0EKPw84NVrdbthG1j9RAjDUGcxEIICor_6HIzv0GLh8-AniEla"
			fmt.Println("######", url)
			logger.Println("search data list baiduUrl : ", url) //真实地址
			redirectRealUrl := retRequest(url)                  //获取调转后的真实地址
			if redirectRealUrl != "" {
				fmt.Println(">>>>>>", redirectRealUrl)
				saveRealUrl(redirectRealUrl, keyword, orderRank)
			} else {
				realUrlTag := "domainUrl|" + strconv.Itoa(orderRank) + "|" + keyword
				req := newRequest(realUrlTag, url)
				p.AddTargetRequestWithParams(req)
			}
		}
		return true
	})

	// times := 0
	if p.GetUrlTag() == "searchList" {
		query.Find("#page-controller .new-pagenav a").EachWithBreak(func(i int, s *goquery.Selection) bool {
			// For each item found, get the band and title
			url, isExsit := s.Attr("href")
			if isExsit {
				// if times == 4 {
				// 	return false
				// }
				// realUrl := "https://www.baidu.com" + url
				realUrl := url
				num++
				numstr := strconv.Itoa(num)
				realUrlTag := "searchListNextPage|" + numstr
				req := newRequest(realUrlTag, realUrl)
				p.AddTargetRequestWithParams(req)
				// times++
			}
			return true
		})
		return
	}

	hasResult := 0
	// fmt.Println("getrul$$$$$", urlTag)
	//if p.GetUrlTag() == urlTag {
	query.Find(".new-nextpage").EachWithBreak(func(i int, s *goquery.Selection) bool {
		url, isExsit := s.Attr("href")
		if isExsit {
			realUrl := url
			num++
			stringnum := strconv.Itoa(num)
			realUrlTag := "searchListNextPage|" + stringnum
			req := newRequest(realUrlTag, realUrl)
			fmt.Println("[info] find next page realUrl:", realUrl, realUrlTag, keyword)
			hasResult = 1
			p.AddTargetRequestWithParams(req)
		}
		return true
	})

	if hasResult == 0 {
		fmt.Println("[error] can not find next page", p.GetRequest().Url, p.GetBodyStr())
	}

}

//保存关键词搜索结果排名
func saveRealUrl(realurl string, keyword string, rank int) {
	u, _ := url.Parse(realurl)
	domain := u.Host
	if domain == "m.goumin.com" && keyword != "" {
		Id, RankSql, IsExist := checkKeywordExist(keyword)
		if IsExist == true {
			if RankSql > rank {
				updateKeywordRank(Id, keyword, rank, realurl, domain)
			}
		} else {
			saveKeywordRankData(keyword, rank, realurl, domain)
		}
	}
}

//获取真实的搜索结果链接
func saveKeyWordRank(p *page.Page, rank int, keyword string) {
	query := p.GetHtmlParser()
	query.Find("noscript").EachWithBreak(func(i int, s *goquery.Selection) bool {
		str := s.Text()
		strArr := strings.Split(str, "content=\"0; url=")
		if len(strArr) > 1 {
			urlArr := strings.Split(strArr[1], "\"")
			fmt.Println("!!!!!!!!!", urlArr[0])
			saveRealUrl(urlArr[0], keyword, rank)
		}
		return true
	})
}
