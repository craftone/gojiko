package ie

import (
	"fmt"
)

type BearerContextCreatedWithinCSRes struct {
	bearerContext
}

// NOT IMPLEMENTED since this type is for S4/S11.
//
// type BearerContextMarkedForRemovalWithinCSRes struct {
// 	BearerContext
// }

type BearerContextCreatedWithinCSResArg struct {
	Ebi          *Ebi
	Cause        *Cause
	PgwDataFteid *Fteid
	ChargingID   *ChargingID
}

func NewBearerContextCreatedWithinCSRes(bcCwCSResArg BearerContextCreatedWithinCSResArg) (*BearerContextCreatedWithinCSRes, error) {
	// BearerContectCreatedWithinCSRes's instance number shoub be 0.
	instance := byte(0)

	bcArg := bearerContextArg{
		ebi:          bcCwCSResArg.Ebi,
		cause:        bcCwCSResArg.Cause,
		pgwDataFteid: bcCwCSResArg.PgwDataFteid,
		chargingID:   bcCwCSResArg.ChargingID,
	}

	if bcArg.cause == nil {
		return nil, fmt.Errorf("Cause is a mondatory IE")
	}
	bcArg.cause.instance = 0

	if bcArg.cause.Value().Type() == CauseTypeAcceptance {
		if bcArg.pgwDataFteid == nil {
			return nil, fmt.Errorf("S5/S8-U PGW F-TEID is a mondatory IE in this condition")
		}
		bcArg.pgwDataFteid.instance = 2
	}

	bc, err := NewBearerContext(instance, bcArg)
	if err != nil {
		return nil, err
	}

	return &BearerContextCreatedWithinCSRes{*bc}, nil
}

func unmarshalBearerContextCreatedWithinCSRes(h header, buf []byte) (*BearerContextCreatedWithinCSRes, error) {
	if h.typeNum != bearerContextNum {
		log.Panic("Invalid type")
	}

	if h.instance != 0 {
		log.Panicf("BearerContectCreatedWithinCSRes's instance number shoub be 0 but %v", h.instance)
	}

	bcCwCSResArg := BearerContextCreatedWithinCSResArg{}
	for len(buf) > 0 {
		msg, tail, err := Unmarshal(buf, CreateSessionRequest)
		if err != nil {
			return nil, err
		}
		buf = tail

		switch msg := msg.(type) {
		case *Ebi:
			bcCwCSResArg.Ebi = msg
		case *Cause:
			bcCwCSResArg.Cause = msg
		case *Fteid:
			if msg.instance == 2 {
				bcCwCSResArg.PgwDataFteid = msg
			} else {
				log.Debugf("Unkown IE : %v", msg)
			}
		case *ChargingID:
			bcCwCSResArg.ChargingID = msg
		default:
			log.Debugf("Unkown IE : %v", msg)
		}
	}

	bc, err := NewBearerContextCreatedWithinCSRes(bcCwCSResArg)
	if err != nil {
		return nil, err
	}
	return bc, nil
}
