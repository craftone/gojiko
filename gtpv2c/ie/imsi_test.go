package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImsi(t *testing.T) {
	imsi := NewImsi("819012345678", 0)
	assert.Equal(t, imsiNum, imsi.header.typeNum)
	assert.Equal(t, "819012345678", imsi.Value)
	assert.Equal(t, tbcd([]byte{0x18, 0x09, 0x21, 0x43, 0x65, 0x87}), imsi.tbcd)
	assert.Equal(t, byte(0), imsi.header.instance)

	imsiMin := NewImsi("123456", 1)
	assert.Equal(t, "123456", imsiMin.Value)
	assert.Equal(t, tbcd([]byte{0x21, 0x43, 0x65}), imsiMin.tbcd)
	assert.Equal(t, byte(1), imsiMin.header.instance)

	imsiMax := NewImsi("123456789012345", 2)
	assert.Equal(t, "123456789012345", imsiMax.Value)
	assert.Equal(t, byte(2), imsiMax.header.instance)
}

func TestImsi_Marshal(t *testing.T) {
	imsiBin := NewImsi("819012345678", 0).Marshal()
	assert.Equal(t, []byte{1, 0, 6, 0, 0x18, 0x09, 0x21, 0x43, 0x65, 0x87}, imsiBin)
}
