package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func testPost() {
	url := "http://sdk.open.api.igexin.com/apiex.htm"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"action":"pushMessageToSingleAction","clientData":"CgASC3B1c2htZXNzYWdlGgAiFmM2Mk5XOUp4d2w5eFJTRlh6SlAwbjQqFm1RZ3FuY0ZDVXk5aG14aHFLQ3ZHTDkyADoKCgASABoAIgItMUIKCAEQABiuTpAKAEIcCK5OEAMYZOIIAOoIBgoAEgAaAPAIAPgIZJAKAEIHCGQQB5AKAEoJZHVyYXRpb249","transmissionContent":"eyJ0eXBlIjo2LCJtYXJrIjoiIiwic2lsZW50IjowLCJ0aXRsZSI6InRlc3QgcHVzaCB0aXRsZSIsImNvbnRlbnQiOiJ0ZXN0IHB1c2ggY29udGVudCEhISIsInVpZCI6MjAzNTgzMH0=","isOffline":true,"offlineExpireTime":43200000,"pushNetWorkType":1,"appId":"mQgqncFCUy9hmxhqKCvGL9","clientId":"aef16d24053b161bb0b35d3a0f775fcd","alias":null,"type":2,"pushType":"TransmissionMsg","version":"3.0.0.0","appkey":"c62NW9Jxwl9xRSFXzJP0n4"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("User-Agent", "GeTui PHP/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func main() {
	fmt.Println("vim-go")
	testPost()
}
