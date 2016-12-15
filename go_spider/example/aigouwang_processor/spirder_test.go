package main

import (
	"github.com/hu17889/go_spider/core/pipeline"
	"github.com/hu17889/go_spider/core/spider"
	"testing"
)

func TestDetailNextPage(t *testing.T) {
	// Spider input:
	//  PageProcesser ;
	//  Task name used in Pipeline for record;
	req_url := "http://bbs.aigou.com/bbs/post/view/552_121328603_1__1_-1.html"
	req := newRequest("threadDetail", req_url)
	spider.NewSpider(NewMyPageProcesser(), "TaskName").
		AddRequest(req).
		AddPipeline(pipeline.NewPipelineConsole()). // Print result on screen
		SetThreadnum(3).                            // Crawl request by three Coroutines
		Run()

}
