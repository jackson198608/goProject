package composite

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	path := "/Users/Snow/img/bigimage_220.jpg"
	path1 := "/Users/Snow/img/watermark/220.png"
	f := NewComposite(path, path1)
	fmt.Println(f.Do())
}
