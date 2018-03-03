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
	configProto byte
	ipcp        *Ipcp
}

type MsToNetwork struct {
	pco
	dnsServerV4Req   bool
	dnsServerV6Req   bool
	ipAllocViaNasSig bool
}

type NetworkToMs struct {
	pco
	dnsServerV4s []*DNSServerV4
	dnsServerV6s []*DNSServerV6
}

func NewMsToNetwork(ipcp *Ipcp, dnsServerV4Req, dnsServerV6Req, ipAllocViaNasSig bool) *MsToNetwork {
	return &MsToNetwork{
		pco{ipcp: ipcp},
		dnsServerV4Req,
		dnsServerV6Req,
		ipAllocViaNasSig,
	}
}

func NewNetworkToMs(ipcp *Ipcp, dnsServerV4s []*DNSServerV4, dnsServerV6s []*DNSServerV6) *NetworkToMs {
	return &NetworkToMs{
		pco{ipcp: ipcp},
		dnsServerV4s,
		dnsServerV6s,
	}
}

func (p MsToNetwork) Marshal() []byte {
	res := make([]byte, 0, 6)
	if p.dnsServerV4Req {
		b := NewDNSServerV4(nil).marshal()
		res = append(res, b...)
	}
	if p.dnsServerV6Req {
		b := NewDNSServerV6(nil).marshal()
		res = append(res, b...)
	}
	if p.ipAllocViaNasSig {
		b := NewIPAllocViaNasSignalling().marshal()
		res = append(res, b...)
	}
	return p.pco.marshal(res)
}

func (p NetworkToMs) Marshal() []byte {
	res := make([]byte, 0, 128)
	for _, dnsServerV4 := range p.dnsServerV4s {
		if dnsServerV4 != nil {
			b := dnsServerV4.marshal()
			res = append(res, b...)
		}
	}
	for _, dnsServerV6 := range p.dnsServerV6s {
		if dnsServerV6 != nil {
			b := dnsServerV6.marshal()
			res = append(res, b...)
		}
	}
	return p.pco.marshal(res)
}

func (p pco) marshal(body []byte) []byte {
	res := make([]byte, 1, 1+16+len(body))
	res[0] = byte(0x80 + p.configProto)
	if p.ipcp != nil {
		ipcpBin := p.ipcp.marshal()
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
			res.ipcp = ipcp
		case dnsServerV4Num:
			res.dnsServerV4Req = true
		case dnsServerV6Num:
			res.dnsServerV6Req = true
		case ipAllocViaNasSigNum:
			res.ipAllocViaNasSig = true
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
			res.ipcp = ipcp
		case dnsServerV4Num:
			dnsServerV4, err := unmarshalDNSServerV4(body)
			if err != nil {
				return err
			}
			res.dnsServerV4s = append(res.dnsServerV4s, dnsServerV4)
		case dnsServerV6Num:
			dnsServerV6, err := unmarshalDNSServerV6(body)
			if err != nil {
				return err
			}
			res.dnsServerV6s = append(res.dnsServerV6s, dnsServerV6)
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
	res := pco{configProto: head & 0x7}

	return res, buf[1:], nil
}

func (p pco) Ipcp() *Ipcp {
	return p.ipcp
}

func (p pco) ConfigProto() byte {
	return p.configProto
}

func (p MsToNetwork) DNSServerV4Req() bool {
	return p.dnsServerV4Req
}

func (p MsToNetwork) DNSServerV6Req() bool {
	return p.dnsServerV6Req
}

func (p MsToNetwork) IPAllocViaNasSig() bool {
	return p.ipAllocViaNasSig
}

func (p NetworkToMs) DNSServerV4s() []*DNSServerV4 {
	return p.dnsServerV4s
}

func (p NetworkToMs) DNSServerV6s() []*DNSServerV6 {
	return p.dnsServerV6s
}
