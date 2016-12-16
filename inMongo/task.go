package inMongo

type Task struct {
	columData string
}

func NewTask(jsonString string) (t *Task) {
	var tR Task
	tR.columData = jsonString
	return &tR
}
