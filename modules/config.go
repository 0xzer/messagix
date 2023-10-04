package modules

import (
	"encoding/json"

	"github.com/0xzer/messagix/types"
)

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

type LSPlatformMessengerSyncParams struct {
	Mailbox string `json:"mailbox,omitempty"`
	Contact string `json:"contact,omitempty"`
	E2Ee    string `json:"e2ee,omitempty"`
}

type InitialCookieConsent struct {
	DeferCookies           bool  `json:"deferCookies,omitempty"`
	InitialConsent         []int `json:"initialConsent,omitempty"`
	NoCookies              bool  `json:"noCookies,omitempty"`
	ShouldShowCookieBanner bool  `json:"shouldShowCookieBanner,omitempty"`
}

type InstagramPasswordEncryption struct {
	KeyID     string `json:"key_id,omitempty"`
	PublicKey string `json:"public_key,omitempty"`
	Version   string `json:"version,omitempty"`
}

type XIGSharedData struct {
	ConfigData *types.XIGConfigData
	Raw    string `json:"raw,omitempty"`
	Native struct {
		Config struct {
			CsrfToken string `json:"csrf_token,omitempty"`
			ViewerID  string    `json:"viewerId,omitempty"`
			Viewer    struct {
				IsBasicAdsOptedIn bool `json:"is_basic_ads_opted_in,omitempty"`
				BasicAdsTier      int  `json:"basic_ads_tier,omitempty"`
			} `json:"viewer,omitempty"`
		} `json:"config,omitempty"`
		SendDeviceIDHeader bool `json:"send_device_id_header,omitempty"`
		ServerChecks       struct {
			Hfe bool `json:"hfe,omitempty"`
		} `json:"server_checks,omitempty"`
		WwwRoutingConfig struct {
			FrontendOnlyRoutes []struct {
				Path        string `json:"path,omitempty"`
				Destination string `json:"destination,omitempty"`
			} `json:"frontend_only_routes,omitempty"`
		} `json:"www_routing_config,omitempty"`
		DeviceID               string `json:"device_id,omitempty"`
		SignalCollectionConfig struct {
			Sid int `json:"sid,omitempty"`
		} `json:"signal_collection_config,omitempty"`
		PrivacyFlowTrigger        any `json:"privacy_flow_trigger,omitempty"`
		PlatformInstallBadgeLinks struct {
			Ios         string `json:"ios,omitempty"`
			Android     string `json:"android,omitempty"`
			WindowsNt10 string `json:"windows_nt_10,omitempty"`
		} `json:"platform_install_badge_links,omitempty"`
		CountryCode    string `json:"country_code,omitempty"`
		ProbablyHasApp bool   `json:"probably_has_app,omitempty"`
	} `json:"native,omitempty"`
}

func (xig *XIGSharedData) ParseRaw() error {
	return json.Unmarshal([]byte(xig.Raw), &xig.ConfigData)
}

type RelayAPIConfigDefaults struct {
	AccessToken   string `json:"accessToken,omitempty"`
	ActorID       string `json:"actorID,omitempty"`
	CustomHeaders struct {
		XIGAppID string `json:"X-IG-App-ID,omitempty"`
		XIGD     string `json:"X-IG-D,omitempty"`
	} `json:"customHeaders,omitempty"`
	EnableNetworkLogger       bool   `json:"enableNetworkLogger,omitempty"`
	FetchTimeout              int    `json:"fetchTimeout,omitempty"`
	GraphBatchURI             string `json:"graphBatchURI,omitempty"`
	GraphURI                  string `json:"graphURI,omitempty"`
	RetryDelays               []int  `json:"retryDelays,omitempty"`
	UseXController            bool   `json:"useXController,omitempty"`
	XhrEncoding               interface{}    `json:"xhrEncoding,omitempty"`
	SubscriptionTopicURI      interface{}    `json:"subscriptionTopicURI,omitempty"`
	WithCredentials           bool   `json:"withCredentials,omitempty"`
	IsProductionEndpoint      bool   `json:"isProductionEndpoint,omitempty"`
	WorkRequestTaggingProduct interface{}    `json:"workRequestTaggingProduct,omitempty"`
	EncryptionKeyParams       interface{}    `json:"encryptionKeyParams,omitempty"`
}