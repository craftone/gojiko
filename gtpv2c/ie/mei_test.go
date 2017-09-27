package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMei(t *testing.T) {
	meiMin, _ := NewMei(1, "012345678901234")
	assert.Equal(t, "012345678901234", meiMin.Value)
	assert.Equal(t, tbcd([]byte{0x10, 0x32, 0x54, 0x76, 0x98, 0x10, 0x32, 0xf4}), meiMin.tbcd)
	assert.Equal(t, byte(1), meiMin.instance)

	meiMax, _ := NewMei(2, "0123456789012345")
	assert.Equal(t, "0123456789012345", meiMax.Value)
	assert.Equal(t, tbcd([]byte{0x10, 0x32, 0x54, 0x76, 0x98, 0x10, 0x32, 0x54}), meiMax.tbcd)
	assert.Equal(t, byte(2), meiMax.instance)

	// shorter than min error
	_, err := NewMei(0, "01234567890123")
	assert.Error(t, err)

	// longer than max error
	_, err = NewMei(0, "01234567890123456")
	assert.Error(t, err)
}

func TestMei_Marshal(t *testing.T) {
	mei, _ := NewMei(0, "0123456789012345")
	meiBin := mei.Marshal()
	assert.Equal(t, []byte{0x4b, 0, 8, 0, 0x10, 0x32, 0x54, 0x76, 0x98, 0x10, 0x32, 0x54}, meiBin)
}

func TestUnmarshal_mei(t *testing.T) {
	meiOrg, _ := NewMei(1, "0123456789012345")
	meiBin := meiOrg.Marshal()
	msg, tail, err := Unmarshal(meiBin, MsToNetwork)
	mei := msg.(*Mei)
	assert.Equal(t, "0123456789012345", mei.Value)
	assert.Equal(t, byte(1), mei.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_meiWithTail(t *testing.T) {
	meiOrg, _ := NewMei(1, "0123456789012345")
	meiBin := meiOrg.Marshal()
	meiBin = append(meiBin, meiBin...)
	msg, tail, err := Unmarshal(meiBin, MsToNetwork)
	mei := msg.(*Mei)
	assert.Equal(t, "0123456789012345", mei.Value)
	assert.Equal(t, byte(1), mei.instance)
	assert.Equal(t, []byte{0x4b, 0, 8, 1, 0x10, 0x32, 0x54, 0x76, 0x98, 0x10, 0x32, 0x54}, tail)
	assert.Nil(t, err)
}
