package main

type Task struct {
	columData string
}

func NewMongoTask(jsonString string) (t *Task) {
	var tR Task
	tR.columData = jsonString
	return &tR
}
