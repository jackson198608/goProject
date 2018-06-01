package wxProgram

import (
	"stathat.com/c/jconfig"
)

type Config struct {
	redisConn   string
	coroutinNum int
}

func loadConfig() {
	config := jconfig.LoadConfig("/etc/pushContentCenterConfig.json")
	c.redisConn = config.GetString("redisConn")
	c.coroutinNum = config.GetInt("coroutinNum")
}
