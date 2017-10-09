package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMsisdn(t *testing.T) {
	msisdn, _ := NewMsisdn(0, "819012345678")
	assert.Equal(t, msisdnNum, msisdn.typeNum)
	assert.Equal(t, "819012345678", msisdn.Value())
	assert.Equal(t, tbcd([]byte{0x18, 0x09, 0x21, 0x43, 0x65, 0x87}), msisdn.tbcd)
	assert.Equal(t, byte(0), msisdn.instance)

	msisdnMin, _ := NewMsisdn(1, "123456")
	assert.Equal(t, "123456", msisdnMin.Value())
	assert.Equal(t, tbcd([]byte{0x21, 0x43, 0x65}), msisdnMin.tbcd)
	assert.Equal(t, byte(1), msisdnMin.instance)

	msisdnMax, _ := NewMsisdn(2, "123456789012345")
	assert.Equal(t, "123456789012345", msisdnMax.Value())
	assert.Equal(t, byte(2), msisdnMax.instance)

	// shorter than min error
	_, err := NewMsisdn(0, "12345")
	assert.Error(t, err)

	// longer than max error
	_, err = NewMsisdn(0, "1234567890123456")
	assert.Error(t, err)
}

func Testmsisdn_Marshal(t *testing.T) {
	msisdn, _ := NewMsisdn(0, "819012345678")
	msisdnBin := msisdn.Marshal()
	assert.Equal(t, []byte{1, 0, 6, 0, 0x18, 0x09, 0x21, 0x43, 0x65, 0x87}, msisdnBin)
}

func TestUnmarshal_msisdn(t *testing.T) {
	msisdnOrg, _ := NewMsisdn(1, "819012345678")
	msisdnBin := msisdnOrg.Marshal()
	msg, tail, err := Unmarshal(msisdnBin, CreateSessionRequest)
	msisdn := msg.(*Msisdn)
	assert.Equal(t, "819012345678", msisdn.Value())
	assert.Equal(t, byte(1), msisdn.instance)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}

func TestUnmarshal_msisdnWithTail(t *testing.T) {
	msisdnOrg, _ := NewMsisdn(1, "819012345678")
	msisdnBin := msisdnOrg.Marshal()
	msisdnBin = append(msisdnBin, msisdnBin...)
	msg, tail, err := Unmarshal(msisdnBin, CreateSessionRequest)
	msisdn := msg.(*Msisdn)
	assert.Equal(t, "819012345678", msisdn.Value())
	assert.Equal(t, byte(1), msisdn.instance)
	assert.Equal(t, []byte{0x4c, 0, 6, 1, 0x18, 0x09, 0x21, 0x43, 0x65, 0x87}, tail)
	assert.Nil(t, err)
}
