package Montnets

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const interfaceUrl = "1"
const userId = "2"
const password = "3"

type Montnets struct {
	phone   string
	message string
}

func NewMontnets(logLevel int, phone string, message string) *Montnets {
	logger.SetLevel(logger.LEVEL(logLevel))

	//check params
	_, err := strconv.Atoi(phone)
	if err != nil {
		logger.Error("phone is not int", phone)
	}

	if message == "" {
		logger.Error("message can not be null", message)
	}

	m := new(Montnets)
	if m == nil {
		logger.Error("new Montnets error")
	}

	m.phone = phone
	m.message = message

	return m

}

func (m *Montnets) send() {

	v := url.Values{}
	v.Set("userId", userId)
	v.Set("password", password)
	v.Set("pszMobis", m.phone)
	v.Set("pszMsg", m.message)
	v.Set("iMobiCount", "1")
	v.Set("pszSubPort", "*")
	v.Set("MsgId", "0")

	body := strings.NewReader(v.Encode()) //把form数据编下码
	defer ioutil.NopCloser(body)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", interfaceUrl, body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //这个一定要加，不加form的值post不过去，被坑了两小时
	req.Header.Set("Host", "61.135.198.131")                            //这个一定要加，不加form的值post不过去，被坑了两小时
	//req.Header.Set("Content-Length", strconv.Itoa(len(bodyStr)))        //这个一定要加，不加form的值post不过去，被坑了两小时
	fmt.Printf("%+v\n", req) //看下发送的结构

	resp, err := client.Do(req) //发送
	defer resp.Body.Close()     //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data), err)

}
