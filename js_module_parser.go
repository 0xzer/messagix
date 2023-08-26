package messagix

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"github.com/0xzer/messagix/modules"
	"golang.org/x/net/html"
)

type ModuleData struct {
	Require [][]interface{} `json:"require,omitempty"`
}

type ScriptTag struct {
	Attributes map[string]string
	Content    string
}

type ModuleParser struct {}

func (m *ModuleParser) load() {
	docStr, _ := os.ReadFile("res.html")
	doc, err := html.Parse(strings.NewReader(string(docStr)))
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