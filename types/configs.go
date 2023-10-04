package types

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/0xzer/messagix/crypto"
)

type MqttWebDeviceID struct {
	ClientID string `json:"clientID,omitempty"`
}

type MQTTWebConfig struct {
	Fbid               string `json:"fbid,omitempty"`
	AppID              int64  `json:"appID,omitempty"`
	Endpoint           string `json:"endpoint,omitempty"`
	PollingEndpoint    string `json:"pollingEndpoint,omitempty"`
	SubscribedTopics   []any  `json:"subscribedTopics,omitempty"`
	Capabilities       int    `json:"capabilities,omitempty"`
	ClientCapabilities int    `json:"clientCapabilities,omitempty"`
	ChatVisibility     bool   `json:"chatVisibility,omitempty"`
	HostNameOverride   string `json:"hostNameOverride,omitempty"`
}

type MQTTConfig struct {
	ProtocolName string
	ProtocolLevel uint8
	ClientId string
	Broker string
	KeepAliveTimeout uint16
	SessionId int64
	AppId int64
	ClientCapabilities int
	Capabilities int
	ChatOn bool
	SubscribedTopics []any
	ConnectionType string
	HostNameOverride string
	Cid string // device id
}

func (m *MQTTConfig) BuildBrokerUrl() string {
	query := &url.Values{}
	query.Add("cid", m.Cid)
	query.Add("sid", strconv.Itoa(int(m.SessionId)))
	
	encodedQuery := query.Encode()
	if !strings.HasSuffix(m.Broker, "=") {
		return m.Broker + encodedQuery
	} else {
		return m.Broker + "&" + encodedQuery
	}
}

type SiteConfig struct {
	ServerRevision string
	AccountId string
	AccountIdInt int64
	Bitmap *crypto.Bitmap
	CSRBitmap *crypto.Bitmap
	HasteSessionId string
	WebSessionId string
	CometReq string
	LsdToken string
	SpinT string
	SpinB string
	SpinR string
	FbDtsg string
	Jazoest string
	Pr string
	HasteSession string
	ConnectionClass string
	VersionId int64
	X_ASDB_ID string
	Locale string
	LgnRnd string
	Trynum string
}