package messagix

import (
	"github.com/0xzer/messagix/byter"
	"github.com/0xzer/messagix/packets"
)

type Payload interface {
	Write() ([]byte, error)
}

type ConnectPayload struct {
	ProtocolName string `lengthType:"uint16"`
	ProtocolLevel uint8
	ConnectFlags uint8
	KeepAliveTime uint16
	ClientId string `lengthType:"uint16"`
	JSONData string `lengthType:"uint16"`
}

func (cp *ConnectPayload) Write() ([]byte, error) {
	return byter.NewWriter().WriteFromStruct(cp)
}

func (c *Client) NewConnectRequest(jsonData string, connectFlags uint8) ([]byte, error) {
	payload := &ConnectPayload{
		ProtocolName: c.configs.mqttConfig.ProtocolName,
		ProtocolLevel: c.configs.mqttConfig.ProtocolLevel,
		ConnectFlags: connectFlags,
		KeepAliveTime: c.configs.mqttConfig.KeepAliveTimeout,
		ClientId: c.configs.mqttConfig.ClientId,
		JSONData: jsonData,
	}

	packet := &packets.ConnectPacket{}
	request := &Request{
		PacketByte: packet.Compress(),
	}
	return request.Write(payload)
}

type PublishPayload struct {
	Topic Topic `lengthType:"uint16"`
	PacketId uint16
	JSONData string `lengthType:"uint16"`
}


func (pb *PublishPayload) Write() ([]byte, error) {
	return byter.NewWriter().WriteFromStruct(pb)
}

func (c *Client) NewPublishRequest(topic Topic, jsonData string, packetId uint16, packetByte byte) ([]byte, error) {
	payload := &PublishPayload{
		Topic: topic,
		PacketId: packetId,
		JSONData: jsonData,
	}
	
	request := &Request{
		PacketByte: packetByte,
	}
	return request.Write(payload)
}