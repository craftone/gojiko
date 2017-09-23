package ie

import "log"
import "errors"
import "encoding/binary"

type Ambr struct {
	header
	UplinkKbps   uint32
	DownlinkKbps uint32
}

func NewAmbr(instance byte, uplinkKbps, downlinkKbps uint32) (*Ambr, error) {
	header, err := newHeader(ambrNum, 8, instance)
	if err != nil {
		return nil, err
	}
	return &Ambr{
		header:       header,
		UplinkKbps:   uplinkKbps,
		DownlinkKbps: downlinkKbps,
	}, nil
}

func (a *Ambr) Marshal() []byte {
	body := make([]byte, 8)
	binary.BigEndian.PutUint32(body[0:4], a.UplinkKbps)
	binary.BigEndian.PutUint32(body[4:8], a.DownlinkKbps)
	return a.header.marshal(body)
}

func unmarshalAmbr(h header, buf []byte) (*Ambr, error) {
	if h.typeNum != ambrNum {
		log.Fatal("Invalud type")
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
