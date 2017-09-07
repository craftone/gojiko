package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCause(t *testing.T) {
	cause := NewCause(0, 1, true, false, true, nil)
	assert.Equal(t, causeNum, cause.header.typeNum)
	assert.Equal(t, byte(1), cause.Value)
	assert.Equal(t, true, cause.Pce)
	assert.Equal(t, false, cause.Bce)
	assert.Equal(t, true, cause.Cs)
	assert.Nil(t, cause.OffendingIe)

	offendingIe := &header{2, 0, 3}
	causeOff := NewCause(1, 2, false, true, false, offendingIe)
	assert.Equal(t, byte(2), causeOff.Value)
	assert.Equal(t, false, causeOff.Pce)
	assert.Equal(t, true, causeOff.Bce)
	assert.Equal(t, false, causeOff.Cs)
	assert.Equal(t, ieTypeNum(2), causeOff.OffendingIe.typeNum)
	assert.Equal(t, byte(3), causeOff.OffendingIe.instance)
}

func TestCause_Marshal(t *testing.T) {
	causeBin := NewCause(1, 2, true, false, true, nil).Marshal()
	assert.Equal(t, []byte{2, 0, 2, 1, 2, 5}, causeBin)

	causeOffBin := NewCause(1, 2, true, true, true, &header{2, 0, 3}).Marshal()
	assert.Equal(t, []byte{2, 0, 6, 1, 2, 7, 2, 0, 0, 3}, causeOffBin)
}

func TestUnmarshal_cause(t *testing.T) {
	causeBin := NewCause(1, 2, true, false, true, nil).Marshal()
	msg, tail, err := Unmarshal(causeBin)
	cause := msg.(*Cause)
	assert.Equal(t, byte(1), cause.instance)
	assert.Equal(t, byte(2), cause.Value)
	assert.Equal(t, true, cause.Pce)
	assert.Equal(t, false, cause.Bce)
	assert.Equal(t, true, cause.Cs)
	assert.Nil(t, cause.OffendingIe)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	causeOffBin := NewCause(1, 2, true, false, true, &header{2, 0, 3}).Marshal()
	msg, tail, err = Unmarshal(causeOffBin)
	cause = msg.(*Cause)
	assert.Equal(t, byte(1), cause.instance)
	assert.Equal(t, byte(2), cause.Value)
	assert.Equal(t, true, cause.Pce)
	assert.Equal(t, false, cause.Bce)
	assert.Equal(t, true, cause.Cs)
	assert.Equal(t, ieTypeNum(2), cause.OffendingIe.typeNum)
	assert.Equal(t, byte(3), cause.OffendingIe.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
