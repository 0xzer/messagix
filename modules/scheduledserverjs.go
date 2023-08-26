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
}

type SchedulerJSRequire struct {
	LSTable *table.LSTable
}

var SchedulerJSDefined = SchedulerJSDefine{}
var SchedulerJSRequired = SchedulerJSRequire{
	LSTable: &table.LSTable{},
}

func SSJSHandle(data interface{}) error {
	box, ok := data.(map[string]interface{})
	if !ok {
		_, ok := data.([]interface{})
		if ok {
			return nil
		}
		log.Fatalf("failed to convert ssjs data to map[string]interface{}")
	}

	var err error
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