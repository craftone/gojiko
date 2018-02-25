package gtpv2c

import (
	"testing"

	"github.com/craftone/gojiko/gtpv2c/ie"

	"github.com/craftone/gojiko/gtp"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBearerResponse_Marshal(t *testing.T) {
	dbRes, err := NewDeleteBearerResponse(gtp.Teid(0x1234), 0x5678, ie.CauseRequestAccepted)
	assert.NoError(t, err)
	dbResBin := dbRes.Marshal()
	assert.Equal(t, []byte{
		0x48,  // First octet
		100,   // DBRes
		0, 14, // Length
		0, 0, 0x12, 0x34, // TEID
		0x00, 0x56, 0x78, // Seq Num
		0,                         // Spare
		0x02, 0, 0x02, 0, 0x10, 0, // Cause
	}, dbResBin)
}

func TestDeleteBearerResponse_Unmarshal(t *testing.T) {
	dbRes, _ := NewDeleteBearerResponse(gtp.Teid(0x1234), 0x5678, ie.CauseRequestAccepted)
	dbResBin := dbRes.Marshal()
	msg, tail, err := Unmarshal(dbResBin)
	assert.Equal(t, []byte{}, tail)
	assert.NoError(t, err)

	dbRes = msg.(*DeleteBearerResponse)
	assert.Equal(t, gtp.Teid(0x1234), dbRes.Teid())
	assert.Equal(t, uint32(0x5678), dbRes.SeqNum())
	assert.Equal(t, ie.CauseRequestAccepted, dbRes.Cause().Value())
}
