work/src$ go install github.com/jackson198608/goProject/go_spider/example/stress
work/src$ ../bin/stress "http://bbs.goumin.com/forum-47-1.html" "bbslist" "/tmp/stress.log" "2" 
work/src$ ../bin/stress "http://mall.goumin.com/mall/list/219" "malllist" "/tmp/stress.log" "1"
work/src$ ../bin/stress "http://www.goumin.com/ask/list-1.html" "asklist" "/tmp/stress.log" "1"