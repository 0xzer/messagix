package messagix

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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
	packetHandler *PacketHandler
	topics []Topic
}

func (c *Client) NewSocketClient() *Socket {
	return &Socket{
		client: c,
		packetHandler: &PacketHandler{
			packetsSent: 1, // starts at 1
			packetChannels: make(map[uint16]chan interface{}, 0),
			packetTimeout: time.Second * 10, // 10 sec timeout if puback is not received
		},
	}
}

func (s *Socket) Connect() error {
	if s.conn != nil {
		s.client.Logger.Err(ErrSocketAlreadyOpen).Msg("Failed to initialize connection to socket")
		return ErrSocketAlreadyOpen
	}

	headers := s.getConnHeaders()
	brokerUrl := s.client.configs.mqttConfig.BuildBrokerUrl()

	dialer := websocket.Dialer{}
	if s.client.proxy != nil {
		dialer.Proxy = s.client.proxy
	}

	s.client.Logger.Debug().Any("broker", brokerUrl).Msg("Dialing socket")
	conn, _, err := dialer.Dial(brokerUrl, headers)
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

	appSettingPublishJSON, err := s.newAppSettingsPublishJSON(s.client.configs.siteConfig.VersionId)
	if err != nil {
		return err
	}

	packetId, err := s.sendPublishPacket(APP_SETTINGS, appSettingPublishJSON, &packets.PublishPacket{QOSLevel: packets.QOS_LEVEL_1})
	if err != nil {
		return fmt.Errorf("failed to send APP_SETTINGS publish packet: %e", err)
	}

	appSettingAck := s.packetHandler.waitForPubACKDetails(packetId)
	if appSettingAck == nil {
		// log.Fatal because if this is never received then why would the other packets work
		// unless later on add some sort of retry-send-packet mechanism
		return fmt.Errorf("failed to get pubAck for packetId: %d", appSettingAck.PacketId)
	}
	
	_, err = s.sendSubscribePacket(FOREGROUND_STATE, packets.QOS_LEVEL_0)
	if err != nil {
		return fmt.Errorf("failed to subscribe to FOREGROUND_STATE topic: %e", err)
	}

	_, err = s.sendSubscribePacket(RESP, packets.QOS_LEVEL_0)
	if err != nil {
		return fmt.Errorf("failed to subscribe to LS_RESP topic: %e", err)
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
	s.client.Logger.Debug().Bytes("bytes", data).Hex("hex", data).Msg("Sending data to socket")

	packetType := data[0] >> 4
	if packetType == packets.PUBLISH || packetType == packets.SUBSCRIBE{
		s.packetHandler.packetsSent++
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

func (s *Socket) sendSubscribePacket(topic Topic, qos packets.QoS) (*Event_SubscribeACK, error) {
	subscribeRequestPayload, packetId, err := s.client.NewSubscribeRequest(topic, qos)
	if err != nil {
		return nil, err
	}

	err = s.sendData(subscribeRequestPayload)
	if err != nil {
		return nil, err
	}

	resp := s.packetHandler.waitForSubACKDetails(packetId)
	if resp == nil {
		return nil, fmt.Errorf("did not receive SubACK packet for packetid: %d", packetId)
	}

	s.client.Logger.Debug().Any("resp", resp).Any("topic", topic).Msg("Successfully subscribed to topic!")
	return resp, nil
}

func (s *Socket) sendPublishPacket(topic Topic, jsonData string, packet *packets.PublishPacket) (uint16, error) {
	publishRequestPayload, packetId, err := s.client.NewPublishRequest(topic, jsonData, packet.Compress())
	if err != nil {
		return packetId, err
	}

	return packetId, s.sendData(publishRequestPayload)
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