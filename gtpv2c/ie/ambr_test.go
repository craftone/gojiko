package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAmbr(t *testing.T) {
	ambr, _ := NewAmbr(0, 4294967, 4294960)
	assert.Equal(t, ambrNum, ambr.typeNum)
	assert.Equal(t, uint32(4294967), ambr.UplinkKbps)
	assert.Equal(t, uint32(4294960), ambr.DownlinkKbps)
}

func TestAmbr_mambrshal(t *testing.T) {
	ambr, _ := NewAmbr(0, 4294967, 4294960)
	ambrBin := ambr.Marshal()
	assert.Equal(t, []byte{0x48, 0, 8, 0, 0, 0x41, 0x89, 0x37, 0, 0x41, 0x89, 0x30}, ambrBin)
}

func TestUnmambrshal_Ambr(t *testing.T) {
	ambrOrg, _ := NewAmbr(0, 4294967, 4294960)
	ambrBin := ambrOrg.Marshal()
	msg, tail, err := Unmarshal(ambrBin, CreateSessionRequest)
	ambr := msg.(*Ambr)
	assert.Equal(t, byte(0), ambr.header.instance)
	assert.Equal(t, uint32(4294967), ambr.UplinkKbps)
	assert.Equal(t, uint32(4294960), ambr.DownlinkKbps)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
