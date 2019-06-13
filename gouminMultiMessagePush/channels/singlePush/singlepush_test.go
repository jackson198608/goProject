package singlePush

import (
	"testing"
	"fmt"
	"io/ioutil"
	"gouminGitlab/common/orm/elasticsearchBase"
	"github.com/olivere/elastic"
)

var p12Bytes []byte

func testConn() ( *elastic.Client) {


	var nodes []string
	nodes = append(nodes, "http://192.168.86.230:9200")
	//nodes = append(nodes, "http://192.168.86.231:9200")
	r,_ := elasticsearchBase.NewClient(nodes)
	esConn,_ :=r.Run()
	return esConn
	// return engine, session
}

func TestDo(t *testing.T) {
	esConn := testConn()

	cBytes, err := ioutil.ReadFile("/etc/pro-lingdang.pem")
	if err != nil {
		return
	}
	p12Bytes = cBytes

	//ios
	//jobStr := "0|60fc3fe0caed4fdeccf2784514054dfd6267c3f47a3609e7e76758ab5b15a548|{\"aps\":{\"alert\":\"\\u60a8\\u6709\\u8bc4\\u8bba\\u672a\\u8bfb,\\u8bf7\\u70b9\\u51fb\\u67e5\\u770b~\",\"sound\":\"default\",\"badge\":1,\"type\":6,\"mark\":\"\"}}"

	//android
	jobStr := "1|ae86b1713f08da8c9b3e7cf4|{\"content\":\"\\u72d7\\u6c11\\u7f51\\u6d4b\\u8bd5\\u6d3b\\u52a8\\u901a\\u77e5\",\"plats\":[1],\"target\":4,\"type\":1,\"alias\":[],\"registrationIds\":[\"ae86b1713f08da8c9b3e7cf4\"],\"androidTitle\":\"\\u72d7\\u6c11\\u7f51\\u6d4b\\u8bd5\\u6d3b\\u52a8\\u901a\\u77e5\",\"androidstyle\":1,\"androidVoice\":1,\"androidShake\":1,\"androidLight\":1,\"unlineTime\":1,\"extras\":\"{\\\"type\\\":1,\\\"mark\\\":5041505,\\\"content\\\":\\\"\\\\u72d7\\\\u6c11\\\\u7f51\\\\u6d4b\\\\u8bd5\\\\u6d3b\\\\u52a8\\\\u901a\\\\u77e5\\\",\\\"uid\\\":\\\"68296\\\",\\\"title\\\":\\\"\\\\u72d7\\\\u6c11\\\\u7f51\\\\u6d4b\\\\u8bd5\\\\u6d3b\\\\u52a8\\\\u901a\\\\u77e5\\\"}\"}"

	m := NewSinglepush(jobStr, esConn,p12Bytes)
	fmt.Println(m.Do())
}
