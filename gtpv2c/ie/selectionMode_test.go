package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSelectionMode(t *testing.T) {
	sm, _ := NewSelectionMode(0, 1)
	assert.Equal(t, selectionModeNum, sm.header.typeNum)
	assert.Equal(t, byte(1), sm.Value)

	_, err := NewSelectionMode(1, 4)
	assert.Error(t, err)
}

func TestSelectionMode_marshal(t *testing.T) {
	sm, _ := NewSelectionMode(0, 0)
	smBin := sm.Marshal()
	assert.Equal(t, []byte{0x80, 0, 1, 0, 0}, smBin)

	sm, _ = NewSelectionMode(0xf, 3)
	smBin = sm.Marshal()
	assert.Equal(t, []byte{0x80, 0, 1, 0xf, 3}, smBin)
}

func TestUnmarshal_selectionMode(t *testing.T) {
	sm, _ := NewSelectionMode(0, 1)
	smBin := sm.Marshal()
	msg, tail, err := Unmarshal(smBin, CreateSessionRequest)
	sm = msg.(*SelectionMode)
	assert.Equal(t, byte(0), sm.header.instance)
	assert.Equal(t, byte(1), sm.Value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
