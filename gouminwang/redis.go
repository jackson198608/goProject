package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackson198608/goProject/common/tools"
	log "github.com/thinkboy/log4go"
	redis "gopkg.in/redis.v4"
)

func connect() (*redis.ClusterClient, error) {
	redisInfo := tools.FormatRedisOption(c.redisConn)
	client, err := tools.GetClusterClient(&redisInfo)
	if err != nil {
		return nil, errors.New("[Error] redis connect error")
	}
	return client, nil
}

func insertArticleDetail(title string, dateline string, author string, content string, sourceUrl string) bool {

	// client := connect(redisConn)
	client, err := connect()
	if err != nil {
		log.Error("[Redis connect error]", err)
		return false
	}
	// jsonStr := {"title":title,"content":content}
	v := map[string]string{
		"title":     title,
		"dateline":  dateline,
		"author":    author,
		"content":   content,
		"sourceUrl": sourceUrl,
	}
	jsonStr, err := json.Marshal(v)
	if err != nil {
		log.Error("[Error] arr encode json error : ", err)
		return false
	}
	err1 := (*client).LPush(c.redisQueueName, jsonStr).Err()
	if err1 != nil {
		fmt.Println("[Error] push str into redis error:  ", jsonStr)
		return false
	}

	client.Close()
	return true
}

// func connect(conn string) (client *redis.Client) {
// 	client = redis.NewClient(&redis.Options{
// 		Addr:     conn,
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	_, err := client.Ping().Result()
// 	if err != nil {
// 		fmt.Println("[Error] redis connect error")
// 	}
// 	return client
// }
