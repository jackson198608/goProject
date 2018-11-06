package composite

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	//path := "/Users/Snow/img/bigimage_220.jpg"
	//path1 := "/Users/Snow/img/watermark/340.png"
	//path := "/Users/Snow/img/15238466413797.gif"
	path1 := "/Users/Snow/img/1.png"
	path := "/Users/Snow/img/IMG_0300.JPG"
	//path :="/Users/Snow/img/201806062118288161.jpg"
	//path := "/Users/Snow/img/15238466874121.gif"
	f := NewComposite(path, path1, "", 0,0)
	fmt.Println(f.Do())
}
