package ie

import (
	"fmt"
)

type BearerContextToBeModifiedWithinMBRes struct {
	bearerContext
}

// NOT IMPLEMENTED since this type is for S4/S11.
//
// type BearerContextToBeRemovedWithinMBRes struct {
// 	BearerContext
// }

type BearerContextToBeModifiedWithinMBResArg struct {
	Ebi        *Ebi
	Cause      *Cause
	ChargingID *ChargingID
}

func NewBearerContextToBeModifiedWithinMBRes(bcTBCwMBResArg BearerContextToBeModifiedWithinMBResArg) (*BearerContextToBeModifiedWithinMBRes, error) {
	// BearerContectToBeModifiedWithinMBReq's instance number shoub be 0.
	instance := byte(0)

	bcArg := bearerContextArg{
		ebi:        bcTBCwMBResArg.Ebi,
		cause:      bcTBCwMBResArg.Cause,
		chargingID: bcTBCwMBResArg.ChargingID,
	}

	if bcArg.ebi == nil {
		return nil, fmt.Errorf("EBI is a mondatory IE")
	}
	if bcArg.cause == nil {
		return nil, fmt.Errorf("Cause is a mondatory IE")
	}

	bc, err := NewBearerContext(instance, bcArg)
	if err != nil {
		return nil, err
	}

	return &BearerContextToBeModifiedWithinMBRes{*bc}, nil
}

func unmarshalBearerContextToBeModifiedWithinMBRes(h header, buf []byte) (*BearerContextToBeModifiedWithinMBRes, error) {
	if h.typeNum != bearerContextNum {
		log.Panic("Invalid type")
	}

	if h.instance != 0 {
		log.Panicf("BearerContectToBeModifiedWithinMBRes's instance number should be 0 but %v", h.instance)
	}

	bcTBCwMBResArg := BearerContextToBeModifiedWithinMBResArg{}
	for len(buf) > 0 {
		msg, tail, err := Unmarshal(buf, ModifyBearerRequest)
		if err != nil {
			return nil, err
		}
		buf = tail

		switch msg := msg.(type) {
		case *Cause:
			bcTBCwMBResArg.Cause = msg
		case *Ebi:
			bcTBCwMBResArg.Ebi = msg
		case *ChargingID:
			bcTBCwMBResArg.ChargingID = msg
		default:
			log.Debugf("Unkown IE : %v", msg)
		}
	}

	bc, err := NewBearerContextToBeModifiedWithinMBRes(bcTBCwMBResArg)
	if err != nil {
		return nil, err
	}
	return bc, nil
}
