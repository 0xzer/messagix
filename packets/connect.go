package packets

type ConnectPacket struct {
	Packet byte
}

func (p *ConnectPacket) Compress() byte {
	return CONNECT << 4
}

func (p *ConnectPacket) Decompress(packetByte byte) error {
	p.Packet = packetByte
	return nil
}

type ConnACKPacket struct {
	Packet byte
}

func (p *ConnACKPacket) Compress() byte {
	return CONNACK << 4
}

func (p *ConnACKPacket) Decompress(packetByte byte) error {
	p.Packet = packetByte
	return nil
}