package messagix

import (
	"encoding/json"
	"fmt"
	"github.com/0xzer/messagix/methods"
	"github.com/0xzer/messagix/socket"
)

type DatabaseManager struct {
	client *Client
	currQueries []socket.DatabaseQuery
}

func (c *Client) NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		client: c,
		currQueries: make([]socket.DatabaseQuery, 0),
	}
}

func (db *DatabaseManager) PublishQueries() error {
	if len(db.currQueries) < 1 {
		return fmt.Errorf("failed to publish database queries in dbmanager: currQueries must have at least 1 query")
	}

	for _, q := range db.currQueries {
		metaData, ok := socket.Queries[q.Database]
		if !ok {
			return fmt.Errorf("invalid database_id: %d", q.Database)
		}
		var t int
		q.EpochId = methods.GenerateEpochId()
		q.Version = db.client.configs.siteConfig.VersionId //db.client.configs.siteConfig.VersionId
		if metaData.SendSyncParams && q.SyncParams == nil{
			params, err := db.client.configs.getSyncParams()
			if err != nil {
				return fmt.Errorf("failed to marshal sync_params into bytes from database query: %e", err)
			}
			t = 1
			q.SyncParams = string(params)
		} else {
			t = 2
			q.LastAppliedCursor = db.client.configs.getLastAppliedCursor(q.Database)
		}

		jsonStr, err := json.Marshal(&q)
		if err != nil {
			return fmt.Errorf("failed to marshal database query into json string: %e", err)
		}
		_, err = db.client.socket.makeLSRequest(jsonStr, t)
		if err != nil {
			return fmt.Errorf("failed to send database query request through socket: %e", err)
		}
		//err := db.client.socket.send
	}
	db.currQueries = make([]socket.DatabaseQuery, 0)
	return nil
}

func (db *DatabaseManager) AddNewQueries(queries []socket.DatabaseQuery) {
	db.currQueries = append(db.currQueries, queries...)
}

func (db *DatabaseManager) AddNewQuery(query socket.DatabaseQuery) {
	db.currQueries = append(db.currQueries, query)
}

/*
	The initial database queries being sent, to "sync" the socket to the client so that the server doesn't send unnecessary updates i'm guessing.
*/
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
	db.AddNewQueries(queries)
}