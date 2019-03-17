package gtpv2c

import (
	"testing"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/stretchr/testify/assert"
)

func TestModifyBearerResponse_Marshal(t *testing.T) {
	mbResArg, err := MakeMBResArg(
		gtp.Teid(0x12345678),    // SgwCtrlTEID
		ie.CauseRequestAccepted, // Cause
		"819012345678", 0x1234,  // MSISDN, ChargingID
		5, 128) // EBI, Recovery
	assert.NoError(t, err)

	mbRes, err := NewModifyBearerResponse(0x1234, mbResArg)
	mbResBin := mbRes.Marshal()
	assert.Equal(t, []byte{
		0x48,  // First octet
		0x23,  // MBRes(35)
		0, 52, // Length
		0x12, 0x34, 0x56, 0x78, // TEID
		0x00, 0x12, 0x34, // Seq Num
		0,                 // Spare
		2, 0, 2, 0, 16, 0, // Cause
		0x4c, 0, 6, 0, 0x18, 0x09, 0x21, 0x43, 0x65, 0x87, // MSISDN
		0x5d, 0, 19, 0, // Bearer Context to be created
		0x49, 0, 1, 0, 5, // EBI
		2, 0, 2, 0, 16, 0, // Cause
		0x5E, 0, 4, 0, 0, 0, 0x12, 0x34, // Charging ID
		0x3, 0, 1, 0, 128, // Recovery
	}, mbResBin)
}

func TestUnmarshal_ModifyBearerResponse(t *testing.T) {
	mbResArg, err := MakeMBResArg(
		gtp.Teid(0x12345678),    // SgwCtrlTEID
		ie.CauseRequestAccepted, // Cause
		"819012345678", 0x1234,  // MSISDN, ChargingID
		5, 128) // EBI, Recovery
	assert.NoError(t, err)

	mbRes, err := NewModifyBearerResponse(0x1234, mbResArg)
	mbResBin := mbRes.Marshal()
	msg, tail, err := Unmarshal(mbResBin)

	mbRes = msg.(*ModifyBearerResponse)
	assert.Equal(t, uint32(0x1234), mbRes.SeqNum())
	assert.Equal(t, gtp.Teid(0x12345678), mbRes.Teid())
	assert.Equal(t, ie.CauseRequestAccepted, mbRes.Cause().Value())
	assert.Equal(t, "819012345678", mbRes.Msisdn().Value())
	assert.Equal(t, ie.CauseRequestAccepted, mbRes.BearerContextTBM().Cause().Value())
	assert.Equal(t, uint32(0x1234), mbRes.BearerContextTBM().ChargingID().Value())
	assert.Equal(t, byte(128), mbRes.Recovery().Value())

	assert.Equal(t, tail, []byte{})
	assert.Nil(t, err)
}
