package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackson198608/gotest/go_spider/core/common/page"
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
		}
		return true
	})
}

// func getArticleImages(query *goquery.Document, content *string) {
// 	query.Find("#page-content .rich_media_content ").EachWithBreak(func(i int, s *goquery.Selection) bool {
// 		var images []string
// 			// *content = s.Text()
// 			*content, _ = s.Html()
// 		}
// 		fmt.Println(*content)
// 		return true
// 	})
// }

func saveArticleDetail(p *page.Page) bool {
	sourceUrl := p.GetRequest().Url
	query := p.GetHtmlParser()
	fmt.Println(p)
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
	fmt.Println(content)
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
