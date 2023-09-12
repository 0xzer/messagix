package socket

type QueryMetadata struct {
	SendSyncParams bool
}

var Queries = map[int64]QueryMetadata{
	1: { SendSyncParams: false },
	2: { SendSyncParams: false },
	5: { SendSyncParams: true },
	16: { SendSyncParams: true },
	26: { SendSyncParams: true },
	28: { SendSyncParams: true },
	95: { SendSyncParams: false },
	104: { SendSyncParams: true },
	140: { SendSyncParams: true },
	141: { SendSyncParams: true },
	142: { SendSyncParams: true },
	143: { SendSyncParams: true },
	196: { SendSyncParams: true },
	198: { SendSyncParams: true },
}

type SyncGroupsTask struct {
	IsAfter                    int    `json:"is_after"`
	ParentThreadKey            int    `json:"parent_thread_key"`
	ReferenceThreadKey         int    `json:"reference_thread_key"`
	ReferenceActivityTimestamp int64  `json:"reference_activity_timestamp"`
	AdditionalPagesToFetch     int    `json:"additional_pages_to_fetch"`
	Cursor                     interface{} `json:"cursor"`
	MessagingTag               interface{}    `json:"messaging_tag"`
	SyncGroup                  int    `json:"sync_group"`
}

func (t *SyncGroupsTask) GetLabel() string {
	return TaskLabels["SyncGroupsTask"]
}

func (t *SyncGroupsTask) Create() (interface{}, interface{}) {
	queueName := "trq"
	return t, queueName
}