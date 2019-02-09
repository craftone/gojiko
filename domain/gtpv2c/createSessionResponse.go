package gtpv2c

import (
	"errors"
	"fmt"
	"net"

	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/craftone/gojiko/domain/gtpv2c/ie/pco"
)

type CreateSessionResponse struct {
	header
	cause               *ie.Cause
	pgwCtrlFteid        *ie.Fteid
	paa                 *ie.Paa
	apnRestriction      *ie.ApnRestriction
	pco                 *ie.PcoNetworkToMs
	bearerContextCeated *ie.BearerContextCreatedWithinCSRes
	recovery            *ie.Recovery
}

type CreateSessionResponseArg struct {
	SgwCtrlTeid         gtp.Teid
	Cause               *ie.Cause
	PgwCtrlFteid        *ie.Fteid
	Paa                 *ie.Paa
	ApnRestriction      *ie.ApnRestriction
	Pco                 *ie.PcoNetworkToMs
	BearerContextCeated *ie.BearerContextCreatedWithinCSRes
	Recovery            *ie.Recovery
}

func NewCreateSessionResponse(seqNum uint32, csResArg CreateSessionResponseArg) (*CreateSessionResponse, error) {
	if err := checkCreateSessionResponseArg(csResArg); err != nil {
		return nil, err
	}

	return &CreateSessionResponse{
		newHeader(CreateSessionResponseNum, false, true, csResArg.SgwCtrlTeid, seqNum),
		csResArg.Cause,
		csResArg.PgwCtrlFteid,
		csResArg.Paa,
		csResArg.ApnRestriction,
		csResArg.Pco,
		csResArg.BearerContextCeated,
		csResArg.Recovery,
	}, nil
}

