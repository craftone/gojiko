package pco

import (
	"fmt"
	"log"
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
	if d.value == nil {
		return []byte{}
	}
	return d.header.marshal(d.value.To4())
}

func unmarshalDnsServerV4(h header, buf []byte) (*DnsServerV4, error) {
	if h.typeNum != dnsServerV4Num {
		log.Fatal("Invalud type")
	}

	if len(buf) == 0 {
		return NewDnsServerV4(nil), nil
	}

	if len(buf) == 4 {
		return NewDnsServerV4(buf[0:4]), nil
	}
	return nil, fmt.Errorf("It should be 4 octets binary : %v", buf)
}
