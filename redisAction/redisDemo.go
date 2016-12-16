package main

import (
	"fmt"
	redis "gopkg.in/redis.v4"
	"time"
	"math/rand"
	// "strconv"
	"reflect"
	"os"
	"io"
	"bufio"
	"strings"
)

var Proxy string = "192.168.86.88:6379"
var ProxyFile string = "/home/wang/www/proxy2.csv"

func testSet(client *redis.Client) {
	err := (*client).Set("zhou", "Google", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := (*client).Get("zhou").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("zhou", val)

	val2, err := (*client).Get("mykey").Result()
	if err == redis.Nil {
		fmt.Println("mykey does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("mykey", val2)
	}
}

func testRpush(client *redis.Client) {
	err := (*client).RPush("list1", "fuck1").Err()
	if err != nil {
		panic(err)
	}
	err = (*client).RPush("list1", "do2").Err()
	if err != nil {
		panic(err)
	}
	err = (*client).RPush("list1", "do3").Err()
	if err != nil {
		panic(err)
	}

	for {
		val := (*client).LPop("list1").Val()
		fmt.Println(val)

		if val == "" {
			fmt.Println("there is no data2")
			break
		}
	}
}

func testLRange(client *redis.Client) {
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
}

func conn(conn string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client

}

func LPop(client *redis.Client) {
    // client := conn(Proxy)
    redisStr := (*client).LPop("proxy").Val()
    fmt.Println(redisStr)
    if redisStr == "" {
        fmt.Println("[notice] got nothing")
        return
    }
    // return redisStr
}

func LRange() string {
	client := conn(Proxy)
	count := (*client).LLen("proxy").Val()
	if count == 0 {
		fmt.Println("[notice] got nothing")
		return ""
	}
	
	fmt.Println(count)
	counts := int(count) //int64类型转为int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	intNum := r.Intn(counts) //取随机数
	int64Nums := int64(intNum) //int类型转为int64
	fmt.Println("type:", reflect.TypeOf(int64Nums))    //查看变量类型
	redisStr := (*client).LRange("proxy", int64Nums, int64Nums).Val() //随机获取代理ip
    if redisStr[0] == "" {
        fmt.Println("[notice] got nothing")
        return ""
    }
    fmt.Println(redisStr[0])
    return redisStr[0]
}

func LRem(ProxyIp string) {
	client := conn(Proxy)
	lRem := (*client).LRem("proxy", 1, ProxyIp).Val()
	fmt.Println(lRem)
}

func LLen(client *redis.Client) int64{
	count := (*client).LLen("proxy").Val()
	// fmt.Println(count)
	if count == 0 {
		fmt.Println("[notice] got nothing")
		return 0
	}
	return count
}

func RPush() {
	client := conn(Proxy)
	f, err := os.Open(ProxyFile)
	if err != nil {
        fmt.Println(err)
    }
    defer f.Close()
    rd := bufio.NewReader(f)
    for {
        line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
        if err != nil || io.EOF == err {
            break
        }
        line = strings.TrimSpace(line) //取得末尾的\n字符
		(*client).RPush("list2", line)
        // fmt.Println(line)
    } 
}

func main() {
	RPush()
	// client := conn("192.168.86.88:6379")
	// rs := LRange()
	// fmt.Println(rs)
	// LRem(rs)
	// if err != nil {
	// 	fmt.Println("[notice] got nothing")

	// }
	// LPop(client)
	//testSet(client)
	// testRpush(client)
}
