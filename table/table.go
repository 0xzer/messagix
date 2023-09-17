package table

import (
	"fmt"
	"log"
)

/*
	Unknown fields = I don't know what type it is supposed to be, because I only see 9 which is undefined
	Trial and error works ^ check console for failed conversation within the decoder
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
	LSRemoveTask []LSRemoveTask
	LSTaskExists []LSTaskExists
	LSDeleteThenInsertContact []LSDeleteThenInsertContact
	LSUpdateTypingIndicator []LSUpdateTypingIndicator
	LSCheckAuthoritativeMessageExists []LSCheckAuthoritativeMessageExists
	LSMoveThreadToInboxAndUpdateParent []LSMoveThreadToInboxAndUpdateParent
	LSUpdateThreadSnippet []LSUpdateThreadSnippet
	LSVerifyThreadExists []LSVerifyThreadExists
	LSBumpThread []LSBumpThread
	LSUpdateParticipantLastMessageSendTimestamp []LSUpdateParticipantLastMessageSendTimestamp
	LSInsertMessage []LSInsertMessage
	LSUpsertGradientColor []LSUpsertGradientColor
	LSUpsertTheme []LSUpsertTheme
	LSInsertStickerAttachment []LSInsertStickerAttachment
}

var SPTable = map[string]string{
	"mciTraceLog": "LSMciTraceLog",
	"executeFirstBlockForSyncTransaction": "LSExecuteFirstBlockForSyncTransaction",
	"updateThreadsRangesV2": "LSUpdateThreadsRangesV2",
	"upsertSyncGroupThreadsRange": "LSUpsertSyncGroupThreadsRange",
	"upsertFolderSeenTimestamp": "LSUpsertFolderSeenTimestamp",
	"setHMPSStatus": "LSSetHMPSStatus",
	"upsertSequenceId": "LSUpsertSequenceId",
	"executeFinallyBlockForSyncTransaction": "LSExecuteFinallyBlockForSyncTransaction",
	"verifyContactRowExists": "LSVerifyContactRowExists",
	"taskExists": "LSTaskExists",
	"removeTask": "LSRemoveTask",
	"deleteThenInsertContact": "LSDeleteThenInsertContact",
	"updateTypingIndicator": "LSUpdateTypingIndicator",
	"checkAuthoritativeMessageExists": "LSCheckAuthoritativeMessageExists",
	"moveThreadToInboxAndUpdateParent": "LSMoveThreadToInboxAndUpdateParent",
	"updateThreadSnippet": "LSUpdateThreadSnippet",
	"verifyThreadExists": "LSVerifyThreadExists",
	"updateReadReceipt": "LSUpdateReadReceipt",
	"setForwardScore": "LSSetForwardScore",
	"bumpThread": "LSBumpThread",
	"updateParticipantLastMessageSendTimestamp": "LSUpdateParticipantLastMessageSendTimestamp",
	"insertMessage": "LSInsertMessage",
	"upsertTheme": "LSUpsertTheme",
	"upsertGradientColor": "LSUpsertGradientColor",
	"insertStickerAttachment": "LSInsertStickerAttachment",
}

func SPToDepMap(sp []string) map[string]string {
	m := make(map[string]string, 0)
	for _, d := range sp {
		depName, ok := SPTable[d]
		if !ok {
			log.Println(fmt.Sprintf("can't convert sp %s to dependency name because it wasn't found in the SPTable", d))
			continue
		}
		m[d] = depName
	}
	return m
}