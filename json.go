package messagix

import (
	"encoding/json"
	"strconv"

	"github.com/0xzer/messagix/modules"
)

type Connect struct {
	AccountId   string `json:"u"` // account id
	SessionId   int64  `json:"s"` // randomly generated sessionid
	ClientCapabilities int    `json:"cp"` // mqttconfig clientCapabilities (3)
	Capabilities         int    `json:"ecp"` // mqttconfig capabilities (10)
	ChatOn      bool   `json:"chat_on"` // mqttconfig chatVisibility (true) - not 100% sure
	Fg          bool   `json:"fg"` // idk what this is
	Cid   string `json:"d"` // cid from html content
    ConnectionType  string `json:"ct"` // connection type? facebook=websocket , insta=cookie_auth
	MqttSid     string `json:"mqtt_sid"` // ""
	AppId       int64  `json:"aid"` // mqttconfig appID (219994525426954)
	SubscribedTopics	[]any  `json:"st"` // mqttconfig subscribedTopics ([])
	Pm          []any  `json:"pm"` // only seen empty array
	Dc          string `json:"dc"` // only seem empty string
	NoAutoFg    bool   `json:"no_auto_fg"` // only seen true
	Gas         any    `json:"gas"` // only seen null
	Pack        []any  `json:"pack"` // only seen empty arr
	HostNameOverride string `json:"php_override"` // mqttconfig hostNameOverride
	P           any    `json:"p"` // only seen null
	UserAgent   string `json:"a"` // user agent
	Aids        any    `json:"aids"` // only seen null
}

func (s *Socket) newConnectJSON() (string, error) {
	payload := &Connect{
		AccountId: modules.SchedulerJSDefined.CurrentUserInitialData.AccountID,
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


type AppSettingsPublish struct {
	LsFdid string `json:"ls_fdid"`
	SchemaVersion   string `json:"ls_sv"`
}

func (s *Socket) newAppSettingsPublishJSON(versionId int64) (string, error) {
	payload := &AppSettingsPublish{
		LsFdid: "",
		SchemaVersion: strconv.Itoa(int(versionId)),
	}

	jsonData, err := json.Marshal(payload)
	return string(jsonData), err
}