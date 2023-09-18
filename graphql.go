package messagix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/0xzer/messagix/graphql"
	"github.com/0xzer/messagix/lightspeed"
	"github.com/0xzer/messagix/table"
	"github.com/0xzer/messagix/types"
	"github.com/google/go-querystring/query"
)

type GraphQLPayload struct {
	Av                   string `url:"av,omitempty"` // not required
	User                 string `url:"__user,omitempty"` // not required
	A                    string `url:"__a,omitempty"` // 1 or 0 wether to include "suggestion_keys" or not in the response - no idea what this is
	Req                  string `url:"__req,omitempty"` // not required
	Hs                   string `url:"__hs,omitempty"` // not required
	Dpr                  string `url:"dpr,omitempty"` // not required
	Ccg                  string `url:"__ccg,omitempty"` // not required
	Rev                  string `url:"__rev,omitempty"` // not required
	S                    string `url:"__s,omitempty"` // not required
	Hsi                  string `url:"__hsi,omitempty"` // not required
	Dyn                  string `url:"__dyn,omitempty"` // not required
	Csr                  string `url:"__csr,omitempty"` // not required
	CometReq             string `url:"__comet_req,omitempty"` // not required but idk what this is
	FbDtsg               string `url:"fb_dtsg,omitempty"`
	Jazoest              string `url:"jazoest,omitempty"` // not required
	Lsd                  string `url:"lsd,omitempty"` // not required
	SpinR                string `url:"__spin_r,omitempty"` // not required
	SpinB                string `url:"__spin_b,omitempty"` // not required
	SpinT                string `url:"__spin_t,omitempty"` // not required
	FbAPICallerClass     string `url:"fb_api_caller_class,omitempty"` // not required
	FbAPIReqFriendlyName string `url:"fb_api_req_friendly_name,omitempty"` // not required
	Variables            string `url:"variables,omitempty"`
	ServerTimestamps     string `url:"server_timestamps,omitempty"` // "true" or "false"
	DocID                string `url:"doc_id,omitempty"`
}

type GraphQL struct {
	client *Client
	lsRequests int
	graphQLRequests int
}

func (g *GraphQL) makeGraphQLRequest(name string, variables interface{}) (*http.Response, []byte, error) {
	graphQLDoc, ok := graphql.GraphQLDocs[name]
	if !ok {
		return nil, nil, fmt.Errorf("could not find graphql doc by the name of: %s", name)
	}

	vBytes, err := json.Marshal(variables)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal graphql variables to json string: %e", err)
	}


	siteConfig := g.client.configs.siteConfig
	payload := &GraphQLPayload{
		Av: siteConfig.AccountId,
		User: siteConfig.AccountId,
		A: "1",
		Req: strconv.Itoa(g.graphQLRequests),
		Hs: siteConfig.HasteSession,
		Dpr: "1",
		Ccg: siteConfig.ConnectionClass,
		Rev: siteConfig.SpinR,
		S: siteConfig.WebSessionId,
		Hsi: siteConfig.HasteSessionId,
		Dyn: siteConfig.Bitmap.CompressedStr,
		Csr: siteConfig.CSRBitmap.CompressedStr,
		CometReq: siteConfig.CometReq,
		FbDtsg: siteConfig.FbDtsg,
		Jazoest: siteConfig.Jazoest,
		Lsd: siteConfig.LsdToken,
		SpinR: siteConfig.SpinR,
		SpinB: siteConfig.SpinB,
		SpinT: siteConfig.SpinT,
		FbAPICallerClass: graphQLDoc.CallerClass,
		FbAPIReqFriendlyName: graphQLDoc.FriendlyName,
		Variables: string(vBytes),
		ServerTimestamps: "true",
		DocID: graphQLDoc.DocId,
	}
	g.graphQLRequests++

	form, err := query.Values(payload)
	if err != nil {
		return nil, nil, err
	}

	payloadBytes := []byte(form.Encode())

	headers := g.client.buildHeaders()
	headers.Add("x-fb-friendly-name", graphQLDoc.FriendlyName)
	headers.Add("sec-fetch-dest", "empty")
	headers.Add("sec-fetch-mode", "cors")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("origin", "https://www.facebook.com")
	headers.Add("referer", "https://www.facebook.com/messages/")

	//g.client.Logger.Debug().Any("headers", headers).Any("payload", string(payloadBytes)).Msg("Sending graphQL request")
	return g.client.MakeRequest("https://www.facebook.com/api/graphql/", "POST", headers, payloadBytes, types.FORM)
}

func (g *GraphQL) makeLSRequest(variables *graphql.LSPlatformGraphQLLightspeedVariables, reqType int) (*table.LSTable, error) {
	strPayload, err := json.Marshal(&variables)
	if err != nil {
		return nil, err
	}
	
	lsVariables := &graphql.LSPlatformGraphQLLightspeedRequestPayload{
		DeviceID: g.client.configs.mqttConfig.Cid,
		IncludeChatVisibility: false,
		RequestID: g.lsRequests,
		RequestPayload: string(strPayload),
		RequestType: reqType,
	}
	g.lsRequests++

	_, respBody, err := g.makeGraphQLRequest("LSGraphQLRequest", &lsVariables)
	if err != nil {
		return nil, err
	}

	os.WriteFile("lsGraphQLResponse.json", respBody, os.ModePerm)
	var graphQLData *graphql.LSPlatformGraphQLLightspeedRequestQuery
	err = json.Unmarshal(respBody, &graphQLData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal LSRequest response bytes into LSPlatformGraphQLLightspeedRequestQuery struct: %e", err)
	}

	lightSpeedRes := graphQLData.Data.Viewer.LightspeedWebRequest

	var lsData *lightspeed.LightSpeedData
	err = json.Unmarshal([]byte(lightSpeedRes.Payload), &lsData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal LSRequest lightspeed payload into lightspeed.LightSpeedData: %e", err)
	}

	dependencies := lightspeed.DependenciesToMap(lightSpeedRes.Dependencies)

	lsTable := &table.LSTable{}
	lsDecoder := lightspeed.NewLightSpeedDecoder(dependencies, lsTable)
	lsDecoder.Decode(lsData.Steps)

	return lsTable, nil
}