package gtpv2c

import (
	"net"
	"testing"

	"github.com/craftone/gojiko/domain/gtpv2c/ie"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/stretchr/testify/assert"
)

func TestModifyBearerRequest_Marshal(t *testing.T) {
	mbReqArg, err := MakeMBReqArg(
		0x0123,            // PGW Ctrl TEID
		"440101234567890", // IMSI
		"440", "10",       // MCC, MNC
		0x1234, 0x01234567, // TAC, ECI
		ie.RatTypeWbEutran, ie.IndicationArg{
			CLII: true,
		},
		net.IPv4(1, 2, 3, 4), gtp.Teid(0x01234567), // SGW Ctrl FTEID
		net.IPv4(5, 6, 7, 8), gtp.Teid(0x76543210), // SGW Data FTEID
		5, 1) // EBI, Recovery
	assert.NoError(t, err)

	mbReq, err := NewModifyBearerRequest(0x1234, mbReqArg)
	mbReqBin := mbReq.Marshal()
	assert.Equal(t, []byte{
		0x48,  // First octet
		0x22,  // MBReq(34)
		0, 93, // Length
		0, 0, 1, 0x23, // TEID
		0x00, 0x12, 0x34, // Seq Num
		0,                                                             // Spare
		0x01, 0, 8, 0, 0x44, 0x10, 0x10, 0x32, 0x54, 0x76, 0x98, 0xf0, // IMSI
		0x56, 0, 13, 0, // ULI header
		0x18, 0x44, 0xf0, 0x01, 0x12, 0x34, 0x44, 0xf0, 0x1, 0x01, 0x23, 0x45, 0x67,
		0x52, 0, 1, 0, 6, // RAT Type
		0x4d, 0, 7, 0, 0, 0, 0, 0x2, 0, 0, 0, // Indication
		0x57, 0, 9, 0, // Sender F-TEID header
		0x86, 0x01, 0x23, 0x45, 0x67, 1, 2, 3, 4,
		0x5d, 0, 18, 0, // Bearer Context to be created
		0x49, 0, 1, 0, 5, // EBI
		0x57, 0, 9, 2, // SGW-DATA FTEID header
		0x84, 0x76, 0x54, 0x32, 0x10, 5, 6, 7, 8,
		0x03, 0, 1, 0, 1, // Recovery
	}, mbReqBin)
}

func TestUnmarshal_ModifyBearerRequest(t *testing.T) {
	mbReqArg, err := MakeMBReqArg(
		0x0123,            // PGW Ctrl TEID
		"440101234567890", // IMSI
		"440", "10",       // MCC, MNC
		0x1234, 0x01234567, // TAC, ECI
		ie.RatTypeWbEutran, ie.IndicationArg{
			CLII: true,
		},
		net.IPv4(1, 2, 3, 4), gtp.Teid(0x01234567), // SGW Ctrl FTEID
		net.IPv4(5, 6, 7, 8), gtp.Teid(0x76543210), // SGW Data FTEID
		5, 1) // EBI, Recovery
	assert.NoError(t, err)

	mbReq, err := NewModifyBearerRequest(0x1234, mbReqArg)
	mbReqBin := mbReq.Marshal()
	msg, tail, err := Unmarshal(mbReqBin)

	mbReq = msg.(*ModifyBearerRequest)
	assert.Equal(t, uint32(0x1234), mbReq.SeqNum())
	assert.Equal(t, "440", mbReq.Uli().Tai().Mcc())
	assert.Equal(t, "10", mbReq.Uli().Tai().Mnc())
	assert.Equal(t, uint16(0x1234), mbReq.Uli().Tai().Tac())
	assert.Equal(t, "440", mbReq.Uli().Ecgi().Mcc())
	assert.Equal(t, "10", mbReq.Uli().Ecgi().Mnc())
	assert.Equal(t, uint32(0x01234567), mbReq.Uli().Ecgi().Eci())
	assert.Equal(t, ie.RatTypeWbEutran, mbReq.RatType().Value())
	assert.Equal(t, net.IPv4(1, 2, 3, 4).To4(), mbReq.SgwCtrlFteid().Ipv4())
	assert.Equal(t, gtp.Teid(0x01234567), mbReq.SgwCtrlFteid().Teid())
	assert.Equal(t, byte(5), mbReq.BearerContextTBM().Ebi().Value())
	assert.Equal(t, net.IPv4(5, 6, 7, 8).To4(), mbReq.BearerContextTBM().SgwDataFteid().Ipv4())
	assert.Equal(t, gtp.Teid(0x76543210), mbReq.BearerContextTBM().SgwDataFteid().Teid())
	assert.Equal(t, byte(1), mbReq.Recovery().Value())

	assert.Equal(t, tail, []byte{})
	assert.Nil(t, err)
}
