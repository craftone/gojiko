package gtpv2c

import (
	"errors"
	"net"

	"github.com/craftone/gojiko/domain/gtp"

	"github.com/craftone/gojiko/domain/gtpv2c/ie"
)

type ModifyBearerRequest struct {
	header
	imsi             *ie.Imsi
	uli              *ie.Uli     // for location information change reporting
	ratType          *ie.RatType // for RAT-change
	indication       *ie.Indication
	sgwCtrlFteid     *ie.Fteid
	bearerContextTBM *ie.BearerContextToBeModifiedWithinMBReq
	recovery         *ie.Recovery
}

type ModifyBearerRequestArg struct {
	PgwCtrlTeid      gtp.Teid
	Imsi             *ie.Imsi
	Uli              *ie.Uli
	RatType          *ie.RatType
	Indication       *ie.Indication
	SgwCtrlFteid     *ie.Fteid
	BearerContextTBM *ie.BearerContextToBeModifiedWithinMBReq
	Recovery         *ie.Recovery
}

func NewModifyBearerRequest(seqNum uint32, mbReqArg ModifyBearerRequestArg) (*ModifyBearerRequest, error) {
	// Actually this is not necessary since MBreq has no mandatory IE.
	// This code exists to be consistant with other implmentations of GTPv2-C messages.
	if err := checkModifyBearerRequestArg(mbReqArg); err != nil {
		return nil, err
	}

	return &ModifyBearerRequest{
		newHeader(ModifyBearerRequestNum, false, true, mbReqArg.PgwCtrlTeid, uint32(seqNum)),
		mbReqArg.Imsi,
		mbReqArg.Uli,
		mbReqArg.RatType,
		mbReqArg.Indication,
		mbReqArg.SgwCtrlFteid,
		mbReqArg.BearerContextTBM,
		mbReqArg.Recovery,
	}, nil
}

// checkModifyBearerRequestArg() ensures mandatory IEs should exists.
func checkModifyBearerRequestArg(mbReqArg ModifyBearerRequestArg) error {
	errMsgs := make([]string, 0)

	// no mandatory IE

	if len(errMsgs) == 0 {
		return nil
	}
	errMsg := ""
	for _, msg := range errMsgs {
		errMsg += msg + " must be specified. "
	}
	return errors.New(errMsg)
}

func (m *ModifyBearerRequest) Marshal() []byte {
	body := make([]byte, 0, 300)

	if m.imsi != nil {
		body = append(body, m.imsi.Marshal()...)
	}
	if m.uli != nil {
		body = append(body, m.uli.Marshal()...)
	}
	if m.ratType != nil {
		body = append(body, m.ratType.Marshal()...)
	}
	if m.indication != nil {
		body = append(body, m.indication.Marshal()...)
	}
	if m.sgwCtrlFteid != nil {
		body = append(body, m.sgwCtrlFteid.Marshal()...)
	}
	if m.bearerContextTBM != nil {
		body = append(body, m.bearerContextTBM.Marshal()...)
	}
	if m.recovery != nil {
		body = append(body, m.recovery.Marshal()...)
	}

	return m.header.marshal(body)
}

func unmarshalModifyBearerRequest(h header, buf []byte) (*ModifyBearerRequest, error) {
	if h.messageType != ModifyBearerRequestNum {
		log.Panic("Invalid messageType")
	}

	mbReqArg := ModifyBearerRequestArg{}
	for len(buf) > 0 {
		msg, tail, err := ie.Unmarshal(buf, ie.ModifyBearerRequest)
		buf = tail
		if err != nil {
			if _, ok := err.(*ie.UnknownIEError); ok {
				log.Println(err)
				continue
			}
			return nil, err
		}

		if msg.Instance() != 0 {
			log.Printf("Unkown IE : %v", msg)
			continue
		}

		switch msg := msg.(type) {
		case *ie.Imsi:
			mbReqArg.Imsi = msg
		case *ie.Uli:
			mbReqArg.Uli = msg
		case *ie.RatType:
			mbReqArg.RatType = msg
		case *ie.Indication:
			mbReqArg.Indication = msg
		case *ie.Fteid:
			mbReqArg.SgwCtrlFteid = msg
		case *ie.BearerContextToBeModifiedWithinMBReq:
			mbReqArg.BearerContextTBM = msg
		case *ie.Recovery:
			mbReqArg.Recovery = msg
		default:
			log.Printf("Unkown IE : %v", msg)
		}
	}
	return NewModifyBearerRequest(gtp.Teid(h.seqNum), mbReqArg)
}

func MakeMBReqArg(
	pgwCtrlTeid gtp.Teid,
	imsi string,
	mcc, mnc string,
	tac uint16, eci uint32,
	ratTypeValue ie.RatTypeValue,
	indicationArg ie.IndicationArg,
	sgwCtrlIPv4 net.IP, sgwCtrlTeid gtp.Teid,
	sgwDataIPv4 net.IP, sgwDataTeid gtp.Teid,
	ebi, recovery byte) (ModifyBearerRequestArg, error) {

	imsiIE, err := ie.NewImsi(0, imsi)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	taiIE, err := ie.NewTai(mcc, mnc, tac)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	ecgiIE, err := ie.NewEcgi(mcc, mnc, eci)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	uliArg := ie.UliArg{
		Tai:  taiIE,
		Ecgi: ecgiIE,
	}
	uliIE, err := ie.NewUli(0, uliArg)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	ratTypeIE, err := ie.NewRatType(0, ratTypeValue)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	// prepare Indication
	indicationIE, err := ie.NewIndication(0, indicationArg)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	sgwCtrlFteidIE, err := ie.NewFteid(0, sgwCtrlIPv4, nil, ie.S5S8SgwGtpCIf, sgwCtrlTeid)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	sgwDataFteidIE, err := ie.NewFteid(1, sgwDataIPv4, nil, ie.S5S8SgwGtpUIf, sgwDataTeid)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	ebiIE, err := ie.NewEbi(0, ebi)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	bearerContextTBMArg := ie.BearerContextToBeModifiedWithinMBReqArg{
		Ebi:          ebiIE,
		SgwDataFteid: sgwDataFteidIE,
	}
	bearerContextTBM, err := ie.NewBearerContextToBeModifiedWithinMBReq(bearerContextTBMArg)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	recoveryIE, err := ie.NewRecovery(0, recovery)
	if err != nil {
		return ModifyBearerRequestArg{}, err
	}

	return ModifyBearerRequestArg{
		PgwCtrlTeid:      pgwCtrlTeid,
		Imsi:             imsiIE,
		Uli:              uliIE,
		RatType:          ratTypeIE,
		Indication:       indicationIE,
		SgwCtrlFteid:     sgwCtrlFteidIE,
		BearerContextTBM: bearerContextTBM,
		Recovery:         recoveryIE,
	}, nil
}

//
// Getters
//

func (m *ModifyBearerRequest) Imsi() *ie.Imsi {
	return m.imsi
}
func (m *ModifyBearerRequest) Uli() *ie.Uli {
	return m.uli
}
func (m *ModifyBearerRequest) RatType() *ie.RatType {
	return m.ratType
}
func (m *ModifyBearerRequest) Indication() *ie.Indication {
	return m.indication
}
func (m *ModifyBearerRequest) SgwCtrlFteid() *ie.Fteid {
	return m.sgwCtrlFteid
}
func (m *ModifyBearerRequest) BearerContextTBM() *ie.BearerContextToBeModifiedWithinMBReq {
	return m.bearerContextTBM
}
func (m *ModifyBearerRequest) Recovery() *ie.Recovery {
	return m.recovery
}
