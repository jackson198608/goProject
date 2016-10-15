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

// test return value to be function,also can be called closure
func test3() F {
	var s string = "say somthing to see the var outside,this is test3 speaking"
	return func() {
		fmt.Println(s)
	}
}

// accept int paramter and return two closure function
func test5(b int) (func(int, int) (bool, int), func(bool) bool) {
	var a bool = true
	var sum int = 0
	return func(x int, y int) (bool, int) {
			fmt.Println("this is the add function in test5 ", sum)
			sum++
			z := (x + y) * b
			return a, z
		}, func(a bool) bool {
			fmt.Println("this is return bool function in test5 ", sum)
			sum++
			return a
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

//take function as paramter
func test4(task func(a int) int) (bool, int) {
	b := task(3)
	return true, b
}

func main() {
	/*
		//按顺序解析()
		test3()()
	*/

	a, b := test5(10)
	_, z := a(1, 5)
	fmt.Println(z)
	b(true)

	/*
		_, a := test4(doTest)
		fmt.Println(a)
	*/

}
