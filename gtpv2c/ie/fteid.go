package ie

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

type Fteid struct {
	header
	IfType IfType
	Ipv4   net.IP
	Ipv6   net.IP
	Value  uint32
}

type IfType byte

const (
	S5S8SgwGtpUIf IfType = 4
	S5S8PgwGtpUIf IfType = 5
	S5S8SgwGtpCIf IfType = 6
	S5S8PgwGtpCIf IfType = 7
)

func NewFteid(instance byte, ipv4, ipv6 net.IP, ifType IfType, value uint32) (*Fteid, error) {
	if ipv4 != nil {
		ipv4 = ipv4.To4()
	}
	if ipv6 != nil {
		ipv6 = ipv6.To16()
	}
	if ipv4 == nil && ipv6 == nil {
		return nil, fmt.Errorf("At least one of V4 or V6 should be set")
	}

	if ifType > 0x3f {
		return nil, fmt.Errorf("Invalud Interface Type")
	}

	length := 5
	if ipv4 != nil {
		length += 4
	}
	if ipv6 != nil {
		length += 16
	}
	header, err := newHeader(fteidNum, uint16(length), instance)
	if err != nil {
		return nil, err
	}
	return &Fteid{
		header, ifType,
		ipv4, ipv6,
		value,
	}, nil
}

func (f *Fteid) Marshal() []byte {
	body := make([]byte, f.length)
	body[0] = setBit(body[0], 7, f.Ipv4 != nil)
	body[0] = setBit(body[0], 6, f.Ipv6 != nil)
	body[0] += byte(f.IfType) & 0x3f
	binary.BigEndian.PutUint32(body[1:5], f.Value)
	offset := 5
	if f.Ipv4 != nil {
		copy(body[5:9], f.Ipv4)
		offset += 4
	}
	if f.Ipv6 != nil {
		copy(body[offset:offset+16], f.Ipv6)
	}
	return f.header.marshal(body)
}

func unmarshalFteid(h header, buf []byte) (*Fteid, error) {
	if h.typeNum != fteidNum {
		log.Fatal("Invalud type")
	}

	if len(buf) < 6 {
		return nil, errors.New("Invalid binary")
	}

	v4flag := getBit(buf[0], 7)
	v6flag := getBit(buf[0], 6)
	ifType := IfType(buf[0] & 0x3f)

	length := 5
	if v4flag {
		length += 4
	}
	if v6flag {
		length += 16
	}
	if len(buf) < length {
		return nil, errors.New("Invalid binary")
	}

	value := binary.BigEndian.Uint32(buf[1:5])

	var ipv4, ipv6 net.IP
	offset := 5
	if v4flag {
		ipv4 = buf[offset : offset+4]
		offset += 4
	}
	if v6flag {
		ipv6 = buf[offset : offset+16]
	}

	fteid, err := NewFteid(h.instance, ipv4, ipv6, ifType, value)
	if err != nil {
		return nil, err
	}
	return fteid, nil
}
