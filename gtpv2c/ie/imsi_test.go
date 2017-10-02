package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImsi(t *testing.T) {
	imsi, _ := NewImsi(0, "819012345678")
	assert.Equal(t, imsiNum, imsi.typeNum)
	assert.Equal(t, "819012345678", imsi.Value)
	assert.Equal(t, tbcd([]byte{0x18, 0x09, 0x21, 0x43, 0x65, 0x87}), imsi.tbcd)
	assert.Equal(t, byte(0), imsi.instance)

	imsiMin, _ := NewImsi(1, "123456")
	assert.Equal(t, "123456", imsiMin.Value)
	assert.Equal(t, tbcd([]byte{0x21, 0x43, 0x65}), imsiMin.tbcd)
	assert.Equal(t, byte(1), imsiMin.instance)

	imsiMax, _ := NewImsi(2, "123456789012345")
	assert.Equal(t, "123456789012345", imsiMax.Value)
	assert.Equal(t, byte(2), imsiMax.instance)

	// shorter than min error
	_, err := NewImsi(0, "12345")
	assert.Error(t, err)

	// longer than max error
	_, err = NewImsi(0, "1234567890123456")
	assert.Error(t, err)
}

func TestImsi_Marshal(t *testing.T) {
	imsi, _ := NewImsi(0, "819012345678")
	imsiBin := imsi.Marshal()
	assert.Equal(t, []byte{1, 0, 6, 0, 0x18, 0x09, 0x21, 0x43, 0x65, 0x87}, imsiBin)
}

func TestUnmarshal_imsi(t *testing.T) {
	imsiOrg, _ := NewImsi(1, "819012345678")
	imsiBin := imsiOrg.Marshal()
	msg, tail, err := Unmarshal(imsiBin, CreateSessionRequest)
	imsi := msg.(*Imsi)
	assert.Equal(t, "819012345678", imsi.Value)
	assert.Equal(t, byte(1), imsi.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_imsiWithTail(t *testing.T) {
	imsiOrg, _ := NewImsi(1, "819012345678")
	imsiBin := imsiOrg.Marshal()
	imsiBin = append(imsiBin, imsiBin...)
	msg, tail, err := Unmarshal(imsiBin, CreateSessionRequest)
	imsi := msg.(*Imsi)
	assert.Equal(t, "819012345678", imsi.Value)
	assert.Equal(t, byte(1), imsi.instance)
	assert.Equal(t, []byte{1, 0, 6, 1, 0x18, 0x09, 0x21, 0x43, 0x65, 0x87}, tail)
	assert.Nil(t, err)
}
