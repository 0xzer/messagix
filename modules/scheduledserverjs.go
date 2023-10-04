package modules

import (
	"log"

	"github.com/0xzer/messagix/table"
)

type SchedulerJSDefine struct {
	MqttWebConfig            MqttWebConfig
	MqttWebDeviceID          MqttWebDeviceID
	WebDevicePerfClassData   WebDevicePerfClassData
	BootloaderConfig         BootLoaderConfig
	CurrentBusinessUser      CurrentBusinessAccount
	SiteData                 SiteData
	SprinkleConfig           SprinkleConfig
	USIDMetadata             USIDMetadata
	WebConnectionClassServerGuess WebConnectionClassServerGuess
	MessengerWebRegion       MessengerWebRegion
	MessengerWebInitData     MessengerWebInitData
	LSD                      LSD
	IntlViewerContext        IntlViewerContext
	IntlCurrentLocale        IntlCurrentLocale
	DTSGInitData             DTSGInitData
	DTSGInitialData          DTSGInitialData
	CurrentUserInitialData   CurrentUserInitialData
	LSPlatformMessengerSyncParams LSPlatformMessengerSyncParams
	ServerNonce ServerNonce
	InitialCookieConsent InitialCookieConsent
	InstagramPasswordEncryption InstagramPasswordEncryption
	XIGSharedData XIGSharedData
	RelayAPIConfigDefaults RelayAPIConfigDefaults
}

type SchedulerJSRequire struct {
	LSTable *table.LSTable
}

var SchedulerJSDefined = SchedulerJSDefine{}
var SchedulerJSRequired = SchedulerJSRequire{
	LSTable: &table.LSTable{},
}

func SSJSHandle(data interface{}) error {
	var err error
	box, ok := data.(map[string]interface{})
	if !ok {
		interfaceData, ok := data.([]interface{})
		if ok {
			err = handleDefine("default_define", interfaceData)
			return err
		}
		log.Fatalf("failed to convert ssjs data to map[string]interface{}")
	}

	for k, v := range box {
		if v == nil {
			continue
		}
		switch k {
			case "__bbox":
				boxMap := v.(map[string]interface{})
				for boxKey, boxData := range boxMap {
					boxDataArr := boxData.([]interface{})
					switch boxKey {
					case "require":
						err = handleRequire("ssjs", boxDataArr)
						continue
					case "define":
						err = handleDefine("ssjs", boxDataArr)
					}
				}
		}
	}
	return err
}