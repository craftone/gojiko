package ie

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type Paa struct {
	header
	value PdnTypeValue
	ipv4  net.IP
	ipv6  net.IP
}

func NewPaa(instance byte, value PdnTypeValue, ipv4, ipv6 net.IP) (*Paa, error) {
	if !value.isValid() {
		return nil, fmt.Errorf("Invalid PDN Type")
	}

	length := 1
	if value == PdnTypeIPv4 {
		if ipv4 == nil || ipv6 != nil {
			return nil, fmt.Errorf("PdnType and IPv4/v6 don't match")
		}
		ipv4 = ipv4.To4()
		length += 4
	} else if value == PdnTypeIPv6 {
		if ipv4 != nil || ipv6 == nil {
			return nil, fmt.Errorf("PdnType and IPv4/v6 don't match")
		}
		ipv6 = ipv6.To16()
		length += 16
	} else { // IPv4v6
		if ipv4 == nil || ipv6 == nil {
			return nil, fmt.Errorf("PdnType and IPv4/v6 don't match")
		}
		ipv4 = ipv4.To4()
		ipv6 = ipv6.To16()
		length += 20
	}

	header, err := newHeader(paaNum, uint16(length), instance)
	if err != nil {
		return nil, err
	}
	return &Paa{
		header,
		value,
		ipv4,
		ipv6,
	}, nil
}

func (p *Paa) Marshal() []byte {
	body := make([]byte, p.length)
	body[0] = byte(p.value)
	offset := 1
	if p.value == PdnTypeIPv6 || p.value == PdnTypeIPv4v6 {
		copy(body[offset:offset+16], p.ipv6)
		offset += 16
	}
	if p.value == PdnTypeIPv4 || p.value == PdnTypeIPv4v6 {
		copy(body[offset:offset+4], p.ipv4)
	}
	return p.header.marshal(body)
}

func unmarshalPaa(h header, buf []byte) (*Paa, error) {
	if h.typeNum != paaNum {
		log.Fatal("Invalud type")
	}

	if len(buf) < 5 {
		return nil, errors.New("Invalid binary")
	}

	pdnType := PdnTypeValue(buf[0] & 0x7)
	if !pdnType.isValid() {
		return nil, errors.New("Invalid PDN Type")
	}

	offset := 1
	var ipv4, ipv6 net.IP
	if pdnType == PdnTypeIPv6 || pdnType == PdnTypeIPv4v6 {
		ipv6 = net.IP(buf[offset : offset+16])
		ipv6 = ipv6.To16()
		offset += 16
	}
	if pdnType == PdnTypeIPv4 || pdnType == PdnTypeIPv4v6 {
		ipv4 = net.IP(buf[offset : offset+4])
		ipv4 = ipv4.To4()
	}

	paa, err := NewPaa(h.instance, pdnType, ipv4, ipv6)
	if err != nil {
		return nil, err
	}
	return paa, nil
}

func (p *Paa) Value() PdnTypeValue {
	return p.value
}

func (p *Paa) IPv4() net.IP {
	return p.ipv4
}

func (p *Paa) IPv6() net.IP {
	return p.ipv6
}
