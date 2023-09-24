package messagix

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"github.com/0xzer/messagix/types"
	"github.com/rs/zerolog"
)

var USER_AGENT = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"

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
	if cookies == nil {
		log.Fatal("cookie struct can not be nil")
	}

	cli := &Client{
		http: &http.Client{
			Transport: &http.Transport{},
		},
		cookies: cookies,
		Logger: logger,
		lsRequests: 0,
		graphQLRequests: 1,
	}

	if proxy != "" {
		err := cli.SetProxy(proxy)
		if err != nil {
			log.Fatalf("failed to set proxy: %e", err)
		}
	}

	socket := cli.NewSocketClient()
	cli.socket = socket

	cli.configs = &Configs{client: cli, needSync: false}

	moduleLoader := &ModuleParser{client: cli}
	moduleLoader.load()

	cli.syncManager = cli.NewSyncManager()
	configSetupErr := cli.configs.SetupConfigs()
	if configSetupErr != nil {
		log.Fatal(configSetupErr)
	}
	
	cli.Account = &Account{client: cli}
	cli.Threads = &Threads{client: cli}
	cli.Messages = &Messages{client: cli}
	
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