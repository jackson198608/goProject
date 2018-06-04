package wxProgram

import (
	"testing"
)

func newTask() *Task{
	appid := "wxe691667e66d1f8eb"
	token := "10_jXOtwQOWGP5WUge_qt5guF6QSsIJau7eWAcPfweypuMa3q8GYgXRdyF3EJC-B_fozr7nIh33mQYT9Fcua3lg3aMYH1iRZyR_b7UlvWmp4pYi4wI5wqsluCFugE28E5t0WNa6hjDQQXuEU6snICFiAJAUOR"
	jobStr := `wxe691667e66d1f8eb|{"touser":"oaBrW5dzQeHbzHFE9a2GZEQI65zE","template_id":"t7OVsxS7ipiiBgAU64HjAiRu82yCjdkXqZc9HydUgDU","form_id":"4f5993a27a227a069d82d7812d20d296","data":{"keyword1":{"value":"\u60a8\u7684\u79c1\u623f\u7167\u88ab\u70b9\u8d5e\u4e86","color":"#000000"},"keyword2":{"value":"\u5f20\u4e09\u7684\u79c1\u623f\u7167","color":"#000000"},"keyword3":{"value":10,"color":"#000000"},"keyword4":{"value":"\u70b9\u51fb\u8fdb\u5165\u5c0f\u7a0b\u5e8f\u67e5\u770b","color":"#000000"}},"page":"pages\/album\/list\/list","color":"#000000","emphasis_keyword":""}`

	t := NewTask(appid,token,jobStr)
	return t
}

func TestParseResult(t *testing.T) {

	f := newTask()
	f.SendRequest()
}
