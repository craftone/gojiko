package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCause(t *testing.T) {
	cause, err := NewCause(0, 1, true, false, true, nil)
	assert.Equal(t, causeNum, cause.header.typeNum)
	assert.Equal(t, byte(1), cause.Value)
	assert.Equal(t, true, cause.Pce)
	assert.Equal(t, false, cause.Bce)
	assert.Equal(t, true, cause.Cs)
	assert.Nil(t, cause.OffendingIe)
	assert.Nil(t, err)

	offendingIe := &header{2, 0, 3}
	cause, _ = NewCause(1, 2, false, true, false, offendingIe)
	assert.Equal(t, byte(2), cause.Value)
	assert.Equal(t, false, cause.Pce)
	assert.Equal(t, true, cause.Bce)
	assert.Equal(t, false, cause.Cs)
	assert.Equal(t, ieTypeNum(2), cause.OffendingIe.typeNum)
	assert.Equal(t, byte(3), cause.OffendingIe.instance)
}

func TestCause_Marshal(t *testing.T) {
	cause, _ := NewCause(1, 2, true, false, true, nil)
	causeBin := cause.Marshal()
	assert.Equal(t, []byte{2, 0, 2, 1, 2, 5}, causeBin)

	causeOff, _ := NewCause(1, 2, true, true, true, &header{2, 0, 3})
	causeOffBin := causeOff.Marshal()
	assert.Equal(t, []byte{2, 0, 6, 1, 2, 7, 2, 0, 0, 3}, causeOffBin)
}

func TestUnmarshal_cause(t *testing.T) {
	causeOrg, _ := NewCause(1, 2, true, false, true, nil)
	causeBin := causeOrg.Marshal()
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

	causeOffOrg, _ := NewCause(1, 2, true, false, true, &header{2, 0, 3})
	causeOffBin := causeOffOrg.Marshal()
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
