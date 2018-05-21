package composite

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	//path := "/Users/Snow/img/bigimage_220.jpg"
	path1 := "/Users/Snow/img/watermark/220.png"
	//path := "/Users/Snow/img/15238466413797.gif"
	path := "/Users/Snow/img/15238466701678.gif"
	//path := "/Users/Snow/img/15238466874121.gif"
	f := NewComposite(path, path1)
	fmt.Println(f.Do())
}
