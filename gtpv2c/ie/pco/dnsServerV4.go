package pco

import (
	"fmt"
	"net"
)

type DNSServerV4 struct {
	header
	value net.IP
}

func NewDNSServerV4(ipv4 net.IP) *DNSServerV4 {
	length := 0
	if ipv4 != nil {
		length = 4
	}
	return &DNSServerV4{
		header{dnsServerV4Num, byte(length)},
		ipv4.To4(),
	}
}

func (d *DNSServerV4) marshal() []byte {
	var body []byte
	if d.value != nil {
		body = d.value.To4()
	}
	return d.header.marshal(body)
}

func unmarshalDNSServerV4(buf []byte) (*DNSServerV4, error) {
	if len(buf) == 0 {
		return NewDNSServerV4(nil), nil
	}

	if len(buf) == 4 {
		return NewDNSServerV4(buf[0:4]), nil
	}
	return nil, fmt.Errorf("It should be 4 octets binary : %v", buf)
}
