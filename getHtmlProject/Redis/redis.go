package Redis

import (
	"github.com/donnie4w/go-logger/logger"
	redis "gopkg.in/redis.v4"
)

func PushTaskData(client *redis.Client, queueName string, tasks interface{}) bool {
	switch realTasks := tasks.(type) {
	case []string:
		logger.Info("this is string task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			err := (*client).RPush(queueName, realTasks[i]).Err()
			if err != nil {
				logger.Error("insert redis error", err)
			}
		}

	case []int:
		logger.Info("this is int task", realTasks)
		for i := 0; i < len(realTasks); i++ {
			err := (*client).RPush(queueName, realTasks[i]).Err()
			if err != nil {
				logger.Error("insert redis error", err)
			}
		}

	default:
		logger.Error("this is not normal format", realTasks)
		return false
	}

	return true
}
