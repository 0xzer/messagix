package modules

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strings"
	"github.com/0xzer/messagix/graphql"
	"github.com/0xzer/messagix/lightspeed"
)

var GraphQLData = &graphql.GraphQLTable{}
var VersionId int64

type EnvJSON struct {
	UseTrustedTypes          bool   `json:"useTrustedTypes,omitempty"`
	IsTrustedTypesReportOnly bool   `json:"isTrustedTypesReportOnly,omitempty"`
	RoutingNamespace         string `json:"routing_namespace,omitempty"`
	Ghlss                    string `json:"ghlss,omitempty"`
	ScheduledCSSJSScheduler  bool   `json:"scheduledCSSJSScheduler,omitempty"`
	UseFbtVirtualModules     bool   `json:"use_fbt_virtual_modules,omitempty"`
	CompatIframeToken        string `json:"compat_iframe_token,omitempty"`
}

type Eqmc struct {
	AjaxURL string `json:"u,omitempty"`
	HasteSessionId string `json:"e,omitempty"`
	S string `json:"s,omitempty"`
	W int    `json:"w,omitempty"`
	FbDtsg string `json:"f,omitempty"`
	L any    `json:"l,omitempty"`
}

type AjaxQueryParams struct {
	A        string `json:"__a"`
	User     string `json:"__user"`
	CometReq string `json:"__comet_req"`
	Jazoest  string `json:"jazoest"`
}

func (e *Eqmc) ParseAjaxURLData() (*AjaxQueryParams, error) {
	u, err := url.Parse(e.AjaxURL)
	if err != nil {
		return nil, err
	}

	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}

	var result AjaxQueryParams

	result.A = params.Get("__a")
	result.User = params.Get("__user")
	result.CometReq = params.Get("__comet_req")
	result.Jazoest = params.Get("jazoest")
	return &result, nil
}

type JSON struct {
	Eqmc Eqmc
	EnvJSON EnvJSON
}

var JSONData = JSON{}

func HandleJSON(data []byte, id string) error {
	var err error
	switch id {
	case "envjson":
		var d EnvJSON
		err = json.Unmarshal(data, &d)
		JSONData.EnvJSON = d
	case "__eqmc":
		var d Eqmc
		err = json.Unmarshal(data, &d)
		JSONData.Eqmc = d
	}
	return err
}

var CsrBitmap = make([]int, 0)
var Bitmap = make([]int, 0)

func interfaceToStructJSON(data interface{}, i interface{}) error {
	b, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &i)
}

func handleDefine(modName string, data []interface{}) error {
	switch modName {
		case "ssjs":
			reflectedMs := reflect.ValueOf(&SchedulerJSDefined).Elem()
			for _, child := range data {
				configData := child.([]interface{})
				config := configData[2]
				configName := configData[0].(string)
				configId := int(configData[3].(float64))
		
				if configId <= 0 {
					continue
				}

				Bitmap = append(Bitmap, configId)
				field := reflectedMs.FieldByName(configName)
				if !field.IsValid() {
					//fmt.Printf("Invalid field name: %s\n", configName)
					continue
				}
				if !field.CanSet() {
					//fmt.Printf("Unsettable field: %s\n", configName)
					continue
				}
		
				err := interfaceToStructJSON(config, field.Addr().Interface())
				if err != nil {
					return err
				}
			}
		}
		return nil
}

func handleRequire(modName string, data []interface{}) error {
	switch modName {
		case "ssjs":
			//reflectedMs := reflect.ValueOf(&SchedulerJSRequired).Elem()
			for _, requireData := range data {
				d := requireData.([]interface{})
				requireType := d[0].(string)
				switch requireType {
					case "CometPlatformRootClient":
						moduleData := d[3].([]interface{})
						for _, v := range moduleData {
							requestsMap, ok := v.([]interface{})
							if !ok {
								continue
							}
							for _, req := range requestsMap {
								var reqData *graphql.GraphQLPreloader
								err := interfaceToStructJSON(req, &reqData)
								if err != nil {
									continue
								}
								if len(reqData.Variables.RequestPayload) > 0 {
									var syncData *graphql.LSPlatformGraphQLLightspeedVariables
									err = json.Unmarshal([]byte(reqData.Variables.RequestPayload), &syncData)
									if err != nil {
										continue
									}
									VersionId = syncData.Version
								}
							}

						}
					case "RelayPrefetchedStreamCache":
						moduleData := d[3].([]interface{})
						//method := d[1].(string)
						//dependencies := d[2].(string)
						parserFunc := parseGraphMethodName(moduleData[0].(string))
						graphQLData := moduleData[1].(map[string]interface{})
						boxData, ok := graphQLData["__bbox"].(map[string]interface{})
						if !ok {
							return fmt.Errorf("could not find __bbox in graphQLData map for parser func: %s", parserFunc)
						}

						result, ok := boxData["result"]
						if !ok {
							return fmt.Errorf("could not find result in __bbox for parser func: %s", parserFunc)
						}

						if parserFunc == "LSPlatformGraphQLLightspeedRequestQuery" {
							handleLightSpeedQLRequest(result)
						} else {
							handleGraphQLData(parserFunc, result)
						}
				}
			}
		}
	return nil
}

func handleLightSpeedQLRequest(data interface{}) {
	var lsData *graphql.LSPlatformGraphQLLightspeedRequestQuery
	err := interfaceToStructJSON(&data, &lsData)
	if err != nil {
		log.Fatalf("failed to parse LightSpeedQLRequest data from html: %e", err)
	}
	
	lsPayload := lsData.Data.Viewer.LightspeedWebRequest.Payload
	dependencies := lightspeed.DependenciesToMap(lsData.Data.Viewer.LightspeedWebRequest.Dependencies)
	decoder := lightspeed.NewLightSpeedDecoder(dependencies, SchedulerJSRequired.LSTable)
	
	var payload lightspeed.LightSpeedData
	err = json.Unmarshal([]byte(lsPayload), &payload)
	if err != nil {
		log.Fatalf("failed to marshal lsPayload into LightSpeedData: %e", err)
	}
	
	decoder.Decode(payload.Steps)
}

func handleGraphQLData(name string, data interface{}) {
	reflectedMs := reflect.ValueOf(GraphQLData).Elem()
	dataField := reflectedMs.FieldByName(name)
	if !dataField.IsValid() {
		log.Println("Not handling GraphQLData for operation:", name)
		return
	}
	
	definition := dataField.Type().Elem()
	newDefinition := reflect.New(definition).Interface()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(fmt.Sprintf("failed to marshal GraphQL operation data %s", name))
		return
	}

	err = json.Unmarshal(jsonBytes, newDefinition)
	if err != nil {
		log.Println(fmt.Sprintf("failed to unmarshal GraphQL operation data %s", name))
		return
	}

	newSlice := reflect.Append(dataField, reflect.Indirect(reflect.ValueOf(newDefinition)))
	dataField.Set(newSlice)
}

func parseGraphMethodName(name string) string {
	var s string
	s = strings.Replace(name, "adp_", "", -1)
	s = strings.Split(s, "RelayPreloader_")[0]
	return s
}