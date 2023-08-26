package table

type LSMciTraceLog struct {
	SomeInt0 int64 `index:"0"`
	MCITraceUnsampledEventTraceId string `index:"1"`
	Unknown2 interface{} `index:"2"`
	SomeInt3 int64 `index:"3"`
	Unknown4 interface{} `index:"4"`
	DatascriptExecute string `index:"5"`
	SomeInt6 int64 `index:"6"`
}

type LSExecuteFirstBlockForSyncTransaction struct {
	DatabaseId int64 `index:"0"`
	EpochId int64 `index:"1"`
	CurrentCursor string `index:"3"`
	SyncStatus int64 `index:"4"`
	SendSyncParams bool `index:"5"`
	MinTimeToSyncTimestampMs int64 `index:"6"` // fix this, use conditionIndex
	CanIgnoreTimestamp bool `index:"7"`
	SyncChannel int64 `index:"8"`
}

type LSExecuteFinallyBlockForSyncTransaction struct {
	Unknown0 bool `index:"0"`
	Unknown1 int64 `index:"1"`
	Unknown2 int64 `index:"2"`
}

type LSSetHMPSStatus struct {
	AccountId int64 `index:"0"`
	Unknown1 int64 `index:"1"`
	Timestamp int64 `index:"2"`
}

type LSUpsertSequenceId struct {
	LastAppliedMailboxSequenceId int64 `index:"0"`
}

type LSSetRegionHint struct {
	Unknown0 int64 `index:"0"`
	RegionHint string `index:"1"`
}