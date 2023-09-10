package messagix

import (
	"encoding/json"
	"log"
	"os"
	"github.com/0xzer/messagix/lightspeed"
	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/packets"
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
		case *Event_PublishResponse:
			s.handlePublishResponseEvent(evt)
		case *Event_PublishACK, *Event_SubscribeACK:
			s.handleACKEvent(evt.(AckEvent))
		case *Event_Ready:
			go s.handleReadyEvent(evt)
		default:
			s.client.Logger.Info().Any("data", data).Msg("sending default event...")
			s.client.eventHandler(resp.ResponseData.Finish())
		}
	}
}

func (s *Socket) handleReadyEvent(data *Event_Ready) {
	appSettingPublishJSON, err := s.newAppSettingsPublishJSON(s.client.configs.siteConfig.VersionId)
	if err != nil {
		log.Fatal(err)
	}

	packetId, err := s.sendPublishPacket(LS_APP_SETTINGS, appSettingPublishJSON, &packets.PublishPacket{QOSLevel: packets.QOS_LEVEL_1}, s.SafePacketId())
	if err != nil {
		log.Fatalf("failed to send APP_SETTINGS publish packet: %e", err)
	}

	appSettingAck := s.responseHandler.waitForPubACKDetails(packetId)
	if appSettingAck == nil {
		log.Fatalf("failed to get pubAck for packetId: %d", appSettingAck.PacketId)
	}

	_, err = s.sendSubscribePacket(LS_FOREGROUND_STATE, packets.QOS_LEVEL_0, true)
	if err != nil {
		log.Fatalf("failed to subscribe to ls_foreground_state: %e", err)
	}

	_, err = s.sendSubscribePacket(LS_RESP, packets.QOS_LEVEL_0, true)
	if err != nil {
		log.Fatalf("failed to subscribe to ls_resp: %e", err)
	}

	err = s.client.Account.ReportAppState(table.FOREGROUND)
	if err != nil {
		log.Fatalf("failed to report app state to foreground (active): %e", err)
	}

	s.client.eventHandler(data.Finish())
}

func (s *Socket) handleACKEvent(ackData AckEvent) {
	packetId := ackData.GetPacketId()
	err := s.responseHandler.updatePacketChannel(uint16(packetId), ackData)
	if err != nil {
		s.client.Logger.Err(err).Any("data", ackData).Any("packetId", packetId).Msg("failed to handle ack event")
		return
	}
}

func (s *Socket) handleErrorEvent(err error) {
	errEvent := &Event_Error{Err: err}
	s.client.eventHandler(errEvent)
}

func (s *Socket) handlePublishResponseEvent(resp *Event_PublishResponse) {
	s.client.Logger.Info().Any("requestId", resp.Data.RequestID).Msg("Got publish response event...")
	packetId := resp.Data.RequestID
	hasPacket := s.responseHandler.hasPacket(uint16(packetId))

	switch resp.Topic {
		case string(LS_RESP):
			if hasPacket {
				err := s.responseHandler.updateRequestChannel(uint16(packetId), resp.Finish())
				if err != nil {
					s.handleErrorEvent(err)
					return
				}
				return
			} else if packetId == 0 {
				if len(resp.Table.LSInsertMessage) > 0 {
					s.client.Logger.Info().Any("table", resp.Table).Msg("Got a new message!")
					os.Exit(1)
				}
				s.client.Logger.Info().Any("packetId", packetId).Any("table", resp.Table).Any("data", resp.Data).Any("topic", resp.Topic).Msg("Got unknown socket event...")
				return
			}
			s.client.Logger.Info().Any("packetId", packetId).Any("data", resp).Msg("Got publish response but was not expecting it for specific packet identifier.")
		default:
			s.client.Logger.Info().Any("packetId", packetId).Any("topic", resp.Topic).Any("data", resp.Data).Msg("Got unknown publish response topic!")
	}
}

// Event_Ready represents the CONNACK packet's response.
//
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
//
// These errors can be for example: failed to send data, failed to read response data and so on.
type Event_Error struct {
	Err error
}

func (e *Event_Error) Finish() ResponseData {
	return e
}

