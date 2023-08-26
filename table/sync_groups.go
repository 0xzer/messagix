package table

type LSTruncateTablesForSyncGroup struct {
	SyncGroup int64 `index:"0"`
}

type LSTruncateThreadRangeTablesForSyncGroup struct {
	ParentThreadKey int64 `index:"0"`
}

type LSUpsertSyncGroupThreadsRange struct {
	SyncGroup int64 `index:"0"`
	ParentThreadKey int64 `index:"1"`
	MinLastActivityTimestampMs int64 `index:"2"`
	HasMoreBefore bool `index:"3"`
	IsLoadingBefore bool `index:"4"`
	MinThreadKey int64 `index:"5"`
}