package ie

import (
	"fmt"
	"log"
)

type BearerContextToBeCreatedWithinCSReq struct {
	BearerContext
}

// NOT IMPLEMENTED since this type is for S4/S11.
//
// type BearerContextToBeRemovedWithinCSReq struct {
// 	BearerContext
// }

type BearerContextToBeCreatedWithinCSReqArg struct {
	Ebi          *Ebi
	BearerQoS    *BearerQoS
	SgwDataFteid *Fteid
}

func NewBearerContextToBeCreatedWithinCSReq(instance byte, bcTBCwCSReqArg BearerContextToBeCreatedWithinCSReqArg) (*BearerContextToBeCreatedWithinCSReq, error) {
	// Enforce instance number 0 since BearerContectToBeCreatedWithinCSReq's
	// instance shoub be 0.
	instance = 0

	bcArg := BearerContextArg{
		Ebi:          bcTBCwCSReqArg.Ebi,
		BearerQoS:    bcTBCwCSReqArg.BearerQoS,
		SgwDataFteid: bcTBCwCSReqArg.SgwDataFteid,
	}

	if bcArg.BearerQoS == nil {
		return nil, fmt.Errorf("Bearer QoS is a mondatory IE")
	}
	bcArg.BearerQoS.instance = 0

	if bcArg.SgwDataFteid == nil {
		return nil, fmt.Errorf("S5/S8-U SGW F-TEID is a mondatory IE in this condition")
	}
	bcArg.SgwDataFteid.instance = 2

	bc, err := NewBearerContext(instance, bcArg)
	if err != nil {
		return nil, err
	}

	return &BearerContextToBeCreatedWithinCSReq{*bc}, nil
}

func unmarshalBearerContextToBeCreatedWithinCSReq(h header, buf []byte) (*BearerContextToBeCreatedWithinCSReq, error) {
	if h.typeNum != bearerContextNum {
		log.Fatal("Invalud type")
	}

	bcTBCwCSReqArg := BearerContextToBeCreatedWithinCSReqArg{}
	for len(buf) > 0 {
		msg, tail, err := Unmarshal(buf, CreateSessionRequest)
		if err != nil {
			return nil, err
		}
		buf = tail

		switch msg := msg.(type) {
		case *Ebi:
			bcTBCwCSReqArg.Ebi = msg
		case *BearerQoS:
			bcTBCwCSReqArg.BearerQoS = msg
		case *Fteid:
			if msg.instance == 2 {
				bcTBCwCSReqArg.SgwDataFteid = msg
			} else {
				log.Printf("Unkown IE : %v", msg)
			}
		default:
			log.Printf("Unkown IE : %v", msg)
		}
	}

	bc, err := NewBearerContextToBeCreatedWithinCSReq(h.instance, bcTBCwCSReqArg)
	if err != nil {
		return nil, err
	}
	return bc, nil
}
