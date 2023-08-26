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
	http *http.Client
	socket *Socket
	taskManager *TaskManager
	eventHandler EventHandler
	configs *Configs

	Logger zerolog.Logger
	cookies *types.Cookies
	proxy Proxy
}

// pass an empty zerolog.Logger{} for no logging
func NewClient(cookies *types.Cookies, logger zerolog.Logger, proxy Proxy) *Client {
	if cookies == nil {
		log.Fatal("cookie struct can not be nil")
	}

	cli := &Client{
		http: &http.Client{
			Transport: &http.Transport{
				Proxy: proxy,
			},
		},
		cookies: cookies,
		proxy: proxy,
		Logger: logger,
	}

	socket := cli.NewSocketClient()
	cli.socket = socket

	cli.configs = &Configs{client: cli}

	moduleLoader := &ModuleParser{}
	moduleLoader.load()


	configSetupErr := cli.configs.SetupConfigs()
	if configSetupErr != nil {
		log.Fatal(configSetupErr)
	}

	taskManager := &TaskManager{client: cli, activeTaskIds: make([]int, 0), currTasks: make([]TaskData, 0)}
	cli.taskManager = taskManager

	cli.taskManager.AddNewTask(&GetContactsTask{Limit: 100})

	log.Println(cli.taskManager.FinalizePayload())
	os.Exit(1)
	return cli
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