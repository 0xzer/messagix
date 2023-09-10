package messagix

import (
	"fmt"
	"log"

	"github.com/0xzer/messagix/table"
	"github.com/0xzer/messagix/tasks"
	"github.com/google/uuid"
)

type Account struct {
	client *Client
}

func (a *Account) GetContacts(limit int64) ([]table.LSVerifyContactRowExists, error) {
	tskm := a.client.NewTaskManager()
	tskm.AddNewTask(&tasks.GetContactsTask{Limit: limit})

	payload, err := tskm.FinalizePayload()
	if err != nil {
		log.Fatal(err)
	}

	packetId, err := a.client.socket.makeLSRequest(payload, 3)
	if err != nil {
		log.Fatal(err)
	}

	resp := a.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return nil, fmt.Errorf("failed to receive response from socket while trying to fetch contacts. packetId: %d", packetId)
	}

	return resp.Table.LSVerifyContactRowExists, nil
}

func (a *Account) GetContactsFull(contactIds []int64) ([]table.LSDeleteThenInsertContact, error) {
	tskm := a.client.NewTaskManager()
	for _, id := range contactIds {
		tskm.AddNewTask(&tasks.GetContactsFullTask{
			ContactId: id,
		})
	}

	payload, err := tskm.FinalizePayload()
	if err != nil {
		log.Fatal(err)
	}

	packetId, err := a.client.socket.makeLSRequest(payload, 3)
	if err != nil {
		log.Fatal(err)
	}

	resp := a.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return nil, fmt.Errorf("failed to receive response from socket while trying to fetch full contact information. packetId: %d", packetId)
	}

	return resp.Table.LSDeleteThenInsertContact, nil
}

func (a *Account) ReportAppState(state table.AppState) error {
	tskm := a.client.NewTaskManager()
	tskm.AddNewTask(&tasks.ReportAppStateTask{AppState: state, RequestId: uuid.NewString()})

	payload, err := tskm.FinalizePayload()
	if err != nil {
		log.Fatal(err)
	}

	packetId, err := a.client.socket.makeLSRequest(payload, 3)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("packetId:", packetId)

	resp := a.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return fmt.Errorf("failed to receive response from socket while trying to report app state. packetId: %d", packetId)
	}

	a.client.Logger.Info().Any("data", resp).Msg("Got report app state resp app state")
	return nil
}