package ipemu

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/craftone/gojiko/domain/gtp"
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

// type IPv4 struct {
// 	version             byte
// 	ihl                 byte
// 	tos                 byte
// 	totalLength         uint16
// 	identification      uint16
// 	variousControlFlags byte
// 	fragmentOffset      uint16
// 	ttl                 byte
// 	headerChecksum      byte
// 	destinationAddress  net.IP
// 	data                []byte
// }

// NewIPv4GPDU creates a GTPv1-U packet : GTPv1 header + IP header + Payload.
// Payload have layers over IP layer (TCP/UDP/ICMP etc).
func (e *IPv4Emulator) NewIPv4GPDU(teid gtp.Teid, tos, ttl byte, payload []byte) ([]byte, error) {
	ipVersion := 4
	ihl := 5
	headerSize := ihl * 4
	maxDataSize := (e.mtu - uint16(headerSize)) & 0xfff8
	if len(payload) > int(maxDataSize) {
		return []byte{}, fmt.Errorf("Too big payload : max payload size is %d but the payload size is %d", maxDataSize, len(payload))
	}
	totalLength := uint16(headerSize + len(payload))

	// make gpdu
	gpdu := make([]byte, totalLength+8)
	gpdu[0] = 0x30 // version:1, ProtocolType:1(GTP), all flags are 0
	gpdu[1] = 0xFF // GTP_TPDU_MSG (0xFF)
	binary.BigEndian.PutUint16(gpdu[2:], totalLength)
	binary.BigEndian.PutUint32(gpdu[4:], uint32(teid))

	// make packet
	packet := gpdu[8:]
	packet[0] = byte(ipVersion<<4 + ihl)
	packet[1] = tos
	binary.BigEndian.PutUint16(packet[2:], totalLength)
	binary.BigEndian.PutUint16(packet[4:], 0)
	binary.BigEndian.PutUint16(packet[6:], 0x4000) // VCS(3bit) + fragmentOffset(15bit)
	packet[8] = ttl
	packet[9] = byte(e.protocol)
	copy(packet[12:16], e.sourceAddr)
	copy(packet[16:20], e.destAddr)
	copy(packet[20:], payload)

	checksum := Checksum(0, packet[0:20])
	binary.BigEndian.PutUint16(packet[10:], checksum)

	return gpdu, nil
}

// Checksum calculates IPv4/UDP/TCP checksum.
func Checksum(initial uint16, data []byte) uint16 {
	sum := uint(initial)
	for i := 0; i < len(data); i += 2 {
		if len(data)-i >= 2 {
			sum += uint(binary.BigEndian.Uint16(data[i : i+2]))
		} else {
			// when data length is odd, 0 is added at the end.
			sum += (uint(data[i]) << 8)
		}
	}

	for sum > 0xFFFF {
		carry := sum >> 16
		sum = sum&0xFFFF + carry
	}
	return ^uint16(sum)
}

func (e *IPv4Emulator) getIdentification() uint16 {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	identification := e.identification
	e.identification++
	return identification
}

// PickOutPayload check and pick up payload (and return the payload)
// from an IP packet (body argument).
// srcAddr, destAddr, etc. are checked ensure matching with IPv4Emulator's parameter.
func (e *IPv4Emulator) PickOutPayload(destPort uint16, body []byte) ([]byte, error) {
	bodyLen := len(body)
	if bodyLen < 20 {
		return nil, fmt.Errorf("Too short packet : length %d", bodyLen)
	}
	ihl := body[0] & 0xf
	proto := Protocol(body[9])
	if proto != e.protocol {
		return nil, fmt.Errorf("Not expected packet : protocol %d", proto)
	}
	srcAddr := net.IP(body[12:16])
	if !e.sourceAddr.Equal(srcAddr) {
		return nil, fmt.Errorf("From unexpected address : %s", srcAddr.String())
	}
	destAddr := net.IP(body[16:20])
	if !e.destAddr.Equal(destAddr) {
		return nil, fmt.Errorf("To unknown address : %s", destAddr.String())
	}
	totalLen := int(binary.BigEndian.Uint16(body[2:4]))
	if bodyLen < totalLen {
		return nil, fmt.Errorf("Too short packet : length %d", bodyLen)
	}
	headerLen := int(ihl) * 4
	if headerLen > totalLen {
		return nil, fmt.Errorf("Header length or/and Total length are invalid")
	}
	udpPacket := body[headerLen:]
	destPortAct := binary.BigEndian.Uint16(udpPacket[2:4])
	if destPortAct != destPort {
		return nil, fmt.Errorf("To unexpected port : expected %d , actual %d", destPort, destPortAct)
	}
	payload := body[headerLen+8:] // 8 is length of UDP header
	return payload, nil
}
