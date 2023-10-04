package messagix

import (
	"fmt"
	"reflect"
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
}

func (c *Configs) SetupConfigs() error {
	schedulerJS := modules.SchedulerJSDefined
	if schedulerJS.SiteData == (modules.SiteData{}) {
		return fmt.Errorf("SetupConfigs was somehow called before modules were initalized")
	}

	authenticated := c.client.IsAuthenticated()
	bitmap, csrBitmap := c.LoadBitmaps()
	c.siteConfig = &types.SiteConfig{
		AccountId: schedulerJS.CurrentUserInitialData.AccountID,
		Bitmap: bitmap,
		CSRBitmap: csrBitmap,
		HasteSessionId: schedulerJS.SiteData.Hsi,
		WebSessionId: methods.GenerateWebsessionID(authenticated),
		LsdToken: schedulerJS.LSD.Token,
		SpinT: strconv.Itoa(schedulerJS.SiteData.SpinT),
		SpinB: schedulerJS.SiteData.SpinB,
		SpinR: strconv.Itoa(schedulerJS.SiteData.SpinR),
		FbDtsg: schedulerJS.DTSGInitialData.Token,
		Pr: strconv.FormatFloat(schedulerJS.SiteData.Pr, 'f', -1, 64),
		HasteSession: schedulerJS.SiteData.HasteSession,
		ConnectionClass: schedulerJS.WebConnectionClassServerGuess.ConnectionClass,
		ServerRevision: strconv.Itoa(int(schedulerJS.SiteData.ServerRevision)),
		VersionId: modules.VersionId,
		Locale: schedulerJS.IntlCurrentLocale.Code,
		X_ASDB_ID: "129477", // hold off on this and check if it ever changes, if so we gotta load the js file and extract it from there
	}

	if authenticated {
		eqmcQuery, err := modules.JSONData.Eqmc.ParseAjaxURLData()
		if err != nil {
			return fmt.Errorf("failed to parse AjaxURLData from eqmc json struct: %e", err)
		}
		accountIdInt, _ := strconv.Atoi(schedulerJS.CurrentUserInitialData.AccountID)

		c.siteConfig.AccountIdInt = int64(accountIdInt)
		c.siteConfig.Jazoest = eqmcQuery.Jazoest
		c.siteConfig.CometReq = eqmcQuery.CometReq
		if c.client.platform == types.Instagram {
			schedulerJS.MqttWebConfig.Endpoint = "wss://edge-chat.instagram.com/chat?"
			schedulerJS.MqttWebConfig.AppID = schedulerJS.MessengerWebInitData.AppID
			userID := c.client.cookies.GetValue("ds_user_id")
			if userID != "" {
				accountIdInt, _ := strconv.Atoi(userID)
				c.siteConfig.AccountIdInt = int64(accountIdInt)
			}
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
			ConnectionType: c.client.socket.getConnectionType(),
			HostNameOverride: "",
			Cid: schedulerJS.MqttWebDeviceID.ClientID,
		}
		c.client.syncManager.syncParams = &modules.SchedulerJSDefined.LSPlatformMessengerSyncParams
		if c.needSync {
			err := c.client.syncManager.UpdateDatabaseSyncParams(
				[]*socket.QueryMetadata{
					{DatabaseId: 1, SendSyncParams: true, LastAppliedCursor: nil, SyncChannel: socket.MailBox},
					{DatabaseId: 2, SendSyncParams: true, LastAppliedCursor: nil, SyncChannel: socket.Contact},
					{DatabaseId: 95, SendSyncParams: true, LastAppliedCursor: nil, SyncChannel: socket.Contact},
				},
			)
			if err != nil {
				return fmt.Errorf("failed to update sync params for databases: 1, 2, 95")
			}

			lsData, err := c.client.syncManager.SyncDataGraphQL([]int64{1,2,95})
			if err != nil {
				return fmt.Errorf("failed to sync data via graphql for databases: 1, 2, 95")
			}

			//c.client.Logger.Info().Any("lsData", lsData).Msg("Synced data through graphql query")
			modules.SchedulerJSRequired.LSTable = lsData
		} else {
			err := c.client.syncManager.SyncTransactions(modules.SchedulerJSRequired.LSTable.LSExecuteFirstBlockForSyncTransaction)
			if err != nil {
				return fmt.Errorf("failed to sync transactions from js module data with syncManager: %e", err)
			}
		}
		c.client.Logger.Info().Any("value", c.siteConfig.Bitmap.CompressedStr).Msg("Loaded __dyn bitmap")
		c.client.Logger.Info().Any("value", c.siteConfig.CSRBitmap.CompressedStr).Msg("Loaded __csr bitmap")
		c.client.Logger.Info().Any("versionId", c.siteConfig.VersionId).Any("appId", schedulerJS.CurrentUserInitialData.AppID).Msg("Loaded versionId & appId")
		c.client.Logger.Info().Any("broker", c.mqttConfig.Broker).Msg("Configs successfully setup!")
		c.client.Logger.Info().
		Any("total_messages", len(modules.SchedulerJSRequired.LSTable.LSUpsertMessage)).
		Any("total_threads", len(modules.SchedulerJSRequired.LSTable.LSDeleteThenInsertThread)).
		Any("total_contacts", len(modules.SchedulerJSRequired.LSTable.LSVerifyContactRowExists)).
		Msg("Account metadata stats")
	} else {
		c.client.Logger.Info().Any("value", c.siteConfig.Bitmap.CompressedStr).Msg("Loaded __dyn bitmap")
		c.client.Logger.Info().Any("value", c.siteConfig.CSRBitmap.CompressedStr).Msg("Loaded __csr bitmap")
		c.client.Logger.Info().Any("platform", c.client.platform).Msg("Configs loaded, but not yet logged in.")
	}
	
	return nil
}

func (c *Configs) ParseFormInputs(inputs []InputTag, reflectedMs reflect.Value) {
	for _, input := range inputs {
		attr := input.Attributes
		key := attr["name"]
		val := attr["value"]

		field := reflectedMs.FieldByNameFunc(func(s string) bool {
			field, _ := reflectedMs.Type().FieldByName(s)
			return field.Tag.Get("name") == key
		})

		if field.IsValid() && field.CanSet() {
			field.SetString(val)
		}
	}
}

// (bitmap, csrBitmap)
func (c *Configs) LoadBitmaps() (*crypto.Bitmap, *crypto.Bitmap) {
	bitmap := crypto.NewBitmap().Update(modules.Bitmap).ToCompressedString()
	csrBitmap := crypto.NewBitmap().Update(modules.CsrBitmap).ToCompressedString()

	return bitmap, csrBitmap
}