package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal_unknown(t *testing.T) {
	recBin := NewRecovery(255, 0).Marshal()
	recBin[0] = 0
	msg, _, err := Unmarshal(recBin)
	assert.Nil(t, msg)
	assert.Error(t, err)
}

func TestUnmarshal_recovery(t *testing.T) {
	recBin := NewRecovery(255, 0).Marshal()
	msg, tail, err := Unmarshal(recBin)
	rec := msg.(*Recovery)
	assert.Equal(t, byte(255), rec.Value)
	assert.Equal(t, byte(0), rec.header.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_recoveryWithTail(t *testing.T) {
	recBin := NewRecovery(255, 0).Marshal()
	recBin = append(recBin, recBin...)
	msg, tail, err := Unmarshal(recBin)
	rec := msg.(*Recovery)
	assert.Equal(t, byte(255), rec.Value)
	assert.Equal(t, byte(0), rec.header.instance)
	assert.Equal(t, []byte{0x3, 0, 1, 0, 0xff}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_imsi(t *testing.T) {
	imsiBin := NewImsi("819012345678", 1).Marshal()
	msg, tail, err := Unmarshal(imsiBin)
	imsi := msg.(*Imsi)
	assert.Equal(t, "819012345678", imsi.Value)
	assert.Equal(t, byte(1), imsi.header.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_imsiWithTail(t *testing.T) {
	imsiBin := NewImsi("819012345678", 1).Marshal()
	imsiBin = append(imsiBin, imsiBin...)
	msg, tail, err := Unmarshal(imsiBin)
	imsi := msg.(*Imsi)
	assert.Equal(t, "819012345678", imsi.Value)
	assert.Equal(t, byte(1), imsi.header.instance)
	assert.Equal(t, []byte{1, 0, 6, 1, 0x18, 0x09, 0x21, 0x43, 0x65, 0x87}, tail)
	assert.Nil(t, err)
}
