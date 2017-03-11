package main

import (
	//"github.com/jackson198608/goProject/tableSplit/pre_forum_post/task"
	"fmt"
)

const dbDsn = "192.168.86.72:3309"
const dbName = "pre_forum_post"
const dbAuth = "root:goumintech"
const numloops = 20

func main() {
	tids := getTask(10)
	fmt.Println(tids)
}
