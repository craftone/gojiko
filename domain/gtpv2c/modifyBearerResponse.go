package gtpv2c

import (
	"errors"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
)

type ModifyBearerResponse struct {
	header
	cause            *ie.Cause
	msisdn           *ie.Msisdn
	bearerContextTBM *ie.BearerContextToBeModifiedWithinMBRes
	recovery         *ie.Recovery
}

type ModifyBearerResponseArg struct {
	SgwCtrlTeid      gtp.Teid
	Cause            *ie.Cause
	Msisdn           *ie.Msisdn
	BearerContextTBM *ie.BearerContextToBeModifiedWithinMBRes
	Recovery         *ie.Recovery
}

func NewModifyBearerResponse(seqNum uint32, mbResArg ModifyBearerResponseArg) (*ModifyBearerResponse, error) {
	if err := checkModifyBearerResponseArg(mbResArg); err != nil {
		return nil, err
	}

	return &ModifyBearerResponse{
		newHeader(ModifyBearerResponseNum, false, true, mbResArg.SgwCtrlTeid, seqNum),
		mbResArg.Cause,
		mbResArg.Msisdn,
		mbResArg.BearerContextTBM,
		mbResArg.Recovery,
	}, nil
}

func checkModifyBearerResponseArg(mbReqArg ModifyBearerResponseArg) error {
	errMsgs := make([]string, 0)

	// Confirm mandatory IEs are exists
	if mbReqArg.Cause == nil {
		errMsgs = append(errMsgs, "Cause")
	}

	if len(errMsgs) == 0 {
		return nil
	}
	errMsg := ""
	for _, msg := range errMsgs {
		errMsg += msg + " must be specified. "
	}
	return errors.New(errMsg)
}

// MakeMBResArg makes a ModifyBearerResponseArg.
func MakeMBResArg(sgwCtrlTEID gtp.Teid, cause ie.CauseValue,
	msisdn string, chargingID uint32,
	ebi, recovory byte) (ModifyBearerResponseArg, error) {

	causeIE, err := ie.NewCause(0, cause, false, false, false, nil)
	if err != nil {
		return ModifyBearerResponseArg{}, err
	}

	msisdnIE, err := ie.NewMsisdn(0, msisdn)
	if err != nil {
		return ModifyBearerResponseArg{}, err
	}

	chargingIDIE, err := ie.NewChargingID(0, chargingID)
	if err != nil {
		return ModifyBearerResponseArg{}, err
	}

	ebiIE, err := ie.NewEbi(0, ebi)
	if err != nil {
		return ModifyBearerResponseArg{}, err
	}

	recoveryIE, err := ie.NewRecovery(0, recovory)
	if err != nil {
		return ModifyBearerResponseArg{}, err
	}

	bearerContextTBMwMBRes, err := ie.NewBearerContextToBeModifiedWithinMBRes(
		ie.BearerContextToBeModifiedWithinMBResArg{
			Ebi:        ebiIE,
			Cause:      causeIE,
			ChargingID: chargingIDIE,
		})

	return ModifyBearerResponseArg{
		SgwCtrlTeid:      sgwCtrlTEID,
		Cause:            causeIE,
		Msisdn:           msisdnIE,
		BearerContextTBM: bearerContextTBMwMBRes,
		Recovery:         recoveryIE,
	}, nil
}

func (m *ModifyBearerResponse) Marshal() []byte {
	body := make([]byte, 0, 300)

	if m.cause != nil {
		body = append(body, m.cause.Marshal()...)
	}
	if m.msisdn != nil {
		body = append(body, m.msisdn.Marshal()...)
	}
	if m.bearerContextTBM != nil {
		body = append(body, m.bearerContextTBM.Marshal()...)
	}
	if m.recovery != nil {
		body = append(body, m.recovery.Marshal()...)
	}

	return m.header.marshal(body)
}

func unmarshalModifyBearerResponse(h header, buf []byte) (*ModifyBearerResponse, error) {
	if h.messageType != ModifyBearerResponseNum {
		log.Panic("Invalid messageType")
	}

	mbResArg := ModifyBearerResponseArg{SgwCtrlTeid: h.teid}
	for len(buf) > 0 {
		msg, tail, err := ie.Unmarshal(buf, ie.ModifyBearerResponse)
		buf = tail
		if err != nil {
			if _, ok := err.(*ie.UnknownIEError); ok {
				log.Println(err)
				continue
			}
			return nil, err
		}

		if _, ok := msg.(*ie.Fteid); ok {
			if msg.Instance() != 1 {
				log.Printf("Instance number (%d) is invalid fot the IE : %v", msg.Instance(), msg)
				continue
			}
		} else if msg.Instance() != 0 {
			log.Printf("Instance number (%d) is invalid fot the IE : %v", msg.Instance(), msg)
			continue
		}

		switch msg := msg.(type) {
		case *ie.Cause:
			mbResArg.Cause = msg
		case *ie.Msisdn:
			mbResArg.Msisdn = msg
		case *ie.BearerContextToBeModifiedWithinMBRes:
			mbResArg.BearerContextTBM = msg
		case *ie.Recovery:
			mbResArg.Recovery = msg
		default:
			log.Printf("Unkown IE : %v", msg)
		}
	}
	return NewModifyBearerResponse(h.seqNum, mbResArg)
}

//
// Getters
//

func (m *ModifyBearerResponse) Cause() *ie.Cause {
	return m.cause
}
func (m *ModifyBearerResponse) Msisdn() *ie.Msisdn {
	return m.msisdn
}
func (m *ModifyBearerResponse) BearerContextTBM() *ie.BearerContextToBeModifiedWithinMBRes {
	return m.bearerContextTBM
}
func (m *ModifyBearerResponse) Recovery() *ie.Recovery {
	return m.recovery
}
