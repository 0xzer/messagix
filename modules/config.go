package modules

type SprinkleConfig struct {
	ParamName       string `json:"param_name,omitempty"`
	ShouldRandomize bool   `json:"should_randomize,omitempty"`
	Version         int    `json:"version,omitempty"`
}

type WebConnectionClassServerGuess struct {
	ConnectionClass string `json:"connectionClass,omitempty"`
}

type WebDevicePerfClassData struct {
	DeviceLevel string `json:"deviceLevel,omitempty"`
	YearClass   any    `json:"yearClass,omitempty"`
}

type USIDMetadata struct {
	BrowserID    string `json:"browser_id,omitempty"`
	PageID       string `json:"page_id,omitempty"`
	TabID        string `json:"tab_id,omitempty"`
	TransitionID int    `json:"transition_id,omitempty"`
	Version      int    `json:"version,omitempty"`
}

type MessengerWebRegion struct {
	Region string `json:"regionNullable,omitempty"`
}