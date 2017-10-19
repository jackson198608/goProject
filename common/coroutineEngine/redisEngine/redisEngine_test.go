package redisEngine

import (
	"fmt"
	redis "gopkg.in/redis.v4"
	"strconv"
	"testing"
)

func jobFunc(c chan int, job string, taskarg ...string) error {
	strconv.Itoa(10)

	fmt.Println("this is jobFunc")
	return nil
}

func TestNew(t *testing.T) {
	r := new(redis.Options)
	NewRedisEngine("test", r, 10, jobFunc)
}

func TestDo(t *testing.T) {
	r := new(redis.Options)
	redisEngine := NewRedisEngine("test", r, 10, jobFunc)
	redisEngine.Do()

}
