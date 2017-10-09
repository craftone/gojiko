package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRecovery(t *testing.T) {
	rec, _ := NewRecovery(0, 1)
	assert.Equal(t, recoveryNum, rec.header.typeNum)
	assert.Equal(t, byte(1), rec.Value())

	_, err := NewRecovery(0x0f, 1)
	assert.Nil(t, err)

	_, err = NewRecovery(0x10, 1)
	assert.Error(t, err)
	_, err = NewRecovery(0xff, 1)
	assert.Error(t, err)
}

func TestRecovery_marshal(t *testing.T) {
	rec, _ := NewRecovery(0, 0)
	recBin := rec.Marshal()
	assert.Equal(t, []byte{3, 0, 1, 0, 0}, recBin)

	rec, _ = NewRecovery(0xf, 255)
	recBin = rec.Marshal()
	assert.Equal(t, []byte{3, 0, 1, 0xf, 255}, recBin)
}

func TestUnmarshal_recovery(t *testing.T) {
	recOrg, _ := NewRecovery(0, 255)
	recBin := recOrg.Marshal()
	msg, tail, err := Unmarshal(recBin, CreateSessionRequest)
	rec := msg.(*Recovery)
	assert.Equal(t, byte(255), rec.Value())
	assert.Equal(t, byte(0), rec.header.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_recoveryWithTail(t *testing.T) {
	recOrg, _ := NewRecovery(0, 255)
	recBin := recOrg.Marshal()
	recBin = append(recBin, recBin...)
	msg, tail, err := Unmarshal(recBin, CreateSessionRequest)
	rec := msg.(*Recovery)
	assert.Equal(t, byte(255), rec.Value())
	assert.Equal(t, byte(0), rec.header.instance)
	assert.Equal(t, []byte{0x3, 0, 1, 0, 0xff}, tail)
	assert.Nil(t, err)
}
