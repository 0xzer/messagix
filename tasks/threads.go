package tasks

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

type ThreadMarkRead struct {
	ThreadId            int64 `json:"thread_id"`
	LastReadWatermarkTs int64 `json:"last_read_watermark_ts"`
	SyncGroup           int64   `json:"sync_group"`
}

func (t *ThreadMarkRead) GetLabel() string {
	return TaskLabels["ThreadMarkRead"]
}

func (t *ThreadMarkRead) Create() (interface{}, interface{}) {
	queueName := strconv.Itoa(int(t.ThreadId))
	return t, queueName
}