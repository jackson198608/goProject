package main

import (
	"fmt"
	"github.com/user/stringutil"
)

func testBasicGrammar() {
	fmt.Println(stringutil.Reverse("!oG ,olleH"));
	//int testing
	var test1 int8 =3;
	var test2 int =5;
	test3 := int(test1)+test2;
	fmt.Println(test3);

	//byte testing
	var b1 byte =3;
	var b2 byte ='a'+b1;
	fmt.Println(b2);
	fmt.Printf("%c\n",b2);

	//rune testing

	//const testing
	const c1 = 3;
	c2 := c1;
	fmt.Printf("%d\n",c2);

	//string testing
	var str1 string = "奥运会就看个金牌数";
	var str2 string;
	str2 = "kobe is mvp";
	fmt.Printf("%s\n",str1);
	fmt.Printf("%s\n",str2);
	var ch byte;
	//print for char
	for i:=0;i<len(str2);i++{
		ch=str2[i];
		fmt.Printf("%c",ch);
	}
	fmt.Printf("\n");

	//print for char
	for i:=0;i<len(str2);i++{
		ch=str2[i];
		fmt.Printf("%c",ch);
	}
	fmt.Printf("\n");

	var chR rune;
	for _,chR = range str1{
		fmt.Printf("%q",chR);
	}		
	fmt.Printf("\n");


	//test pointer
	var a *int;
	var i int = 4;
	if(a == nil){
		fmt.Println("this is the mempty pointer");
	}else{
		fmt.Println("this is not empty");
	}

	a=&i;
	if(a == nil){
		fmt.Println("this is the mempty pointer");
	}else{
		fmt.Println("this is not empty");
	}

	//array testing
	var players [6]int=[6]int{2, 3, 5, 7, 11, 13};
	fmt.Printf("%d\n",players[1]);

	//slice testing
	//var mySlice []int = players[:3]
}

func testFurtherGrammar(){
	var a [4]int=[4]int{1,2,3,4};
	fmt.Println("a=",a," len=",len(a)," cap=",cap(a));
	var b []int = a[0:1];
	fmt.Println("b=",b," len=",len(b)," cap=",cap(b));
	b=append(b,11,12);
	fmt.Println("a=",a," len=",len(a)," cap=",cap(a));
	fmt.Println("b=",b," len=",len(b)," cap=",cap(b));


	b=append(b,1,2,3,4,5,6,7,8);
	fmt.Println("a=",a," len=",len(a)," cap=",cap(a));
	fmt.Println("b=",b," len=",len(b)," cap=",cap(b));

	c := []int{1,2,3};
	fmt.Println("c=",c," len=",len(c)," cap=",cap(c));
}

func testSlice(a []int) int{
	a[3]=100;
	fmt.Println(a);
	return 1;
}

func main(){
	//testBasicGrammar();
	//testFurtherGrammar();
	/*testSlice*/
	//a := []int{1,2,3,4,5};
	//b := testSlice(a);
	//fmt.Println(b);

	/*testInterface*/	
}
