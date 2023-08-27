package messagix

import (
	"bytes"
	"io"
	"net/http"

	"github.com/0xzer/messagix/types"
)

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