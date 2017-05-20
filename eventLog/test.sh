work$ bin/dianping "/data/dianping/" "http://bbs.goumin.com/forum-34-1.html" "gouminList" "/tmp/1.log" 1000 1 "北京"
work$ bin/dianping "/data/dianping/" "http://www.dianping.com/search/category/2/95/g25147r2580" "shopList" "/tmp/1.log" 1000 1 "北京"

work/src$ go install github.com/jackson198608/gotest/go_spider/example/dianping
work/src$ go get github.com/jackson198608/gotest/go_spider/example/dianping
work$ ../bin/dianping "/data/dianping/" "http://www.dianping.com/search/category/2/95/g25147r2580" "shopList" "/tmp/1.log" 1000 1 "北京"

work/src$ go install github.com/jackson198608/gotest/redisAction
work/src$ ../bin/redisAction
mlog.StraceInst().Println("start crawl : " + req.GetUrl())

work/src$ go install github.com/jackson198608/gotest/bbsPosttable
work/src$ ../bin/bbsPosttable

work/src$ go get gopkg.in/mgo.v2


work/src$ go install github.com/jackson198608/gotest/exportData/
work/src$ ../bin/exportData "mall" "select * from address limit 10"
work/src$ ../bin/exportData "mall" "select * from \`order\` limit 10"

work/src$ go install github.com/jackson198608/goProject/eventLog/
work/src$ ../bin/eventLog create
work/src$ ../bin/eventLog do
    nohup eventLog create >/tmp/eventlog.log 2>&1 &