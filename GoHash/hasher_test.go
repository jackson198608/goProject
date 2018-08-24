package GoHash

import (
	"fmt"
	//"github.com/jackson198608/goProject/GoHash/data"
	"testing"
	//"hash/crc32"
	//"github.com/jackson198608/goProject/go_spider/core/common/util"
	"strconv"
)

func newH()(*Map){
	h := NewHasher(64,nil)
	return h

}

func TestHasher(t *testing.T){
	//num := 10
	prefix := "golang_"
	m := "100232378"
	//var a = []byte(m)
	////b := Crc32Hash.
	//fmt.Println(crc32.ChecksumIEEE(a))
	//fmt.Println("========")
	//c := newC()
	//c.Add(m)
	//n,err := c.Get("3")
	//fmt.Println(n)
	//fmt.Println(err)
	//var targets []string
	//var t []string = make([]string,10)
	targets := make([]string,10)
	for i:=0;i<10 ;i++  {
		name := prefix+strconv.Itoa(i)
		targets[i] = name
	}
	//fmt.Println(targets)
	h := newH()

	h.Add(targets)
	//fmt.Println(h.keys)

	v := h.Get(m)
	fmt.Println(v)



}
