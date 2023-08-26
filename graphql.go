package messagix

import "net/http"

type GraphQLPayload struct {
	Av                   string `json:"av,omitempty"` // not required
	User                 string `json:"__user,omitempty"` // not required
	A                    string `json:"__a,omitempty"` // 1 or 0 wether to include "suggestion_keys" or not in the response - no idea what this is
	Req                  string `json:"__req,omitempty"` // not required
	Hs                   string `json:"__hs,omitempty"` // not required
	Dpr                  string `json:"dpr,omitempty"` // not required
	Ccg                  string `json:"__ccg,omitempty"` // not required
	Rev                  string `json:"__rev,omitempty"` // not required
	S                    string `json:"__s,omitempty"` // not required
	Hsi                  string `json:"__hsi,omitempty"` // not required
	Dyn                  string `json:"__dyn,omitempty"` // not required
	Csr                  string `json:"__csr,omitempty"` // not required
	CometReq             string `json:"__comet_req,omitempty"` // not required but idk what this is
	FbDtsg               string `json:"fb_dtsg,omitempty"`
	Jazoest              string `json:"jazoest,omitempty"` // not required
	Lsd                  string `json:"lsd,omitempty"` // not required
	SpinR                string `json:"__spin_r,omitempty"` // not required
	SpinB                string `json:"__spin_b,omitempty"` // not required
	SpinT                string `json:"__spin_t,omitempty"` // not required
	FbAPICallerClass     string `json:"fb_api_caller_class,omitempty"` // not required
	FbAPIReqFriendlyName string `json:"fb_api_req_friendly_name,omitempty"` // not required
	Variables            string `json:"variables,omitempty"`
	ServerTimestamps     string `json:"server_timestamps,omitempty"` // "true" or "false"
	DocID                string `json:"doc_id,omitempty"`
}

func (c *Client) getGraphQLHeaders() http.Header {
	h := http.Header{}
	h.Add("user-agent", USER_AGENT)
	
	return h
}