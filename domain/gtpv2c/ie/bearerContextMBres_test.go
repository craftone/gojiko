package ie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBearerContextModifiedWithinMBRes(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	cause, _ := NewCause(0, CauseRequestAccepted, false, false, false, nil)
	chargingID, _ := NewChargingID(0, 0x00558776)

	bcTBMwMBResArg := BearerContextToBeModifiedWithinMBResArg{
		Ebi:        ebi,
		Cause:      cause,
		ChargingID: chargingID,
	}
	bcTBMwMBRes, err := NewBearerContextToBeModifiedWithinMBRes(bcTBMwMBResArg)

	assert.Equal(t, bearerContextNum, bcTBMwMBRes.typeNum)
	assert.Nil(t, err)

	// check Mandatories

	bcTBMwMBResArg.Ebi = nil
	_, err = NewBearerContextToBeModifiedWithinMBRes(bcTBMwMBResArg)
	assert.Error(t, err)
	bcTBMwMBResArg.Ebi = ebi

	bcTBMwMBResArg.Cause = nil
	_, err = NewBearerContextToBeModifiedWithinMBRes(bcTBMwMBResArg)
	assert.Error(t, err)
	bcTBMwMBResArg.Cause = cause

	bcTBMwMBResArg.ChargingID = nil
	_, err = NewBearerContextToBeModifiedWithinMBRes(bcTBMwMBResArg)
	assert.NoError(t, err)
}

func TestBearerContextModifiedWithinMBRes_Marshal(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	cause, _ := NewCause(0, CauseRequestAccepted, false, false, false, nil)
	chargingID, _ := NewChargingID(0, 0x00558776)

	bcTBMwMBResArg := BearerContextToBeModifiedWithinMBResArg{
		Ebi:        ebi,
		Cause:      cause,
		ChargingID: chargingID,
	}
	bcTBMwMBRes, _ := NewBearerContextToBeModifiedWithinMBRes(bcTBMwMBResArg)
	bcTBMwMBResBin := bcTBMwMBRes.Marshal()

	assert.Equal(t, []byte{
		0x5D, 0, 0x13, 0, //Header
		0x49, 0, 1, 0, 5, //EBI
		0x02, 0, 2, 0, 0x10, 0, //Cause
		0x5E, 0, 4, 0, 0, 0x55, 0x87, 0x76, //Charging ID
	}, bcTBMwMBResBin)
}

func TestUnmarshal_BearerContextModifiedWithinMBRes(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	cause, _ := NewCause(0, CauseRequestAccepted, false, false, false, nil)
	chargingID, _ := NewChargingID(0, 0x00558776)

	bcTBMwMBResArg := BearerContextToBeModifiedWithinMBResArg{
		Ebi:        ebi,
		Cause:      cause,
		ChargingID: chargingID,
	}
	bcTBMwMBRes, _ := NewBearerContextToBeModifiedWithinMBRes(bcTBMwMBResArg)
	bcTBMwMBResBin := bcTBMwMBRes.Marshal()

	msg, tail, err := Unmarshal(bcTBMwMBResBin, ModifyBearerResponse)
	bcTBMwMBRes = msg.(*BearerContextToBeModifiedWithinMBRes)
	assert.Equal(t, bearerContextNum, bcTBMwMBRes.typeNum)
	assert.Equal(t, ebi, bcTBMwMBRes.Ebi())
	assert.Equal(t, cause, bcTBMwMBRes.Cause())
	assert.Equal(t, chargingID, bcTBMwMBRes.ChargingID())
	assert.Nil(t, err)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
