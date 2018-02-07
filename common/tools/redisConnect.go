package tools

import (
	redis "gopkg.in/redis.v4"
	"strings"
)

func FormatRedisOption(redisConn string) redis.ClusterOptions {
	var redisStr []string
	redisConns := strings.Split(redisConn, ",")
	for i, _ := range redisConns {
		redisStr = append(redisStr, redisConns[i])
	}
	redisInfo := redis.ClusterOptions{
		Addrs: redisStr,
	}
	return redisInfo
}

func GetClusterClient(redisInfo *redis.ClusterOptions) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(redisInfo)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetClient(redisInfo *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(redisInfo)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
