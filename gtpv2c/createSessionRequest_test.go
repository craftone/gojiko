package gtpv2c

import (
	"net"
	"testing"

	"github.com/craftone/gojiko/gtp"
	"github.com/stretchr/testify/assert"
)

func TestCreateSessionRequest_Marshal(t *testing.T) {
	csReqArg, err := MakeCSReqArg(
		"440101234567890",  // IMSI
		"819012345678",     // MSISDN
		"1234567812345678", // MEI
		"440", "10",        // MCC, MNC
		net.IPv4(1, 2, 3, 4), gtp.Teid(0x01234567), // Sgw Ctrl FTEID
		net.IPv4(5, 6, 7, 8), gtp.Teid(0x76543210), // SGW Data FTEID
		"apn.jp", 5, 1) // APN, EBI, Recovery
	assert.NoError(t, err)

	csReq, err := NewCreateSessionRequest(0x1234, csReqArg)
	csReqBin := csReq.Marshal()
	assert.Equal(t, []byte{
		0x48,   // First octet
		0x20,   // CSReq(32)
		0, 220, // Length
		0, 0, 0, 0, // TEID
		0x00, 0x12, 0x34, // Seq Num
		0,          // Spare
		1, 0, 8, 0, // IMSI header
		0x44, 0x10, 0x10, 0x32, 0x54, 0x76, 0x98, 0xf0,
		0x4c, 0, 6, 0, // MSISDN header
		0x18, 0x09, 0x21, 0x43, 0x65, 0x87,
		0x4b, 0, 8, 0, // MEI header
		0x21, 0x43, 0x65, 0x87, 0x21, 0x43, 0x65, 0x87,
		0x56, 0, 13, 0, // ULI header
		0x18, 0x44, 0xf0, 0x01, 0, 0, 0x44, 0xf0, 0x1, 0, 0, 0, 0,
		0x53, 0, 3, 0, 0x44, 0xf0, 0x01, // Serving Network
		0x52, 0, 1, 0, 6, // RAT Type
		0x4d, 0, 3, 0, 0, 0, 0, // Indication
		0x57, 0, 9, 0, // Sender F-TEID header
		0x86, 0x01, 0x23, 0x45, 0x67, 1, 2, 3, 4,
		0x47, 0, 6, 0, 0x61, 0x70, 0x6e, 0x2e, 0x6a, 0x70, // APN
		0x80, 0, 1, 0, 1, // Selection Mode
		0x63, 0, 1, 0, 1, // PDN Type
		0x4f, 0, 5, 0, 1, 0, 0, 0, 0, // PAA
		0x7f, 0, 1, 0, 0, // Max APN Restriction
		0x48, 0, 8, 0, 0, 0x41, 0x89, 0x37, 0, 0x41, 0x89, 0x37, // APN AMBR
		0x4e, 0, 26, 0, 0x80, // PCO header
		0x80, 0x21, 0x10, // IPCP header
		1,        // Code : Configure-Request
		0,        // Identifier : 0
		00, 0x0c, // Length: 12
		0x81,       // Option : 129 Primary DNS
		6,          // Length : 6
		0, 0, 0, 0, // 0.0.0.0
		0x83,       // Option : 131 Secondary DNS
		6,          // Length : 6
		0, 0, 0, 0, // 0.0.0.0
		0, 0x0d, 0, // IPv4-DNS-Server
		0, 0x0a, 0, // IP address allocation via NAS signalling
		0x5d, 0, 44, 0, // Bearer Context to be created
		0x49, 0, 1, 0, 5, // EBI
		0x50, 0, 22, 0, // Bearer QoS header
		0x7c, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0x57, 0, 9, 2, // SGW-DATA FTEID header
		0x84, 0x76, 0x54, 0x32, 0x10, 5, 6, 7, 8,
		0x03, 0, 1, 0, 1, // Recovery
	}, csReqBin)
}

func TestUnmarshal_CreateSessionRequest(t *testing.T) {
	csReqArg, _ := MakeCSReqArg(
		"440101234567890",  // IMSI
		"819012345678",     // MSISDN
		"1234567812345678", // MEI
		"440", "10",        // MCC, MNC
		net.IPv4(1, 2, 3, 4), gtp.Teid(0x01234567), // Sgw Ctrl FTEID
		net.IPv4(5, 6, 7, 8), gtp.Teid(0x76543210), // SGW Data FTEID
		"apn.jp", 5, 1) // APN, EBI, Recovery

	csReq, err := NewCreateSessionRequest(0x1234, csReqArg)
	csReqBin := csReq.Marshal()
	msg, tail, err := Unmarshal(csReqBin)

	csReq = msg.(*CreateSessionRequest)
	assert.Equal(t, uint32(0x1234), csReq.SeqNum())
	assert.Equal(t, "440101234567890", csReq.Imsi().Value())
	assert.Equal(t, "819012345678", csReq.Msisdn().Value())
	assert.Equal(t, "1234567812345678", csReq.Mei().Value())
	assert.Equal(t, "440", csReq.Uli().Tai().Mcc())
	assert.Equal(t, "10", csReq.Uli().Tai().Mnc())
	assert.Equal(t, "440", csReq.Uli().Ecgi().Mcc())
	assert.Equal(t, "10", csReq.Uli().Ecgi().Mnc())
	assert.Equal(t, "440", csReq.ServingNetwork().Mcc())
	assert.Equal(t, "10", csReq.ServingNetwork().Mnc())
	assert.Equal(t, byte(6), csReq.RatType().Value())
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), csReq.SgwCtrlFteid().Ipv4())
	assert.Equal(t, gtp.Teid(0x01234567), csReq.SgwCtrlFteid().Teid())
	assert.Equal(t, "apn.jp", csReq.Apn().Value())
	assert.Equal(t, byte(5), csReq.BearerContextTBC().Ebi().Value())
	assert.Equal(t, net.IPv4(5, 6, 7, 8).To4(), csReq.BearerContextTBC().SgwDataFteid().Ipv4())
	assert.Equal(t, gtp.Teid(0x76543210), csReq.BearerContextTBC().SgwDataFteid().Teid())
	assert.Equal(t, byte(1), csReq.Recovery().Value())

	assert.Equal(t, tail, []byte{})
	assert.Nil(t, err)
}
