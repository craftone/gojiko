package pco

import (
	"fmt"
	"net"
)

type DnsServerV4 struct {
	header
	value net.IP
}

func NewDnsServerV4(ipv4 net.IP) *DnsServerV4 {
	length := 0
	if ipv4 != nil {
		length = 4
	}
	return &DnsServerV4{
		header{dnsServerV4Num, byte(length)},
		ipv4.To4(),
	}
}

func (d *DnsServerV4) marshal() []byte {
	var body []byte
	if d.value != nil {
		body = d.value.To4()
	}
	return d.header.marshal(body)
}

func unmarshalDnsServerV4(buf []byte) (*DnsServerV4, error) {
	if len(buf) == 0 {
		return NewDnsServerV4(nil), nil
	}

	if len(buf) == 4 {
		return NewDnsServerV4(buf[0:4]), nil
	}
	return nil, fmt.Errorf("It should be 4 octets binary : %v", buf)
}
