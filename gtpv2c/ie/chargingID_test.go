package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChargingID(t *testing.T) {
	_, err := NewChargingID(100, 0x5605238)
	assert.Error(t, err)
}

func TestChargingID_Marshal(t *testing.T) {
	chargingID, _ := NewChargingID(1, 0x5605238)
	chargingIDBin := chargingID.Marshal()
	assert.Equal(t, []byte{0x5e, 0, 4, 1, 0x5, 0x60, 0x52, 0x38}, chargingIDBin)
}

func TestUnmarshal_ChargingID(t *testing.T) {
	chargingID, _ := NewChargingID(1, 0x5605238)
	chargingIDBin := chargingID.Marshal()
	msg, tail, err := Unmarshal(chargingIDBin, MsToNetwork)
	chargingID = msg.(*ChargingID)
	assert.Equal(t, byte(1), chargingID.instance)
	assert.Equal(t, uint32(0x5605238), chargingID.Value)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	chargingIDBin2 := chargingIDBin[0:4]
	msg, tail, err = Unmarshal(chargingIDBin2, MsToNetwork)
	assert.Error(t, err)
}
