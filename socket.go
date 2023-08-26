package messagix

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/0xzer/messagix/packets"
	"github.com/gorilla/websocket"
)

var (
	ErrSocketClosed      = errors.New("socket is closed")
	ErrSocketAlreadyOpen = errors.New("socket is already open")
)

type Socket struct {
	client *Client
	conn *websocket.Conn

	packetsSent uint16

	topics []Topic
}

func (c *Client) NewSocketClient() *Socket {
	return &Socket{
		client: c,
		packetsSent: 0,
	}
}

func (s *Socket) Connect() error {
	if s.conn != nil {
		s.client.Logger.Err(ErrSocketAlreadyOpen).Msg("Failed to initialize connection to socket")
		return ErrSocketAlreadyOpen
	}

	headers := s.getConnHeaders()
	brokerUrl := s.client.configs.mqttConfig.BuildBrokerUrl()

	s.client.Logger.Debug().Any("broker", brokerUrl).Msg("Dialing socket")
	conn, _, err := websocket.DefaultDialer.Dial(brokerUrl, headers)
	if err != nil {
		return fmt.Errorf("failed to dial socket: %s", err.Error())
	}

	conn.SetCloseHandler(s.CloseHandler)
	
	s.conn = conn

    err = s.sendConnectPacket()
    if err != nil {
        return fmt.Errorf("failed to send CONNECT packet to socket: %s", err.Error())
    }

	go s.beginReadStream()

	appSettingPublishJSON, err := s.newAppSettingsPublishJSON()
	if err != nil {
		return err
	}

	publishPacketByte := &packets.PublishPacket{
		QOSLevel: packets.QOS_LEVEL_1,
	}

	appSettingPublishPayload, err := s.client.NewPublishRequest(APP_SETTINGS, appSettingPublishJSON, s.packetsSent+1, publishPacketByte.Compress())
	if err != nil {
		return err
	}

	err = s.sendData(appSettingPublishPayload)
	if err != nil {
		return err
	}
	return nil
}

func (s *Socket) beginReadStream() {
	for {
		messageType, p, err := s.conn.ReadMessage()
		if err != nil {
			s.handleErrorEvent(fmt.Errorf("error reading message from websocket: %s", err.Error()))
			return
		}

		switch messageType {
			case websocket.TextMessage:
				s.client.Logger.Debug().Any("data", p).Bytes("bytes", p).Msg("Received TextMessage")
			case websocket.BinaryMessage:
				s.handleBinaryMessage(p)
		}
	}
}

func (s *Socket) sendData(data []byte) error {
	s.client.Logger.Debug().Any("data", data).Hex("hex", data).Msg("Sending data to socket")

	packetType := data[0] >> 4
	if packetType == packets.PUBLISH {
		s.packetsSent += 1
	}
	
	err := s.conn.WriteMessage(websocket.BinaryMessage, data)
    if err != nil {
        e := fmt.Errorf("error sending data to websocket: %s", err.Error())
		s.handleErrorEvent(e)
		return e
    }
    return nil
}

func (s *Socket) sendConnectPacket() error {
	connectAdditionalData, err := s.newConnectJSON()
	if err != nil {
		return err
	}

	connectFlags := CreateConnectFlagByte(ConnectFlags{CleanSession: true, Username: true})
	connectPayload, err := s.client.NewConnectRequest(connectAdditionalData, connectFlags)
	if err != nil {
		return err
	}
	return s.sendData(connectPayload)
}

func (s *Socket) CloseHandler(code int, text string) error {
	s.conn = nil

	if s.client.eventHandler != nil {
		socketCloseEvent := &Event_SocketClosed{
			Code: code,
			Text: text,
		}
		s.client.eventHandler(socketCloseEvent)
	}
	
	return nil
}

func (s *Socket) setTopics(topics []Topic) {
	s.topics = topics
}

func (s *Socket) getConnHeaders() http.Header {
	h := http.Header{}

	h.Add("cookie", s.client.cookies.ToString())
	h.Add("user-agent", USER_AGENT)
	h.Add("origin", "https://www.facebook.com")

	return h
}