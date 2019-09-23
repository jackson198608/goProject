package basepush

import (
	"testing"
	"fmt"
	"io/ioutil"
	"github.com/jackson198608/goProject/common/tools"
	"gopkg.in/redis.v4"
)

var p12Bytes []byte

func testConn() ( *redis.ClusterClient) {
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
	redisConn, _ := testConn()

	cBytes, err := ioutil.ReadFile("/etc/pro-lingdang.pem")
	if err != nil {
		return
	}
	p12Bytes = cBytes

	//ios
	// jobStr := "0|60fc3fe0caed4fdeccf2784514054dfd6267c3f47a3609e7e76758ab5b15a548|{\"aps\":{\"alert\":\"\\u60a8\\u6709\\u8bc4\\u8bba\\u672a\\u8bfb,\\u8bf7\\u70b9\\u51fb\\u67e5\\u770b~\",\"sound\":\"default\",\"badge\":1,\"type\":6,\"mark\":\"\"}}"

	//android
	// jobStr := "1|2efb35600238a8a30099161a|{\"content\":\"\\u72d7\\u6c11\\u7f51\\u6d4b\\u8bd5\\u6d3b\\u52a8\\u901a\\u77e5\",\"plats\":[1],\"target\":4,\"type\":1,\"alias\":[],\"registrationIds\":[\"ae86b1713f08da8c9b3e7cf4\"],\"androidTitle\":\"\\u72d7\\u6c11\\u7f51\\u6d4b\\u8bd5\\u6d3b\\u52a8\\u901a\\u77e5\",\"androidstyle\":1,\"androidVoice\":1,\"androidShake\":1,\"androidLight\":1,\"unlineTime\":1,\"extras\":\"{\\\"type\\\":1,\\\"mark\\\":5041505,\\\"content\\\":\\\"\\\\u72d7\\\\u6c11\\\\u7f51\\\\u6d4b\\\\u8bd5\\\\u6d3b\\\\u52a8\\\\u901a\\\\u77e5\\\",\\\"uid\\\":\\\"68296\\\",\\\"title\\\":\\\"\\\\u72d7\\\\u6c11\\\\u7f51\\\\u6d4b\\\\u8bd5\\\\u6d3b\\\\u52a8\\\\u901a\\\\u77e5\\\"}\"}"
	//1:android , 
	//2:huawei, 
	//3:xiaomi  7Gg2j8/yzoWfTJQh+Mr5x0BU7gQJRhE1xtQyJkXmVLRQEX7HgRaoQ3tSDKQ+vMnZ , 
	//4:oppo  CN_872d386c6d9ea3546cfe8f109da6e584, 
	//5:vivo,  15682454150111531282227
	//6:meizu   SO76c0e7e414c790176424e6e5b58525e64660c7c4746
	jobStr := "6|SO76c0e7e414c790176424e6e5b58525e64660c7c4746|{\"type\":1,\"mark\":\"4847131\",\"content\":\"铃铛签到\",\"uid\":-1,\"title\":\"铃铛签到\"}"
	

	m := Newpush(jobStr, redisConn,p12Bytes)
	fmt.Println(m.Do())
}

func TestPush(t *testing.T) {
	fmt.Println("push test")
}
