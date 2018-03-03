package gtpv2c

import (
	"testing"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/stretchr/testify/assert"
)

func TestNewEchoRequest(t *testing.T) {
	er, err := NewEchoRequest(1, 2)
	assert.Equal(t, byte(2), er.version)
	assert.Equal(t, false, er.piggybackingFlag)
	assert.Equal(t, false, er.teidFlag)
	assert.Equal(t, EchoRequestNum, er.messageType)
	assert.Equal(t, uint32(1), er.seqNum)
	assert.Equal(t, gtp.Teid(0), er.teid)
	assert.Equal(t, byte(2), er.Recovery().Value())
	assert.NoError(t, err)
}

func TestEchoRequest_Marshal(t *testing.T) {
	er, _ := NewEchoRequest(1, 2)
	assert.Equal(t, []byte{0x40, 1, 0, 9, 0, 0, 1, 0, 3, 0, 1, 0, 2}, er.Marshal())

	er, _ = NewEchoRequest(0x123456, 0xff)
	assert.Equal(t, []byte{0x40, 1, 0, 9, 0x12, 0x34, 0x56, 0, 3, 0, 1, 0, 0xff}, er.Marshal())
}

func TestNewEchoResponse(t *testing.T) {
	er, err := NewEchoResponse(1, 2)
	assert.Equal(t, byte(2), er.version)
	assert.Equal(t, false, er.piggybackingFlag)
	assert.Equal(t, false, er.teidFlag)
	assert.Equal(t, EchoResponseNum, er.messageType)
	assert.Equal(t, uint32(1), er.seqNum)
	assert.Equal(t, gtp.Teid(0), er.teid)
	assert.Equal(t, byte(2), er.Recovery().Value())
	assert.NoError(t, err)
}

func TestEchoResponse_Marshal(t *testing.T) {
	er, _ := NewEchoResponse(1, 2)
	assert.Equal(t, []byte{0x40, 2, 0, 9, 0, 0, 1, 0, 3, 0, 1, 0, 2}, er.Marshal())

	er, _ = NewEchoResponse(0x123456, 0xff)
	assert.Equal(t, []byte{0x40, 2, 0, 9, 0x12, 0x34, 0x56, 0, 3, 0, 1, 0, 0xff}, er.Marshal())
}

func TestUnmarshal_echoRequest(t *testing.T) {
	er, _ := NewEchoRequest(1, 2)
	erBin := er.Marshal()
	msg, _, err := Unmarshal(erBin)
	er = msg.(*EchoRequest)
	assert.Equal(t, uint32(1), er.seqNum)
	assert.Equal(t, byte(2), er.Recovery().Value())
	assert.Nil(t, err)
}

func TestUnmarshal_echoResponse(t *testing.T) {
	er, _ := NewEchoResponse(1, 2)
	erBin := er.Marshal()
	msg, _, err := Unmarshal(erBin)
	er = msg.(*EchoResponse)
	assert.Equal(t, uint32(1), er.seqNum)
	assert.Equal(t, byte(2), er.Recovery().Value())
	assert.Nil(t, err)
}
