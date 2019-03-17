package ie

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBearerContextToBeModifiedWithinMBReq(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	sgwDataFteid, _ := NewFteid(2, net.IPv4(1, 2, 3, 4), nil, S5S8SgwGtpUIf, 0x4075C31B)

	bcTBMwMBReqArg := BearerContextToBeModifiedWithinMBReqArg{
		Ebi:          ebi,
		SgwDataFteid: sgwDataFteid,
	}
	bcTBMwMBReq, err := NewBearerContextToBeModifiedWithinMBReq(bcTBMwMBReqArg)

	assert.Equal(t, bearerContextNum, bcTBMwMBReq.typeNum)
	assert.Nil(t, err)

	// check Mondatory
	bcTBMwMBReqArg.Ebi = nil
	_, err = NewBearerContextToBeModifiedWithinMBReq(bcTBMwMBReqArg)
	assert.Error(t, err)
	bcTBMwMBReqArg.Ebi = ebi

	// check Mondatory
	bcTBMwMBReqArg.SgwDataFteid = nil
	_, err = NewBearerContextToBeModifiedWithinMBReq(bcTBMwMBReqArg)
	assert.Error(t, err)
}

func TestBearerContextToBeModifiedWithinMBReq_Marshal(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	sgwDataFteid, _ := NewFteid(0, net.IPv4(1, 2, 3, 4), nil, S5S8SgwGtpUIf, 0x4075C31B)

	bcTBMwMBReqArg := BearerContextToBeModifiedWithinMBReqArg{
		Ebi:          ebi,
		SgwDataFteid: sgwDataFteid,
	}
	bcTBMwMBReq, _ := NewBearerContextToBeModifiedWithinMBReq(bcTBMwMBReqArg)
	bcTBMwMBReqBin := bcTBMwMBReq.Marshal()

	assert.Equal(t, []byte{
		0x5D, 0, 0x12, 0, //Header
		0x49, 0, 1, 0, 5, //EBI
		0x57, 0, 9, 2, 0x84, 0x40, 0x75, 0xc3, 0x1b, 1, 2, 3, 4, //SgwDataFTEID
	}, bcTBMwMBReqBin)
}

func TestUnmarshal_BearerContextToBeModifiedWithinMBReq(t *testing.T) {
	ebi, _ := NewEbi(0, 5)
	sgwDataFteid, _ := NewFteid(2, net.IPv4(1, 2, 3, 4), nil, S5S8SgwGtpUIf, 0x4075C31B)

	bcTBMwMBReqArg := BearerContextToBeModifiedWithinMBReqArg{
		Ebi:          ebi,
		SgwDataFteid: sgwDataFteid,
	}
	bcTBMwMBReq, _ := NewBearerContextToBeModifiedWithinMBReq(bcTBMwMBReqArg)
	bcTBMwMBReqBin := bcTBMwMBReq.Marshal()

	msg, tail, err := Unmarshal(bcTBMwMBReqBin, ModifyBearerRequest)
	bcTBMwMBReq = msg.(*BearerContextToBeModifiedWithinMBReq)
	assert.Equal(t, bearerContextNum, bcTBMwMBReq.typeNum)
	assert.Equal(t, ebi, bcTBMwMBReq.Ebi())
	assert.Equal(t, sgwDataFteid, bcTBMwMBReq.SgwDataFteid())
	assert.Nil(t, err)
	assert.Equal(t, []byte{}, tail)
	assert.Nil(t, err)
}
