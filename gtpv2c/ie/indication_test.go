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
	assert.Equal(t, true, idc.DAF())
	assert.Equal(t, false, idc.DTF())
	assert.Equal(t, true, idc.ISRAU())
	assert.Equal(t, false, idc.CCRSI())
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
	assert.Equal(t, []byte{0x4d, 0, 7, 1, 0xaa, 0xaa, 2, 0, 0, 0, 0}, idcBin)
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
	msg, tail, err := Unmarshal(idcBin, CreateSessionRequest)
	idc = msg.(*Indication)
	assert.Equal(t, true, idc.DAF())
	assert.Equal(t, false, idc.DTF())
	assert.Equal(t, true, idc.HI())
	assert.Equal(t, false, idc.DFI())
	assert.Equal(t, true, idc.OI())
	assert.Equal(t, false, idc.ISRSI())
	assert.Equal(t, true, idc.ISRAI())
	assert.Equal(t, false, idc.SGWCI())
	assert.Equal(t, true, idc.SQCI())
	assert.Equal(t, false, idc.UIMSI())
	assert.Equal(t, true, idc.CFSI())
	assert.Equal(t, false, idc.CRSI())
	assert.Equal(t, true, idc.P())
	assert.Equal(t, false, idc.PT())
	assert.Equal(t, true, idc.SI())
	assert.Equal(t, false, idc.MSV())
	assert.Equal(t, true, idc.ISRAU())
	assert.Equal(t, false, idc.CCRSI())
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
