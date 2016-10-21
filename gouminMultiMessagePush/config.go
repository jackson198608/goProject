package main

type Config struct {
	redisConn string
	mongoConn string
}

type redisData struct {
	pushStr   string
	insertStr string
}
