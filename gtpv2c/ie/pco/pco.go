package pco

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type pcoTypeNum uint16

const (
	dnsServerV6Num      pcoTypeNum = 0x0003
	ipAllocViaNasSigNum pcoTypeNum = 0x000A
	dnsServerV4Num      pcoTypeNum = 0x000D
	ipcpNum             pcoTypeNum = 0x8021
	lcpNum              pcoTypeNum = 0xC021
	papNum              pcoTypeNum = 0xC023
	chapNum             pcoTypeNum = 0xC223
)

type pco struct {
	ConfigProto byte
	Ipcp        *Ipcp
}

type MsToNetwork struct {
	pco
	DNSServerV4Req   bool
	DNSServerV6Req   bool
	IPAllocViaNasSig bool
}

type NetworkToMs struct {
	pco
	DNSServerV4s []*DNSServerV4
	DNSServerV6s []*DNSServerV6
}

func NewMsToNetwork(ipcp *Ipcp, dnsServerV4Req, dnsServerV6Req, ipAllocViaNasSig bool) *MsToNetwork {
	return &MsToNetwork{
		pco{Ipcp: ipcp},
		dnsServerV4Req,
		dnsServerV6Req,
		ipAllocViaNasSig,
	}
}

func NewNetworkToMs(ipcp *Ipcp, dnsServerV4s []*DNSServerV4, dnsServerV6s []*DNSServerV6) *NetworkToMs {
	return &NetworkToMs{
		pco{Ipcp: ipcp},
		dnsServerV4s,
		dnsServerV6s,
	}
}

func (p MsToNetwork) Marshal() []byte {
	res := make([]byte, 0, 6)
	if p.DNSServerV4Req {
		b := NewDNSServerV4(nil).marshal()
		res = append(res, b...)
	}
	if p.DNSServerV6Req {
		b := NewDNSServerV6(nil).marshal()
		res = append(res, b...)
	}
	if p.IPAllocViaNasSig {
		b := NewIPAllocViaNasSignalling().marshal()
		res = append(res, b...)
	}
	return p.pco.marshal(res)
}

func (p NetworkToMs) Marshal() []byte {
	res := make([]byte, 0, 128)
	for _, dnsServerV4 := range p.DNSServerV4s {
		if dnsServerV4 != nil {
			b := dnsServerV4.marshal()
			res = append(res, b...)
		}
	}
	for _, dnsServerV6 := range p.DNSServerV6s {
		if dnsServerV6 != nil {
			b := dnsServerV6.marshal()
			res = append(res, b...)
		}
	}
	return p.pco.marshal(res)
}

func (p pco) marshal(body []byte) []byte {
	res := make([]byte, 1, 1+16+len(body))
	res[0] = byte(0x80 + p.ConfigProto)
	if p.Ipcp != nil {
		ipcpBin := p.Ipcp.marshal()
		res = append(res, ipcpBin...)
	}
	// copy body
	res = append(res, body...)
	return res
}

func unmarshalContainer(buf []byte, f func(header, []byte) error) error {
	for len(buf) > 0 {
		if len(buf) < 3 {
			return fmt.Errorf("Too short container : %v", buf)
		}
		h := header{
			typeNum: pcoTypeNum(binary.BigEndian.Uint16(buf[0:2])),
			length:  buf[2],
		}
		if len(buf)-3 < int(h.length) {
			return fmt.Errorf("Too short container body : %v", buf)
		}
		body := buf[3 : 3+int(h.length)]
		buf = buf[3+int(h.length):]

		err := f(h, body)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnmarshalMsToNetowrk(buf []byte) (*MsToNetwork, []byte, error) {
	pco, tail, err := unmarshal(buf)
	if err != nil {
		return &MsToNetwork{}, buf, err
	}

	res := &MsToNetwork{pco: pco}
	err = unmarshalContainer(tail, func(h header, body []byte) error {
		switch h.typeNum {
		case ipcpNum:
			ipcp, err := unmarshalIpcp(body)
			if err != nil {
				return err
			}
			res.Ipcp = ipcp
		case dnsServerV4Num:
			res.DNSServerV4Req = true
		case dnsServerV6Num:
			res.DNSServerV6Req = true
		case ipAllocViaNasSigNum:
			res.IPAllocViaNasSig = true
		}
		return nil
	})
	if err != nil {
		return res, buf, err
	}
	return res, []byte{}, nil
}

func UnmarshalNetowrkToMs(buf []byte) (*NetworkToMs, []byte, error) {
	pco, tail, err := unmarshal(buf)
	if err != nil {
		return &NetworkToMs{}, buf, err
	}

	res := &NetworkToMs{pco: pco}
	unmarshalContainer(tail, func(h header, body []byte) error {
		switch h.typeNum {
		case ipcpNum:
			ipcp, err := unmarshalIpcp(body)
			if err != nil {
				return err
			}
			res.Ipcp = ipcp
		case dnsServerV4Num:
			dnsServerV4, err := unmarshalDNSServerV4(body)
			if err != nil {
				return err
			}
			res.DNSServerV4s = append(res.DNSServerV4s, dnsServerV4)
		case dnsServerV6Num:
			dnsServerV6, err := unmarshalDNSServerV6(body)
			if err != nil {
				return err
			}
			res.DNSServerV6s = append(res.DNSServerV6s, dnsServerV6)
		}
		return nil
	})
	return res, []byte{}, nil
}

func unmarshal(buf []byte) (pco, []byte, error) {
	if len(buf) < 4 {
		return pco{}, buf, errors.New("It needs at least 4 bytes")
	}

	head := buf[0]
	if head&0x80 == 0 {
		return pco{}, buf, errors.New("MSB of the first octet must be 1")
	}
	res := pco{ConfigProto: head & 0x7}

	return res, buf[1:], nil
}
