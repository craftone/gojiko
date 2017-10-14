package ie

import "fmt"

type bearerContext struct {
	header
	bearerContextArg
}

type bearerContextArg struct {
	ebi          *Ebi
	bearerQoS    *BearerQoS
	sgwDataFteid *Fteid
	cause        *Cause
	pgwDataFteid *Fteid
	chargingID   *ChargingID
}

func NewBearerContext(instance byte, bcArg bearerContextArg) (*bearerContext, error) {
	header, err := newHeader(bearerContextNum, 0, instance)
	if err != nil {
		return nil, err
	}

	// Check and Set Instance num for EBI
	if bcArg.ebi == nil {
		return nil, fmt.Errorf("EBI (EPS Bearer ID) is a mondatory IE")
	}
	bcArg.ebi.instance = 0

	return &bearerContext{
		header,
		bcArg,
	}, nil
}

func (b *bearerContext) Marshal() []byte {
	body := make([]byte, 0, 44) // 44 octet depends on experience
	if b.ebi != nil {
		body = append(body, b.ebi.Marshal()...)
	}
	if b.cause != nil {
		body = append(body, b.cause.Marshal()...)
	}
	if b.bearerQoS != nil {
		body = append(body, b.bearerQoS.Marshal()...)
	}
	if b.sgwDataFteid != nil {
		body = append(body, b.sgwDataFteid.Marshal()...)
	}
	if b.pgwDataFteid != nil {
		body = append(body, b.pgwDataFteid.Marshal()...)
	}
	if b.chargingID != nil {
		body = append(body, b.chargingID.Marshal()...)
	}
	return b.header.marshal(body)
}

// Getters

func (b *bearerContext) Ebi() *Ebi {
	return b.ebi
}
func (b *bearerContext) BearerQoS() *BearerQoS {
	return b.bearerQoS
}
func (b *bearerContext) SgwDataFteid() *Fteid {
	return b.sgwDataFteid
}
func (b *bearerContext) Cause() *Cause {
	return b.cause
}
func (b *bearerContext) PgwDataFteid() *Fteid {
	return b.pgwDataFteid
}
func (b *bearerContext) ChargingID() *ChargingID {
	return b.chargingID
}
