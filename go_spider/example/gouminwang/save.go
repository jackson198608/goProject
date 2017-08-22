package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
	"os"
	"path"
	"regexp"
	"strings"
)

func getArticleTitle(query *goquery.Document, title *string) {
	query.Find("#page-content .rich_media_title").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			*title = strings.Trim(s.Text(), "\n ")
		}
		return true
	})
}

func getArticleDateline(query *goquery.Document, dateline *string) {
	query.Find("#page-content #meta_content .rich_media_meta_text").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			*dateline = s.Text()
		}
		return true
	})
}

func getArticleAuthor(query *goquery.Document, content *string) {
	query.Find("#page-content #meta_content .rich_media_meta_nickname").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			*content = s.Text()
		}
		return true
	})
}

func getArticleContent(query *goquery.Document, content *string) {
	query.Find("#page-content .rich_media_content ").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			// *content = s.Text()
			*content, _ = s.Html()
			// *content = strings.Replace(*content, "<img.*?src=""([^""]*)"".*?>", ".", -1)
			// *content = strings.Replace(*content, "?wx_fmt=", ".", -1)
			// *content = strings.Replace(*content, "data-src", "src", -1)
			// *content = strings.Replace(*content, "http://mmbiz.qpic.cn", imgUrl, -1)
			// *content = strings.Replace(*content, "https://mmbiz.qlogo.cn", imgUrl, -1)
			// re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
			re, _ := regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
			*content = re.ReplaceAllString(*content, "")
			re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
			*content = re.ReplaceAllString(*content, "")
			// var str = "slfjslkdfjsldkfj<img src='dfd' />"
			// str := string(*content)
			lenstr := len(*content)
			contentBytes := []byte(*content)

			start := 0
			end := 0

			for {
				a := strings.Index(string(contentBytes), "<img")
				if a < 0 {
					contentBytes[start] = '<'
					contentBytes[end] = '>'
					break
				}

				contentBytes[a] = '['
				for i := a; i < lenstr; i++ {
					if contentBytes[i] == '>' {
						contentBytes[i] = ']'
						start = a
						end = i
						break
					}
				}
			}

			*content = string(contentBytes)

			re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
			*content = re.ReplaceAllString(*content, "\n")
			re, _ = regexp.Compile("\\s{2,}")
			*content = re.ReplaceAllString(*content, "\n")

			lenstr = len(*content)
			contentBytes = []byte(*content)
			for {
				a := strings.Index(string(contentBytes), "[img")
				if a < 0 {
					break
				}
				contentBytes[a] = '<'
				for i := a; i < lenstr; i++ {
					if contentBytes[i] == ']' {
						contentBytes[i] = '>'
						break
					}
				}
			}

			*content = string(contentBytes)
			fmt.Println("[info]", *content)
		}
		return true
	})
}

func getArticleImages(p *page.Page) {
	query := p.GetHtmlParser()
	query.Find("#page-content .rich_media_content img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		var image *string = new(string)
		*image, _ = s.Attr("data-src")
		*image = strings.Replace(*image, "?wx_fmt=", ".", -1)
		fmt.Println(*image)
		//生成对图片的抓取任务
		req := newImageRequest("articleImage", *image)
		p.AddTargetRequestWithParams(req)
		return true
	})
}

func saveImage(p *page.Page) bool {

	url := p.GetRequest().Url
	//get fullpath
	abPath := getPathFromUrl(url)
	fmt.Println("^^^^^")
	fullPath := saveDir + abPath
	fmt.Println(fullPath)
	fullDirPath := path.Dir(fullPath)
	err := os.MkdirAll(fullDirPath, 0777)
	if err != nil {
		logger.Println("[error]create dir error:", err, " ", fullDirPath, " ", url)
		return false
	}

	//save file
	result, err := os.Create(fullPath)
	if err != nil {
		logger.Println("[error]create file error:", err, " ", fullPath, " ", url)
		return false
	}
	logger.Println("[info] save image:", url)
	logger.Println("[info] save in:", fullPath)
	logger.Println("[info] save len:", len(p.GetBodyStr()))
	result.WriteString(p.GetBodyStr())
	logger.Println("[info] save iamge WriteString:", len(p.GetBodyStr()))
	result.Close()

	return true
}

func getPathFromUrl(url string) string {
	//find first ? position
	var i int = 0
	len := len(url)
	for i = 0; i < len; i++ {
		if url[i] == '?' {
			break
		}
	}
	if i == (len - 1) {
		//have no ?
		return url
	} else {
		path := url[7:i]
		return string(path)
	}

}

func saveArticleDetail(p *page.Page) bool {
	sourceUrl := p.GetRequest().Url
	query := p.GetHtmlParser()
	// fmt.Println(p)
	// 文章标题
	var title *string = new(string)
	getArticleTitle(query, title)
	logger.Println("[info] article title: ", *title)

	var dateline *string = new(string)
	getArticleDateline(query, dateline)

	var author *string = new(string)
	getArticleAuthor(query, author)
	logger.Println("[info] article author: ", *author)

	var content *string = new(string)
	getArticleContent(query, content)

	// getArticleImages(p)
	var status bool = false
	if *title != "" {
		status = insertArticleDetail(*title, *dateline, *author, *content)
	} else {
		logger.Println("[info]again find article detail: ", sourceUrl)
		realUrlTag := "articleDetail"
		req := newRequest(realUrlTag, sourceUrl)
		p.AddTargetRequestWithParams(req)
	}
	return status
}
