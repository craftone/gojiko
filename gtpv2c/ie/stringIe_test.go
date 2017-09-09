package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApn(t *testing.T) {
	apn := NewApn(0, "example.com")
	assert.Equal(t, apnNum, apn.typeNum)
	assert.Equal(t, "example.com", apn.Value)
}

func TestApn_Marshal(t *testing.T) {
	apn := NewApn(1, "example.com")
	apnBin := apn.Marshal()
	assert.Equal(t, []byte{0x47, 0, 0x0b, 1, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d}, apnBin)
}

func TestUnmarshal_apn(t *testing.T) {
	apnBin := NewApn(1, "example.com").Marshal()
	msg, tail, err := Unmarshal(apnBin)
	apn := msg.(*Apn)
	assert.Equal(t, byte(1), apn.instance)
	assert.Equal(t, "example.com", apn.Value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
