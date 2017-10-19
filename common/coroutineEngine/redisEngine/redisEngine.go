package redisEngine

import (
	redis "gopkg.in/redis.v4"
)

type RedisEngine struct {
	queueName   string
	redisInfo   *redis.Options
	coroutinNum int
	workFun     func(c chan int, job string, taskarg ...string) error
}

func NewRedisEngine(queueName string, redisInfo *redis.Options, coroutinNum int, workFun func(c chan int, job string, taskarg ...string) error) *RedisEngine {
	if (queueName == "") || (redisInfo == nil) || (coroutinNum <= 0) || (workFun == nil) {
		return nil
	}

	r := new(RedisEngine)
	if r == nil {
		return nil
	}

	r.queueName = queueName
	r.redisInfo = redisInfo
	r.coroutinNum = coroutinNum
	r.workFun = workFun

	return r

}

func (r *RedisEngine) Do() error {
	c := make(chan int, 1)
	r.workFun(c, "", "hello")

	return nil
}
