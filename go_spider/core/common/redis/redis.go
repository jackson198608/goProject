package redis

import (
    redis "gopkg.in/redis.v5"
)
var Proxy string = "192.168.86.88:6379"
type Redis struct {
}
//根据参数 count 的值，移除列表中与参数 value 相等的元素。
func LRem(ProxyIp string){
    client := conn(Proxy)
    (*client).LRem("proxy", 1, ProxyIp).Val()
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
