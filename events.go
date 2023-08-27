package messagix

import (
	"log"
	"os"

	"github.com/0xzer/messagix/packets"
)

func (s *Socket) handleBinaryMessage(data []byte) {
	s.client.Logger.Debug().Hex("hex-data", data).Bytes("bytes", data).Any("str", string(data)).Msg("Received BinaryMessage")
	if s.client.eventHandler == nil {
		return
	}

	resp := &Response{}
	err := resp.Read(data)
	//log.Println(resp)
	if err != nil {
		s.handleErrorEvent(err)
	} else {
		switch evt := resp.ResponseData.(type) {
		case *packets.PublishACK:
			s.handlePublishACK(int16(evt.PacketId))
			os.Exit(1)
		default:
			log.Println(resp.ResponseData)
			s.client.eventHandler(resp.ResponseData)
		}
	}
}

func (s *Socket) handlePublishACK(packetId int16) {
	log.Println("packetid:",packetId)
	os.Exit(1)
}

func (s *Socket) handleErrorEvent(err error) {
	errEvent := &Event_Error{Err: err}
	s.client.eventHandler(errEvent)
}

// Event_Ready represents the CONNACK packet's response.
// The library provides the raw parsed data, so handle connection codes as needed for your application.
type Event_Ready struct {
	IsNewSession   bool
	ConnectionCode ConnectionCode
}

// Event_Error is emitted whenever the library encounters/receives an error.
// These errors can be for example: failed to send data, failed to read response data and so on.
type Event_Error struct {
	Err error
}

// Event_SocketClosed is emitted whenever the websockets CloseHandler() is called.
// This provides great flexability because the user can then decide whether the client should reconnect or not.
type Event_SocketClosed struct {
	Code int
	Text string
}