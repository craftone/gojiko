package gtpv2c

import (
	"errors"
	"log"
	"net"

	"github.com/craftone/gojiko/gtpv2c/ie/pco"

	"github.com/craftone/gojiko/gtp"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

type CreateSessionRequest struct {
	header
	imsi             *ie.Imsi
	msisdn           *ie.Msisdn
	mei              *ie.Mei
	uli              *ie.Uli
	servingNetwork   *ie.ServingNetwork
	ratType          *ie.RatType
	indication       *ie.Indication
	senderFteid      *ie.Fteid
	apn              *ie.Apn
	selectionMode    *ie.SelectionMode
	pdnType          *ie.PdnType
	paa              *ie.Paa
	apnRestriction   *ie.ApnRestriction
	apnAmbr          *ie.Ambr
	pco              *ie.PcoMsToNetwork
	bearerContextTBC *ie.BearerContextToBeCreatedWithinCSReq
	recovery         *ie.Recovery
}

type CreateSessionRequestArg struct {
	Imsi             *ie.Imsi
	Msisdn           *ie.Msisdn
	Mei              *ie.Mei
	Uli              *ie.Uli
	ServingNetwork   *ie.ServingNetwork
	RatType          *ie.RatType
	Indication       *ie.Indication
	SenderFteid      *ie.Fteid
	Apn              *ie.Apn
	SelectionMode    *ie.SelectionMode
	PdnType          *ie.PdnType
	Paa              *ie.Paa
	ApnRestriction   *ie.ApnRestriction
	ApnAmbr          *ie.Ambr
	Pco              *ie.PcoMsToNetwork
	BearerContextTBC *ie.BearerContextToBeCreatedWithinCSReq
	Recovery         *ie.Recovery
}

func NewCreateSessionRequest(seqNum uint32, csReqArg CreateSessionRequestArg) (*CreateSessionRequest, error) {
	if err := checkCreateSessionRequestArg(csReqArg); err != nil {
		return nil, err
	}

	return &CreateSessionRequest{
		newHeader(createSessionRequestNum, false, true, 0, seqNum),
		csReqArg.Imsi,
		csReqArg.Msisdn,
		csReqArg.Mei,
		csReqArg.Uli,
		csReqArg.ServingNetwork,
		csReqArg.RatType,
		csReqArg.Indication,
		csReqArg.SenderFteid,
		csReqArg.Apn,
		csReqArg.SelectionMode,
		csReqArg.PdnType,
		csReqArg.Paa,
		csReqArg.ApnRestriction,
		csReqArg.ApnAmbr,
		csReqArg.Pco,
		csReqArg.BearerContextTBC,
		csReqArg.Recovery,
	}, nil
}

func checkCreateSessionRequestArg(csReqArg CreateSessionRequestArg) error {
	errMsgs := make([]string, 0)
	if csReqArg.Imsi == nil {
		errMsgs = append(errMsgs, "IMSI")
	}
	if csReqArg.Msisdn == nil {
		errMsgs = append(errMsgs, "MSISDN")
	}
	if csReqArg.Mei == nil {
		errMsgs = append(errMsgs, "MEI")
	}
	if csReqArg.Uli == nil {
		errMsgs = append(errMsgs, "ULI")
	}
	if csReqArg.ServingNetwork == nil {
		errMsgs = append(errMsgs, "Serving Network")
	}
	if csReqArg.RatType == nil {
		errMsgs = append(errMsgs, "RatType")
	}
	if csReqArg.Indication == nil {
		errMsgs = append(errMsgs, "Indication")
	}
	if csReqArg.SenderFteid == nil {
		errMsgs = append(errMsgs, "Sender F-TEID")
	}
	if csReqArg.Apn == nil {
		errMsgs = append(errMsgs, "APN")
	}
	if csReqArg.SelectionMode == nil {
		errMsgs = append(errMsgs, "Selection Mode")
	}
	if csReqArg.PdnType == nil {
		errMsgs = append(errMsgs, "PDN Type")
	}
	if csReqArg.Paa == nil {
		errMsgs = append(errMsgs, "PAA")
	}
	if csReqArg.ApnRestriction == nil {
		errMsgs = append(errMsgs, "Max APN Restriction")
	}
	if csReqArg.ApnAmbr == nil {
		errMsgs = append(errMsgs, "APN-AMBR")
	}
	if csReqArg.Pco == nil {
		errMsgs = append(errMsgs, "PCO")
	}
	if csReqArg.BearerContextTBC == nil {
		errMsgs = append(errMsgs, "Bearer Context to be created")
	}
	if csReqArg.Recovery == nil {
		errMsgs = append(errMsgs, "Recovery")
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

// MakeCSReqArg makes a CreateSessionRequestArg.
// - MCC and MNC are used for ULI and Serving Network.
// - RAT Type is always EUTRAN(6).
// - Selection Mode is always MS provided APN,subscr not verified (0x01).
// - PDN Type is always IPv4.
// - PAA is always IPv4:0.0.0.0.
// - Max APN Restriction is always No Existing Contexts or Restriction (0).
// - APN AMBR is always up:4294967 kbps, down:4294967 kbps.
func MakeCSReqArg(imsi, msisdn, mei, mcc, mnc string,
	sgwCtrlIPv4 net.IP, sgwCtrlTeid gtp.Teid,
	sgwDataIPv4 net.IP, sgwDataTeid gtp.Teid,
	apn string, ebi, recovery byte) (CreateSessionRequestArg, error) {
	imsiIE, err := ie.NewImsi(0, imsi)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	msisdnIE, err := ie.NewMsisdn(0, msisdn)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	meiIE, err := ie.NewMei(0, mei)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	taiIE, err := ie.NewTai(mcc, mnc, 0)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}
	ecgiIE, err := ie.NewEcgi(mcc, mnc, 0)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}
	uliArg := ie.UliArg{
		Tai:  taiIE,
		Ecgi: ecgiIE,
	}
	uliIE, err := ie.NewUli(0, uliArg)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	servingNetworkIE, err := ie.NewServingNetwork(0, mcc, mnc)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	ratTypeIE, err := ie.NewRatType(0, 6) // EUTRAN(6)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	// Indication is always all 0.
	indicationIE, err := ie.NewIndication(0, ie.IndicationArg{})
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	senderFteidIE, err := ie.NewFteid(0, sgwCtrlIPv4, nil, ie.S5S8SgwGtpCIf, sgwCtrlTeid)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	apnIE, err := ie.NewApn(0, apn)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	selectionModeIE, err := ie.NewSelectionMode(0, 1)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	pdnTypeIE, err := ie.NewPdnType(0, ie.PdnTypeIPv4)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	paaIE, err := ie.NewPaa(0, ie.PdnTypeIPv4, net.IPv4(0, 0, 0, 0), nil)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	apnRestrictionIE, err := ie.NewApnRestriction(0, 0)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	apnAmbrIE, err := ie.NewAmbr(0, 4294967, 4294967)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	ipcpIE := pco.NewIpcp(pco.ConfigureRequest, 0, net.IPv4(0, 0, 0, 0), net.IPv4(0, 0, 0, 0))
	pcoMsToNetworkIE := pco.NewMsToNetwork(ipcpIE, true, false, true)
	pcoIE, err := ie.NewPcoMsToNetwork(0, pcoMsToNetworkIE)

	ebiIE, err := ie.NewEbi(0, ebi)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	bearerQoSIE, err := ie.NewBearerQoS(0, ie.BearerQoSArg{
		Pci:         true,
		Pl:          15,
		Pvi:         false,
		Label:       9,
		UplinkMBR:   0,
		DownlinkMBR: 0,
		UplinkGBR:   0,
		DownlinkGBR: 0,
	})
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	sgwDataFteidIE, err := ie.NewFteid(2, sgwDataIPv4, nil, ie.S5S8SgwGtpUIf, sgwDataTeid)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	bearerContextTBCIE, err := ie.NewBearerContextToBeCreatedWithinCSReq(
		ie.BearerContextToBeCreatedWithinCSReqArg{ebiIE, bearerQoSIE, sgwDataFteidIE})

	recoveryIE, err := ie.NewRecovery(0, recovery)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	return CreateSessionRequestArg{
		Imsi:             imsiIE,
		Msisdn:           msisdnIE,
		Mei:              meiIE,
		Uli:              uliIE,
		ServingNetwork:   servingNetworkIE,
		RatType:          ratTypeIE,
		Indication:       indicationIE,
		SenderFteid:      senderFteidIE,
		Apn:              apnIE,
		SelectionMode:    selectionModeIE,
		PdnType:          pdnTypeIE,
		Paa:              paaIE,
		ApnRestriction:   apnRestrictionIE,
		ApnAmbr:          apnAmbrIE,
		Pco:              pcoIE,
		BearerContextTBC: bearerContextTBCIE,
		Recovery:         recoveryIE,
	}, nil
}

func (c *CreateSessionRequest) Marshal() []byte {
	body := make([]byte, 0, 300)

	if c.imsi != nil {
		body = append(body, c.imsi.Marshal()...)
	}
	if c.msisdn != nil {
		body = append(body, c.msisdn.Marshal()...)
	}
	if c.mei != nil {
		body = append(body, c.mei.Marshal()...)
	}
	if c.uli != nil {
		body = append(body, c.uli.Marshal()...)
	}
	if c.servingNetwork != nil {
		body = append(body, c.servingNetwork.Marshal()...)
	}
	if c.ratType != nil {
		body = append(body, c.ratType.Marshal()...)
	}
	if c.indication != nil {
		body = append(body, c.indication.Marshal()...)
	}
	if c.senderFteid != nil {
		body = append(body, c.senderFteid.Marshal()...)
	}
	if c.apn != nil {
		body = append(body, c.apn.Marshal()...)
	}
	if c.selectionMode != nil {
		body = append(body, c.selectionMode.Marshal()...)
	}
	if c.pdnType != nil {
		body = append(body, c.pdnType.Marshal()...)
	}
	if c.paa != nil {
		body = append(body, c.paa.Marshal()...)
	}
	if c.apnRestriction != nil {
		body = append(body, c.apnRestriction.Marshal()...)
	}
	if c.apnAmbr != nil {
		body = append(body, c.apnAmbr.Marshal()...)
	}
	if c.pco != nil {
		body = append(body, c.pco.Marshal()...)
	}
	if c.bearerContextTBC != nil {
		body = append(body, c.bearerContextTBC.Marshal()...)
	}
	if c.recovery != nil {
		body = append(body, c.recovery.Marshal()...)
	}

	return c.header.marshal(body)
}

func unmarshalCreateSessionRequest(h header, buf []byte) (*CreateSessionRequest, error) {
	if h.messageType != createSessionRequestNum {
		log.Fatal("Invalud messageType")
	}
	log.Printf("%v\n", buf)
	csReqArg := CreateSessionRequestArg{}
	for len(buf) > 0 {
		msg, tail, err := ie.Unmarshal(buf, ie.CreateSessionRequest)
		if err != nil {
			return nil, err
		}
		buf = tail

		if msg.Instance() != 0 {
			log.Printf("Unkown IE : %v", msg)
			continue
		}

		switch msg := msg.(type) {
		case *ie.Imsi:
			csReqArg.Imsi = msg
		case *ie.Msisdn:
			csReqArg.Msisdn = msg
		case *ie.Mei:
			csReqArg.Mei = msg
		case *ie.Uli:
			csReqArg.Uli = msg
		case *ie.ServingNetwork:
			csReqArg.ServingNetwork = msg
		case *ie.RatType:
			csReqArg.RatType = msg
		case *ie.Indication:
			csReqArg.Indication = msg
		case *ie.Fteid:
			csReqArg.SenderFteid = msg
		case *ie.Apn:
			csReqArg.Apn = msg
		case *ie.SelectionMode:
			csReqArg.SelectionMode = msg
		case *ie.PdnType:
			csReqArg.PdnType = msg
		case *ie.Paa:
			csReqArg.Paa = msg
		case *ie.ApnRestriction:
			csReqArg.ApnRestriction = msg
		case *ie.Ambr:
			csReqArg.ApnAmbr = msg
		case *ie.PcoMsToNetwork:
			csReqArg.Pco = msg
		case *ie.BearerContextToBeCreatedWithinCSReq:
			csReqArg.BearerContextTBC = msg
		case *ie.Recovery:
			csReqArg.Recovery = msg
		default:
			log.Printf("Unkown IE : %v", msg)
		}
	}
	return NewCreateSessionRequest(h.seqNum, csReqArg)
}

//
// Getters
//

func (c *CreateSessionRequest) Imsi() *ie.Imsi {
	return c.imsi
}
func (c *CreateSessionRequest) Msisdn() *ie.Msisdn {
	return c.msisdn
}
func (c *CreateSessionRequest) Mei() *ie.Mei {
	return c.mei
}
func (c *CreateSessionRequest) Uli() *ie.Uli {
	return c.uli
}
func (c *CreateSessionRequest) ServingNetwork() *ie.ServingNetwork {
	return c.servingNetwork
}
func (c *CreateSessionRequest) RatType() *ie.RatType {
	return c.ratType
}
func (c *CreateSessionRequest) Indication() *ie.Indication {
	return c.indication
}
func (c *CreateSessionRequest) SenderFteid() *ie.Fteid {
	return c.senderFteid
}
func (c *CreateSessionRequest) Apn() *ie.Apn {
	return c.apn
}
func (c *CreateSessionRequest) SelectionMode() *ie.SelectionMode {
	return c.selectionMode
}
func (c *CreateSessionRequest) PdnType() *ie.PdnType {
	return c.pdnType
}
func (c *CreateSessionRequest) Paa() *ie.Paa {
	return c.paa
}
func (c *CreateSessionRequest) ApnRestriction() *ie.ApnRestriction {
	return c.apnRestriction
}
func (c *CreateSessionRequest) ApnAmbr() *ie.Ambr {
	return c.apnAmbr
}
func (c *CreateSessionRequest) Pco() *ie.PcoMsToNetwork {
	return c.pco
}
func (c *CreateSessionRequest) BearerContextTBC() *ie.BearerContextToBeCreatedWithinCSReq {
	return c.bearerContextTBC
}
func (c *CreateSessionRequest) Recovery() *ie.Recovery {
	return c.recovery
}
