package table

/*
	Unknown fields = I don't know what type it is supposed to be, because I only see 9 which is undefined
*/

type LSTable struct {
	LSMciTraceLog []LSMciTraceLog
	LSExecuteFirstBlockForSyncTransaction []LSExecuteFirstBlockForSyncTransaction
	LSTruncateMetadataThreads []LSTruncateMetadataThreads
	LSTruncateThreadRangeTablesForSyncGroup []LSTruncateThreadRangeTablesForSyncGroup
	LSUpsertSyncGroupThreadsRange []LSUpsertSyncGroupThreadsRange
	LSUpsertInboxThreadsRange []LSUpsertInboxThreadsRange
	LSUpdateThreadsRangesV2 []LSUpdateThreadsRangesV2
	LSUpsertFolderSeenTimestamp []LSUpsertFolderSeenTimestamp
	LSSetHMPSStatus []LSSetHMPSStatus
	LSTruncateTablesForSyncGroup []LSTruncateTablesForSyncGroup
	LSDeleteThenInsertThread []LSDeleteThenInsertThread
	LSAddParticipantIdToGroupThread []LSAddParticipantIdToGroupThread
	LSClearPinnedMessages []LSClearPinnedMessages
	LSWriteThreadCapabilities []LSWriteThreadCapabilities
	LSUpsertMessage []LSUpsertMessage
	LSSetForwardScore []LSSetForwardScore
	LSSetMessageDisplayedContentTypes []LSSetMessageDisplayedContentTypes
	LSUpdateReadReceipt []LSUpdateReadReceipt
	LSInsertNewMessageRange []LSInsertNewMessageRange
	LSDeleteExistingMessageRanges []LSDeleteExistingMessageRanges
	LSUpsertSequenceId []LSUpsertSequenceId
	LSVerifyContactRowExists []LSVerifyContactRowExists
	LSThreadsRangesQuery []LSThreadsRangesQuery
	LSSetRegionHint []LSSetRegionHint
	LSExecuteFinallyBlockForSyncTransaction []LSExecuteFinallyBlockForSyncTransaction
}