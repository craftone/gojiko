package ipemu

import (
	"net"
	"testing"

	"github.com/craftone/gojiko/gtp"

	"github.com/stretchr/testify/assert"
)

func TestIPv4Emulator_NewIPv4GPDU(t *testing.T) {
	//version: 4, ihl:5, tos:00, totalLen:0073=115, id:0000,
	//fragment: 4000,  ttl:40, proto:11, crc:b861,
	//saddr: c0,a8,00,01 = 192.168.0.1
	//darrd: c0,a8,00,c7 = 192.168.0.199,
	//data: 95bytes
	//  sample data from : https://en.wikipedia.org/wiki/IPv4_header_checksum
	ipv4Emu := NewIPv4Emulator(UDP, net.IPv4(192, 168, 0, 1), net.IPv4(192, 168, 0, 199), 1500)
	teid := gtp.Teid(0x12345678)
	data := make([]byte, 95)
	tos := byte(0)
	ttl := byte(0x40)
	packet, err := ipv4Emu.NewIPv4GPDU(teid, tos, ttl, data)
	assert.NoError(t, err)
	expected := make([]byte, 12+20+95)
	expectedHeader := []byte{
		0x20,      // GTP version:1, all flags are 0
		0xFF,      // GTP_TPDU_MSG (0xFF)
		0x00, 119, // totalLen: 4+115
		0x12, 0x34, 0x56, 0x78, // teid
		0, 0, 0, 0, // seqNum, etc
		0x45,      // version: 4, ihl: 5
		0x00,      // tos: 0,
		0x00, 115, // totalLen : 115
		0x00, 0x00, // id:0
		0x40, 0x00, // fragment: 0x4000
		0x40,       // ttl
		0x11,       // protocol
		0xb8, 0x61, // checksum
		0xc0, 0xa8, 0x00, 0x01, //source address
		0xc0, 0xa8, 0x00, 0xc7, //destination address
	}
	copy(expected[0:], expectedHeader)
	assert.Equal(t, expected, packet)

	// boundary condition
	data2 := make([]byte, 1480)
	_, err = ipv4Emu.NewIPv4GPDU(teid, tos, ttl, data2)
	assert.NoError(t, err)

	// too big (should be fragmented) data
	data3 := make([]byte, 1481)
	_, err = ipv4Emu.NewIPv4GPDU(teid, tos, ttl, data3)
	assert.Error(t, err)
}
