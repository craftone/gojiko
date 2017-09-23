package pco

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type pcoTypeNum uint16

const (
	dnsServerV6Num pcoTypeNum = 0x0003
	dnsServerV4Num pcoTypeNum = 0x000D
	ipcpNum        pcoTypeNum = 0x8021
	lcpNum         pcoTypeNum = 0xC021
	papNum         pcoTypeNum = 0xC023
	chapNum        pcoTypeNum = 0xC223
)

type Pco struct {
	ConfigProto  byte
	DnsServerV4s []*DnsServerV4
	DnsServerV6s []*DnsServerV6
}

type header struct {
	typeNum pcoTypeNum
	length  byte
}

func newHeader(typeNum pcoTypeNum, length byte) header {
	return header{typeNum, length}
}

func (h *header) marshal(body []byte) []byte {
	bodyLen := byte(len(body))
	res := make([]byte, 3+bodyLen)
	// make header
	binary.BigEndian.PutUint16(res[0:2], uint16(h.typeNum))
	res[2] = bodyLen
	// copy body
	copy(res[3:], body)
	return res
}

func (p Pco) Marshal() []byte {
	res := make([]byte, 1, 255)
	res[0] = byte(0x80 + p.ConfigProto)
	for _, dnsServerV4 := range p.DnsServerV4s {
		if dnsServerV4 != nil {
			b := dnsServerV4.marshal()
			res = append(res, b...)
		}
	}
	for _, dnsServerV6 := range p.DnsServerV6s {
		if dnsServerV6 != nil {
			b := dnsServerV6.marshal()
			res = append(res, b...)
		}
	}

	return res
}

func Unmarshal(buf []byte) (Pco, []byte, error) {
	if len(buf) < 4 {
		return Pco{}, buf, errors.New("It needs at least 4 bytes")
	}

	head := buf[0]
	if head&0x80 == 0 {
		return Pco{}, buf, errors.New("MSB of the first octet must be 1")
	}
	res := Pco{ConfigProto: head & 0x7}

	offset := 1
	for len(buf) > offset {
		if len(buf)-offset < 3 {
			return Pco{}, buf, fmt.Errorf("Invalid binary %v", buf)
		}
		h := header{
			typeNum: pcoTypeNum(binary.BigEndian.Uint16(buf[offset : offset+2])),
			length:  buf[offset+2],
		}
		offset += 3
		if len(buf)-offset < int(h.length) {
			return Pco{}, buf, fmt.Errorf("Invalid binary %v", buf)
		}
		body := buf[offset : offset+int(h.length)]
		offset += int(h.length)

		switch h.typeNum {
		case dnsServerV4Num:
			dnsServerV4, err := unmarshalDnsServerV4(h, body)
			if err != nil {
				return Pco{}, buf, fmt.Errorf("Invalid DnsServerV4 binary : %v", err)
			}
			res.DnsServerV4s = append(res.DnsServerV4s, dnsServerV4)
		case dnsServerV6Num:
			dnsServerV6, err := unmarshalDnsServerV6(h, body)
			if err != nil {
				return Pco{}, buf, fmt.Errorf("Invalid DnsServerV6 binary : %v", err)
			}
			res.DnsServerV6s = append(res.DnsServerV6s, dnsServerV6)
		}
	}
	return res, buf[offset:], nil
}
