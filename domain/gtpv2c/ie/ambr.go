package ie

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Ambr struct {
	header
	uplinkKbps   uint32
	downlinkKbps uint32
}

func NewAmbr(instance byte, uplinkKbps, downlinkKbps uint32) (*Ambr, error) {
	header, err := newHeader(ambrNum, 8, instance)
	if err != nil {
		return nil, err
	}
	return &Ambr{
		header:       header,
		uplinkKbps:   uplinkKbps,
		downlinkKbps: downlinkKbps,
	}, nil
}

func (a *Ambr) Marshal() []byte {
	body := make([]byte, 8)
	binary.BigEndian.PutUint32(body[0:4], a.uplinkKbps)
	binary.BigEndian.PutUint32(body[4:8], a.downlinkKbps)
	return a.header.marshal(body)
}

func unmarshalAmbr(h header, buf []byte) (*Ambr, error) {
	if h.typeNum != ambrNum {
		log.Panic("Invalid type")
	}

	if len(buf) != 8 {
		return nil, errors.New("invalid length")
	}

	uplinkKbps := binary.BigEndian.Uint32(buf[0:4])
	downlinkKbps := binary.BigEndian.Uint32(buf[4:8])
	ambr, err := NewAmbr(h.instance, uplinkKbps, downlinkKbps)
	if err != nil {
		return nil, err
	}
	return ambr, nil
}

func (a *Ambr) UplinkKbps() uint32 {
	return a.uplinkKbps
}

func (a *Ambr) DownlinkKbps() uint32 {
	return a.downlinkKbps
}

func (a *Ambr) String() string {
	return fmt.Sprintf("Uplink AMBR: %d kbps, Downlink AMBR: %d kbps", a.uplinkKbps, a.downlinkKbps)
}
