package messagix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/0xzer/messagix/cookies"
	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/packets"
	"github.com/0xzer/messagix/types"
	"github.com/gorilla/websocket"
)

var (
	connectionTypes = map[types.Platform]string{
		types.Instagram: "cookie_auth",
		types.Facebook: "websocket",
	}
	ErrSocketClosed      = errors.New("socket is closed")
	ErrSocketAlreadyOpen = errors.New("socket is already open")
)

var handshakeBytes = []byte{192, 0} // pingreq packet

type Socket struct {
	client *Client
	conn *websocket.Conn
	responseHandler *ResponseHandler
	mu *sync.Mutex
    packetsSent uint16
	handshakeInterval *time.Ticker
}

func (c *Client) NewSocketClient() *Socket {
	return &Socket{
		client: c,
		responseHandler: &ResponseHandler{
			client: c,
			requestChannels: make(map[uint16]chan interface{}, 0),
			packetChannels: make(map[uint16]chan interface{}, 0),
			packetTimeout: time.Second * 10, // 10 sec timeout if puback is not received
		},
		mu: &sync.Mutex{},
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
				go s.handleBinaryMessage(p)
		}
	}
}

func (s *Socket) startHandshakeInterval() {
    if s.handshakeInterval != nil {
        s.handshakeInterval.Stop()
    }

    s.handshakeInterval = time.NewTicker(15 * time.Second)
	s.client.Logger.Info().Msg("Starting handshakeInterval...")
    go func() {
        for range s.handshakeInterval.C {
            s.sendHandshake()
			s.client.Logger.Info().Msg("Sent handshake to socket")
        }
    }()
}

func (s *Socket) sendHandshake() {
	err := s.sendData(handshakeBytes)
	if err != nil {
		log.Fatalf("failed to send handshake data: %e", err)
	}

	// TO-DO implement channels to handle the server not responding with a PINGRESP packet
}

func (s *Socket) sendData(data []byte) error {
	//s.client.Logger.Debug().Any("hex", debug.BeautifyHex(data)).Msg("Sending data to socket")
	err := s.conn.WriteMessage(websocket.BinaryMessage, data)
    if err != nil {
        e := fmt.Errorf("error sending data to websocket: %s", err.Error())
		s.handleErrorEvent(e)
		return e
    }

    return nil
}

func (s *Socket) SafePacketId() uint16 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.packetsSent++
	if s.packetsSent == 0 || s.packetsSent > 65535 {
		s.packetsSent = 1
	}
	return s.packetsSent
}

func (s *Socket) sendConnectPacket() error {
	connectAdditionalData, err := s.newConnectJSON()
	if err != nil {
		return err
	}

	connectFlags := packets.CreateConnectFlagByte(packets.ConnectFlags{CleanSession: true, Username: true})
	connectPayload, err := s.client.NewConnectRequest(connectAdditionalData, connectFlags)
	if err != nil {
		return err
	}
	return s.sendData(connectPayload)
}

func (s *Socket) sendSubscribePacket(topic Topic, qos packets.QoS, wait bool) (*Event_SubscribeACK, error) {
	subscribeRequestPayload, packetId, err := s.client.NewSubscribeRequest(topic, qos)
	if err != nil {
		return nil, err
	}

	err = s.sendData(subscribeRequestPayload)
	if err != nil {
		return nil, err
	}

	var resp *Event_SubscribeACK
	if wait {
		resp = s.responseHandler.waitForSubACKDetails(packetId)
		if resp == nil {
			return nil, fmt.Errorf("did not receive SubACK packet for packetid: %d", packetId)
		}
	}

	return resp, nil
}

func (s *Socket) sendPublishPacket(topic Topic, jsonData string, packet *packets.PublishPacket, packetId uint16) (uint16, error) {
	publishRequestPayload, packetId, err := s.client.NewPublishRequest(topic, jsonData, packet.Compress(), packetId)
	if err != nil {
		return packetId, err
	}

	s.responseHandler.addPacketChannel(packetId)
	return packetId, s.sendData(publishRequestPayload)
}

type SocketLSRequestPayload struct {
	AppId string `json:"app_id"`
	Payload string `json:"payload"`
	RequestId int `json:"request_id"`
	Type int `json:"type"`
}

func (s *Socket) makeLSRequest(payload []byte, t int) (uint16, error) {
	packetId := s.SafePacketId()
	lsPayload := &SocketLSRequestPayload{
		AppId: modules.SchedulerJSDefined.CurrentUserInitialData.AppID,
		Payload: string(payload),
		RequestId: int(packetId),
		Type: t,
	}

	jsonPayload, err := json.Marshal(lsPayload)
	if err != nil {
		return 0, err
	}

	_, err = s.sendPublishPacket(LS_REQ, string(jsonPayload), &packets.PublishPacket{QOSLevel: packets.QOS_LEVEL_1}, packetId)
	if err != nil {
		return 0, err
	}

	return packetId, nil
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

func (s *Socket) getConnHeaders() http.Header {
	h := http.Header{}

	h.Add("cookie", cookies.CookiesToString(s.client.cookies))
	h.Add("user-agent", USER_AGENT)
	h.Add("origin", s.client.getEndpoint("base_url"))

	return h
}

func (s *Socket) getConnectionType() string {
	return connectionTypes[s.client.platform]
}