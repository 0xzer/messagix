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

	Logger zerolog.Logger

	http *http.Client
	socket *Socket
	taskManager *TaskManager
	eventHandler EventHandler
	graphQl *GraphQL
	configs *Configs

	cookies *types.Cookies
	proxy Proxy
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


	configSetupErr := cli.configs.SetupConfigs()
	if configSetupErr != nil {
		log.Fatal(configSetupErr)
	}

	cli.Account = &Account{client: cli}

	taskManager := &TaskManager{client: cli, activeTaskIds: make([]int, 0), currTasks: make([]TaskData, 0)}
	cli.taskManager = taskManager

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

// Sets the topics the client should subscribe to
func (c *Client) SetTopics(topics []Topic) {
	c.socket.setTopics(topics)
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