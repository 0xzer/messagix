package messagix

import "encoding/json"

type Topic string

const (
	UNKNOWN_TOPIC Topic = "unknown"
	APP_SETTINGS Topic = "/ls_app_settings"
)

type AppSettingsPublish struct {
	LsFdid string `json:"ls_fdid"`
	SchemaVersion   string `json:"ls_sv"` /* unsure if this is language-code based (but it is 6512074225573706 all the time)*/
}

func (s *Socket) newAppSettingsPublishJSON() (string, error) {
	payload := &AppSettingsPublish{
		LsFdid: "",
		SchemaVersion: "6512074225573706",
	}

	jsonData, err := json.Marshal(payload)
	return string(jsonData), err
}