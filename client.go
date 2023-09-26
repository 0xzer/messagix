package messagix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/0xzer/messagix/types"
	"github.com/google/go-querystring/query"
	"github.com/rs/zerolog"
)

var USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
var ErrRedirectAttempted = errors.New("redirect attempted")

type EventHandler func(evt interface{})
type Proxy func(*http.Request) (*url.URL, error)
type Client struct {
	Account *Account
	Threads *Threads
	Messages *Messages

	Logger zerolog.Logger

	http *http.Client
	socket *Socket
	eventHandler EventHandler
	configs *Configs
	syncManager *SyncManager

	cookies *types.Cookies
	proxy Proxy

	lsRequests int
	graphQLRequests int
}

// pass an empty zerolog.Logger{} for no logging
func NewClient(cookies *types.Cookies, logger zerolog.Logger, proxy string) *Client {

	cli := &Client{
		http: &http.Client{
			Transport: &http.Transport{},
		},
		cookies: cookies,
		Logger: logger,
		lsRequests: 0,
		graphQLRequests: 1,
	}

	cli.Account = &Account{client: cli}
	cli.Threads = &Threads{client: cli}
	cli.Messages = &Messages{client: cli}
	cli.configs = &Configs{client: cli, needSync: false}
	if proxy != "" {
		err := cli.SetProxy(proxy)
		if err != nil {
			log.Fatalf("failed to set proxy: %e", err)
		}
	}

	if cookies == nil || cookies.Xs == "" {
		return cli
	}

	socket := cli.NewSocketClient()
	cli.socket = socket

	moduleLoader := &ModuleParser{client: cli}
	moduleLoader.load("https://www.facebook.com/messages")

	cli.syncManager = cli.NewSyncManager()
	configSetupErr := cli.configs.SetupConfigs()
	if configSetupErr != nil {
		log.Fatal(configSetupErr)
	}
	
	return cli
}

func (c *Client) SetProxy(proxy string) error {
	proxyParsed, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	proxyUrl := http.ProxyURL(proxyParsed)
	c.http.Transport = &http.Transport{
		Proxy: proxyUrl,
	}
	c.proxy = proxyUrl
	c.Logger.Debug().Any("addr", proxyParsed.Host).Msg("Proxy Updated")
	return nil
}

func (c *Client) SetEventHandler(handler EventHandler) {
	c.eventHandler = handler
}

func (c *Client) Connect() error {
	return c.socket.Connect()
}

func (c *Client) SaveSession(path string) error {
	jsonBytes, err := json.Marshal(c.cookies)
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonBytes, os.ModePerm)
}

var cookieConsentUrl = "https://www.facebook.com/cookie/consent/"
func (c *Client) sendCookieConsent(jsDatr string) error {
	h := c.buildHeaders()
	h.Add("sec-fetch-dest", "empty")
	h.Add("sec-fetch-mode", "cors")
	h.Add("sec-fetch-site", "same-origin") // header is required, otherwise they dont send the csr bitmap data in the response. lets also include the other headers to be sure
	h.Add("sec-fetch-user", "?1")
	h.Add("host", "www.facebook.com")
	h.Add("upgrade-insecure-requests", "1")
	h.Add("origin", "https://www.facebook.com")
	h.Add("cookie", "_js_datr="+jsDatr)
	h.Add("referer", "https://www.facebook.com/login")

	payload := c.NewHttpQuery()
	payload.AcceptOnlyEssential = "false"
	form, err := query.Values(payload)
	if err != nil {
		return err
	}

	payloadBytes := []byte(form.Encode())
	req, _, err := c.MakeRequest(cookieConsentUrl, "POST", h, payloadBytes, types.FORM)
	if err != nil {
		return err
	}

	datr := c.findCookie(req.Cookies(), "datr")
	if datr == nil {
		return fmt.Errorf("consenting to cookies failed, could not find datr cookie in set-cookie header")
	}

	c.cookies = &types.Cookies{
		Datr: datr.Value,
		Wd: "901x1156",
	}

	return nil
}

func (c *Client) enableRedirects() {
	c.http.CheckRedirect = nil
}

func (c *Client) disableRedirects() {
	c.http.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return ErrRedirectAttempted
	}
}