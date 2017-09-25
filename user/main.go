package main

import (
	"fmt"
)

var driverName string = "mysql"

var dbAuth string = "dog123:dog123"

var dbDsn string = "192.168.86.193:3307"

var dbName string = "test"

func main() {

	fmt.Println("***")
	getUser()
}
