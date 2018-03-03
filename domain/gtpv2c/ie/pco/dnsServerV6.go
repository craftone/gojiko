package pco

import (
	"fmt"
	"net"
)

type DNSServerV6 struct {
	header
	value net.IP
}

func NewDNSServerV6(ipv6 net.IP) *DNSServerV6 {
	length := 0
	if ipv6 != nil {
		length = 4
	}
	return &DNSServerV6{
		header{dnsServerV6Num, byte(length)},
		ipv6.To16(),
	}
}

func (d *DNSServerV6) marshal() []byte {
	var body []byte
	if d.value != nil {
		body = d.value.To16()
	}
	return d.header.marshal(body)
}

func unmarshalDNSServerV6(buf []byte) (*DNSServerV6, error) {
	if len(buf) == 0 {
		return NewDNSServerV6(nil), nil
	}

	if len(buf) == 16 {
		return NewDNSServerV6(buf[0:16]), nil
	}
	return nil, fmt.Errorf("It should be 16 octets binary : %v", buf)
}

func (d *DNSServerV6) Value() net.IP {
	return d.value
}
