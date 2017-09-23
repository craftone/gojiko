package pco

import (
	"fmt"
	"log"
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
	if d.value == nil {
		return []byte{}
	}
	return d.header.marshal(d.value.To16())
}

func unmarshalDnsServerV6(h header, buf []byte) (*DnsServerV6, error) {
	if h.typeNum != dnsServerV6Num {
		log.Fatal("Invalud type")
	}

	if len(buf) == 0 {
		return NewDnsServerV6(nil), nil
	}

	if len(buf) == 16 {
		return NewDnsServerV6(buf[0:16]), nil
	}
	return nil, fmt.Errorf("It should be 16 octets binary : %v", buf)
}
