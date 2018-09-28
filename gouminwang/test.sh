work/src$ go install github.com/jackson198608/goProject/go_spider/example/gouminwang
work/src$ ../bin/gouminwang "/data/gouminwang/" "http://mp.weixin.qq.com/profile?src=3&timestamp=1502248336&ver=1&signature=fQiCz0fa59OB5qMVPTRjp9OrKK39L8W8jW*ypNWocuzQ4O52NzCyrm641yYqwBpX9nQm5mWV*OTbPZYwFHvSmA==" "articleList|0" "/tmp/1.log" 1
work/src$ ../bin/gouminwang "/data/gouminwang/" "http://www.dianping.com/search/category/2/10/g110" "articleList|0" "/tmp/1.log" 1
work/src$ ../bin/gouminwang "/data/gouminwang/" "http://weixin.sogou.com/weixin?type=2&query=狗与爱的世界&ie=utf8&s_from=input&_sug_=n&_sug_type_=&w=01019900&sut=2593&sst0=1502269098653&lkt=1%2C1502269098550%2C1502269098550" "articleList|0" "/tmp/1.log" 1
work/src$ ../bin/gouminwang "http://weixin.sogou.com/weixin?type=2&query=狗与爱的世界&ie=utf8&s_from=input&_sug_=n&_sug_type_=&w=01019900&sut=2593&sst0=1502269098653&lkt=1%2C1502269098550%2C1502269098550" "articleList|0"

work/src$ ../bin/gouminwang /etc/gouminwang.json

