package main

import "fmt"

//define type to be simple usage
type F func()
type FS []F

// basic usage
func test1() {
	c := FS{func() { fmt.Print("hello ") }, func() { fmt.Println("douzifly") }}
	for i := 0; i < len(c); i++ {
		c[i]()
	}
	fmt.Println("vim-go")
}

func do() {
	fmt.Println("do")
}

// test return value to be function
func test3() F {
	return func() {
		fmt.Println("return by test3")
	}
}

// test slice function
func test2() {
	var d []F = make([]F, 3, 3)
	d[0] = func() {
		for i := 0; i < 3; i++ {
			fmt.Println("test")
		}
	}

	d[1] = func() {
		a := 0
		b := 1
		c := a + b
		fmt.Println(c)
	}
	d[2] = do
	d[0]()
	d[1]()
	d[2]()
}

func doTest(a int) int {
	fmt.Println("I am doTest")
	return a + 1
}

func test4(task func(a int) int) (bool, int) {
	b := task(3)
	return true, b
}

func main() {
	//按顺序解析()
	//test3()()
	_, a := test4(doTest)
	fmt.Println(a)

}
