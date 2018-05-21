package compress

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	//path := "/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png"
	//path := "/Users/Snow/img/15238466413797.gif"
	path := "/Users/Snow/img/15238466701678.gif"
	//path := "/Users/Snow/img/15238466874121.gif"
	width := 200
	height := 200
	f := NewCompress(path, width, height)
	fmt.Println(f.Do())
}
