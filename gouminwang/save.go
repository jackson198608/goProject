package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/goProject/go_spider/core/common/page"
	// "io/ioutil"
	// "net/http"
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
	query.Find("#page-content #meta_content #post-date").EachWithBreak(func(i int, s *goquery.Selection) bool {
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
			*content = strings.Replace(*content, "?wx_fmt=", ".", -1)
			*content = strings.Replace(*content, "&amp;tp=webp&amp;wxfrom=5&amp;wx_lazy=1", "", -1)
			*content = strings.Replace(*content, "&wxfrom=5&wx_lazy=1&tp=webp", "", -1)
			*content = strings.Replace(*content, "?tp=webp&wxfrom=5&wx_lazy=1", "", -1)
			*content = strings.Replace(*content, "?tp=webp&wxfrom=5", "", -1)
			*content = strings.Replace(*content, "&tp=webp&wxfrom=5&wx_lazy=1", "", -1)
			*content = strings.Replace(*content, "&tp=webp&wxfrom=5", "", -1)
			*content = strings.Replace(*content, "&wxfrom=5&tp=webp", "", -1)
			*content = strings.Replace(*content, "&amp;wxfrom=5&amp;tp=webp", "", -1)
			*content = strings.Replace(*content, "&amp;wxfrom=5", "", -1)
			*content = strings.Replace(*content, "&wxfrom=5", "", -1)
			*content = strings.Replace(*content, "&amp;tp=webp", "", -1)
			*content = strings.Replace(*content, "&tp=webp", "", -1)
			// *content = strings.Replace(*content, "data-src", "src", -1)
			*content = strings.Replace(*content, "http://mmbiz.qpic.cn", imgUrl+"/mmbiz.qpic.cn", -1)
			*content = strings.Replace(*content, "https://mmbiz.qpic.cn", imgUrl+"/mmbiz.qpic.cn", -1)
			*content = strings.Replace(*content, "https://mmbiz.qlogo.cn", imgUrl+"/mmbiz.qlogo.cn", -1)
			re, _ := regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
			*content = re.ReplaceAllString(*content, "")
			re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
			*content = re.ReplaceAllString(*content, "")

			/*lenstr := len(*content)
			contentBytes := []byte(*content)

			start := 0
			end := 0
			a := 0
			for {
				// aa := strings.IndexRune(string(contentBytes), "<img")
				a = strings.Index(string(contentBytes), "<img")
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
			b := strings.Index(string(contentBytes), "[img")
			if b > 0 {
				contentBytes[b] = '<'
				for j := b; j < lenstr; j++ {
					if contentBytes[j] == ']' {
						contentBytes[j] = '>'
						break
					}
				}
			}
			*content = string(contentBytes)
			// fmt.Println("[info]", *content)


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
			*content = string(contentBytes)*/
			*content = replaceHtmlImg(content)
			fmt.Println("[info]", *content)
		}
		return true
	})
}

func replaceHtmlImg(content *string) string {
	lenstr := len(*content)
	contentBytes := []byte(*content)

	start := 0
	end := 0
	a := 0
	for {
		// aa := strings.IndexRune(string(contentBytes), "<img")
		a = strings.Index(string(contentBytes), "<img")
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
	b := strings.Index(string(contentBytes), "[img")
	if b > 0 {
		contentBytes[b] = '<'
		for j := b; j < lenstr; j++ {
			if contentBytes[j] == ']' {
				contentBytes[j] = '>'
				break
			}
		}
	}
	*content = string(contentBytes)
	// fmt.Println("[info]", *content)

	// *content = replaceHtmlImg(*content)
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
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

	scontent := string(contentBytes)
	return scontent
}

func getArticleImages(p *page.Page) {
	query := p.GetHtmlParser()
	query.Find("#page-content .rich_media_content img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		var image *string = new(string)
		*image, _ = s.Attr("data-src")
		*image = strings.Replace(*image, "?wx_fmt=", ".", -1)
		*image = strings.Replace(*image, "&amp;tp=webp&amp;wxfrom=5&amp;wx_lazy=1", "", -1)
		*image = strings.Replace(*image, "&wxfrom=5&wx_lazy=1&tp=webp", "", -1)
		*image = strings.Replace(*image, "?tp=webp&wxfrom=5&wx_lazy=1", "", -1)
		*image = strings.Replace(*image, "?tp=webp&wxfrom=5", "", -1)
		*image = strings.Replace(*image, "&tp=webp&wxfrom=5&wx_lazy=1", "", -1)
		*image = strings.Replace(*image, "&tp=webp&wxfrom=5", "", -1)
		*image = strings.Replace(*image, "&wxfrom=5&tp=webp", "", -1)
		*image = strings.Replace(*image, "&amp;wxfrom=5&amp;tp=webp", "", -1)
		*image = strings.Replace(*image, "&amp;wxfrom=5", "", -1)
		*image = strings.Replace(*image, "&wxfrom=5", "", -1)
		*image = strings.Replace(*image, "&amp;tp=webp", "", -1)
		*image = strings.Replace(*image, "&tp=webp", "", -1)
		// fmt.Println(*image)
		//生成对图片的抓取任务
		// req := newImageRequest("shopImage", "http://www.testing.com:89/imgBridge.php?url="+*image)
		req := newImageRequest("shopImage", *image)
		p.AddTargetRequestWithParams(req)
		return true
	})
}

func saveImage(p *page.Page) bool {

	url := p.GetRequest().Url
	fmt.Println("^^^^^" + url)
	//get fullpath
	// realurl := strings.Split(url, "imageUrl=")
	// fmt.Println(realurl[0])
	abPath := getPathFromUrl(url)
	fullPath := saveDir + abPath
	fullDirPath := path.Dir(fullPath)
	fmt.Println(fullPath)
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

	getArticleImages(p)
	var status bool = false
	if *title != "" {
		status = insertArticleDetail(*title, *dateline, *author, *content, sourceUrl)
	} else {
		logger.Println("[info]again find article detail: ", sourceUrl)
		realUrlTag := "articleDetail"
		req := newRequest(realUrlTag, sourceUrl)
		p.AddTargetRequestWithParams(req)
	}
	return status
}
