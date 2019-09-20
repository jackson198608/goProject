package oppo

import (
	// "errors"
	"fmt"
	"testing"
	"gopkg.in/redis.v4"
	"github.com/jackson198608/goProject/common/tools"
)

func testConn() ( *redis.ClusterClient) {
	redisstr := "192.168.86.82:6380,192.168.86.68:6381,192.168.86.82:6382,192.168.86.82:6383,192.168.86.82:6384,192.168.86.82:6385"
	redisInfo := tools.FormatRedisOption(redisstr)
	redisConn,_ := tools.GetClusterClient(&redisInfo)
	return redisConn
}

func TestToken(t *testing.T) {
	str := `{"code":0,"data":{"auth_token":"0f901777-0e23-4483-bd1d-cef3154d660c","create_time":1568808043026},"message":"Success"}`
	token := ""
	jsonStr := ""
	redisConn := testConn()
	op := NewPush(token, jsonStr, redisConn)
	err, authInfo := op.ParseToken(str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(authInfo.Data.AuthToken)
}