package socket

type SyncChannel int64
const (
	MailBox SyncChannel = 1
	Contact SyncChannel = 2
)

type QueryMetadata struct {
	DatabaseId int64
	SendSyncParams bool
	LastAppliedCursor interface{}
	SyncParams interface{}
	SyncChannel 
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

func (t *SyncGroupsTask) Create() (interface{}, interface{}, bool) {
	queueName := "trq"
	return t, queueName, false
}