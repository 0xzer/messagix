package messagix

import (
	"errors"
	"time"
)

type PacketHandler struct {
    packetsSent uint16
	packetChannels map[uint16]chan interface{}

    packetTimeout time.Duration
}

func (p *PacketHandler) SetPacketTimeout(d time.Duration) {
    p.packetTimeout = d
}

func (p *PacketHandler) addPacketChannel(packetId uint16) {
    p.packetChannels[packetId] = make(chan interface{}, 1) // buffered channel with capacity of 1
}

func (p *PacketHandler) updatePacketChannel(packetId uint16, ackData interface{}) error {
    if ch, ok := p.packetChannels[packetId]; ok {
        ch <- ackData
        return nil
    }

    return errors.New("failed to update packet channel for packetId %d because there was no packet channel found")
}

func (p *PacketHandler) waitForPubACKDetails(packetId uint16) *Event_PublishACK {
    if ch, ok := p.packetChannels[packetId]; ok {
        select {
        case ackEvent := <-ch:
            p.deletePacketDetails(packetId)
            data, ok := ackEvent.(*Event_PublishACK)
            if !ok {
                return nil
            }
            return data
        case <-time.After(p.packetTimeout):
            p.deletePacketDetails(packetId)
            return nil
        }
    }
    return nil
}

func (p *PacketHandler) waitForSubACKDetails(packetId uint16) *Event_SubscribeACK {
    if ch, ok := p.packetChannels[packetId]; ok {
        select {
        case ackEvent := <-ch:
            p.deletePacketDetails(packetId)
            data, ok := ackEvent.(*Event_SubscribeACK)
            if !ok {
                return nil
            }
            return data
        case <-time.After(p.packetTimeout):
            p.deletePacketDetails(packetId)
            return nil
        }
    }
    return nil
}

func (p *PacketHandler) deletePacketDetails(packetId uint16) {
    if ch, ok := p.packetChannels[packetId]; ok {
        close(ch)
        delete(p.packetChannels, packetId)
    }
}