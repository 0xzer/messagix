package messagix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/types"
	"golang.org/x/net/html"
)

type ModuleData struct {
	Require [][]interface{} `json:"require,omitempty"`
}

type ScriptTag struct {
	Attributes map[string]string
	Content    string
}

type ModuleParser struct {
	client *Client
}

func (m *ModuleParser) fetchMainSite() []byte { // just log.fatal if theres an error because the library should not be able to continue then
	headers := m.client.buildHeaders()
	headers.Add("connection", "keep-alive")
	headers.Add("host", "www.facebook.com")
	headers.Add("sec-fetch-dest", "document")
	headers.Add("sec-fetch-mode", "navigate")
	headers.Add("sec-fetch-site", "none") // header is required, otherwise they dont send the csr bitmap data in the response. lets also include the other headers to be sure
	headers.Add("sec-fetch-user", "?1")
	headers.Add("upgrade-insecure-requests", "1")
	_, responseBody, err := m.client.MakeRequest("https://www.facebook.com/messages", "GET", headers, nil, types.NONE)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to fetch main messages site: %e", err))
	}

	return responseBody
}

func (m *ModuleParser) load() {
	htmlData := m.fetchMainSite()
	os.WriteFile("res.html", htmlData, os.ModePerm)
	doc, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		log.Fatalf("failed to parse doc string: %e", err)
	}
	scriptTags := m.findScriptTags(doc)
	for _, tag := range scriptTags {
		id := tag.Attributes["id"]
		switch id {
		case "envjson", "__eqmc":
			modules.HandleJSON([]byte(tag.Content), id)
		default:
			if tag.Content == "" {
				continue
			}
			var data *ModuleData
			err := json.Unmarshal([]byte(tag.Content), &data)
			if err != nil {
				os.WriteFile("test.json", []byte(tag.Content), os.ModePerm)
				log.Fatalf("failed to unmarshal content to moduleData: %e", err)
			}

			req := data.Require
		    for _, mod := range req {
				m.handleModule(mod)
			}
		}
	}
}

func (m *ModuleParser) handleModule(data []interface{}) {
	modName := data[0].(string)
	modData := data[3].([]interface{})
	switch modName {
		case "ScheduledServerJS", "ScheduledServerJSWithCSS":
			method := data[1].(string)
			for _, d := range modData {
				switch method {
				case "handle":
					err := modules.SSJSHandle(d)
					if err != nil {
						log.Fatalf("failed to handle scheduledserverjs module: %e", err)
					}
				}
			}
		case "Bootloader":
			method := data[1].(string)
			for _, d := range modData {
				switch method {
					case "handlePayload":
						err := modules.HandlePayload(d, &modules.SchedulerJSDefined.BootloaderConfig)
						if err != nil {
							log.Fatalf("failed to handle Bootloader_handlePayload call: %e", err)
						}
						//debug.Debug().Any("csrBitmap", modules.CsrBitmap).Msg("handlePayload")
				}
			}
		}
}

func (m *ModuleParser) findScriptTags(n *html.Node) []ScriptTag {
	var tags []ScriptTag
	if n.Type == html.ElementNode && n.Data == "script" {
		attributes := make(map[string]string)
		for _, a := range n.Attr {
			attributes[a.Key] = a.Val
		}
		content := ""
		if n.FirstChild != nil {
			content = n.FirstChild.Data
		}
		tags = append(tags, ScriptTag{Attributes: attributes, Content: strings.Replace(content, ",BootloaderConfig", ",", -1)})
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tags = append(tags, m.findScriptTags(c)...)
	}
	return tags
}