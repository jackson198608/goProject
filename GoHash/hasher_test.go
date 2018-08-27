package GoHash

import (
	"fmt"
	//"github.com/jackson198608/goProject/GoHash/data"
	"testing"
	//"hash/crc32"
	//"github.com/jackson198608/goProject/go_spider/core/common/util"
)

func newH()(*Map){
	h := NewHasher(64,nil)
	return h

}

func TestHasher(t *testing.T){
	num := 10  //表数量
	prefix := "lingdang_"  //表前缀
	m := "100232378"
	////var a = []byte(m)
	////fmt.Println(crc32.ChecksumIEEE(a))
	////fmt.Println("========")
	//
	//targets := make([]string,10)
	//for i:=0;i<10 ;i++  {
	//	name := prefix+strconv.Itoa(i)
	//	targets[i] = name
	//}
	////fmt.Println(targets)
	//h := newH()
	//
	//h.Add(targets)
	////fmt.Println(h.keys)
	//
	//v := h.Get(m)
	//fmt.Println(v)
	n := CreateHash(prefix,num,m)
	fmt.Println(n)


}
