package socket

import (
	"strconv"
	"github.com/0xzer/messagix/table"
)

type SendMessageTask struct {
	ThreadId int64 `json:"thread_id"`
	Otid string `json:"otid"`
	Source table.ThreadSourceType `json:"source"`
	SendType table.SendType `json:"send_type"`
	SyncGroup int64 `json:"sync_group"`
	Text string `json:"text,omitempty"`
	InitiatingSource table.InitiatingSource `json:"initiating_source"`
	SkipUrlPreviewGen int32 `json:"skip_url_preview_gen"` // 0 or 1
	TextHasLinks int32 `json:"text_has_links"` // 0 or 1
}

func (t *SendMessageTask) GetLabel() string {
	return TaskLabels["SendMessageTask"]
}

func (t *SendMessageTask) Create() (interface{}, interface{}) {
	queueName := strconv.Itoa(int(t.ThreadId))
	return t, queueName
}

type ThreadMarkReadTask struct {
	ThreadId            int64 `json:"thread_id"`
	LastReadWatermarkTs int64 `json:"last_read_watermark_ts"`
	SyncGroup           int64   `json:"sync_group"`
}

func (t *ThreadMarkReadTask) GetLabel() string {
	return TaskLabels["ThreadMarkRead"]
}

func (t *ThreadMarkReadTask) Create() (interface{}, interface{}) {
	queueName := strconv.Itoa(int(t.ThreadId))
	return t, queueName
}

type FetchMessagesTask struct {
	ThreadKey int64 `json:"thread_key"`
	Direction int64 `json:"direction"` // 0
	ReferenceTimestampMs int64 `json:"reference_timestamp_ms"`
	ReferenceMessageId string `json:"reference_message_id"`
	SyncGroup int64 `json:"sync_group"` // 1
	Cursor string `json:"cursor"`
}

func (t *FetchMessagesTask) GetLabel() string {
	return TaskLabels["FetchMessagesTask"]
}

func (t *FetchMessagesTask) Create() (interface{}, interface{}) {
	threadStr := strconv.Itoa(int(t.ThreadKey))
	queueName := "mrq." + threadStr
	return t, queueName
}