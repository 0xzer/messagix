package messagix

import (
	"encoding/json"
	"strconv"
)

type Topic string

const (
	UNKNOWN_TOPIC Topic = "unknown"
	LS_APP_SETTINGS Topic = "/ls_app_settings"
	LS_FOREGROUND_STATE Topic = "/ls_foreground_state"
	LS_REQ Topic = "/ls_req"
	LS_RESP Topic = "/ls_resp"
)

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