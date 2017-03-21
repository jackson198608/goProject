# 设计思想
- 经常有需求是从redis读task，然后new task，再do
- 考虑做一个公共化modules，初始化传入task,queueName.
- 自动启动守护进程loop


# 项目定制化改动
- 做一些针对帖子分表项目的改变
- 比如 执行100个队列
- 每个队列对应个帖子分表
- 从push，到pop都要从不同的tid中获取了

