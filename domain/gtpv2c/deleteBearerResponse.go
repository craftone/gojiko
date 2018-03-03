package gtpv2c

import (
	"fmt"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
)

type DeleteBearerResponse struct {
	header
	cause    *ie.Cause
	lbi      *ie.Ebi
	recovery *ie.Recovery
}

func NewDeleteBearerResponse(teid gtp.Teid, seqNum uint32, cause ie.CauseValue, ebi byte, rec byte) (*DeleteBearerResponse, error) {
	causeIE, err := ie.NewCause(0, cause, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	ebiIE, err := ie.NewEbi(0, ebi)
	if err != nil {
		return nil, err
	}

	recoveryIE, err := ie.NewRecovery(0, rec)
	if err != nil {
		return nil, err
	}

	return &DeleteBearerResponse{
		newHeader(DeleteBearerResponseNum, false, true, teid, seqNum),
		causeIE,
		ebiIE,
		recoveryIE,
	}, nil
}

func (d *DeleteBearerResponse) Marshal() []byte {
	body := make([]byte, 0, 16)
	body = append(body, d.cause.Marshal()...)
	body = append(body, d.lbi.Marshal()...)
	body = append(body, d.recovery.Marshal()...)
	return d.header.marshal(body)
}

func unmarshalDeleteBearerResponse(h header, buf []byte) (*DeleteBearerResponse, error) {
	if h.messageType != DeleteBearerResponseNum {
		panic(fmt.Sprintf("Invalid message Type : %d", h.messageType))
	}

	var causeIE *ie.Cause
	var ebiIE *ie.Ebi
	var recoveryIE *ie.Recovery
	for len(buf) > 0 {
		msg, tail, err := ie.Unmarshal(buf, ie.DeleteBearerResponse)
		buf = tail
		if err != nil {
			if _, ok := err.(*ie.UnknownIEError); ok {
				log.Println(err)
				continue
			}
			return nil, err
		}

		if msg.Instance() != 0 {
			log.Printf("Unkown IE : %#v", msg)
		}

		switch msg := msg.(type) {
		case *ie.Cause:
			causeIE = msg
		case *ie.Ebi:
			ebiIE = msg
		case *ie.Recovery:
			recoveryIE = msg
		default:
			log.Printf("Unkown IE : %#v", msg)
		}
	}
	if causeIE == nil {
		return nil, fmt.Errorf("No Cause Delete Bearer Request message")
	}
	return NewDeleteBearerResponse(h.Teid(), h.seqNum, causeIE.Value(), ebiIE.Value(), recoveryIE.Value())
}

// Getters

func (d *DeleteBearerResponse) Cause() *ie.Cause {
	return d.cause
}

func (d *DeleteBearerResponse) Lbi() *ie.Ebi {
	return d.lbi
}

func (d *DeleteBearerResponse) Recovery() *ie.Recovery {
	return d.recovery
}
