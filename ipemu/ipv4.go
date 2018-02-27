package ipemu

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/craftone/gojiko/gtp"
)

type Protocol byte

const (
	ICMP Protocol = 1
	UDP  Protocol = 17
)

type IPv4Emulator struct {
	identification uint16
	protocol       Protocol
	sourceAddr     net.IP
	destAddr       net.IP
	mtu            uint16
	mtx            sync.Mutex
}

func NewIPv4Emulator(proto Protocol, saddr, daddr net.IP, mtu uint16) *IPv4Emulator {
	return &IPv4Emulator{
		identification: 0,
		protocol:       proto,
		sourceAddr:     saddr.To4(),
		destAddr:       daddr.To4(),
		mtu:            mtu,
	}
}

type IPv4 struct {
	version             byte
	ihl                 byte
	tos                 byte
	totalLength         uint16
	identification      uint16
	variousControlFlags byte
	fragmentOffset      uint16
	ttl                 byte
	headerChecksum      byte
	destinationAddress  net.IP
	data                []byte
}

func (e *IPv4Emulator) NewIPv4GPDU(teid gtp.Teid, tos, ttl byte, data []byte) ([]byte, error) {
	ipVersion := 4
	ihl := 5
	headerSize := ihl * 4
	maxDataSize := (e.mtu - uint16(headerSize)) & 0xfff8
	if len(data) > int(maxDataSize) {
		return []byte{}, fmt.Errorf("Too big data : max data size is %d but the data size is %d", maxDataSize, len(data))
	}
	totalLength := uint16(headerSize + len(data))

	// make gpdu
	gpdu := make([]byte, totalLength+12)
	gpdu[0] = 0x20 // version:1, all flags are 0
	gpdu[1] = 0xFF // GTP_TPDU_MSG (0xFF)
	binary.BigEndian.PutUint16(gpdu[2:], 4+totalLength)
	binary.BigEndian.PutUint32(gpdu[4:], uint32(teid))

	// make packet
	packet := gpdu[12:]
	packet[0] = byte(ipVersion<<4 + ihl)
	packet[1] = tos
	binary.BigEndian.PutUint16(packet[2:], totalLength)
	binary.BigEndian.PutUint16(packet[4:], 0)
	binary.BigEndian.PutUint16(packet[6:], 0x4000) // VCS(3bit) + fragmentOffset(15bit)
	packet[8] = ttl
	packet[9] = byte(e.protocol)
	copy(packet[12:16], e.sourceAddr)
	copy(packet[16:20], e.destAddr)
	copy(packet[20:], data)

	// checksum
	crc := 0
	for _, i := range []int{0, 2, 4, 6, 8, 12, 14, 16, 18} {
		crc += int(binary.BigEndian.Uint16(packet[i : i+2]))
	}

	for {
		carry := crc >> 16
		if carry == 0 {
			break
		}
		crc = crc & 0xFFFF
		crc += carry
	}
	checksum := ^uint16(crc)
	binary.BigEndian.PutUint16(packet[10:], checksum)

	return gpdu, nil
}

func (e *IPv4Emulator) getIdentification() uint16 {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	identification := e.identification
	e.identification++
	return identification
}
