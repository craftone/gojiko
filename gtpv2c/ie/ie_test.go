package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_recovery(t *testing.T) {
	recBin := NewRecovery(255, 0).Marshal()
	msg, tail, err := Unmarshal(recBin)
	rec := msg.(*Recovery)
	assert.Equal(t, byte(255), rec.Value)
	assert.Equal(t, byte(0), rec.header.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_recoveryWitTail(t *testing.T) {
	recBin := NewRecovery(255, 0).Marshal()
	recBin = append(recBin, recBin...)
	msg, tail, err := Unmarshal(recBin)
	rec := msg.(*Recovery)
	assert.Equal(t, byte(255), rec.Value)
	assert.Equal(t, byte(0), rec.header.instance)
	assert.Equal(t, []byte{0x3, 0, 1, 0, 0xff}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_unknown(t *testing.T) {
	recBin := NewRecovery(255, 0).Marshal()
	recBin[0] = 0
	msg, _, err := Unmarshal(recBin)
	assert.Nil(t, msg)
	assert.Error(t, err)
}
