package ie

import (
	"fmt"
)

type BearerContextToBeModifiedWithinMBReq struct {
	bearerContext
}

// NOT IMPLEMENTED since this type is for S4/S11.
//
// type BearerContextToBeRemovedWithinMBReq struct {
// 	BearerContext
// }

type BearerContextToBeModifiedWithinMBReqArg struct {
	Ebi          *Ebi
	SgwDataFteid *Fteid
}

func NewBearerContextToBeModifiedWithinMBReq(bcTBCwMBReqArg BearerContextToBeModifiedWithinMBReqArg) (*BearerContextToBeModifiedWithinMBReq, error) {
	// BearerContectToBeModifiedWithinMBReq's instance number shoub be 0.
	instance := byte(0)

	bcArg := bearerContextArg{
		ebi:          bcTBCwMBReqArg.Ebi,
		sgwDataFteid: bcTBCwMBReqArg.SgwDataFteid,
	}

	if bcArg.sgwDataFteid == nil {
		return nil, fmt.Errorf("S5/S8-U SGW F-TEID is a mondatory IE in this condition")
	}
	bcArg.sgwDataFteid.instance = 2

	bc, err := NewBearerContext(instance, bcArg)
	if err != nil {
		return nil, err
	}

	return &BearerContextToBeModifiedWithinMBReq{*bc}, nil
}

func unmarshalBearerContextToBeModifiedWithinMBReq(h header, buf []byte) (*BearerContextToBeModifiedWithinMBReq, error) {
	if h.typeNum != bearerContextNum {
		log.Panic("Invalid type")
	}

	if h.instance != 0 {
		log.Panicf("BearerContectToBeModifiedWithinMBReq's instance number shoub be 0 but %v", h.instance)
	}

	bcTBCwMBReqArg := BearerContextToBeModifiedWithinMBReqArg{}
	for len(buf) > 0 {
		msg, tail, err := Unmarshal(buf, ModifyBearerRequest)
		if err != nil {
			return nil, err
		}
		buf = tail

		switch msg := msg.(type) {
		case *Ebi:
			bcTBCwMBReqArg.Ebi = msg
		case *Fteid:
			if msg.instance == 2 {
				bcTBCwMBReqArg.SgwDataFteid = msg
			} else {
				log.Debugf("Unkown IE : %v", msg)
			}
		default:
			log.Debugf("Unkown IE : %v", msg)
		}
	}

	bc, err := NewBearerContextToBeModifiedWithinMBReq(bcTBCwMBReqArg)
	if err != nil {
		return nil, err
	}
	return bc, nil
}
