package database

import (
	"encoding/json"
	"log"

	"github.com/0xzer/messagix/methods"
)

/*
	type 1 = sync by cursor
	type 2 = sync by sync_params
*/

type SubscribeDatabase struct {
	AppID     string `json:"app_id"`
	Payload   interface{} `json:"payload"`
	RequestID int    `json:"request_id"`
	Type      int    `json:"type"`
}

type SubscribePayload struct {
	Database int64 `json:"database"`
	EpochID           int64  `json:"epoch_id"`
	FailureCount      interface{}   `json:"failure_count"`
	LastAppliedCursor interface{} `json:"last_applied_cursor"`
	SyncParams        interface{}    `json:"sync_params"`
	Version           int64  `json:"version"`
}

func (s *SubscribeDatabase) Create() ([]byte, error) {
	return json.Marshal(&s)
}

func NewSubscribeDatabase(appId string, database int64, lastAppliedCursor interface{}, syncParams interface{}, version int64, requestId int64, t int64) ([]byte, error) {
	payload := &SubscribePayload{
		Database: database,
		EpochID: methods.GenerateEpochId(),
		LastAppliedCursor: lastAppliedCursor,
		SyncParams: syncParams,
		Version: version,
	}
	payloadStr, err := json.Marshal(&payload)
	if err != nil {
		log.Fatal(err)
	}
	db := &SubscribeDatabase{
		AppID: appId,
		Payload: string(payloadStr),
		RequestID: int(requestId),
		Type: int(t),
	}
	return json.Marshal(&db)
}