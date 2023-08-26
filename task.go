package messagix

import (
	"github.com/0xzer/messagix/methods"
)

type TaskData struct {
	FailureCount interface{} `json:"failure_count,omitempty"`
	Label string `json:"label,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
	QueueName interface{} `json:"queue_name,omitempty"`
	TaskId int64 `json:"task_id,omitempty"`
}

type TaskPayload struct {
	EpochId int64 `json:"epoch_id,omitempty"`
	Tasks []TaskData `json:"tasks,omitempty"`
	VersionId int64 `json:"version_id,omitempty"`
}

type TaskManager struct {
	client *Client
	activeTaskIds []int
	currTasks []TaskData
}

func (tm *TaskManager) FinalizePayload() *TaskPayload {
	p := &TaskPayload{
		EpochId: methods.GenerateEpochId(),
		Tasks: tm.currTasks,
		VersionId: tm.client.configs.siteConfig.VersionId,
	}
	tm.currTasks = make([]TaskData, 0)
	return p
}

func (tm *TaskManager) AddNewTask(task Task) {
	payload, queueName := task.create()
	label := task.getLabel()
	tm.client.Logger.Debug().Any("label", label).Any("payload", payload).Any("queueName", queueName).Msg("Creating task")
	taskData := TaskData{
		FailureCount: nil,
		Label: label,
		Payload: payload,
		QueueName: queueName,
		TaskId: tm.GetTaskId(),
	}
	
	tm.currTasks = append(tm.currTasks, taskData)
}

func (tm *TaskManager) GetTaskId() int64 {
	if len(tm.currTasks) == 0 {
		return 0
	}
	return int64(len(tm.currTasks))
}

type Task interface {
	getLabel() string
	create() (interface{}, interface{}) // payload, queue_name
}

type GetContactsTask struct {
	Limit int64 `json:"limit,omitempty"`
}

func (t *GetContactsTask) getLabel() string {
	return "452"
}

func (t *GetContactsTask) create() (interface{}, interface{}) {
	queueName := []string{"search_contacts", methods.GenerateTimestampString()}
	return t.Limit, queueName
}