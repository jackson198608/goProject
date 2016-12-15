package pipeline

import (
    "github.com/jackson198608/gotest/go_spider/core/common/com_interfaces"
    "github.com/jackson198608/gotest/go_spider/core/common/page_items"
)

type PipelineConsole struct {
}

func NewPipelineConsole() *PipelineConsole {
    return &PipelineConsole{}
}

func (this *PipelineConsole) Process(items *page_items.PageItems, t com_interfaces.Task) {
    println("----------------------------------------------------------------------------------------------")
    println("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
    println("Crawled result : ")
    for key, value := range items.GetAll() {
        println(key + "\t:\t" + value)
    }
}
