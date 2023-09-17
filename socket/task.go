package socket

/*
	type 3 = task
*/

var TaskLabels = map[string]string{
	"GetContactsTask": "452",
	"SendMessageTask": "46",
	"ThreadMarkRead": "21",
	"GetContactsFullTask": "207",
	"ReportAppStateTask": "123",
	"SyncGroupsTask": "145", 
}

type Task interface {
	GetLabel() string
	Create() (interface{}, interface{}) // payload, queue_name
}

type TaskData struct {
	FailureCount interface{} `json:"failure_count"`
	Label string `json:"label,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
	QueueName interface{} `json:"queue_name,omitempty"`
	TaskId int64 `json:"task_id"`
}