package messagix

import (
	"log"
	"strconv"
	"github.com/0xzer/messagix/crypto"
	"github.com/0xzer/messagix/methods"
	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/types"
)

type Configs struct {
	client *Client
	mqttConfig *types.MQTTConfig
	siteConfig *types.SiteConfig
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
	}

	return nil
}