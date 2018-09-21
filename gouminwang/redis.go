package main

import (
	"encoding/json"
	"fmt"
	redis "gopkg.in/redis.v4"
)

func insertArticleDetail(title string, dateline string, author string, content string, sourceUrl string) bool {

	client := connect(redisConn)
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
		logger.Println("[Error] arr encode json error : ", err)
		return false
	}
	err1 := (*client).LPush(redisQueueName, jsonStr).Err()
	if err1 != nil {
		fmt.Println("[Error] push str into redis error:  ", jsonStr)
		return false
	}

	client.Close()
	return true
}

func connect(conn string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("[Error] redis connect error")
	}
	return client
}
