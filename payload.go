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
	JSONData string `lengthType:""`
}


func (pb *PublishPayload) Write() ([]byte, error) {
	return byter.NewWriter().WriteFromStruct(pb)
}

func (c *Client) NewPublishRequest(topic Topic, jsonData string, packetByte byte, packetId uint16) ([]byte, uint16, error) {
	payload := &PublishPayload{
		Topic: topic,
		PacketId: packetId,
		JSONData: jsonData,
	}
	
	c.socket.responseHandler.addPacketChannel(packetId)
	c.socket.responseHandler.addRequestChannel(packetId)

	request := &Request{
		PacketByte: packetByte,
	}
	reqBytes, err := request.Write(payload)
	if err != nil {
		c.socket.responseHandler.deleteDetails(packetId, PacketChannel)
		c.socket.responseHandler.deleteDetails(packetId, RequestChannel)
		return nil, 0, err
	}
	return reqBytes, packetId, nil
}

type SubscribePayload struct {
	PacketId uint16
	Topic Topic `lengthType:"uint16"`
	QoSLevel packets.QoS
}

func (sb *SubscribePayload) Write() ([]byte, error) {
	return byter.NewWriter().WriteFromStruct(sb)
}

func (c *Client) NewSubscribeRequest(topic Topic, qos packets.QoS) ([]byte, uint16, error) {
	packetByte := &packets.SubscribePacket{}
	packetId := c.socket.SafePacketId()
	c.socket.responseHandler.addPacketChannel(packetId)
	payload := &SubscribePayload{
		PacketId: packetId,
		Topic: topic,
		QoSLevel: qos,
	}
	
	request := &Request{
		PacketByte: packetByte.Compress(),
	}
	reqBytes, err := request.Write(payload)
	if err != nil {
		c.socket.responseHandler.deleteDetails(packetId, PacketChannel)
		return nil, 0, err
	}

	return reqBytes, packetId, nil
}