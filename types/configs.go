package types

import (
	"net/url"
	"strconv"

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
	Cid string
}

func (m *MQTTConfig) BuildBrokerUrl() string {
	query := &url.Values{}
	query.Add("cid", m.Cid)
	query.Add("sid", strconv.Itoa(int(m.SessionId)))
	
	return m.Broker + "&" + query.Encode()
}

type GraphQLPayload struct {
	Av                   string `json:"av,omitempty"` // user id
	User                 string `json:"__user,omitempty"` // user id
	A                    string `json:"__a,omitempty"` // always 1?
	ReqId                string `json:"__req,omitempty"`
	HasteSession         string `json:"__hs,omitempty"`
	Pr                   string `json:"dpr,omitempty"`
	ConnectionClass      string `json:"__ccg,omitempty"`
	Revision       		 string `json:"__rev,omitempty"`
	WebSessionId         string `json:"__s,omitempty"`
	HasteSessionId       string `json:"__hsi,omitempty"`
	CompressedBitmap     string `json:"__dyn,omitempty"`
	CompressedCsrBitmap  string `json:"__csr,omitempty"`
	CometReq             string `json:"__comet_req,omitempty"`
	FbDtsg               string `json:"fb_dtsg,omitempty"`
	Jazoest              string `json:"jazoest,omitempty"`
	LsdToken             string `json:"lsd,omitempty"`
	SpinR                string `json:"__spin_r,omitempty"`
	SpinB                string `json:"__spin_b,omitempty"`
	SpinT                string `json:"__spin_t,omitempty"`
	FbAPICallerClass     string `json:"fb_api_caller_class,omitempty"`
	FbAPIReqFriendlyName string `json:"fb_api_req_friendly_name,omitempty"`
	Variables            interface{} `json:"variables,omitempty"`
	ServerTimestamps     bool   `json:"server_timestamps,omitempty"` // "true" or "false"
	DocID                string `json:"doc_id,omitempty"`
}

type SiteConfig struct {
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
}