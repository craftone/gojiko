package gtpv2c

import (
	"testing"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSessionRequest_Marshal(t *testing.T) {
	dbReq, err := NewDeleteSessionRequest(gtp.Teid(0x1234), 0x5678, byte(5))
	assert.NoError(t, err)
	dbReqBin := dbReq.Marshal()
	assert.Equal(t, []byte{
		0x48,  // First octet
		36,    // DSReq(36)
		0, 13, // Length
		0, 0, 0x12, 0x34, // TEID
		0x00, 0x56, 0x78, // Seq Num
		0,                // Spare
		0x49, 0, 1, 0, 5, // EBI
	}, dbReqBin)
}

func TestDeleteSessionRequest_Unmarshal(t *testing.T) {
	dsReq, _ := NewDeleteSessionRequest(gtp.Teid(0x1234), 0x5678, byte(5))
	dsReqBin := dsReq.Marshal()
	msg, tail, err := Unmarshal(dsReqBin)
	assert.Equal(t, []byte{}, tail)
	assert.NoError(t, err)

	dsReq = msg.(*DeleteSessionRequest)
	assert.Equal(t, gtp.Teid(0x1234), dsReq.Teid())
	assert.Equal(t, uint32(0x5678), dsReq.SeqNum())
	assert.Equal(t, byte(5), dsReq.Lbi().Value())
}
