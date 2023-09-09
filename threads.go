package messagix

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/0xzer/messagix/methods"
	"github.com/0xzer/messagix/table"
	"github.com/0xzer/messagix/tasks"
)


type Threads struct {
	client *Client
}

type MessageBuilder struct {
	client *Client
	payload *tasks.SendMessageTask
	readPayload *tasks.ThreadMarkRead
}

func (t *Threads) NewMessageBuilder(threadId int64) *MessageBuilder {
	return &MessageBuilder{
		client: t.client,
		payload: &tasks.SendMessageTask{
			ThreadId: threadId,
			SkipUrlPreviewGen: 0,
			TextHasLinks: 0,
		},
		readPayload: &tasks.ThreadMarkRead{
			ThreadId: threadId,
		},
	}
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

	log.Println("packetId:", packetId)

	resp := m.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return nil, fmt.Errorf("failed to receive response from socket while trying to send message. packetId: %d", packetId)
	}

	m.client.Logger.Debug().Any("resp", resp).Msg("got response!")
	os.Exit(1)
	return nil, nil
}