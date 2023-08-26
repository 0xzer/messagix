package modules

type MqttWebConfig struct {
	AppID              int64  `json:"appID,omitempty"`
	Capabilities       int    `json:"capabilities,omitempty"`
	ChatVisibility     bool   `json:"chatVisibility,omitempty"`
	ClientCapabilities int    `json:"clientCapabilities,omitempty"`
	Endpoint           string `json:"endpoint,omitempty"`
	Fbid               string `json:"fbid,omitempty"`
	HostNameOverride   string `json:"hostNameOverride,omitempty"`
	PollingEndpoint    string `json:"pollingEndpoint,omitempty"`
	SubscribedTopics   []any  `json:"subscribedTopics,omitempty"`
}

type MqttWebDeviceID struct {
	ClientID string `json:"clientID,omitempty"`
}