package table

type LSTaskExists struct {
	TaskId int64 `index:"0"`
}

type LSRemoveTask struct {
	TaskId int64 `index:"0"`
}

type LSUpdateTaskQueueName struct {
	QueueNameTaskId string `index:"0"`
	QueueName string `index:"1"`
}

type LSUpdateTaskValue struct {
	QueueNameTaskId string `index:"0"`
	/*
		b.taskValue.split(a[1]).join(a[2]) // b = curr obj value
	*/
	TaskValue1 string `index:"1"`
	TaskValue2 string `index:"2"`
}

type LSMailboxTaskCompletionApiOnTaskCompletion struct {
	TaskId int64 `index:"0"`
	Success bool `index:"1"`
}