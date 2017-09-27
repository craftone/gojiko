package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServingNetwork(t *testing.T) {
	sn, err := NewServingNetwork(0, "440", "10")
	assert.Equal(t, "440", sn.Mcc)
	assert.Equal(t, "10", sn.Mnc)
	assert.Nil(t, err)

	_, err = NewServingNetwork(0, "4400", "10")
	assert.Error(t, err)
}

func TestServingNetwork_Marshal(t *testing.T) {
	sn, _ := NewServingNetwork(1, "440", "10")
	snBin := sn.Marshal()
	assert.Equal(t, []byte{0x53, 0, 3, 1, 0x44, 0xf0, 0x01}, snBin)
}

func TestUnmarshal_servingNetwork(t *testing.T) {
	sn, _ := NewServingNetwork(1, "440", "10")
	snBin := sn.Marshal()
	msg, tail, err := Unmarshal(snBin, MsToNetwork)
	sn = msg.(*ServingNetwork)
	assert.Equal(t, byte(1), sn.instance)
	assert.Equal(t, "440", sn.Mcc)
	assert.Equal(t, "10", sn.Mnc)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
