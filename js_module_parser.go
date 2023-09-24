package messagix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/types"
	"golang.org/x/net/html"
)

var versionPattern = regexp.MustCompile(`__d\("LSVersion"[^)]+\)\{e\.exports="(\d+)"\}`)

type ModuleData struct {
	Require [][]interface{} `json:"require,omitempty"`
}

type LinkTag struct {
	Attributes map[string]string
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
	doc, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		log.Fatalf("failed to parse doc string: %e", err)
	}
	scriptTags := m.findScriptTags(doc)
	for _, tag := range scriptTags {
		rel, ok := tag.Attributes["rel"]
		if ok {
			log.Println(rel)
		}
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
				log.Fatalf("failed to unmarshal content to moduleData: %e", err)
			}

			req := data.Require
		    for _, mod := range req {
				m.handleModule(mod)
			}
		}
	}

	// on certain occasions, the server does not return the lightspeed data or version
	// when this is the case, the server "preloads" the js files in the link tags, so we need to loop through them until we can find the "LSVersion" module and extract the exported version string
	if modules.VersionId == 0 {
		m.client.configs.needSync = true
		m.client.Logger.Info().Msg("Setting configs.needSync to true")
		var doneCrawling bool
		linkTags := m.findLinkTags(doc)
		for _, tag := range linkTags {
			if doneCrawling {
				break
			}
			as := tag.Attributes["as"]
			href := tag.Attributes["href"]

			switch as {
			case "script":
				doneCrawling, err = m.crawlJavascriptFile(href)
				if err != nil {
					log.Fatalf("failed to crawl js file %s: %e", href, err)
				}
			}
		}
	}
}

func (m *ModuleParser) crawlJavascriptFile(href string) (bool, error) {
	_, jsContent, err := m.client.MakeRequest(href, "GET", http.Header{}, nil, types.NONE)
	if err != nil {
		return false, err
	}

	if err != nil {
		log.Fatal(err)
	}
	
	versionMatches := versionPattern.FindStringSubmatch(string(jsContent))
	if len(versionMatches) > 0 {
		versionInt, err := strconv.ParseInt(versionMatches[1], 10, 64)
		if err != nil {
			return false, err
		}
		modules.VersionId = versionInt
		return true, nil
	}
	return false, nil
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
		/*
		add later if needed for the gkx data
		case "HasteSupportData":
			log.Println("got haste support data!")
			m.client.Logger.Debug().Any("data", modData).Msg("Got haste support data")
			os.Exit(1)
		}
		*/
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

func (m *ModuleParser) findLinkTags(n *html.Node) []LinkTag {
	var tags []LinkTag
	if n.Type == html.ElementNode && n.Data == "link" {
		attributes := make(map[string]string)
		for _, a := range n.Attr {
			attributes[a.Key] = a.Val
		}
		tags = append(tags, LinkTag{Attributes: attributes})
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tags = append(tags, m.findLinkTags(c)...)
	}
	return tags
}