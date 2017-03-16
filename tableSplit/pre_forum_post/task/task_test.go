package task

import (
	"fmt"
	"testing"
)

var dbAuth string = "dog123:dog123"

var dbDsn string = "210.14.154.198:3306"
var dbName string = "new_dog123"

func buildArgs() []string {
	args := make([]string, 0, 3)
	args = append(args, dbAuth)
	args = append(args, dbDsn)
	args = append(args, dbName)
	return args
}

func TestGetPids(t *testing.T) {
	args := buildArgs()
	buildArgs()
	//task := NewTask(2730146, dbAuth, dbDsn, dbName)
	task := NewTask(0, "2730142", args)
	if task != nil {
		fmt.Println(task.pids)
	}

}

func TestDo(t *testing.T) {
	//task := NewTask(2730146, dbAuth, dbDsn, dbName)
	args := buildArgs()
	task := NewTask(0, "1035", args)
	if task != nil {
		fmt.Println(task.pids)
	}
	task.Do()

}

func checkError(err error, errstr string) {
	if err != nil {
		fmt.Println("[err]", errstr, err)
	}
}
