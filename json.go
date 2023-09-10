package messagix

import (
	"encoding/json"
)

type Connect struct {
	AccountId   string `json:"u"` // account id
	SessionId   int64  `json:"s"` // randomly generated sessionid
	ClientCapabilities int    `json:"cp"` // mqttconfig clientCapabilities (3)
	Capabilities         int    `json:"ecp"` // mqttconfig capabilities (10)
	ChatOn      bool   `json:"chat_on"` // mqttconfig chatVisibility (true) - not 100% sure
	Fg          bool   `json:"fg"` // idk what this is
	Cid   string `json:"d"` // cid from html content
    ConnectionType  string `json:"ct"` // connection type? websocket
	MqttSid     string `json:"mqtt_sid"` // ""
	AppId       int64  `json:"aid"` // mqttconfig appID (219994525426954)
	SubscribedTopics	[]any  `json:"st"` // mqttconfig subscribedTopics ([])
	Pm          []any  `json:"pm"` // only seen empty array
	Dc          string `json:"dc"` // only seem empty string
	NoAutoFg    bool   `json:"no_auto_fg"` // only seen true
	Gas         any    `json:"gas"` // only seen null
	Pack        []any  `json:"pack"` // only seen empty arr
	HostNameOverride string `json:"php_override"` // mqttconfig hostNameOverride ("") - not 100% sure
	P           any    `json:"p"` // only seen null
	UserAgent   string `json:"a"` // user agent
	Aids        any    `json:"aids"` // only seen null
}

func (s *Socket) newConnectJSON() (string, error) {
	payload := &Connect{
		AccountId: s.client.cookies.AccountId,
		SessionId: s.client.configs.mqttConfig.SessionId,
		ClientCapabilities: s.client.configs.mqttConfig.ClientCapabilities,
		Capabilities: s.client.configs.mqttConfig.Capabilities,
		ChatOn: s.client.configs.mqttConfig.ChatOn,
		Fg: false,
		ConnectionType: s.client.configs.mqttConfig.ConnectionType,
		MqttSid: "",
		AppId: s.client.configs.mqttConfig.AppId,
		SubscribedTopics: s.client.configs.mqttConfig.SubscribedTopics,
		Pm: make([]any, 0),
		Dc: "",
		NoAutoFg: true,
		Gas: nil,
		Pack: make([]any, 0),
		HostNameOverride: s.client.configs.mqttConfig.HostNameOverride,
		P: nil,
		UserAgent: USER_AGENT,
		Aids: nil,
		Cid: s.client.configs.mqttConfig.ClientId,
	}
	
	jsonData, err := json.Marshal(payload)
	return string(jsonData), err
}