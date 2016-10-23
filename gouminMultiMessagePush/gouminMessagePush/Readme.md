# 1.安装方法
- cd gouminMessagePush
- sh install

# 2.简介
## 2.1 msgConfig.json
### 2.1.1 httpTimeOut
- http connection timeout settting

### 2.1.2 currentNum
- 进程执行任务的最大并发数

### 2.1.3 redisConn
- redis 连接参数 ip:port

### 2.1.4 mongoConn
- mongoDb 连接参数 ip:port


# 3.进程执行方式
- 循环5秒
- 检查redis任务，执行插入mongo，并发送push


# 4.任务数据
## 4.1 格式
- phoneType|pushJsonString^mongoInserJsonString
- phoneType   0:ios  1:android

## 4.2 demo数据

```
127.0.0.1:6379> LPOP MessageCenter
"0|fb71306452499efc778cc77d1be6614b8e1753e79b7acd8348ce6cb47abd4dc2|{\"aps\":{\"alert\":\"task\",\"sound\":\"default\",\"badge\":1,\"type\":6,\"mark\":\"\"}}^{\"uid\":1895167,\"type\":1,\"makr\":281,\"isnew\":0,\"from\":0,\"channel\":1,\"channel_types\":2,\"title\":\"\xe7\x8b\x97\xe7\x8b\x97\xe7\x9a\x84\xe5\xaf\x82\xe5\xaf\x9e\xe9\x83\xbd\xe5\xb8\x82\xe4\xb9\x8b\xe6\xae\x87\",\"content\":\"\xe5\xb0\x8f\xe7\x9f\xad\xe8\x85\xbf\xe5\x9c\xa8\xe5\xae\xb6\xe6\xbc\x82\xe7\xa7\xbb\xef\xbc\x8c\xe5\xb1\x81\xe8\x82\xa1\xe5\xb7\xae\xe7\x82\xb9\xe6\xb2\xa1\xe7\x94\xa9\xe6\x8e\x89\xe4\xba\x86~\",\"image\":\"/messagepush/day_161020/20161020_7a50e50.jpg\",\"url_type\":1,\"url\":\"4346101\",\"created\":\"2016-10-20 14:12:28\",\"modified\":\"0000-00-00 00:00:00\"}"
```



