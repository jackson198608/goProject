# 设计思想
- 经常有需求是从redis读task，然后new task，再do
- 考虑做一个公共化modules，初始化传入task,queueName.
- 自动启动守护进程loop

