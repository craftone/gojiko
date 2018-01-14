package domain

import (
	"net"
	"sync"
	"time"

	"github.com/craftone/gojiko/gtp"

	"github.com/craftone/gojiko/config"
	gsc "github.com/craftone/gojiko/domain/gtpSessionCmd"
	"github.com/sirupsen/logrus"

	"github.com/craftone/gojiko/gtpv2c"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

type GtpSessionStatus byte

const (
	GssIdle GtpSessionStatus = iota
	GssSendingCSReq
	GssConnected
)

type GtpSession struct {
	id     SessionID
	status GtpSessionStatus
	mtx    sync.RWMutex

	cmdReqChan           chan gsc.Cmd
	cmdResChan           chan gsc.Res
	toCtrlSenderChan     chan UDPpacket
	fromCtrlReceiverChan chan UDPpacket
	toDataSenderChan     chan UDPpacket
	fromDataReceiverChan chan UDPpacket

	sgwCtrl      *SgwCtrl
	sgwCtrlFTEID *ie.Fteid
	sgwDataFTEID *ie.Fteid
	pgwCtrlFTEID *ie.Fteid // at first, ZERO FTEID
	pgwDataFTEID *ie.Fteid // at first, nil

	pgwCtrlAddr net.UDPAddr
	pgwDataAddr net.UDPAddr
	sgwCtrlAddr net.UDPAddr
	sgwDataAddr net.UDPAddr

	imsi           *ie.Imsi
	msisdn         *ie.Msisdn
	ebi            *ie.Ebi
	paa            *ie.Paa // at first, 0.0.0.0
	apn            *ie.Apn
	ambr           *ie.Ambr
	ratType        *ie.RatType
	servingNetwork *ie.ServingNetwork
	pdnType        *ie.PdnType
}

// this function is for GoRoutine
func gtpSessionRoutine(session *GtpSession) {
	myLog := log.WithField("SessionID", session.id)
	myLog.Debug("Start a GTP session goroutine")

	for msg := range session.cmdReqChan {
		myLog.Debugf("Received CMD : %v", msg)

		switch cmd := msg.(type) {
		case gsc.CreateSessionReq:
			err := session.procCreateSession(cmd, myLog)
			if err != nil {
				log.Error(err)
			}
		}
	}
	myLog.Debug("Stop a GTP session goroutine")
}

func (session *GtpSession) procCreateSession(cmd gsc.CreateSessionReq, myLog *logrus.Entry) error {
	session.status = GssSendingCSReq
	seqNum := session.sgwCtrl.nextSeqNum()

	recoveryIE, err := ie.NewRecovery(0, session.sgwCtrl.recovery)
	if err != nil {
		return err
	}
	bearerContextTBCIE, err := ie.NewBearerContextToBeCreatedWithinCSReq(
		ie.BearerContextToBeCreatedWithinCSReqArg{
			Ebi:          session.ebi,
			BearerQoS:    cmd.BearerQoS,
			SgwDataFteid: session.sgwDataFTEID,
		})
	if err != nil {
		return err
	}

	csReqArg := gtpv2c.CreateSessionRequestArg{
		Imsi:             session.imsi,
		Msisdn:           session.msisdn,
		Mei:              cmd.Mei,
		Uli:              cmd.Uli,
		ServingNetwork:   session.servingNetwork,
		RatType:          session.ratType,
		Indication:       cmd.Indication,
		SgwCtrlFteid:     session.sgwCtrlFTEID,
		Apn:              session.apn,
		SelectionMode:    cmd.SelectionMode,
		PdnType:          session.pdnType,
		Paa:              session.paa,
		ApnRestriction:   cmd.ApnRestriction,
		ApnAmbr:          session.ambr,
		Pco:              cmd.Pco,
		BearerContextTBC: bearerContextTBCIE,
		Recovery:         recoveryIE,
	}

	csReq, err := gtpv2c.NewCreateSessionRequest(seqNum, csReqArg)
	if err != nil {
		return err
	}
	csReqBin := csReq.Marshal()

	// Send a CSReq packet to the PGW
	raddr := session.pgwCtrlAddr
	session.toCtrlSenderChan <- UDPpacket{raddr, csReqBin}

	var res gsc.Res
	afterChan := time.After(config.Gtpv2cTimeoutDuration())
retry:
	select {
	case recv := <-session.fromCtrlReceiverChan:
		myLog.Debugf("received packet from %v body: %v", recv.raddr, recv.body)

		// Ensure received packet has sent from correct PGW address
		if !recv.raddr.IP.Equal(session.pgwCtrlAddr.IP) ||
			recv.raddr.Port != session.pgwCtrlAddr.Port {
			myLog.Debugf("Received invalid GTPv2-C packet from unkown address : %v , expected : %v", recv.raddr, session.pgwCtrlAddr)
			goto retry
		}

		// Unmarchal received packet
		msg, _, err := gtpv2c.Unmarshal(recv.body)
		if err != nil {
			myLog.Debugf("Received invalid GTPv2-C packet")
			goto retry
		}
		myLog.Debugf("received GTPv2-C packet : %v", msg)

		// Ensure received packete is a Create Session Response
		csres, ok := msg.(*gtpv2c.CreateSessionResponse)
		if !ok {
			myLog.Debugf("Received packet is not a Create Session Response message.")
			goto retry
		}

		causeType, causeMsg := ie.CauseDetail(csres.Cause().Value())
		switch causeType {
		case ie.CauseTypeAcceptance:
			// Set PGW's F-TEIDs into the session
			session.pgwCtrlFTEID = csres.PgwCtrlFteid()
			session.pgwDataFTEID = csres.BearerContextCeated().PgwDataFteid()
			// Set PDN Address Allocation into the session
			session.paa = csres.Paa()

			res = gsc.Res{Code: gsc.ResOK, Msg: causeMsg}
		case ie.CauseTypeRetryableRejection:
			res = gsc.Res{Code: gsc.ResRetryableNG, Msg: causeMsg}
		default:
			res = gsc.Res{Code: gsc.ResNG, Msg: causeMsg}
		}

	case <-afterChan:
		myLog.Error("The Create Session Response timed out")
		res = gsc.Res{Code: gsc.ResTimeout, Msg: "Timeout"}
	}
	session.cmdResChan <- res

	return nil
}

// setter & getter

func (s *GtpSession) ID() SessionID {
	return s.id
}

func (s *GtpSession) Status() GtpSessionStatus {
	return s.status
}

func (s *GtpSession) SgwCtrlFTEID() (net.IP, gtp.Teid) {
	return s.sgwCtrlFTEID.Ipv4(), s.sgwCtrlFTEID.Teid()
}

func (s *GtpSession) SgwDataFTEID() (net.IP, gtp.Teid) {
	return s.sgwDataFTEID.Ipv4(), s.sgwDataFTEID.Teid()
}

func (s *GtpSession) PgwCtrlFTEID() (net.IP, gtp.Teid) {
	return s.pgwCtrlFTEID.Ipv4(), s.pgwCtrlFTEID.Teid()
}

func (s *GtpSession) PgwDataFTEID() (net.IP, gtp.Teid) {
	return s.pgwDataFTEID.Ipv4(), s.pgwDataFTEID.Teid()
}

func (s *GtpSession) Imsi() string {
	return s.imsi.Value()
}

func (s *GtpSession) Msisdn() string {
	return s.msisdn.Value()
}

func (s *GtpSession) Ebi() byte {
	return s.ebi.Value()
}

func (s *GtpSession) Paa() net.IP {
	return s.paa.IPv4()
}

func (s *GtpSession) Apn() string {
	return s.apn.Value()
}

func (s *GtpSession) Ambr() string {
	return s.ambr.String()
}

func (s *GtpSession) RatType() string {
	return s.ratType.String()
}

func (s *GtpSession) ServingNetwork() string {
	return s.servingNetwork.String()
}

func (s *GtpSession) PdnType() string {
	return s.pdnType
}
