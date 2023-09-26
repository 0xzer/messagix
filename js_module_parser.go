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
	"github.com/0xzer/messagix/methods"
	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/types"
	"golang.org/x/net/html"
)

var jsDatrPattern = regexp.MustCompile(`"_js_datr","([^"]+)"`)
var versionPattern = regexp.MustCompile(`__d\("LSVersion"[^)]+\)\{e\.exports="(\d+)"\}`)

type ModuleData struct {
	Define    [][]interface{} `json:"define,omitempty"`
	Instances [][]interface{} `json:"instances,omitempty"`
	Markup    [][]interface{} `json:"markup,omitempty"`
	Elements  [][]interface{} `json:"elements,omitempty"`
	Require   [][]interface{} `json:"require,omitempty"`
	Contexts  [][]interface{} `json:"contexts,omitempty"`
}

type AttributeMap map[string]string
type ScriptTag struct {
	Attributes AttributeMap
	Content    string
}

type LinkTag struct {
	Attributes AttributeMap
}

type FormTag struct {
	Attributes AttributeMap
	Inputs     []InputTag
}

type InputTag struct {
	Attributes AttributeMap
}

type ModuleParser struct {
	client *Client
	FormTags []FormTag
	LoginInputs []InputTag
	JSDatr string
}

func (m *ModuleParser) fetchPageData(page string) []byte { // just log.fatal if theres an error because the library should not be able to continue then
	headers := m.client.buildHeaders()
	headers.Add("connection", "keep-alive")
	headers.Add("host", "www.facebook.com")
	headers.Add("sec-fetch-dest", "document")
	headers.Add("sec-fetch-mode", "navigate")
	headers.Add("sec-fetch-site", "none") // header is required, otherwise they dont send the csr bitmap data in the response. lets also include the other headers to be sure
	headers.Add("sec-fetch-user", "?1")
	headers.Add("upgrade-insecure-requests", "1")
	_, responseBody, err := m.client.MakeRequest(page, "GET", headers, nil, types.NONE)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to fetch page data for module parser: %e", err))
	}

	return responseBody
}

func (m *ModuleParser) load(page string) {
	htmlData := m.fetchPageData(page)
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
				tagStr := string(tag.Content)
				if strings.HasPrefix(tagStr, "requireLazy") {
					m.requireLazyModule(tagStr)
					continue
				}
				//log.Println(fmt.Sprintf("failed to unmarshal content to moduleData: %e | content=%s", err, string(tag.Content)))
				continue
			}

			req := data.Require
		    for _, mod := range req {
				m.handleModule(mod)
			}
		}
	}

	loggedIn := modules.SchedulerJSDefined.CurrentUserInitialData.AccountID != "0"
	// on certain occasions, the server does not return the lightspeed data or version
	// when this is the case, the server "preloads" the js files in the link tags, so we need to loop through them until we can find the "LSVersion" module and extract the exported version string
	if modules.VersionId == 0  && loggedIn {
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

	formTags := m.findFormTags(doc)
	m.FormTags = formTags
	if !loggedIn {
		loginNode := m.findNodeByID(doc, "loginform")
		loginInputs := m.findInputTags(loginNode)
		m.LoginInputs = loginInputs
	}
}

type BigPipe struct {
	JSMods *ModuleData `json:"jsmods,omitempty"`
}

