package ie

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBearerContextToBeCreatedWithinCSReq(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	bearerQoS, _ := NewBearerQoS(0, BearerQoSArg{
		Pci:         true,
		Pl:          15,
		Pvi:         false,
		Label:       9,
		UplinkMBR:   0,
		DownlinkMBR: 0,
		UplinkGBR:   0,
		DownlinkGBR: 0,
	})
	sgwDataFteid, _ := NewFteid(2, net.IPv4(1, 2, 3, 4), nil, S5S8SgwGtpUIf, 0x4075C31B)

	bcTBCwCSReqArg := BearerContextToBeCreatedWithinCSReqArg{
		Ebi:          ebi,
		BearerQoS:    bearerQoS,
		SgwDataFteid: sgwDataFteid,
	}
	bcTBCwCSReq, err := NewBearerContextToBeCreatedWithinCSReq(bcTBCwCSReqArg)

	assert.Equal(t, bearerContextNum, bcTBCwCSReq.typeNum)
	assert.Nil(t, err)

	bcTBCwCSReqArg.Ebi = nil
	_, err = NewBearerContextToBeCreatedWithinCSReq(bcTBCwCSReqArg)
	assert.Error(t, err)

	bcTBCwCSReqArg.Ebi = ebi
	bcTBCwCSReqArg.BearerQoS = nil
	_, err = NewBearerContextToBeCreatedWithinCSReq(bcTBCwCSReqArg)
	assert.Error(t, err)

	bcTBCwCSReqArg.BearerQoS = bearerQoS
	bcTBCwCSReqArg.SgwDataFteid = nil
	_, err = NewBearerContextToBeCreatedWithinCSReq(bcTBCwCSReqArg)
	assert.Error(t, err)
}

func TestBearerContextToBeCreatedWithinCSReq_Marshal(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	bearerQoS, _ := NewBearerQoS(0, BearerQoSArg{
		Pci:         true,
		Pl:          15,
		Pvi:         false,
		Label:       9,
		UplinkMBR:   0,
		DownlinkMBR: 0,
		UplinkGBR:   0,
		DownlinkGBR: 0,
	})
	sgwDataFteid, _ := NewFteid(0, net.IPv4(1, 2, 3, 4), nil, S5S8SgwGtpUIf, 0x4075C31B)

	bcTBCwCSReqArg := BearerContextToBeCreatedWithinCSReqArg{
		Ebi:          ebi,
		BearerQoS:    bearerQoS,
		SgwDataFteid: sgwDataFteid,
	}
	bcTBCwCSReq, _ := NewBearerContextToBeCreatedWithinCSReq(bcTBCwCSReqArg)
	bcTBCwCSReqBin := bcTBCwCSReq.Marshal()

	assert.Equal(t, []byte{
		0x5D, 0, 0x2c, 0, //Header
		0x49, 0, 1, 0, 5, //EBI
		0x50, 0, 0x16, 0, 0x7c, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, //BearerQoS
		0x57, 0, 9, 2, 0x84, 0x40, 0x75, 0xc3, 0x1b, 1, 2, 3, 4, //SgwDataFTEID
	}, bcTBCwCSReqBin)
}

func TestUnmarshal_BearerContextToBeCreatedWithinCSReq(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	bearerQoS, _ := NewBearerQoS(0, BearerQoSArg{
		Pci:         true,
		Pl:          15,
		Pvi:         false,
		Label:       9,
		UplinkMBR:   0,
		DownlinkMBR: 0,
		UplinkGBR:   0,
		DownlinkGBR: 0,
	})
	sgwDataFteid, _ := NewFteid(2, net.IPv4(1, 2, 3, 4), nil, S5S8SgwGtpUIf, 0x4075C31B)

	bcTBCwCSReqArg := BearerContextToBeCreatedWithinCSReqArg{
		Ebi:          ebi,
		BearerQoS:    bearerQoS,
		SgwDataFteid: sgwDataFteid,
	}
	bcTBCwCSReq, _ := NewBearerContextToBeCreatedWithinCSReq(bcTBCwCSReqArg)
	bcTBCwCSReqBin := bcTBCwCSReq.Marshal()

	msg, tail, err := Unmarshal(bcTBCwCSReqBin, CreateSessionRequest)
	bcTBCwCSReq = msg.(*BearerContextToBeCreatedWithinCSReq)
	assert.Equal(t, bearerContextNum, bcTBCwCSReq.typeNum)
	assert.Equal(t, ebi, bcTBCwCSReq.Ebi())
	assert.Equal(t, bearerQoS, bcTBCwCSReq.BearerQoS())
	assert.Equal(t, sgwDataFteid, bcTBCwCSReq.SgwDataFteid())
	assert.Nil(t, err)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)

	//ensure no refference to the buffer
	copy(bcTBCwCSReqBin, make([]byte, len(bcTBCwCSReqBin)))
	assert.Equal(t, bearerContextNum, bcTBCwCSReq.typeNum)
	assert.Equal(t, ebi, bcTBCwCSReq.Ebi())
	assert.Equal(t, bearerQoS, bcTBCwCSReq.BearerQoS())
	assert.Equal(t, sgwDataFteid, bcTBCwCSReq.SgwDataFteid())
}
