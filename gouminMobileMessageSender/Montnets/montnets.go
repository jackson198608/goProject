package Montnets

import (
	"errors"
	"github.com/donnie4w/go-logger/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const interfaceUrl = ""
const userId = ""
const password = ""

type Montnets struct {
	phone   string
	message string
}

func NewMontnets(logLevel logger.LEVEL, phone string, message string) *Montnets {
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

func (m *Montnets) Send() error {

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

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Host", "61.135.198.131")
	//req.Header.Set("Content-Length", strconv.Itoa(len(bodyStr)))

	resp, err := client.Do(req) //发送
	if err != nil {
		logger.Error("connect to remote err", err)
		return errors.New("connect to remote err")
	}
	defer resp.Body.Close() //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	logger.Debug("this data we got is :", string(data))
	return nil
}
