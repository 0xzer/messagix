package messagix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/textproto"
	"reflect"

	"github.com/0xzer/messagix/methods"
	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/types"
	"github.com/google/go-querystring/query"
)

type MediaType string
const (
	IMAGE_JPEG MediaType = "image/jpeg"
	VIDEO_MP4 MediaType = "video/mp4"
)

type MercuryUploadMedia struct {
	Filename string
	MediaType MediaType
	MediaData []byte
}

var MERCURY_UPLOAD_URL = "https://www.facebook.com/ajax/mercury/upload.php?"
func (c *Client) SendMercuryUploadRequest(medias []*MercuryUploadMedia) ([]*types.MercuryUploadResponse, error) {
	responses := make([]*types.MercuryUploadResponse, 0)
	for _, media := range medias {
		urlQueries := c.NewHttpQuery()
		queryValues, err := query.Values(urlQueries)
		if err != nil {
			return nil, fmt.Errorf("failed to convert HttpQuery into query.Values for mercury upload: %e", err)
		}

		payloadBytes := queryValues.Encode()
		url := MERCURY_UPLOAD_URL + payloadBytes
		payload, contentType := c.NewMercuryMediaPayload(media)

		h := c.buildHeaders()
		h.Add("content-type", contentType)
		h.Add("origin", "https://www.facebook.com")
		h.Add("referer", "https://www.facebook.com/messages")
		h.Add("sec-fetch-dest", "empty")
		h.Add("sec-fetch-mode", "cors")
		h.Add("sec-fetch-site", "same-origin") // header is required
		
		_, respBody, err := c.MakeRequest(url, "POST", h, payload, types.NONE)
		if err != nil {
			return nil, fmt.Errorf("failed to send MercuryUploadRequest: %e", err)
		}

		resp, err := c.parseMercuryResponse(respBody)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mercury response: %e", err)
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (c *Client) parseMercuryResponse(respBody []byte) (*types.MercuryUploadResponse, error) {
	if len(respBody) < 9 {
		return nil, fmt.Errorf("mercury upload response body was less than 9 in size")
	}

	jsonData := respBody[9:]
	var mercuryResponse *types.MercuryUploadResponse
	if err := json.Unmarshal(jsonData, &mercuryResponse); err != nil {
		return nil, err
	}

	err := c.parseMetadata(mercuryResponse)
	if err != nil {
		return nil, err
	}

	return mercuryResponse, nil
}

func (c *Client) parseMetadata(response *types.MercuryUploadResponse) error {
	var err error

	switch metadata := response.Payload.Metadata.(type) {
	case []interface{}:
		var realMetadata types.ImageMetadata
		err = modules.InterfaceToStructJSON(metadata[0], &realMetadata)
		response.Payload.Metadata = &realMetadata
	case map[string]interface{}:
		var realMetadata types.VideoMetadata
		err = modules.InterfaceToStructJSON(metadata["0"], &realMetadata)
		response.Payload.Metadata = &realMetadata
	default:
		return fmt.Errorf("got invalid metadata type, cannot proceed with type assertion: %v", reflect.TypeOf(metadata))
	}

	return err
}

// returns payloadBytes, multipart content-type header
func (c *Client) NewMercuryMediaPayload(media *MercuryUploadMedia) ([]byte, string) {
	var mercuryPayload bytes.Buffer
	writer := multipart.NewWriter(&mercuryPayload)

	err := writer.SetBoundary("----WebKitFormBoundary" + methods.RandStr(16))
	if err != nil {
		log.Fatalf("Failed to set boundary: %v", err)
	}
	
	partHeader := textproto.MIMEHeader{
		"Content-Disposition": []string{`form-data; name="farr"; filename="` + media.Filename + `"`},
		"Content-Type": []string{string(media.MediaType)},
	}

	mediaPart, err := writer.CreatePart(partHeader)
	if err != nil {
		log.Fatalf("Failed to create multipart writer: %v", err)
	}

	_, err = mediaPart.Write(media.MediaData)
	if err != nil {
		log.Fatalf("Failed to write data to multipart section: %v", err)
	}

	err = writer.Close()
	if err != nil {
		log.Fatalf("Failed to close multipart writer: %v", err)
	}

	return mercuryPayload.Bytes(), writer.FormDataContentType()
}