func (m *ModuleParser) requireLazyModule(data string) {
	moduleSplit := strings.Split(data[12:], "],")
	var moduleNames []string
	err := json.Unmarshal([]byte(moduleSplit[0] + "]"), &moduleNames)
	if err != nil {
		log.Fatalf("failed to get module names from requireLazy module: %e", err)
	}

	for _, mName := range moduleNames {
		switch mName {
		case "__bigPipe":
			if len(data) < 5000 {
				continue
			}
			handleData := strings.Split(strings.Split(data, "bigPipe.onPageletArrive(")[1], ");requireLazy")[0]
			handleData = methods.PreprocessJSObject(handleData[:len(handleData)-6])
			var bigPipeData *BigPipe
			err = json.Unmarshal([]byte(handleData), &bigPipeData)
			if err != nil {
				log.Fatalf("failed to unmarshal BigPipe data")
			}

			for _, d := range bigPipeData.JSMods.Define {
				err := modules.SSJSHandle(d)
				if err != nil {
					log.Fatalf("failed to handle serverjs module: %e", err)
				}
			}

			if m.client.cookies == nil {
				jsDatrMatches := jsDatrPattern.FindStringSubmatch(handleData)
				if len(jsDatrMatches) > 1 {
					m.JSDatr = jsDatrMatches[1]
				}
			}
		case "ServerJS":
			handleData := "{" + strings.Split(strings.Split(data, ".handle({")[1], ");requireLazy")[0]
			var moduleData *ModuleData
			err = json.Unmarshal([]byte(handleData), &moduleData)
			if err != nil {
				log.Fatalf("failed to unmarshal handleData[0] into struct *ModuleData")
			}

			for _, d := range moduleData.Define {
				err := modules.SSJSHandle(d)
				if err != nil {
					log.Fatalf("failed to handle serverjs module: %e", err)
				}
			}
		case "HasteSupportData":
			handleSplit := strings.Split(data, ".handle({")
			if len(handleSplit) <= 2 {
				//log.Println("skipping HasteSupportData because it contains irrelevant stuff")
				continue
			}
			handleData := "{" + strings.Split(strings.Split(data, ".handle({")[2], ");requireLazy")[0]
			handleData = handleData[:len(handleData)-5]
			var moduleData *ModuleData
			err = json.Unmarshal([]byte(handleData), &moduleData)
			if err != nil {
				log.Fatalf("failed to unmarshal handleData[0] into struct *ModuleData")
			}
			for _, d := range moduleData.Define {
				err := modules.SSJSHandle(d)
				if err != nil {
					log.Fatalf("failed to handle hastesupportdata module: %e", err)
				}
			}
		case "bootstrapWebSession":
			continue
		default:
			//log.Println(fmt.Sprintf("Skipping requireLazy module with name: %s", mName))
			continue
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

type NodeProcessor func(*html.Node) interface{}

func (m *ModuleParser) processNode(n *html.Node, tag string, processor NodeProcessor) interface{} {
	if n.Type == html.ElementNode && n.Data == tag {
		return processor(n)
	}
	return nil
}

func (m *ModuleParser) recursiveSearch(n *html.Node, tag string, processor NodeProcessor) []interface{} {
	var result []interface{}

	if item := m.processNode(n, tag, processor); item != nil {
		result = append(result, item)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, m.recursiveSearch(c, tag, processor)...)
	}

	return result
}

func (m *ModuleParser) findTags(tag string, processor NodeProcessor, n *html.Node) []interface{} {
	return m.recursiveSearch(n, tag, processor)
}

func (m *ModuleParser) findScriptTags(n *html.Node) []ScriptTag {
	processor := func(n *html.Node) interface{} {
		attributes := make(AttributeMap)
		for _, a := range n.Attr {
			attributes[a.Key] = a.Val
		}
		content := ""
		if n.FirstChild != nil {
			content = n.FirstChild.Data
		}
		return ScriptTag{Attributes: attributes, Content: strings.Replace(content, ",BootloaderConfig", ",", -1)}
	}

	tags := m.findTags("script", processor, n)
	scriptTags := make([]ScriptTag, len(tags))
	for i, t := range tags {
		scriptTags[i] = t.(ScriptTag)
	}
	return scriptTags
}

func (m *ModuleParser) findLinkTags(n *html.Node) []LinkTag {
	processor := func(n *html.Node) interface{} {
		attributes := make(AttributeMap)
		for _, a := range n.Attr {
			attributes[a.Key] = a.Val
		}
		return LinkTag{Attributes: attributes}
	}

	tags := m.findTags("link", processor, n)
	linkTags := make([]LinkTag, len(tags))
	for i, t := range tags {
		linkTags[i] = t.(LinkTag)
	}
	return linkTags
}

// For Form Tags:
func (m *ModuleParser) findFormTags(n *html.Node) []FormTag {
	processor := func(n *html.Node) interface{} {
		formAttributes := make(AttributeMap)
		for _, a := range n.Attr {
			formAttributes[a.Key] = a.Val
		}
		var formInputs []InputTag
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "input" {
				inputAttributes := make(AttributeMap)
				for _, a := range c.Attr {
					inputAttributes[a.Key] = a.Val
				}
				formInputs = append(formInputs, InputTag{Attributes: inputAttributes})
			}
		}
		return FormTag{Attributes: formAttributes, Inputs: formInputs}
	}

	tags := m.findTags("form", processor, n)
	formTags := make([]FormTag, len(tags))
	for i, t := range tags {
		formTags[i] = t.(FormTag)
	}
	return formTags
}

func (m *ModuleParser) findInputTags(n *html.Node) []InputTag {
	processor := func(n *html.Node) interface{} {
		attributes := make(AttributeMap)
		for _, a := range n.Attr {
			attributes[a.Key] = a.Val
		}
		return InputTag{Attributes: attributes}
	}

	tags := m.findTags("input", processor, n)
	inputTags := make([]InputTag, len(tags))
	for i, t := range tags {
		inputTags[i] = t.(InputTag)
	}
	return inputTags
}

func (m *ModuleParser) findNodeByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "id" && a.Val == id {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := m.findNodeByID(c, id)
		if result != nil {
			return result
		}
	}
	return nil
}