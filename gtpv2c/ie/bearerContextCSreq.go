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

func NewBearerContextToBeCreatedWithinCSReq(bcTBCwCSReqArg BearerContextToBeCreatedWithinCSReqArg) (*BearerContextToBeCreatedWithinCSReq, error) {
	// BearerContectToBeCreatedWithinCSReq's instance number shoub be 0.
	instance := byte(0)

	bcArg := BearerContextArg{
		ebi:          bcTBCwCSReqArg.Ebi,
		bearerQoS:    bcTBCwCSReqArg.BearerQoS,
		sgwDataFteid: bcTBCwCSReqArg.SgwDataFteid,
	}

	if bcArg.bearerQoS == nil {
		return nil, fmt.Errorf("Bearer QoS is a mondatory IE")
	}
	bcArg.bearerQoS.instance = 0

	if bcArg.sgwDataFteid == nil {
		return nil, fmt.Errorf("S5/S8-U SGW F-TEID is a mondatory IE in this condition")
	}
	bcArg.sgwDataFteid.instance = 2

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

	if h.instance != 0 {
		log.Fatalf("BearerContectToBeCreatedWithinCSReq's instance number shoub be 0 but %v", h.instance)
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

	bc, err := NewBearerContextToBeCreatedWithinCSReq(bcTBCwCSReqArg)
	if err != nil {
		return nil, err
	}
	return bc, nil
}
