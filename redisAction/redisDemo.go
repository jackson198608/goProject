package main

import (
	"fmt"
	redis "gopkg.in/redis.v4"
)

func testSet(client *redis.Client) {
	err := (*client).Set("zhou", "Google", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := (*client).Get("zhou").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("zhou", val)

	val2, err := (*client).Get("mykey").Result()
	if err == redis.Nil {
		fmt.Println("mykey does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("mykey", val2)
	}
}

func testRpush(client *redis.Client) {
	err := (*client).RPush("list1", "fuck1").Err()
	if err != nil {
		panic(err)
	}
	err = (*client).RPush("list1", "do2").Err()
	if err != nil {
		panic(err)
	}
	err = (*client).RPush("list1", "do3").Err()
	if err != nil {
		panic(err)
	}

	for {
		val := (*client).LPop("list1").Val()
		fmt.Println(val)

		if val == "" {
			fmt.Println("there is no data2")
			break
		}
	}
}

func testLRange(client *redis.Client) {
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
	(*client).RPush("list2", "do1")
}

func conn(conn string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     conn,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client

}

func main() {
	client := conn("210.14.154.198:6379")
	//testSet(client)
	testRpush(client)
}