func checkCreateSessionResponseArg(csReqArg CreateSessionResponseArg) error {
	errMsgs := make([]string, 0)

	// Confirm mandatory IEs are exists
	if csReqArg.Cause == nil {
		errMsgs = append(errMsgs, "Cause")
	}
	if csReqArg.BearerContextCeated == nil {
		errMsgs = append(errMsgs, "Bearer Context Created")
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf("Some mandatory IEs are missing : %v", errMsgs)
	}

	// Confirm conditional IEs are exists in CSres Acceptance condition
	if csReqArg.Cause.Value().Type() == ie.CauseTypeAcceptance {
		if csReqArg.PgwCtrlFteid == nil {
			errMsgs = append(errMsgs, "PgwCtrlFteid")
		}
		if csReqArg.Paa == nil {
			errMsgs = append(errMsgs, "Paa")
		}
		if csReqArg.ApnRestriction == nil {
			errMsgs = append(errMsgs, "ApnRestriction")
		}
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

// MakeCSResArg makes a CreateSessionResponseArg.
// - PDN Type is always IPv4.
// - APN Restriction is always No Existing Contexts or Restriction (0).
// - APN AMBR is always up:4294967 kbps, down:4294967 kbps.
// - ChargingID is always 0x12345678.
func MakeCSResArg(sgwCtrlTEID gtp.Teid, cause ie.CauseValue,
	pgwCtrlIPv4 net.IP, pgwCtrlTeid gtp.Teid,
	pgwDataIPv4 net.IP, pgwDataTeid gtp.Teid,
	pdnIPv4Addr, priDNSIPv4Addr, secDNSIPv4Addr net.IP,
	ebi byte) (CreateSessionResponseArg, error) {

	causeIE, err := ie.NewCause(0, cause, false, false, false, nil)
	if err != nil {
		return CreateSessionResponseArg{}, err
	}

	pgwCtrlFteidIE, err := ie.NewFteid(1, pgwCtrlIPv4, nil, ie.S5S8PgwGtpCIf, pgwCtrlTeid)
	if err != nil {
		return CreateSessionResponseArg{}, err
	}

	paaIE, err := ie.NewPaa(0, ie.PdnTypeIPv4, pdnIPv4Addr, nil)
	if err != nil {
		return CreateSessionResponseArg{}, err
	}

	apnRestrictionIE, err := ie.NewApnRestriction(0, 0)
	if err != nil {
		return CreateSessionResponseArg{}, err
	}

	ipcpIE := pco.NewIpcp(pco.ConfigureNack, 0, priDNSIPv4Addr, secDNSIPv4Addr)
	dnsServerV4s := []*pco.DNSServerV4{
		pco.NewDNSServerV4(priDNSIPv4Addr),
		pco.NewDNSServerV4(secDNSIPv4Addr),
	}
	pcoNetworkToMsIE := pco.NewNetworkToMs(ipcpIE, dnsServerV4s, nil)
	pcoIE, err := ie.NewPcoNetworkToMs(0, pcoNetworkToMsIE)

	ebiIE, err := ie.NewEbi(0, ebi)
	if err != nil {
		return CreateSessionResponseArg{}, err
	}

	pgwDataFteidIE, err := ie.NewFteid(2, pgwDataIPv4, nil, ie.S5S8PgwGtpUIf, pgwDataTeid)
	if err != nil {
		return CreateSessionResponseArg{}, err
	}

	chargingIDIE, err := ie.NewChargingID(0, 0x12345678)
	if err != nil {
		return CreateSessionResponseArg{}, err
	}

	bearerContextCIE, err := ie.NewBearerContextCreatedWithinCSRes(
		ie.BearerContextCreatedWithinCSResArg{
			Ebi:          ebiIE,
			Cause:        causeIE,
			PgwDataFteid: pgwDataFteidIE,
			ChargingID:   chargingIDIE,
		})

	return CreateSessionResponseArg{
		SgwCtrlTeid:         sgwCtrlTEID,
		Cause:               causeIE,
		PgwCtrlFteid:        pgwCtrlFteidIE,
		Paa:                 paaIE,
		ApnRestriction:      apnRestrictionIE,
		Pco:                 pcoIE,
		BearerContextCeated: bearerContextCIE,
	}, nil
}

func (c *CreateSessionResponse) Marshal() []byte {
	body := make([]byte, 0, 300)

	if c.cause != nil {
		body = append(body, c.cause.Marshal()...)
	}
	if c.pgwCtrlFteid != nil {
		body = append(body, c.pgwCtrlFteid.Marshal()...)
	}
	if c.paa != nil {
		body = append(body, c.paa.Marshal()...)
	}
	if c.apnRestriction != nil {
		body = append(body, c.apnRestriction.Marshal()...)
	}
	if c.pco != nil {
		body = append(body, c.pco.Marshal()...)
	}
	if c.bearerContextCeated != nil {
		body = append(body, c.bearerContextCeated.Marshal()...)
	}
	if c.recovery != nil {
		body = append(body, c.recovery.Marshal()...)
	}

	return c.header.marshal(body)
}

func unmarshalCreateSessionResponse(h header, buf []byte) (*CreateSessionResponse, error) {
	if h.messageType != CreateSessionResponseNum {
		log.Panic("Invalid messageType")
	}

	csResArg := CreateSessionResponseArg{SgwCtrlTeid: h.teid}
	for len(buf) > 0 {
		msg, tail, err := ie.Unmarshal(buf, ie.CreateSessionResponse)
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
			csResArg.Cause = msg
		case *ie.Fteid:
			csResArg.PgwCtrlFteid = msg
		case *ie.Paa:
			csResArg.Paa = msg
		case *ie.ApnRestriction:
			csResArg.ApnRestriction = msg
		case *ie.PcoNetworkToMs:
			csResArg.Pco = msg
		case *ie.BearerContextCreatedWithinCSRes:
			csResArg.BearerContextCeated = msg
		case *ie.Recovery:
			csResArg.Recovery = msg
		default:
			log.Printf("Unkown IE : %v", msg)
		}
	}
	return NewCreateSessionResponse(h.seqNum, csResArg)
}

//
// Getters
//

func (c *CreateSessionResponse) Cause() *ie.Cause {
	return c.cause
}
func (c *CreateSessionResponse) PgwCtrlFteid() *ie.Fteid {
	return c.pgwCtrlFteid
}
func (c *CreateSessionResponse) Paa() *ie.Paa {
	return c.paa
}
func (c *CreateSessionResponse) ApnRestriction() *ie.ApnRestriction {
	return c.apnRestriction
}
func (c *CreateSessionResponse) Pco() *ie.PcoNetworkToMs {
	return c.pco
}
func (c *CreateSessionResponse) BearerContextCeated() *ie.BearerContextCreatedWithinCSRes {
	return c.bearerContextCeated
}
func (c *CreateSessionResponse) Recovery() *ie.Recovery {
	return c.recovery
}
