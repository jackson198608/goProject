package wxProgram

import (
	//"fmt"
	"strings"
	"strconv"
	"net/http"
	"fmt"
	"github.com/jackson198608/goProject/common/http/abuyunHttpClient"
	"encoding/json"
	"gopkg.in/redis.v4"
	"github.com/donnie4w/go-logger/logger"
	"github.com/jackson198608/goProject/common/tools"
	//"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
	"github.com/jackson198608/goProject/common/coroutineEngine/redisEngine"
)

var c Config = Config{
	"127.0.0.1:6379",
	1,
}      //redis info

type Task struct {
	phoneType int
	appId string
	TaskJson string
}
type accessTokenBody struct {
	Access_token string
	Expires_in int
}

const _token = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&"
const _sendUrl = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?"

const card_accestoken_key = "card_access_token_"

//type Result struct {
//	XMLName    xml.Name `xml:"returnsms"`  // 指定最外层的标签为returnsms
//	Returnstatus string  `xml:"returnstatus"` // 读取returnstatus，并将结果保存到returnstatus变量中
//	Message string `xml:"message"`  //读取错误信息
//}



/**
实例化
 */
func NewTask(redisString string) (t *Task){
	var tR Task
	result := strings.Split(redisString,"|")
	if len(result) != 3 {
		return nil
	}
	tR.phoneType,_ = strconv.Atoi(result[0])
	tR.appId = result[1]
	tR.TaskJson = result[2]
	return &tR
}
/**
发送请求
 */
func (p *Task) SendRequest() error {
	if p.appId == "" {
		logger.Error("[appid empty] ", p.appId, p.TaskJson)
		return nil
	}
	//获取secret
	secret := getSecret(p.appId)
	if secret == ""{
		logger.Error("[secret empty] ", p.appId, p.TaskJson)
		return nil
	}
	var mongoConnInfo []string
	var mysqlInfo []string
	redisInfo := tools.FormatRedisOption(c.redisConn)
	logger.Info("start work")
	r, err := redisEngine.NewRedisEngine("", &redisInfo, mongoConnInfo, mysqlInfo, c.coroutinNum, 1, jobFuc)
	if err != nil {
		logger.Error("[NewRedisEngine] ", err)
	}
	//获取AccessToken
	accessToken := p.getAccessToken(p.appId,secret,r)
	//发送请求
	if accessToken == "" {
		return nil
	}
	//请求微信
	err := p.requestWeixin(accessToken,p.TaskJson)
	if err != nil {
		logger.Error("[request weixin fail] ", err, p.appId, p.TaskJson)
	}
	return nil
}
/**
do request weixin
 */
func (p *Task)requestWeixin(accesstoken string,messqge string) error{
	var target = _sendUrl+"access_token="+accesstoken
	var h http.Header = make(http.Header)
	abuyun := p.setAbuyun()
	statusCode, _, body, err := abuyun.SendPostRequest(target,h,messqge,true)
	if err != nil {
		return err
	}
	fmt.Println(statusCode)
	fmt.Println("------------")
	if statusCode == 200 {
		fmt.Println(333333)
		fmt.Println(body)
	}
	return nil
}

/**
get secret form config by appid
 */
func getSecret(appId string) string{
	configs := make(map[string]string)
	secret := configs[appId]
	if secret != "" {
		return secret
	}
	//return ""
	return "feb421cb4a2cb7d0f9cb9a10fac15593"
}
/**
get access_token
 */
func (p *Task)getAccessToken(appId string,secret string,redisInfo *redis.Client) string{
	//查询redis缓存，获取accesstoken
	key := card_accestoken_key + appId
	acTo := redisInfo.Get(key).Val()
	if acTo != "" {
		return acTo
	}

	target := _token+"appid="+appId+"&secret="+secret
	var h http.Header = make(http.Header)
	abuyun := p.setAbuyun()
	statusCode, _, body, err := abuyun.SendRequest(target,h,"",true)
	if err != nil {
		return err.Error()
	}
	if statusCode == 200 {
		var at accessTokenBody
		if err:=json.Unmarshal([]byte(body),&at);err==nil{  //解析json，获取accesstoken
			//fmt.Println(at.Access_token)
			//更新缓存

			return at.Access_token
		}
	}
	return ""
}
func (p *Task) setAbuyun() *abuyunHttpClient.AbuyunProxy {
	var abuyun *abuyunHttpClient.AbuyunProxy = abuyunHttpClient.NewAbuyunProxy("", "", "")

	if abuyun == nil {
		fmt.Println("create abuyun error")
		return nil
	}
	return abuyun
}