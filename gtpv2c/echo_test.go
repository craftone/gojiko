package gtpv2c

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEchoRequest(t *testing.T) {
	er := NewEchoRequest(1, 2)
	assert.Equal(t, byte(2), er.header.version)
	assert.Equal(t, false, er.header.piggybackingFlag)
	assert.Equal(t, false, er.header.teidFlag)
	assert.Equal(t, echoRequest, er.header.messageType)
	assert.Equal(t, uint32(1), er.header.seqNum)
	assert.Equal(t, uint32(0), er.header.teid)
	assert.Equal(t, byte(2), er.recovery.Value)
}

func TestEchoRequest_Marshal(t *testing.T) {
	er := NewEchoRequest(1, 2)
	assert.Equal(t, []byte{0x40, 1, 0, 9, 0, 0, 1, 0, 3, 0, 1, 0, 2}, er.Marshal())

	er = NewEchoRequest(0x123456, 0xff)
	assert.Equal(t, []byte{0x40, 1, 0, 9, 0x12, 0x34, 0x56, 0, 3, 0, 1, 0, 0xff}, er.Marshal())
}

func TestNewEchoResponse(t *testing.T) {
	er := NewEchoResponse(1, 2)
	assert.Equal(t, byte(2), er.header.version)
	assert.Equal(t, false, er.header.piggybackingFlag)
	assert.Equal(t, false, er.header.teidFlag)
	assert.Equal(t, echoResponse, er.header.messageType)
	assert.Equal(t, uint32(1), er.header.seqNum)
	assert.Equal(t, uint32(0), er.header.teid)
	assert.Equal(t, byte(2), er.recovery.Value)
}

func TestEchoResponse_Marshal(t *testing.T) {
	er := NewEchoResponse(1, 2)
	assert.Equal(t, []byte{0x40, 2, 0, 9, 0, 0, 1, 0, 3, 0, 1, 0, 2}, er.Marshal())

	er = NewEchoResponse(0x123456, 0xff)
	assert.Equal(t, []byte{0x40, 2, 0, 9, 0x12, 0x34, 0x56, 0, 3, 0, 1, 0, 0xff}, er.Marshal())
}
