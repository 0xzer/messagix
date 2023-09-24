package messagix

import (
	"fmt"
	"log"
	"github.com/0xzer/messagix/socket"
	"github.com/0xzer/messagix/table"
)

type Messages struct {
	client *Client
}

func (m *Messages) SendReaction(threadId int64, messageId string, reaction string) (*table.LSTable, error) {
	tskm := m.client.NewTaskManager()
	tskm.AddNewTask(&socket.SendReactionTask{
		ThreadKey: threadId,
		MessageID: messageId,
		ActorID: m.client.configs.siteConfig.AccountIdInt,
		Reaction: reaction,
		SendAttribution: table.MESSENGER_INBOX_IN_THREAD,
	})

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
		return nil, fmt.Errorf("failed to receive response from socket after sending reaction. packetId: %d", packetId)
	}
	resp.Finish()

	return &resp.Table, nil
}