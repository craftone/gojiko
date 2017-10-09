package ie

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBearerContextCreatedWithinCSRes(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	cause, _ := NewCause(0, CauseRequestAccepted, false, false, false, nil)
	pgwDataFteid, _ := NewFteid(2, net.IPv4(1, 2, 3, 4), nil, S5S8PgwGtpUIf, 0x86440002)
	chargingID, _ := NewChargingID(0, 0x00558776)

	bcCwCSResArg := BearerContextCreatedWithinCSResArg{
		Ebi:          ebi,
		Cause:        cause,
		PgwDataFteid: pgwDataFteid,
		ChargingID:   chargingID,
	}
	bcCwCSRes, err := NewBearerContextCreatedWithinCSRes(bcCwCSResArg)

	assert.Equal(t, bearerContextNum, bcCwCSRes.typeNum)
	assert.Nil(t, err)

	bcCwCSResArg.Ebi = nil
	_, err = NewBearerContextCreatedWithinCSRes(bcCwCSResArg)
	assert.Error(t, err)

	bcCwCSResArg.Ebi = ebi
	bcCwCSResArg.Cause = nil
	_, err = NewBearerContextCreatedWithinCSRes(bcCwCSResArg)
	assert.Error(t, err)

	bcCwCSResArg.Cause = cause
	bcCwCSResArg.PgwDataFteid = nil
	_, err = NewBearerContextCreatedWithinCSRes(bcCwCSResArg)
	assert.Error(t, err)

	bcCwCSResArg.PgwDataFteid = pgwDataFteid
	bcCwCSResArg.ChargingID = nil
	_, err = NewBearerContextCreatedWithinCSRes(bcCwCSResArg)
	assert.NoError(t, err)
}

func TestBearerContextCreatedWithinCSRes_Marshal(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	cause, _ := NewCause(0, CauseRequestAccepted, false, false, false, nil)
	pgwDataFteid, _ := NewFteid(2, net.IPv4(1, 2, 3, 4), nil, S5S8PgwGtpUIf, 0x86440002)
	chargingID, _ := NewChargingID(0, 0x00558776)

	bcCwCSResArg := BearerContextCreatedWithinCSResArg{
		Ebi:          ebi,
		Cause:        cause,
		PgwDataFteid: pgwDataFteid,
		ChargingID:   chargingID,
	}
	bcCwCSRes, _ := NewBearerContextCreatedWithinCSRes(bcCwCSResArg)
	bcCwCSResBin := bcCwCSRes.Marshal()

	assert.Equal(t, []byte{
		0x5D, 0, 0x20, 0, //Header
		0x49, 0, 1, 0, 5, //EBI
		0x02, 0, 2, 0, 0x10, 0, //Cause
		0x57, 0, 9, 2, 0x85, 0x86, 0x44, 0x00, 0x02, 1, 2, 3, 4, //PgwDataFTEID
		0x5E, 0, 4, 0, 0, 0x55, 0x87, 0x76, //BearerQoS
	}, bcCwCSResBin)
}

func TestUnmarshal_BearerContextCreatedWithinCSRes(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	cause, _ := NewCause(0, CauseRequestAccepted, false, false, false, nil)
	pgwDataFteid, _ := NewFteid(2, net.IPv4(1, 2, 3, 4), nil, S5S8PgwGtpUIf, 0x86440002)
	chargingID, _ := NewChargingID(0, 0x00558776)

	bcCwCSResArg := BearerContextCreatedWithinCSResArg{
		Ebi:          ebi,
		Cause:        cause,
		PgwDataFteid: pgwDataFteid,
		ChargingID:   chargingID,
	}
	bcCwCSRes, _ := NewBearerContextCreatedWithinCSRes(bcCwCSResArg)
	bcCwCSResBin := bcCwCSRes.Marshal()

	msg, tail, err := Unmarshal(bcCwCSResBin, CreateSessionResponse)
	bcCwCSRes = msg.(*BearerContextCreatedWithinCSRes)
	assert.Equal(t, bearerContextNum, bcCwCSRes.typeNum)
	assert.Equal(t, ebi, bcCwCSRes.ebi)
	assert.Equal(t, cause, bcCwCSRes.cause)
	assert.Equal(t, pgwDataFteid, bcCwCSRes.pgwDataFteid)
	assert.Equal(t, chargingID, bcCwCSRes.chargingID)
	assert.Nil(t, err)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
