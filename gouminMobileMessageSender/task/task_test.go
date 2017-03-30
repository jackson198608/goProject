package task

import (
	"testing"
	//"github.com/donnie4w/go-logger/logger"
)

func TestDo(t *testing.T) {
	task := NewTask(0, "18210091845|[铃铛宠物] 您好，1391（验证码不要告诉别人哦，5分钟内有效）小铃铛等你很久了，快去登录领豆吧！可以兑换商品抽取大奖呢！【狗民网】")
	if task != nil {
		task.Do()
	}
}
