package tools

import (
	redis "gopkg.in/redis.v4"
)

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
