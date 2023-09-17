package messagix

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/0xzer/messagix/byter"
	"github.com/0xzer/messagix/packets"
)

type ResponseData interface {
	Finish() ResponseData
	SetIdentifier(identifier int16)
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
	qosLevel := (r.PacketByte >> 1) & 0x03

	responseHandlerFunc, ok := responseMap[packetType]
	if !ok {
		return fmt.Errorf("could not find response func handler for packet type %d", packetType)
	}
	r.ResponseData = responseHandlerFunc()

	if packetType == packets.PUBLISH && qosLevel == 1 {
		identifier, err := reader.ReadInteger(reflect.Uint16, 2, "big")
		if err != nil {
			log.Fatalf("failed to read int16 message identifier from publish response packet with qos level 1: %e", err)
		}
		log.Println("got qos_level 1:", r.PacketByte)
		log.Println("got message identifier:", identifier)
		log.Println("data left:", reader.Buff.Len())
		os.Exit(1)
	}

	offset := len(data) - reader.Buff.Len()
	remainingData := data[offset:]
	return byter.NewReader(remainingData).ReadToStruct(r.ResponseData)
}