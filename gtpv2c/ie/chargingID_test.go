package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChargingID(t *testing.T) {
	_, err := NewChargingID(100, 0x00558776)
	assert.Error(t, err)
}

func TestChargingID_Marshal(t *testing.T) {
	chargingID, _ := NewChargingID(1, 0x00558776)
	chargingIDBin := chargingID.Marshal()
	assert.Equal(t, []byte{0x5e, 0, 4, 1, 0x00, 0x55, 0x87, 0x76}, chargingIDBin)
}

func TestUnmarshal_ChargingID(t *testing.T) {
	chargingID, _ := NewChargingID(1, 0x00558776)
	chargingIDBin := chargingID.Marshal()
	msg, tail, err := Unmarshal(chargingIDBin, CreateSessionRequest)
	chargingID = msg.(*ChargingID)
	assert.Equal(t, byte(1), chargingID.instance)
	assert.Equal(t, uint32(0x00558776), chargingID.Value())
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	chargingIDBin2 := chargingIDBin[0:4]
	_, _, err = Unmarshal(chargingIDBin2, CreateSessionRequest)
	assert.Error(t, err)

	//ensure no refference to the buffer
	copy(chargingIDBin, make([]byte, len(chargingIDBin)))
	assert.Equal(t, byte(1), chargingID.instance)
	assert.Equal(t, uint32(0x00558776), chargingID.Value())
}
