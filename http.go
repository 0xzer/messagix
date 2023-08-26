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
	headers.Add("content-type", string(contentType))
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