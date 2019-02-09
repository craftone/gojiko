package pco

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Ipcp struct {
	header
	code       ipcpCode
	identifier byte
	priDNS     net.IP
	secDNS     net.IP
}

type ipcpCode byte
type ipcpOption byte

const (
	ConfigureRequest ipcpCode = 1
	ConfigureAck     ipcpCode = 2
	ConfigureNack    ipcpCode = 3
	ConfigureReject  ipcpCode = 4
	TerminalRequest  ipcpCode = 5
	TerminalAck      ipcpCode = 6
	CodeReject       ipcpCode = 7
)

const (
	primaryDNSServerAddress    ipcpOption = 129
	primaryNBNSServerAddress   ipcpOption = 130
	secondaryDNSServerAddress  ipcpOption = 131
	secondaryNBNSServerAddress ipcpOption = 132
)

// NewIpcp makes Ipcp struct.
// Since it is not an important function, a provisional implementation of
// only DNS related.
func NewIpcp(code ipcpCode, identifier byte, priDNS, secDNS net.IP) *Ipcp {
	return &Ipcp{
		header{ipcpNum, 16},
		code,
		identifier,
		priDNS.To4(),
		secDNS.To4(),
	}
}

func (i *Ipcp) marshal() []byte {
	body := make([]byte, i.length)

	body[0] = byte(i.code)
	body[1] = byte(i.identifier)

	// length
	binary.BigEndian.PutUint16(body[2:4], 16)

	// Primary DNS Server
	body[4] = byte(primaryDNSServerAddress)
	body[5] = 6
	copy(body[6:10], i.priDNS)

	// Secondary DNS Server
	body[10] = byte(secondaryDNSServerAddress)
	body[11] = 6
	copy(body[12:16], i.secDNS)

	return i.header.marshal(body)
}

func unmarshalIpcp(buf []byte) (*Ipcp, error) {
	if len(buf) < 4 {
		return nil, fmt.Errorf("Too short IPCP binary : %v", buf)
	}
	code := ipcpCode(buf[0])
	identifier := buf[1]
	length := binary.BigEndian.Uint16(buf[2:4])
	if len(buf) != int(length) {
		return nil, fmt.Errorf("Invalid length IPCP binary : %v", buf)
	}

	priDNS := net.IPv4(0, 0, 0, 0)
	secDNS := net.IPv4(0, 0, 0, 0)

	tail := buf[4:]
	for len(tail) > 0 {
		if len(tail) < 2 {
			return nil, fmt.Errorf("Too short IPCP configuration option : %v", tail)
		}
		option := ipcpOption(tail[0])
		optionLength := tail[1]
		if len(tail) < int(optionLength) {
			return nil, fmt.Errorf("Too short IPCP configuration option : %v", tail)
		}
		switch option {
		case primaryDNSServerAddress:
			ipv4buf := make([]byte, 4)
			copy(ipv4buf, tail[2:6])
			priDNS = net.IP(ipv4buf)
		case secondaryDNSServerAddress:
			ipv4buf := make([]byte, 4)
			copy(ipv4buf, tail[2:6])
			secDNS = net.IP(ipv4buf)
		}
		tail = tail[optionLength:]
	}

	return NewIpcp(code, identifier, priDNS, secDNS), nil
}

func (i *Ipcp) Code() ipcpCode {
	return i.code
}

func (i *Ipcp) Identifier() byte {
	return i.identifier
}
func (i *Ipcp) PriDNS() net.IP {
	return i.priDNS
}

func (i *Ipcp) SecDNS() net.IP {
	return i.secDNS
}
