package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRatType(t *testing.T) {
	rt, _ := NewRatType(0, 1)
	assert.Equal(t, ratTypeNum, rt.header.typeNum)
	assert.Equal(t, byte(1), rt.Value())

	_, err := NewRatType(0x0f, 1)
	assert.Nil(t, err)

	_, err = NewRatType(0x10, 1)
	assert.Error(t, err)
	_, err = NewRatType(0xff, 1)
	assert.Error(t, err)
}

func TestRatType_marshal(t *testing.T) {
	rt, _ := NewRatType(0, 0)
	rtBin := rt.Marshal()
	assert.Equal(t, []byte{0x52, 0, 1, 0, 0}, rtBin)

	rt, _ = NewRatType(0xf, 255)
	rtBin = rt.Marshal()
	assert.Equal(t, []byte{0x52, 0, 1, 0xf, 255}, rtBin)
}

func TestUnmarshal_RatType(t *testing.T) {
	rtOrg, _ := NewRatType(0, 255)
	rtBin := rtOrg.Marshal()
	msg, tail, err := Unmarshal(rtBin, CreateSessionRequest)
	rt := msg.(*RatType)
	assert.Equal(t, byte(255), rt.Value())
	assert.Equal(t, byte(0), rt.header.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestRatType_String(t *testing.T) {
	rt, _ := NewRatType(0, 0)
	assert.Equal(t, "<reserved> (0)", rt.String())
	rt, _ = NewRatType(0, 1)
	assert.Equal(t, "UTRAN (1)", rt.String())
	rt, _ = NewRatType(0, 6)
	assert.Equal(t, "EUTRAN (6)", rt.String())
	rt, _ = NewRatType(0, 7)
	assert.Equal(t, "Virtual (7)", rt.String())
	rt, _ = NewRatType(0, 255)
	assert.Equal(t, "<reserved> (255)", rt.String())
}
