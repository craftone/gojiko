package gtpv2c

import (
	"net"
	"testing"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/stretchr/testify/assert"
)

func TestCreateSessionResponse_Marshal(t *testing.T) {
	csResArg, err := MakeCSResArg(
		gtp.Teid(0x12345678),                       // SgwCtrlTEID
		ie.CauseRequestAccepted,                    // Cause
		net.IPv4(1, 2, 3, 4), gtp.Teid(0x01234567), // PGW Ctrl FTEID
		net.IPv4(5, 6, 7, 8), gtp.Teid(0x76543210), // PGW Data FTEID
		net.IPv4(9, 10, 11, 12),                    // PDN Allocated IP address
		net.IPv4(8, 8, 8, 8), net.IPv4(8, 8, 4, 4), // PriDNS, SecDNS
		5) // EBI
	assert.NoError(t, err)

	csRes, err := NewCreateSessionResponse(0x1234, csResArg)
	csResBin := csRes.Marshal()
	assert.Equal(t, []byte{
		0x48,   // First octet
		0x21,   // CSRes(33)
		0, 115, // Length
		0x12, 0x34, 0x56, 0x78, // TEID
		0x00, 0x12, 0x34, // Seq Num
		0,                 // Spare
		2, 0, 2, 0, 16, 0, // Cause
		0x57, 0, 9, 1, // PGW Ctrl F-TEID header
		0x87, 0x01, 0x23, 0x45, 0x67, 1, 2, 3, 4,
		0x4f, 0, 5, 0, 1, 9, 10, 11, 12, // PAA
		0x7f, 0, 1, 0, 0, // Max APN Restriction
		0x4e, 0, 34, 0, 0x80, // PCO header
		0x80, 0x21, 0x10, // IPCP header
		3,        // Code : Configure-Response
		0,        // Identifier : 0
		00, 0x0c, // Length: 12
		0x81,       // Option : 129 Primary DNS
		6,          // Length : 6
		8, 8, 8, 8, // 8.8.8.8
		0x83,       // Option : 131 Secondary DNS
		6,          // Length : 6
		8, 8, 4, 4, // 8.8.4.4
		0x00, 0x0d, 4, 8, 8, 8, 8, //PriDNS
		0x00, 0x0d, 4, 8, 8, 4, 4, //PriDNS
		0x5d, 0, 32, 0, // Bearer Context to be created
		0x49, 0, 1, 0, 5, // EBI
		2, 0, 2, 0, 16, 0, // Cause
		0x57, 0, 9, 2, // PGW-DATA FTEID header
		0x85, 0x76, 0x54, 0x32, 0x10, 5, 6, 7, 8,
		0x5E, 0, 4, 0, 0x12, 0x34, 0x56, 0x78, // Charging ID
	}, csResBin)
}

func TestUnmarshal_CreateSessionResponse(t *testing.T) {
	csResArg, err := MakeCSResArg(
		gtp.Teid(0x12345678),                       // SgwCtrlTEID
		ie.CauseRequestAccepted,                    // Cause
		net.IPv4(1, 2, 3, 4), gtp.Teid(0x01234567), // PGW Ctrl FTEID
		net.IPv4(5, 6, 7, 8), gtp.Teid(0x76543210), // PGW Data FTEID
		net.IPv4(9, 10, 11, 12),                    // PDN Allocated IP address
		net.IPv4(8, 8, 8, 8), net.IPv4(8, 8, 4, 4), // PriDNS, SecDNS
		5) // EBI
	csResArg.Recovery, _ = ie.NewRecovery(0, 128)

	csRes, err := NewCreateSessionResponse(0x1234, csResArg)
	csResBin := csRes.Marshal()
	msg, tail, err := Unmarshal(csResBin)

	csRes = msg.(*CreateSessionResponse)
	assert.Equal(t, uint32(0x1234), csRes.SeqNum())
	assert.Equal(t, gtp.Teid(0x12345678), csRes.Teid())
	assert.Equal(t, ie.CauseRequestAccepted, csRes.Cause().Value())
	assert.Equal(t, ie.CauseRequestAccepted, csRes.BearerContextCeated().Cause().Value())
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), csRes.PgwCtrlFteid().Ipv4())
	assert.Equal(t, gtp.Teid(0x01234567), csRes.PgwCtrlFteid().Teid())
	assert.Equal(t, net.IPv4(5, 6, 7, 8).To4(), csRes.BearerContextCeated().PgwDataFteid().Ipv4())
	assert.Equal(t, gtp.Teid(0x76543210), csRes.BearerContextCeated().PgwDataFteid().Teid())
	assert.Equal(t, net.IPv4(9, 10, 11, 12).To4(), csRes.Paa().IPv4())
	assert.Equal(t, net.IPv4(8, 8, 8, 8).To4(), csRes.Pco().DNSServerV4s()[0].Value())
	assert.Equal(t, net.IPv4(8, 8, 4, 4).To4(), csRes.Pco().DNSServerV4s()[1].Value())
	assert.Equal(t, byte(5), csRes.BearerContextCeated().Ebi().Value())
	assert.Equal(t, byte(128), csRes.Recovery().Value())

	assert.Equal(t, tail, []byte{})
	assert.Nil(t, err)

	// Some unkown type IE.

	csResBin[0x77] = 0xff
	msg, tail, err = Unmarshal(csResBin)

	csRes = msg.(*CreateSessionResponse)
	assert.Equal(t, uint32(0x1234), csRes.SeqNum())
	assert.Equal(t, gtp.Teid(0x12345678), csRes.Teid())
	assert.Equal(t, ie.CauseRequestAccepted, csRes.Cause().Value())
	assert.Equal(t, ie.CauseRequestAccepted, csRes.BearerContextCeated().Cause().Value())
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), csRes.PgwCtrlFteid().Ipv4())
	assert.Equal(t, gtp.Teid(0x01234567), csRes.PgwCtrlFteid().Teid())
	assert.Equal(t, net.IPv4(5, 6, 7, 8).To4(), csRes.BearerContextCeated().PgwDataFteid().Ipv4())
	assert.Equal(t, gtp.Teid(0x76543210), csRes.BearerContextCeated().PgwDataFteid().Teid())
	assert.Equal(t, net.IPv4(9, 10, 11, 12).To4(), csRes.Paa().IPv4())
	assert.Equal(t, net.IPv4(8, 8, 8, 8).To4(), csRes.Pco().DNSServerV4s()[0].Value())
	assert.Equal(t, net.IPv4(8, 8, 4, 4).To4(), csRes.Pco().DNSServerV4s()[1].Value())
	assert.Equal(t, byte(5), csRes.BearerContextCeated().Ebi().Value())
	assert.Nil(t, csRes.Recovery())

	assert.Equal(t, tail, []byte{})
	assert.Nil(t, err)
}
