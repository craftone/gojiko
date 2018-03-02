package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApn(t *testing.T) {
	apn, _ := NewApn(0, "apn-example.com")
	assert.Equal(t, apnNum, apn.typeNum)
	assert.Equal(t, "apn-example.com", apn.String())

	errorApns := []string{
		"", "apn-example.", "apn-example.c*m", "-apn-example.com", "apn-example-.com",
	}
	for _, errorApn := range errorApns {
		_, err := NewApn(0, errorApn)
		assert.Error(t, err)
	}
}

func TestApn_Marshal(t *testing.T) {
	apn, _ := NewApn(1, "example.com")
	apnBin := apn.Marshal()
	assert.Equal(t, []byte{0x47, 0, 0x0c, 1, 7, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 3, 0x63, 0x6f, 0x6d}, apnBin)
}

func TestUnmarshal_apn(t *testing.T) {
	apnOrg, _ := NewApn(1, "example.com")
	apnBin := apnOrg.Marshal()
	msg, tail, err := Unmarshal(apnBin, CreateSessionRequest)
	apn := msg.(*Apn)
	assert.Equal(t, byte(1), apn.instance)
	assert.Equal(t, "example.com", apn.String())
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	//ensure no refference to the buffer
	copy(apnBin, make([]byte, len(apnBin)))
	assert.Equal(t, byte(1), apn.instance)
	assert.Equal(t, "example.com", apn.String())
}
