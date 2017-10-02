package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBearerQoS(t *testing.T) {
	bearerQoSArg := BearerQoSArg{
		Pci:         true,
		Pl:          15,
		Pvi:         false,
		Label:       9,
		UplinkMBR:   0,
		DownlinkMBR: 0,
		UplinkGBR:   0,
		DownlinkGBR: 0,
	}
	bearerQoS, err := NewBearerQoS(1, bearerQoSArg)
	assert.Equal(t, bearerQoSNum, bearerQoS.typeNum)
	assert.Equal(t, true, bearerQoS.Pci)
	assert.Equal(t, byte(15), bearerQoS.Pl)
	assert.Equal(t, false, bearerQoS.Pvi)
	assert.Equal(t, byte(9), bearerQoS.Label)
	assert.Nil(t, err)

	bearerQoSArg.Pl = 16
	_, err = NewBearerQoS(1, bearerQoSArg)
	assert.Error(t, err)

	bearerQoSArg.Pl = 15
	bearerQoSArg.UplinkMBR = MaxBitrate
	_, err = NewBearerQoS(1, bearerQoSArg)
	assert.NoError(t, err)

	bearerQoSArg.UplinkMBR = MaxBitrate + 1
	_, err = NewBearerQoS(1, bearerQoSArg)
	assert.Error(t, err)
}

func TestBearerQoS_Marshal(t *testing.T) {
	bearerQoS, _ := NewBearerQoS(1, BearerQoSArg{
		Pci:         true,
		Pl:          15,
		Pvi:         false,
		Label:       9,
		UplinkMBR:   0,
		DownlinkMBR: 0,
		UplinkGBR:   0,
		DownlinkGBR: 0,
	})
	bearerQoSBin := bearerQoS.Marshal()
	assert.Equal(t, []byte{
		0x50, 0, 0x16, 1,
		0x7c, 0x09,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}, bearerQoSBin)
}

func TestUnmarshal_BearerQoS(t *testing.T) {
	bearerQoS, _ := NewBearerQoS(1, BearerQoSArg{
		Pci:         true,
		Pl:          15,
		Pvi:         false,
		Label:       9,
		UplinkMBR:   0,
		DownlinkMBR: 0,
		UplinkGBR:   0,
		DownlinkGBR: 0,
	})
	bearerQoSBin := bearerQoS.Marshal()
	msg, tail, err := Unmarshal(bearerQoSBin, CreateSessionRequest)
	bearerQoS = msg.(*BearerQoS)
	assert.Equal(t, bearerQoSNum, bearerQoS.typeNum)
	assert.Equal(t, true, bearerQoS.Pci)
	assert.Equal(t, byte(15), bearerQoS.Pl)
	assert.Equal(t, false, bearerQoS.Pvi)
	assert.Equal(t, byte(9), bearerQoS.Label)
	assert.Nil(t, err)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
