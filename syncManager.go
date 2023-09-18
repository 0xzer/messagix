package messagix

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/0xzer/messagix/graphql"
	"github.com/0xzer/messagix/methods"
	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/socket"
	"github.com/0xzer/messagix/table"
)

type SyncManager struct {
	client *Client
	store map[int64]*socket.QueryMetadata
	syncParams *modules.LSPlatformMessengerSyncParams
}

func (c *Client) NewSyncManager() *SyncManager {
	return &SyncManager{
		client: c,
		store: map[int64]*socket.QueryMetadata{
			1: { SendSyncParams: false, SyncChannel: socket.MailBox },
			2: { SendSyncParams: false, SyncChannel: socket.Contact },
			5: { SendSyncParams: true },
			16: { SendSyncParams: true },
			26: { SendSyncParams: true },
			28: { SendSyncParams: true },
			95: { SendSyncParams: false, SyncChannel: socket.Contact },
			104: { SendSyncParams: true },
			140: { SendSyncParams: true },
			141: { SendSyncParams: true },
			142: { SendSyncParams: true },
			143: { SendSyncParams: true },
			196: { SendSyncParams: true },
			198: { SendSyncParams: true },
		},
		syncParams: &modules.LSPlatformMessengerSyncParams{},
	}
}

func (sm *SyncManager) EnsureSyncedSocket(databases []int64) error {
	for _, db := range databases {
		database, ok := sm.store[db]
		if !ok {
			return fmt.Errorf("could not find sync store for database: %d", db)
		}

		_, err := sm.SyncSocketData(db, database)
		if err != nil {
			return fmt.Errorf("failed to ensure database is synced through socket: (databaseId=%d)", db)
		}
	}

	return nil
}

func (sm *SyncManager) SyncSocketData(databaseId int64, db *socket.QueryMetadata) (*table.LSTable, error) {
	var t int
	payload := &socket.DatabaseQuery{
		Database: databaseId,
		Version: sm.client.configs.siteConfig.VersionId,
		EpochId: methods.GenerateEpochId(),
	}

	if db.SendSyncParams {
		t = 1
		payload.SyncParams = sm.getSyncParams(db.SyncChannel)
	} else {
		t = 2
		payload.LastAppliedCursor = db.LastAppliedCursor
	}

	jsonPayload, err := json.Marshal(&payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DatabaseQuery struct into json bytes (databaseId=%d): %e", databaseId, err)
	}

	packetId, err := sm.client.socket.makeLSRequest(jsonPayload, t)
	if err != nil {
		return nil, fmt.Errorf("failed to make lightspeed socket request with DatabaseQuery byte payload (databaseId=%d): %e", databaseId, err)
	}

	resp := sm.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return nil, fmt.Errorf("timed out while waiting for sync response from socket (databaseId=%d)", databaseId)
	}
	resp.Finish()

	block := resp.Table.LSExecuteFirstBlockForSyncTransaction[0]
	nextCursor, currentCursor := block.NextCursor, block.CurrentCursor

	if nextCursor == currentCursor {
		return &resp.Table, nil
	}

	// Update the last applied cursor to the next cursor and recursively fetch again
	db.LastAppliedCursor = nextCursor
	db.SendSyncParams = block.SendSyncParams
	db.SyncChannel = socket.SyncChannel(block.SyncChannel)
	err = sm.SyncTransactions(resp.Table.LSExecuteFirstBlockForSyncTransaction) // Also sync the transaction with the store map because the db param is just a copy of the map entry
	if err != nil {
		return nil, err
	}

	return sm.SyncSocketData(databaseId, db)
}

func (sm *SyncManager) SyncDataGraphQL(dbs []int64) (*table.LSTable, error) {
	var tableData *table.LSTable
	for _, db := range dbs {
		database, ok := sm.store[db]
		if !ok {
			return nil, fmt.Errorf("could not find sync store for database: %d", db)
		}
		
		variables := &graphql.LSPlatformGraphQLLightspeedVariables{
			Database: int(db),
			LastAppliedCursor: database.LastAppliedCursor,
			Version: sm.client.configs.siteConfig.VersionId,
			EpochID: 0,
		}
		if database.SendSyncParams {
			variables.SyncParams = sm.getSyncParams(database.SyncChannel)
		}

		lsTable, err := sm.client.graphQl.makeLSRequest(variables, 1)
		if err != nil {
			return nil, err
		}

		if db == 1 {
			tableData = lsTable
		}
		sm.SyncTransactions(lsTable.LSExecuteFirstBlockForSyncTransaction)
	}

	return tableData, nil
}

func (sm *SyncManager) SyncTransactions(transactions []table.LSExecuteFirstBlockForSyncTransaction) error {
	for _, transaction := range transactions {
		database, ok := sm.store[transaction.DatabaseId]
		if !ok {
			return fmt.Errorf("failed to update database %d by block transaction", transaction.DatabaseId)
		}
	
		database.LastAppliedCursor = transaction.NextCursor
		database.SendSyncParams = transaction.SendSyncParams
		database.SyncChannel = socket.SyncChannel(transaction.SyncChannel)
		//sm.client.Logger.Info().Any("new_cursor", database.LastAppliedCursor).Any("syncChannel", database.SyncChannel).Any("sendSyncParams", database.SendSyncParams).Any("database_id", transaction.DatabaseId).Msg("Updated database by transaction...")
	}

	return nil
}

func (sm *SyncManager) UpdateDatabaseSyncParams(dbs []*socket.QueryMetadata) error {
	for _, db := range dbs {
		database, ok := sm.store[db.DatabaseId]
		if !ok {
			return fmt.Errorf("failed to update sync params for database: %d", db.DatabaseId)
		}
		database.SendSyncParams = db.SendSyncParams
		database.SyncChannel = db.SyncChannel
	}
	return nil
}

func (sm *SyncManager) UpdateDatabaseCursor(dbs []*socket.QueryMetadata) error {
	for _, db := range dbs {
		database, ok := sm.store[db.DatabaseId]
		if !ok {
			return fmt.Errorf("failed to update cursor for database: %d", db.DatabaseId)
		}
		database.LastAppliedCursor = db.LastAppliedCursor
	}
	return nil
}

func (sm *SyncManager) getSyncParams(ch socket.SyncChannel) interface{} {
	switch ch {
	case socket.MailBox:
		return sm.syncParams.Mailbox
	case socket.Contact:
		return sm.syncParams.Contact
	default:
		log.Fatalf("Unknown syncChannel: %d", ch)
		return nil
	}
}

func (sm *SyncManager) getCursor(db int64) interface{} {
	database, ok := sm.store[db]
	if !ok {
		return nil
	}
	return database.LastAppliedCursor
}
/*
func (db *DatabaseManager) AddInitQueries() {
	queries := []socket.DatabaseQuery{
		{Database: 1},
		{Database: 2},
		{Database: 95},
		{Database: 16},
		{Database: 26},
		{Database: 28},
		{Database: 95},
		{Database: 104},
		{Database: 140},
		{Database: 141},
		{Database: 142},
		{Database: 143},
		{Database: 196},
		{Database: 198},
		
	}
}
*/