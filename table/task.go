package table

type LSTaskExists struct {
	TaskId int64 `index:"0"`
}

type LSRemoveTask struct {
	TaskId int64 `index:"0"`
}