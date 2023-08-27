package messagix

import (
	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/table"
)

func (s *Socket) handleBinaryMessage(data []byte) {
	s.client.Logger.Debug().Hex("hex-data", data).Bytes("bytes", data).Msg("Received BinaryMessage")
	if s.client.eventHandler == nil {
		return
	}

	resp := &Response{}
	err := resp.Read(data)
	if err != nil {
		s.handleErrorEvent(err)
	} else {
		switch evt := resp.ResponseData.(type) {
		case *Event_PublishACK, *Event_SubscribeACK:
			s.handleACKEvent(evt.(AckEvent))
		default:
			s.client.eventHandler(resp.ResponseData.Finish())
		}
	}
}

func (s *Socket) handleACKEvent(ackData AckEvent) {
	packetId := ackData.GetPacketId()
	err := s.packetHandler.updatePacketChannel(uint16(packetId), ackData)
	if err != nil {
		s.client.Logger.Info().Any("data", ackData).Any("packetId", packetId).Msg(err.Error())
		return
	}

	s.client.Logger.Info().Any("data", ackData).Any("packetId", packetId).Msg("Updated packet channel!")
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
	CurrentUser modules.CurrentUserInitialData `skip:"1"`
	Threads []table.LSDeleteThenInsertThread `skip:"1"`
	Messages []table.LSUpsertMessage `skip:"1"`
	Contacts []table.LSVerifyContactRowExists `skip:"1"`
}

func (e *Event_Ready) Finish() ResponseData {
	e.CurrentUser = modules.SchedulerJSDefined.CurrentUserInitialData
	e.Threads = modules.SchedulerJSRequired.LSTable.LSDeleteThenInsertThread
	e.Messages = modules.SchedulerJSRequired.LSTable.LSUpsertMessage
	e.Contacts = modules.SchedulerJSRequired.LSTable.LSVerifyContactRowExists
	return e
}

// Event_Error is emitted whenever the library encounters/receives an error.
// These errors can be for example: failed to send data, failed to read response data and so on.
type Event_Error struct {
	Err error
}

func (e *Event_Error) Finish() ResponseData {
	return e
}

// Event_SocketClosed is emitted whenever the websockets CloseHandler() is called.
// This provides great flexability because the user can then decide whether the client should reconnect or not.
type Event_SocketClosed struct {
	Code int
	Text string
}

func (e *Event_SocketClosed) Finish() ResponseData {
	return e
}

type AckEvent interface{
	GetPacketId() uint16
}

// Event_PublishACK is never emitted, it only handles the acknowledgement after a PUBLISH packet has been sent.
type Event_PublishACK struct {
	PacketId uint16
}

func (pb *Event_PublishACK) GetPacketId() uint16 {
	return pb.PacketId
}

func (pb *Event_PublishACK) Finish() ResponseData {
	return pb
}

// Event_SubscribeACK is never emitted, it only handles the acknowledgement after a SUBSCRIBE packet has been sent.
type Event_SubscribeACK struct {
	PacketId uint16
	QoSLevel uint8 // 0, 1, 2, 128
}

func (pb *Event_SubscribeACK) GetPacketId() uint16 {
	return pb.PacketId
}

func (pb *Event_SubscribeACK) Finish() ResponseData {
	return pb
}