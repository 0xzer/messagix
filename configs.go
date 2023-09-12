package messagix

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/0xzer/messagix/crypto"
	"github.com/0xzer/messagix/methods"
	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/socket"
	"github.com/0xzer/messagix/types"
)

type Configs struct {
	client *Client
	needSync bool
	mqttConfig *types.MQTTConfig
	siteConfig *types.SiteConfig
	syncCursors map[int64]string
}

func (c *Configs) SetupConfigs() error {
	schedulerJS := modules.SchedulerJSDefined
	if schedulerJS.SiteData == (modules.SiteData{}) {
		log.Fatalf("SetupConfigs was somehow called before modules were initalized")
	}
	
	c.mqttConfig = &types.MQTTConfig{
		ProtocolName: "MQIsdp",
		ProtocolLevel: 3,
		ClientId: "mqttwsclient",
		Broker: schedulerJS.MqttWebConfig.Endpoint,
		KeepAliveTimeout: 15,
		SessionId: methods.GenerateSessionId(),
		AppId: schedulerJS.MqttWebConfig.AppID,
		ClientCapabilities: schedulerJS.MqttWebConfig.ClientCapabilities,
		Capabilities: schedulerJS.MqttWebConfig.Capabilities,
		ChatOn: schedulerJS.MqttWebConfig.ChatVisibility,
		SubscribedTopics: schedulerJS.MqttWebConfig.SubscribedTopics,
		ConnectionType: "websocket",
		HostNameOverride: "",
		Cid: schedulerJS.MqttWebDeviceID.ClientID,
	}

	bitmap := crypto.NewBitmap().Update(modules.Bitmap).ToCompressedString()
	csrBitmap := crypto.NewBitmap().Update(modules.CsrBitmap).ToCompressedString()

	eqmcQuery, err := modules.JSONData.Eqmc.ParseAjaxURLData()
	if err != nil {
		log.Fatalf("failed to parse AjaxURLData from eqmc json struct: %e", err)
	}

	c.siteConfig = &types.SiteConfig{
		AccountId: schedulerJS.CurrentUserInitialData.AccountID,
		Bitmap: bitmap,
		CSRBitmap: csrBitmap,
		HasteSessionId: schedulerJS.SiteData.Hsi,
		WebSessionId: methods.GenerateWebsessionID(),
		CometReq: eqmcQuery.CometReq,
		LsdToken: schedulerJS.LSD.Token,
		SpinT: strconv.Itoa(schedulerJS.SiteData.SpinT),
		SpinB: schedulerJS.SiteData.SpinB,
		SpinR: strconv.Itoa(schedulerJS.SiteData.SpinR),
		FbDtsg: schedulerJS.DTSGInitialData.Token,
		Jazoest: eqmcQuery.Jazoest,
		Pr: strconv.Itoa(schedulerJS.SiteData.Pr),
		HasteSession: schedulerJS.SiteData.HasteSession,
		ConnectionClass: schedulerJS.WebConnectionClassServerGuess.ConnectionClass,
		VersionId: modules.VersionId,
		Locale: schedulerJS.IntlCurrentLocale.Code,
		X_ASDB_ID: "129477", // hold off on this and check if it ever changes, if so we gotta load the js file and extract it from there
	}

	c.client.graphQl = &GraphQL{client: c.client, lsRequests: 0, graphQLRequests: 1}
	if c.needSync {
		lsData, err := c.client.graphQl.makeLSRequest(1, 0, 1)
		if err != nil {
			log.Fatalf("failed to sync lightspeed data: %e", err)
		}

		modules.SchedulerJSRequired.LSTable = lsData
	}

	for _, block := range modules.SchedulerJSRequired.LSTable.LSExecuteFirstBlockForSyncTransaction {
		c.syncCursors[block.DatabaseId] = block.CurrentCursor
	}

	c.client.Logger.Info().Any("value", c.siteConfig.Bitmap.CompressedStr).Msg("Loaded __dyn bitmap")
	c.client.Logger.Info().Any("value", c.siteConfig.CSRBitmap.CompressedStr).Msg("Loaded __csr bitmap")
	c.client.Logger.Info().Any("value", c.siteConfig.VersionId).Msg("Loaded versionId")
	c.client.Logger.Info().Any("broker", c.mqttConfig.Broker).Msg("Configs successfully setup!")
	c.client.Logger.Info().
	Any("total_messages", len(modules.SchedulerJSRequired.LSTable.LSUpsertMessage)).
	Any("total_threads", len(modules.SchedulerJSRequired.LSTable.LSDeleteThenInsertThread)).
	Any("total_contacts", len(modules.SchedulerJSRequired.LSTable.LSVerifyContactRowExists)).
	Msg("Account metadata stats")

	c.client.Logger.Info().Any("length", len(modules.SchedulerJSRequired.LSTable.LSExecuteFirstBlockForSyncTransaction)).Any("executeFirstBlockForSyncTransaction", modules.SchedulerJSRequired.LSTable.LSExecuteFirstBlockForSyncTransaction).Msg("sync information")
	c.client.Logger.Info().Any("list", c.syncCursors).Msg("Cursors")
	return nil
}

func (c *Configs) getSyncParams() ([]byte, error) {
	return json.Marshal(&socket.SyncParams{
		Locale: c.siteConfig.Locale,
	})
}

func (c *Configs) getLastAppliedCursor(database int64) string {
	return c.syncCursors[database]
}

func (c *Configs) updateLastAppliedCursor(database int64, cursor string) {
	c.syncCursors[database] = cursor
}