package messagix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/0xzer/messagix/graphql"
	"github.com/0xzer/messagix/lightspeed"
	"github.com/0xzer/messagix/table"
	"github.com/0xzer/messagix/types"
	"github.com/google/go-querystring/query"
)

func (c *Client) makeGraphQLRequest(name string, variables interface{}) (*http.Response, []byte, error) {
	graphQLDoc, ok := graphql.GraphQLDocs[name]
	if !ok {
		return nil, nil, fmt.Errorf("could not find graphql doc by the name of: %s", name)
	}

	vBytes, err := json.Marshal(variables)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal graphql variables to json string: %e", err)
	}


	payload := c.NewHttpQuery()
	payload.FbAPICallerClass = graphQLDoc.CallerClass
	payload.FbAPIReqFriendlyName = graphQLDoc.FriendlyName
	payload.Variables = string(vBytes)
	payload.ServerTimestamps = "true"
	payload.DocID = graphQLDoc.DocId

	c.graphQLRequests++

	form, err := query.Values(payload)
	if err != nil {
		return nil, nil, err
	}

	payloadBytes := []byte(form.Encode())

	headers := c.buildHeaders()
	headers.Add("x-fb-friendly-name", graphQLDoc.FriendlyName)
	headers.Add("sec-fetch-dest", "empty")
	headers.Add("sec-fetch-mode", "cors")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("origin", "https://www.facebook.com")
	headers.Add("referer", "https://www.facebook.com/messages/")

	//g.client.Logger.Debug().Any("headers", headers).Any("payload", string(payloadBytes)).Msg("Sending graphQL request")
	return c.MakeRequest("https://www.facebook.com/api/graphql/", "POST", headers, payloadBytes, types.FORM)
}

func (c *Client) makeLSRequest(variables *graphql.LSPlatformGraphQLLightspeedVariables, reqType int) (*table.LSTable, error) {
	strPayload, err := json.Marshal(&variables)
	if err != nil {
		return nil, err
	}
	
	lsVariables := &graphql.LSPlatformGraphQLLightspeedRequestPayload{
		DeviceID: c.configs.mqttConfig.Cid,
		IncludeChatVisibility: false,
		RequestID: c.lsRequests,
		RequestPayload: string(strPayload),
		RequestType: reqType,
	}
	c.lsRequests++

	_, respBody, err := c.makeGraphQLRequest("LSGraphQLRequest", &lsVariables)
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