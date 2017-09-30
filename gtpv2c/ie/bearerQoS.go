package ie

import "log"
import "errors"
import "fmt"
import "encoding/binary"

type BearerQoS struct {
	header
	BearerQoSArg
}

type BearerQoSArg struct {
	Pci         bool
	Pl          byte
	Pvi         bool
	Label       byte
	UplinkMBR   uint64
	DownlinkMBR uint64
	UplinkGBR   uint64
	DownlinkGBR uint64
}

const MaxBitrate = 10 * 1000 * 1000 * 1000

func NewBearerQoS(instance byte, bearerQoSArg BearerQoSArg) (*BearerQoS, error) {
	if bearerQoSArg.Pl > 0xf {
		return nil, fmt.Errorf("Invalid Bearer QoS's PL : %v", bearerQoSArg.Pl)
	}
	if bearerQoSArg.UplinkMBR > MaxBitrate ||
		bearerQoSArg.DownlinkMBR > MaxBitrate ||
		bearerQoSArg.UplinkGBR > MaxBitrate ||
		bearerQoSArg.DownlinkGBR > MaxBitrate {
		return nil, fmt.Errorf("A bitrate should be 0~%d", MaxBitrate)
	}

	header, err := newHeader(bearerQoSNum, 22, instance)
	if err != nil {
		return nil, err
	}
	return &BearerQoS{
		header,
		bearerQoSArg,
	}, nil
}

func putBitrate(body []byte, bitrate uint64) {
	body[0] = byte(bitrate >> 32)
	binary.BigEndian.PutUint32(body[1:5], uint32(bitrate))
}

func (p *BearerQoS) Marshal() []byte {
	body := make([]byte, 22)
	body[0] = setBit(body[0], 6, p.Pci)
	body[0] += (p.Pl << 2)
	body[0] = setBit(body[0], 0, p.Pvi)
	body[1] = p.Label

	putBitrate(body[2:7], p.UplinkMBR)
	putBitrate(body[7:12], p.DownlinkMBR)
	putBitrate(body[12:17], p.UplinkMBR)
	putBitrate(body[17:22], p.DownlinkMBR)
	return p.header.marshal(body)
}

func getBitrate(body []byte) uint64 {
	var res uint64
	res = uint64(body[0]) << 32
	res += uint64(binary.BigEndian.Uint32(body[1:5]))
	return res
}

func unmarshalBearerQoS(h header, buf []byte) (*BearerQoS, error) {
	if h.typeNum != bearerQoSNum {
		log.Fatal("Invalud type")
	}

	if len(buf) != 22 {
		return nil, errors.New("Invalid binary length")
	}

	pci := getBit(buf[0], 6)
	pl := (buf[0] >> 2) & 0xf
	pvi := getBit(buf[0], 0)
	label := buf[1]
	uplinkMBR := getBitrate(buf[2:7])
	downlinkMBR := getBitrate(buf[7:12])
	uplinkGBR := getBitrate(buf[12:17])
	downlinkGBR := getBitrate(buf[17:22])
	bearerQoSArg := BearerQoSArg{
		pci,
		pl,
		pvi,
		label,
		uplinkMBR,
		downlinkMBR,
		uplinkGBR,
		downlinkGBR,
	}
	bq, err := NewBearerQoS(h.instance, bearerQoSArg)
	if err != nil {
		return nil, err
	}
	return bq, nil
}
