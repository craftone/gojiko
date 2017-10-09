package ie

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/craftone/gojiko/gtp"
)

type Fteid struct {
	header
	ifType IfType
	ipv4   net.IP
	ipv6   net.IP
	teid   gtp.Teid
}

type IfType byte

const (
	S5S8SgwGtpUIf IfType = 4
	S5S8PgwGtpUIf IfType = 5
	S5S8SgwGtpCIf IfType = 6
	S5S8PgwGtpCIf IfType = 7
)

func NewFteid(instance byte, ipv4, ipv6 net.IP, ifType IfType, teid gtp.Teid) (*Fteid, error) {
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
		teid,
	}, nil
}

func (f *Fteid) Marshal() []byte {
	body := make([]byte, f.length)
	body[0] = setBit(body[0], 7, f.ipv4 != nil)
	body[0] = setBit(body[0], 6, f.ipv6 != nil)
	body[0] += byte(f.ifType) & 0x3f
	binary.BigEndian.PutUint32(body[1:5], uint32(f.teid))
	offset := 5
	if f.ipv4 != nil {
		copy(body[5:9], f.ipv4)
		offset += 4
	}
	if f.ipv6 != nil {
		copy(body[offset:offset+16], f.ipv6)
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

	teid := gtp.Teid(binary.BigEndian.Uint32(buf[1:5]))

	var ipv4, ipv6 net.IP
	offset := 5
	if v4flag {
		ipv4 = buf[offset : offset+4]
		offset += 4
	}
	if v6flag {
		ipv6 = buf[offset : offset+16]
	}

	fteid, err := NewFteid(h.instance, ipv4, ipv6, ifType, teid)
	if err != nil {
		return nil, err
	}
	return fteid, nil
}

func (f *Fteid) IfType() IfType {
	return f.ifType
}
func (f *Fteid) Ipv4() net.IP {
	return f.ipv4
}
func (f *Fteid) Ipv6() net.IP {
	return f.ipv6
}
func (f *Fteid) Teid() gtp.Teid {
	return f.teid
}
