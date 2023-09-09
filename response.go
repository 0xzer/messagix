package messagix

import (
	"fmt"

	"github.com/0xzer/messagix/byter"
	"github.com/0xzer/messagix/packets"
)

type ResponseData interface {
	Finish() ResponseData
}
type responseHandler func() (ResponseData)
var responseMap = map[uint8]responseHandler{
	packets.CONNACK: func() ResponseData {return &Event_Ready{}},
	packets.PUBACK: func() ResponseData {return &Event_PublishACK{}},
	packets.SUBACK: func() ResponseData {return &Event_SubscribeACK{}},
	packets.PUBLISH: func() ResponseData {return &Event_PublishResponse{}},
}

type Response struct {
	PacketByte uint8
    RemainingLength uint32 `vlq:"true"`
	ResponseData ResponseData
}

func (r *Response) Read(data []byte) error {
	reader := byter.NewReader(data)
	err := reader.ReadToStruct(r)
	if err != nil {
		return err
	}

	packetType := r.PacketByte >> 4 // parse the packet type by the leftmost 4 bits
	responseHandlerFunc, ok := responseMap[packetType]
	if !ok {
		return fmt.Errorf("could not find response func handler for packet type %d", packetType)
	}
	r.ResponseData = responseHandlerFunc()

	offset := len(data) - reader.Buff.Len()
	remainingData := data[offset:]
	return byter.NewReader(remainingData).ReadToStruct(r.ResponseData)
}