package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApnRestriction(t *testing.T) {
	ar, _ := NewApnRestriction(0, 1)
	assert.Equal(t, apnRestrictionNum, ar.typeNum)
	assert.Equal(t, byte(1), ar.Value)
}

func TestApnRestriction_marshal(t *testing.T) {
	ar, _ := NewApnRestriction(0, 0)
	arBin := ar.Marshal()
	assert.Equal(t, []byte{0x7f, 0, 1, 0, 0}, arBin)
}

func TestUnmarshal_ApnRestriction(t *testing.T) {
	arOrg, _ := NewApnRestriction(0, 255)
	arBin := arOrg.Marshal()
	msg, tail, err := Unmarshal(arBin, CreateSessionRequest)
	ar := msg.(*ApnRestriction)
	assert.Equal(t, byte(255), ar.Value)
	assert.Equal(t, byte(0), ar.header.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
