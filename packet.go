package messagix

import "github.com/0xzer/messagix/packets"

type PacketHandler struct {
	packetIdHandler map[uint16]*packets.PublishACK
}

func (s *Socket) storePacketDetails(packetId uint16, packet *packets.PublishACK) {
    s.packetIdHandler[packetId] = packet
}

func (s *Socket) getPacketDetails(packetId uint16) *packets.PublishACK {
    return s.packetIdHandler[packetId]
}

func (s *Socket) deletePacketDetails(packetId uint16) {
    delete(s.packetIdHandler, packetId)
}