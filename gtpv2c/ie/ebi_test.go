package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEbi(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	assert.Equal(t, ebiNum, ebi.typeNum)
	assert.Equal(t, byte(5), ebi.value)
	ebi, _ = NewEbi(0, 15)
	assert.Equal(t, ebiNum, ebi.typeNum)
	assert.Equal(t, byte(15), ebi.value)

	_, err := NewEbi(0, 0)
	assert.Error(t, err)
	_, err = NewEbi(0, 1)
	assert.Error(t, err)
	_, err = NewEbi(0, 4)
	assert.Error(t, err)
	_, err = NewEbi(0, 16)
	assert.Error(t, err)
	_, err = NewEbi(0, 255)
	assert.Error(t, err)
}

func TestEbi_mebishal(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	ebiBin := ebi.Marshal()
	assert.Equal(t, []byte{0x49, 0, 1, 0, 5}, ebiBin)
}

func TestUnmebishal_Ebi(t *testing.T) {
	ebiOrg, _ := NewEbi(0, 5)
	ebiBin := ebiOrg.Marshal()
	msg, tail, err := Unmarshal(ebiBin, CreateSessionRequest)
	ebi := msg.(*Ebi)
	assert.Equal(t, byte(5), ebi.Value())
	assert.Equal(t, byte(0), ebi.header.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	//ensure no refference to the buffer
	copy(ebiBin, make([]byte, len(ebiBin)))
	assert.Equal(t, byte(5), ebi.Value())
	assert.Equal(t, byte(0), ebi.header.instance)
}
