package pco

import (
	"fmt"
	"net"
)

type DnsServerV6 struct {
	header
	value net.IP
}

func NewDnsServerV6(ipv6 net.IP) *DnsServerV6 {
	length := 0
	if ipv6 != nil {
		length = 4
	}
	return &DnsServerV6{
		header{dnsServerV6Num, byte(length)},
		ipv6.To16(),
	}
}

func (d *DnsServerV6) marshal() []byte {
	var body []byte
	if d.value != nil {
		body = d.value.To16()
	}
	return d.header.marshal(body)
}

func unmarshalDnsServerV6(buf []byte) (*DnsServerV6, error) {
	if len(buf) == 0 {
		return NewDnsServerV6(nil), nil
	}

	if len(buf) == 16 {
		return NewDnsServerV6(buf[0:16]), nil
	}
	return nil, fmt.Errorf("It should be 16 octets binary : %v", buf)
}
