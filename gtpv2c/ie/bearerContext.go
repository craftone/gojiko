package ie

import "fmt"

type BearerContext struct {
	header
	BearerContextArg
}

type BearerContextArg struct {
	ebi          *Ebi
	bearerQoS    *BearerQoS
	sgwDataFteid *Fteid
	cause        *Cause
	pgwDataFteid *Fteid
	chargingID   *ChargingID
}

func NewBearerContext(instance byte, bcArg BearerContextArg) (*BearerContext, error) {
	header, err := newHeader(bearerContextNum, 0, instance)
	if err != nil {
		return nil, err
	}

	// Check and Set Instance num for EBI
	if bcArg.ebi == nil {
		return nil, fmt.Errorf("EBI (EPS Bearer ID) is a mondatory IE")
	}
	bcArg.ebi.instance = 0

	return &BearerContext{
		header,
		bcArg,
	}, nil
}

func (b *BearerContext) Marshal() []byte {
	body := make([]byte, 0, 44) // 44 octet depends on experience
	if b.ebi != nil {
		body = append(body, b.ebi.Marshal()...)
	}
	if b.bearerQoS != nil {
		body = append(body, b.bearerQoS.Marshal()...)
	}
	if b.sgwDataFteid != nil {
		body = append(body, b.sgwDataFteid.Marshal()...)
	}
	if b.cause != nil {
		body = append(body, b.cause.Marshal()...)
	}
	if b.pgwDataFteid != nil {
		body = append(body, b.pgwDataFteid.Marshal()...)
	}
	if b.chargingID != nil {
		body = append(body, b.chargingID.Marshal()...)
	}
	return b.header.marshal(body)
}
