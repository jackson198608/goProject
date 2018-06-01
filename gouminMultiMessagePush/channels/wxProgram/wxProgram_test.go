package wxProgram

import (
	"testing"
)

func newTask() *Task{
	msg := `{"touser":"oaBrW5dzQeHbzHFE9a2GZEQI65zE","template_id":"t7OVsxS7ipiiBgAU64HjAiRu82yCjdkXqZc9HydUgDU","form_id":"8c34de924242cc60543743952cde8210","data":{"keyword1":{"value":"\u60a8\u7684\u79c1\u623f\u7167\u88ab\u70b9\u8d5e\u4e86","color":"#000000"},"keyword2":{"value":"\u5f20\u4e09\u7684\u79c1\u623f\u7167","color":"#000000"},"keyword3":{"value":12,"color":"#000000"},"keyword4":{"value":"\u70b9\u51fb\u8fdb\u5165\u5c0f\u7a0b\u5e8f\u67e5\u770b","color":"#000000"}},"page":"pages\/album\/list\/list","color":"#000000","emphasis_keyword":"keyword1.DATA"}`
	redisString := "2|meowStar|"+msg
	t := NewTask(redisString)
	return t
}

func TestParseResult(t *testing.T) {

	f := newTask()
	f.SendRequest()
}
