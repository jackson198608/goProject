package compress

import (
	"fmt"
	"testing"
)

func TestParseJson(t *testing.T) {
	jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200}"
	c := NewCompress(jobStr)
	fmt.Println(c.parseJson())
}

func TestDo(t *testing.T) {
	jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200}"
	f := NewCompress(jobStr)
	fmt.Println(f.Do())
}
