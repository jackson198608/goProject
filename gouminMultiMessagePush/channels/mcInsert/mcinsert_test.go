package mcInsert

import (
"testing"
"fmt"
"io/ioutil"
"gouminGitlab/common/orm/elasticsearchBase"
"github.com/olivere/elastic"
	"github.com/jackson198608/goProject/common/tools"
	"gopkg.in/redis.v4"
)

var p12Bytes []byte

func testConn() ( *redis.ClusterClient,*elastic.Client) {
	redisstr := "192.168.86.82:6380,192.168.86.68:6381,192.168.86.82:6382,192.168.86.82:6383,192.168.86.82:6384,192.168.86.82:6385"
	redisInfo := tools.FormatRedisOption(redisstr)
	redisConn,_ := tools.GetClusterClient(&redisInfo)

	var nodes []string
	nodes = append(nodes, "http://192.168.86.230:9200")
	//nodes = append(nodes, "http://192.168.86.231:9200")
	r,_ := elasticsearchBase.NewClient(nodes)
	esConn,_ :=r.Run()
	return redisConn,esConn
	// return engine, session
}

func TestDo(t *testing.T) {
	redisConn,esConn := testConn()

	cBytes, err := ioutil.ReadFile("/etc/pro-lingdang.pem")
	if err != nil {
		return
	}
	p12Bytes = cBytes

	jobStr := "{\"uid\":2189280,\"type\":1,\"mark\":651,\"isnew\":0,\"from\":0,\"channel\":1,\"channel_types\":1,\"title\":\"\\u6d4b\\u8bd5\\u53d1\\u9001\\u6d88\\u606f\",\"content\":\"\\u6d4b\\u8bd5\\u53d1\\u9001\\u6d88\\u606f\",\"image\":\"\",\"url_type\":0,\"url\":\"mall.goumin.com\",\"created\":\"2019-06-04 16:43:57\",\"modified\":\"\"}"
	m := NewTask(jobStr)
	fmt.Println(m.Insert(esConn,redisConn))
}
