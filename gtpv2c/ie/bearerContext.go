package ie

import "fmt"

type BearerContext struct {
	header
	BearerContextArg
}

type BearerContextArg struct {
	Ebi          *Ebi
	BearerQoS    *BearerQoS
	SgwDataFteid *Fteid
	Cause        *Cause
	PgwDataFteid *Fteid
	ChargingID   *ChargingID
}

func NewBearerContext(instance byte, bcArg BearerContextArg) (*BearerContext, error) {
	header, err := newHeader(bearerContextNum, 0, instance)
	if err != nil {
		return nil, err
	}

	// Check and Set Instance num for EBI
	if bcArg.Ebi == nil {
		return nil, fmt.Errorf("EBI (EPS Bearer ID) is a mondatory IE")
	}
	bcArg.Ebi.instance = 0

	return &BearerContext{
		header,
		bcArg,
	}, nil
}

func (b *BearerContext) Marshal() []byte {
	body := make([]byte, 0, 44) // 44 octet depends on experience
	if b.Ebi != nil {
		body = append(body, b.Ebi.Marshal()...)
	}
	if b.BearerQoS != nil {
		body = append(body, b.BearerQoS.Marshal()...)
	}
	if b.SgwDataFteid != nil {
		body = append(body, b.SgwDataFteid.Marshal()...)
	}
	if b.Cause != nil {
		body = append(body, b.Cause.Marshal()...)
	}
	if b.PgwDataFteid != nil {
		body = append(body, b.PgwDataFteid.Marshal()...)
	}
	if b.ChargingID != nil {
		body = append(body, b.ChargingID.Marshal()...)
	}
	return b.header.marshal(body)
}

// func unmarshalBearerContext(h header, buf []byte, msgType MsgType) (*BearerContext, error) {
// 	if h.typeNum != bearerContextNum {
// 		log.Fatal("Invalud type")
// 	}

// 	bcArg := BearerContextArg{}
// 	msg, tail, err := Unmarshal(buf, MsToNetwork)
// 	if err != nil {
// 		return nil, err
// 	}

// 	switch msg := msg.(type) {
// 	case *Ebi:
// 		bcArg.Ebi = msg
// 	case *BearerQoS:
// 		bcArg.BearerQoS = msg
// 	default:
// 		log.Printf("Unkown IE : %v", msg)
// 	}

// 	bc, err := NewBearerContext(h.instance, bcArg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bc, nil
// }
