package redisEngine

import (
	"fmt"
	redis "gopkg.in/redis.v4"
	"testing"
)

func jobFunc(job string, mysqlConns []*xorm.Engine, mgoConns []*mgo.Session, taskarg []string) error {

	fmt.Println("this is jobFunc")
	return nil
}

func newRedis() *RedisEngine {
	r := new(redis.Options)
	NewRedisEngine("test", r, 10, jobFunc)
}

func TestDo(t *testing.T) {
	r := new(redis.Options)
	redisEngine := NewRedisEngine("test", r, 10, jobFunc)
	redisEngine.Do()

}
