package gtpv2c

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_echoRequest(t *testing.T) {
	erBin := NewEchoRequest(1, 2).Marshal()
	msg, _, err := Unmarshal(erBin)
	er := msg.(*EchoRequest)
	assert.Equal(t, uint32(1), er.header.seqNum)
	assert.Equal(t, byte(2), er.recovery.Value)
	assert.Nil(t, err)
}
