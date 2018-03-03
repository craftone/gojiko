package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServingNetwork(t *testing.T) {
	sn, err := NewServingNetwork(0, "440", "10")
	assert.Equal(t, "440", sn.Mcc())
	assert.Equal(t, "10", sn.Mnc())
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
	msg, tail, err := Unmarshal(snBin, CreateSessionRequest)
	sn = msg.(*ServingNetwork)
	assert.Equal(t, byte(1), sn.instance)
	assert.Equal(t, "440", sn.Mcc())
	assert.Equal(t, "10", sn.Mnc())
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestServingNetwork_String(t *testing.T) {
	sn, _ := NewServingNetwork(1, "440", "10")
	assert.Equal(t, "MCC: 440, MNC: 10", sn.String())
	sn, _ = NewServingNetwork(1, "449", "111")
	assert.Equal(t, "MCC: 449, MNC: 111", sn.String())
}
