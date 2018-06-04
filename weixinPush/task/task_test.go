package task

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/jackson198608/goProject/common/tools"
	mgo "gopkg.in/mgo.v2"
	"testing"
	"gouminGitlab/common/weixin/accessToken/accesstokenManager"
	"fmt"
)

const dbAuth = "dog123:dog123"
const dbDsn = "192.168.86.193:3307"
const dbName = "new_dog123"
const mongoConn = "192.168.86.192:27017" //"192.168.86.193:27017,192.168.86.193:27018,192.168.86.193:27019"

//var weixinAccessTokens *accesstokenManager.Manager

func newtask() (*Task, error) {
	connStr := tools.GetMysqlDsn(dbAuth, dbDsn, dbName)
	engine, err := xorm.NewEngine("mysql", connStr)
	if err != nil {
		return nil, err
	}

	engines := []*xorm.Engine{engine}

	//get mongo session
	session, err := mgo.Dial(mongoConn)
	if err != nil {
		return nil, err
	}

	sessions := []*mgo.Session{session}

	jobStr := `wxe691667e66d1f8eb|{"formid":"t7OVsxS7ipiiBgAU64HjAiRu82yCjdkXqZc9HydUgDU","data":{"message":"\u72d7\u6c11\u7f51","color":"#000000"},"template_id":"t7OVsxS7ipiiBgAU64HjAiRu82yCjdkXqZc9HydUgDU","openid":"t7OVsxS7ipiiBgAU64HjAiRu82yCjdkXqZc9HydUgDU"}`

	t, err := NewTask(jobStr, engines, sessions)
	return t, err

}

func TestNewTask(t *testing.T) {
	secret := "wxe691667e66d1f8eb:feb421cb4a2cb7d0f9cb9a10fac15593,wx304f38886b8bba74:da6db4abedb0b199d4cbf37fab7c33b9"
	weixinAccessTokens := accesstokenManager.NewManager(secret)
	fmt.Println(weixinAccessTokens)
	fmt.Println(999999999)
	task, err := newtask()
	if task == nil {
		t.Log("task create error", err)
		t.Fail()
	}
	a := task.Do(weixinAccessTokens)
	fmt.Println(a)
}
