package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var Type string = "0"
var FileName string = "/tmp/5.csv"

func main() {
	FileName = os.Args[1]
	Type = os.Args[2]

	fmt.Println("vim-go: ", FileName, " ", Type)

	do()
}

func do() {
	f, err := os.Open(FileName)
	if err != nil {
		panic(err)

	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, " ", "", -1)
		line1 := line[0 : len(line)-1]
		command := "spider /data/dianping1/ " + string(line1)
		command1 := " shopDetail /tmp/4.log 1 8 广州 1"
		command2 := command + command1
		fmt.Println(command2)
		/*
			cmd := exec.Command(command2)
			err = cmd.Run()
			fmt.Println(err)
		*/
	}
}
