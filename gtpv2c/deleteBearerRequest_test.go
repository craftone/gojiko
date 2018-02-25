package gtpv2c

import (
	"testing"

	"github.com/craftone/gojiko/gtp"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBearerRequest_Marshal(t *testing.T) {
	dbReq, err := NewDeleteBearerRequest(gtp.Teid(0x1234), 0x5678, byte(5))
	assert.NoError(t, err)
	dbReqBin := dbReq.Marshal()
	assert.Equal(t, []byte{
		0x48,  // First octet
		99,    // DBReq(99)
		0, 13, // Length
		0, 0, 0x12, 0x34, // TEID
		0x00, 0x56, 0x78, // Seq Num
		0,                // Spare
		0x49, 0, 1, 0, 5, // EBI
	}, dbReqBin)
}

func TestDeleteBearerRequest_Unmarshal(t *testing.T) {
	dbReq, _ := NewDeleteBearerRequest(gtp.Teid(0x1234), 0x5678, byte(5))
	dbReqBin := dbReq.Marshal()
	msg, tail, err := Unmarshal(dbReqBin)
	assert.Equal(t, []byte{}, tail)
	assert.NoError(t, err)

	dbReq = msg.(*DeleteBearerRequest)
	assert.Equal(t, gtp.Teid(0x1234), dbReq.Teid())
	assert.Equal(t, uint32(0x5678), dbReq.SeqNum())
	assert.Equal(t, byte(5), dbReq.Lbi().Value())
}
