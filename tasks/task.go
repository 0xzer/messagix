package tasks

/*
	type 1 = sync by cursor
	type 2 = sync by sync_params
	type 3 = task
*/

var TaskLabels = map[string]string{
	"GetContactsTask": "452",
	"SendMessageTask": "46",
	"ThreadMarkRead": "21",
	"GetContactsFullTask": "207",
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

type TaskPayload struct {
	EpochId int64 `json:"epoch_id"`
	DataTraceId string `json:"data_trace_id,omitempty"`
	Tasks []TaskData `json:"tasks,omitempty"`
	VersionId string `json:"version_id"`
}