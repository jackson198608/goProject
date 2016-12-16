package main

import (
    "fmt"
    redis "gopkg.in/redis.v5"
    "time"
    "math/rand"
    "os"
    "io"
    "bufio"
    "strings"
    // "strconv"
)
var Proxy string = "192.168.86.88:6379"
var ProxyFile string = "/home/wang/www/proxy.csv" //代理批量数据文件

//从redis中随机读取出一个ip代理
func LRange() string {
    client := conn(Proxy)
    count := (*client).LLen("proxy").Val()
    if count == 0 {
        fmt.Println("[notice] got nothing")
        return ""
    }
    
    counts := int(count) //int64类型转为int
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    intNum := r.Intn(counts) //取随机数
    int64Nums := int64(intNum) //int类型转为int64
    // fmt.Println("type:", reflect.TypeOf(int64Nums))    //查看变量类型
    redisStr := (*client).LRange("proxy", int64Nums, int64Nums).Val() //随机获取代理ip
    if redisStr[0] == "" {
        fmt.Println("[notice] got nothing")
        return ""
    }
    return redisStr[0]
}

func LPop() {
    client := conn(Proxy)
    redisStr := (*client).LPop("proxy").Val()
    fmt.Println(redisStr)
    if redisStr == "" {
        fmt.Println("[notice] got nothing")
        return
    }
    // return redisStr
}

//根据参数 count 的值，移除列表中与参数 value 相等的元素。
func LRem(ProxyIp string) {
    client := conn(Proxy)
    (*client).LRem("proxy", 1, ProxyIp).Val()
    // fmt.Println(lRem)
}

//将一个或多个值 value 插入到列表 key 的表尾(最右边)。
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
        (*client).RPush("proxy", line)
        // fmt.Println(line)
    } 
}


func conn(conn string) (client *redis.Client) {
    client = redis.NewClient(&redis.Options{
        Addr:     conn,
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    // pong, err := client.Ping().Result()
    // fmt.Println(pong, err)
    // Output: PONG <nil>
    return client

}

// func main() {
//     // client := conn("192.168.86.88:6379")
//     //testSet(client)
//     LPop()
// }