// Event_SocketClosed is emitted whenever the websockets CloseHandler() is called.
//
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

// Event_PublishResponse is never emitted, instead we will convert this into seperate events
//
// It will also be used for handling the responses after calling a function like GetContacts through the requestId
type Event_PublishResponse struct {
	Topic string `lengthType:"uint16" endian:"big"`
	Data PublishResponseData `jsonString:"1"`
	Table table.LSTable
}

type PublishResponseData struct {
	RequestID int64 `json:"request_id,omitempty"`
	Payload   string `json:"payload,omitempty"`
	Sp     []string `json:"sp,omitempty"` // dependencies
	Target int      `json:"target,omitempty"`
}

func (pb *Event_PublishResponse) Finish() ResponseData {
	pb.Table = table.LSTable{}
	var lsData *lightspeed.LightSpeedData
	err := json.Unmarshal([]byte(pb.Data.Payload), &lsData)
	if err != nil {
		log.Println("failed to unmarshal PublishResponseData JSON payload into lightspeed.LightSpeedData struct: %e", err)
		return pb
	}

	dependencies := table.SPToDepMap(pb.Data.Sp)
	decoder := lightspeed.NewLightSpeedDecoder(dependencies, &pb.Table)
	decoder.Decode(lsData.Steps)
	return pb
}

// Event_NewMessage is emitted whenever the client receives an LSInsertMessage dependency in Event_PublishResponse's table field.
type Event_NewMessage struct {
	Message []table.LSInsertMessage // slice of LSInsertMessage because there could potentially be more?
	MciTraceLog []table.LSMciTraceLog
	ExecuteFirstBlockForSyncTransaction []table.LSExecuteFirstBlockForSyncTransaction
	CheckAuthoritativeMessageExists []table.LSCheckAuthoritativeMessageExists
	VerifyThreadExists []table.LSVerifyThreadExists
	BumpThread []table.LSBumpThread
	MoveThreadToInboxAndUpdateParent []table.LSMoveThreadToInboxAndUpdateParent
	WriteThreadCapabilities []table.LSWriteThreadCapabilities
	UpdateThreadSnippet []table.LSUpdateThreadSnippet
	UpdateReadReceipt []table.LSUpdateReadReceipt
	SetForwardScore []table.LSSetForwardScore
	SetMessageDisplayedContentTypes []table.LSSetMessageDisplayedContentTypes
	UpdateParticipantLastMessageSendTimestamp []table.LSUpdateParticipantLastMessageSendTimestamp
	UpdateTypingIndicator []table.LSUpdateTypingIndicator
	UpsertSequenceId []table.LSUpsertSequenceId
	ExecuteFinallyBlockForSyncTransaction []table.LSExecuteFinallyBlockForSyncTransaction
}

func (nm *Event_NewMessage) Finish(table table.LSTable) *Event_NewMessage {
	nm.Message = table.LSInsertMessage
	nm.MciTraceLog = table.LSMciTraceLog
	nm.ExecuteFirstBlockForSyncTransaction = table.LSExecuteFirstBlockForSyncTransaction
	nm.CheckAuthoritativeMessageExists = table.LSCheckAuthoritativeMessageExists
	nm.VerifyThreadExists = table.LSVerifyThreadExists
	nm.BumpThread = table.LSBumpThread
	nm.MoveThreadToInboxAndUpdateParent = table.LSMoveThreadToInboxAndUpdateParent
	nm.WriteThreadCapabilities = table.LSWriteThreadCapabilities
	nm.UpdateThreadSnippet = table.LSUpdateThreadSnippet
	nm.UpdateReadReceipt = table.LSUpdateReadReceipt
	nm.SetForwardScore = table.LSSetForwardScore
	nm.SetMessageDisplayedContentTypes = table.LSSetMessageDisplayedContentTypes
	nm.UpdateParticipantLastMessageSendTimestamp = table.LSUpdateParticipantLastMessageSendTimestamp
	nm.UpdateTypingIndicator = table.LSUpdateTypingIndicator
	nm.UpsertSequenceId = table.LSUpsertSequenceId
	nm.ExecuteFinallyBlockForSyncTransaction = table.LSExecuteFinallyBlockForSyncTransaction
	return nm
}