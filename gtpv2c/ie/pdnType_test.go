package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPdnType(t *testing.T) {
	pt, _ := NewPdnType(0, PdnTypeIPv4)
	assert.Equal(t, pdnTypeNum, pt.typeNum)
	assert.Equal(t, PdnTypeIPv4, pt.Value)

	pt, _ = NewPdnType(0, PdnTypeIPv6)
	assert.Equal(t, pdnTypeNum, pt.typeNum)
	assert.Equal(t, PdnTypeIPv6, pt.Value)

	pt, _ = NewPdnType(0, PdnTypeIPv4v6)
	assert.Equal(t, pdnTypeNum, pt.typeNum)
	assert.Equal(t, PdnTypeIPv4v6, pt.Value)

	_, err := NewPdnType(1, 0)
	assert.Error(t, err)

	_, err = NewPdnType(1, 4)
	assert.Error(t, err)
}

func TestPdnType_marshal(t *testing.T) {
	pt, _ := NewPdnType(0, PdnTypeIPv4)
	ptBin := pt.Marshal()
	assert.Equal(t, []byte{0x63, 0, 1, 0, 1}, ptBin)
}

func TestUnmarshal_pdnType(t *testing.T) {
	pt, _ := NewPdnType(0, PdnTypeIPv4)
	ptBin := pt.Marshal()
	msg, tail, err := Unmarshal(ptBin, CreateSessionRequest)
	pt = msg.(*PdnType)
	assert.Equal(t, byte(0), pt.instance)
	assert.Equal(t, PdnTypeIPv4, pt.Value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
