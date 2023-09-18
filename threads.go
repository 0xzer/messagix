package messagix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/0xzer/messagix/methods"
	"github.com/0xzer/messagix/table"
	"github.com/0xzer/messagix/socket"
)

type Threads struct {
	client *Client
}

func (t *Threads) FetchMessages(ThreadId int64, ReferenceTimestampMs int64, ReferenceMessageId string, Cursor string) (*table.LSTable, error) {
	tskm := t.client.NewTaskManager()
	tskm.AddNewTask(&socket.FetchMessagesTask{ThreadKey: ThreadId, Direction: 0, ReferenceTimestampMs: ReferenceTimestampMs, ReferenceMessageId: ReferenceMessageId, SyncGroup: 1, Cursor: Cursor})

	payload, err := tskm.FinalizePayload()
	if err != nil {
		log.Fatal(err)
	}

	packetId, err := t.client.socket.makeLSRequest(payload, 3)
	if err != nil {
		log.Fatal(err)
	}


	resp := t.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return nil, fmt.Errorf("failed to receive response from socket while trying to fetch messages. (packetId=%d, thread_key=%d, cursor=%s, reference_message_id=%s, reference_timestamp_ms=%d)", packetId, ThreadId, Cursor, ReferenceMessageId, ReferenceTimestampMs)
	}

	return &resp.Table, nil
}

type MessageBuilder struct {
	client *Client
	payload *socket.SendMessageTask
	readPayload *socket.ThreadMarkReadTask
}

func (t *Threads) NewMessageBuilder(threadId int64) *MessageBuilder {
	return &MessageBuilder{
		client: t.client,
		payload: &socket.SendMessageTask{
			ThreadId: threadId,
			SkipUrlPreviewGen: 0,
			TextHasLinks: 0,
		},
		readPayload: &socket.ThreadMarkReadTask{
			ThreadId: threadId,
		},
	}
}

func (m *MessageBuilder) SetReplyMetadata(replyMetadata *socket.ReplyMetaData) *MessageBuilder {
	m.payload.ReplyMetaData = replyMetadata
	return m
}

func (m *MessageBuilder) SetSource(source table.ThreadSourceType) *MessageBuilder {
	m.payload.Source = source
	return m
}

func (m *MessageBuilder) SetInitiatingSource(initatingSource table.InitiatingSource) *MessageBuilder {
	m.payload.InitiatingSource = initatingSource
	return m
}

func (m *MessageBuilder) SetSendType(sendType table.SendType) *MessageBuilder {
	m.payload.SendType = sendType
	return m
}


func (m *MessageBuilder) SetSyncGroup(syncGroup int64) *MessageBuilder {
	m.payload.SyncGroup = syncGroup
	m.readPayload.SyncGroup = syncGroup
	return m
}

func (m *MessageBuilder) SetSkipUrlPreviewGen() *MessageBuilder {
	m.payload.SkipUrlPreviewGen = 1
	return m
}

func (m *MessageBuilder) SetTextHasLinks() *MessageBuilder {
	m.payload.TextHasLinks = 1
	return m
}

func (m *MessageBuilder) SetText(text string) *MessageBuilder {
	m.payload.Text = text
	return m
}

func (m *MessageBuilder) SetLastReadWatermarkTs(ts int64) *MessageBuilder {
	m.readPayload.LastReadWatermarkTs = ts
	return m
}

func (m *MessageBuilder) Execute() (error, error){
	tskm := m.client.NewTaskManager()
	m.payload.Otid = strconv.Itoa(int(methods.GenerateEpochId()))
	tskm.AddNewTask(m.payload)
	tskm.AddNewTask(m.readPayload)
	tskm.setTraceId(methods.GenerateTraceId())

	payload, err := tskm.FinalizePayload()
	if err != nil {
		log.Fatal(err)
	}

	packetId, err := m.client.socket.makeLSRequest(payload, 3)
	if err != nil {
		log.Fatal(err)
	}

	resp := m.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return nil, fmt.Errorf("failed to receive response from socket while trying to send message. packetId: %d", packetId)
	}

	return nil, nil
}