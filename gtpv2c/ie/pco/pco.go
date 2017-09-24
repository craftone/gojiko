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

type pco struct {
	ConfigProto byte
}

type PcoMsToNetwork struct {
	pco
	DnsServerV4Req bool
	DnsServerV6Req bool
}

type PcoNetworkToMs struct {
	pco
	DnsServerV4s []*DnsServerV4
	DnsServerV6s []*DnsServerV6
}

func (p PcoMsToNetwork) Marshal() []byte {
	res := make([]byte, 0, 6)
	if p.DnsServerV4Req {
		b := NewDnsServerV4(nil).marshal()
		res = append(res, b...)
	}
	if p.DnsServerV6Req {
		b := NewDnsServerV6(nil).marshal()
		res = append(res, b...)
	}
	return p.pco.marshal(res)
}

func (p PcoNetworkToMs) Marshal() []byte {
	res := make([]byte, 0, 128)
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
	return p.pco.marshal(res)
}

func (p pco) marshal(body []byte) []byte {
	res := make([]byte, 1+len(body))
	res[0] = byte(0x80 + p.ConfigProto)
	// copy body
	copy(res[1:], body)
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

func UnmarshalMsToNetowrk(buf []byte) (PcoMsToNetwork, []byte, error) {
	pco, tail, err := unmarshal(buf)
	if err != nil {
		return PcoMsToNetwork{}, buf, err
	}

	res := PcoMsToNetwork{pco: pco}
	err = unmarshalContainer(tail, func(h header, body []byte) error {
		switch h.typeNum {
		case dnsServerV4Num:
			res.DnsServerV4Req = true
		case dnsServerV6Num:
			res.DnsServerV6Req = true
		}
		return nil
	})
	if err != nil {
		return res, buf, err
	}
	return res, []byte{}, nil
}

func UnmarshalNetowrkToMs(buf []byte) (PcoNetworkToMs, []byte, error) {
	pco, tail, err := unmarshal(buf)
	if err != nil {
		return PcoNetworkToMs{}, buf, err
	}

	res := PcoNetworkToMs{pco: pco}
	unmarshalContainer(tail, func(h header, body []byte) error {
		switch h.typeNum {
		case dnsServerV4Num:
			dnsServerV4, err := unmarshalDnsServerV4(body)
			if err != nil {
				return err
			}
			res.DnsServerV4s = append(res.DnsServerV4s, dnsServerV4)
		case dnsServerV6Num:
			dnsServerV6, err := unmarshalDnsServerV6(body)
			if err != nil {
				return err
			}
			res.DnsServerV6s = append(res.DnsServerV6s, dnsServerV6)
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
