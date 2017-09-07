package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecovery_new(t *testing.T) {
	rec := NewRecovery(0, 1)
	assert.Equal(t, recoveryNum, rec.header.typeNum)
	assert.Equal(t, byte(1), rec.Value)
}

func TestRecovery_marshal(t *testing.T) {
	var rec []byte
	rec = NewRecovery(0, 0).Marshal()
	assert.Equal(t, []byte{3, 0, 1, 0, 0}, rec)

	rec = NewRecovery(0xf, 255).Marshal()
	assert.Equal(t, []byte{3, 0, 1, 0xf, 255}, rec)
}

func TestUnmarshal_recovery(t *testing.T) {
	recBin := NewRecovery(0, 255).Marshal()
	msg, tail, err := Unmarshal(recBin)
	rec := msg.(*Recovery)
	assert.Equal(t, byte(255), rec.Value)
	assert.Equal(t, byte(0), rec.header.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_recoveryWithTail(t *testing.T) {
	recBin := NewRecovery(0, 255).Marshal()
	recBin = append(recBin, recBin...)
	msg, tail, err := Unmarshal(recBin)
	rec := msg.(*Recovery)
	assert.Equal(t, byte(255), rec.Value)
	assert.Equal(t, byte(0), rec.header.instance)
	assert.Equal(t, []byte{0x3, 0, 1, 0, 0xff}, tail)
	assert.Nil(t, err)
}
