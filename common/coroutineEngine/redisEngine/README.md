# introduction
- this package is design for the situation when you want to  parallel your code with redis queue 

# Installation

	go get github.com/jackson198608/goProject/common/coroutineEngine/redisEngine	

# Quick Start

## Create redisEngine 
- engine needs to Know: queueName ,redisInfo, cocurrentNum, workFunc
- the most important is the workFunc 
- workFunc will be invoked in coroutine which created inside the engine
- @todo change demo code to the real

```Go
	r := new(redis.Options)
	NewRedisEngine("test", r, 10, jobFunc)
```

## do the job

```Go
	r.Do()
```

