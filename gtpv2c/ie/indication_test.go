package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIndication(t *testing.T) {
	ida := IndicationArg{
		true, false, true, false,
		true, false, true, false,
		true, false, true, false,
		true, false, true, false,
		true, false,
	}
	idc, err := NewIndication(1, ida)
	assert.Equal(t, indicationNum, idc.typeNum)
	assert.Equal(t, true, idc.DAF)
	assert.Equal(t, false, idc.DTF)
	assert.Equal(t, true, idc.ISRAU)
	assert.Equal(t, false, idc.CCRSI)
	assert.Nil(t, err)
}

func TestIndication_Marshal(t *testing.T) {
	ida := IndicationArg{
		true, false, true, false,
		true, false, true, false,
		true, false, true, false,
		true, false, true, false,
		true, false,
	}
	idc, _ := NewIndication(1, ida)
	idcBin := idc.Marshal()
	assert.Equal(t, []byte{0x4d, 0, 3, 1, 0xaa, 0xaa, 2}, idcBin)
}

func TestUnmarshal_indication(t *testing.T) {
	ida := IndicationArg{
		true, false, true, false,
		true, false, true, false,
		true, false, true, false,
		true, false, true, false,
		true, false,
	}
	idc, _ := NewIndication(1, ida)
	idcBin := idc.Marshal()
	msg, tail, err := Unmarshal(idcBin, MsToNetwork)
	idc = msg.(*Indication)
	assert.Equal(t, ida, idc.IndicationArg)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
