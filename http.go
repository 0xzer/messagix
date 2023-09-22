package messagix

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/0xzer/messagix/types"
)

type HttpQuery struct {
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

func (c *Client) NewHttpQuery() *HttpQuery {
	siteConfig := c.configs.siteConfig
	return &HttpQuery{
		Av: siteConfig.AccountId,
		User: siteConfig.AccountId,
		A: "1",
		Req: strconv.Itoa(c.graphQLRequests),
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
	}
}

func (c *Client) MakeRequest(url string, method string, headers http.Header, payload []byte, contentType types.ContentType) (*http.Response, []byte, error) {
	newRequest, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, err
	}

	if contentType != types.NONE {
		headers.Add("content-type", string(contentType))
	}

	newRequest.Header = headers

	response, err := c.http.Do(newRequest)
	if err !=  nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return response, responseBody, nil
}

// 129477
func (c *Client) buildHeaders() http.Header {
	w, _ := c.cookies.GetViewports()

	headers := http.Header{}
	headers.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	headers.Add("accept-language", "en-US,en;q=0.9")
	headers.Add("cache-control", "0")
	headers.Add("dpr", "1")
	headers.Add("sec-ch-prefers-color-scheme", "light")
	headers.Add("sec-ch-ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	headers.Add("sec-ch-ua-full-version-list", "\"Google Chrome\";v=\"113.0.5672.92\", \"Chromium\";v=\"113.0.5672.92\", \"Not-A.Brand\";v=\"24.0.0.0\"")
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-ch-ua-model", "")
	headers.Add("sec-ch-ua-platform", "Linux")
	headers.Add("sec-ch-ua-platform-version", "6.4.10")
	headers.Add("user-agent", USER_AGENT)
	headers.Add("viewport-width", w)

	if c.configs.siteConfig != nil {
		headers.Add("x-asbd-id", c.configs.siteConfig.X_ASDB_ID)
		headers.Add("x-fb-lsd", c.configs.siteConfig.LsdToken)
	}

	if c.cookies != nil {
		headers.Add("cookie", c.cookies.ToString())
	}

	return headers
}