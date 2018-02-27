package task

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	jobStr := "{\"path\":\"/Users/Snow/img/data/attachment/forum/201802/12/1518412896316.png\",\"width\":200,\"height\":200}"
	f := NewTask(jobStr)
	fmt.Println(f.Do())
}
