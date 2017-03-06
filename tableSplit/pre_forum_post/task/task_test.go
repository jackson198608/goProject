package task

import (
	"fmt"
	"testing"
)

var dbAuth string = "root:goumintech"

var dbDsn string = "192.168.86.72:3309"
var dbName string = "test_dz2"

func TestGetPids(t *testing.T) {
	//task := NewTask(2730146, dbAuth, dbDsn, dbName)
	task := NewTask(2730142, dbAuth, dbDsn, dbName)
	if task != nil {
		fmt.Println(task.pids)
	}

}

func TestDo(t *testing.T) {
	//task := NewTask(2730146, dbAuth, dbDsn, dbName)
	task := NewTask(2730142, dbAuth, dbDsn, dbName)
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
