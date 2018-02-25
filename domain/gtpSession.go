package domain

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/craftone/gojiko/gtp"

	"github.com/craftone/gojiko/config"
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
	repo   *GtpSessionRepo
	mtx    sync.RWMutex

	cmdReqChan           chan gtpSessionCmd
	cmdResChan           chan GscRes
	receiveCSresChan     chan *gtpv2c.CreateSessionResponse
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
	mei            *ie.Mei
	ebi            *ie.Ebi
	paa            *ie.Paa // at first, 0.0.0.0
	apn            *ie.Apn
	ambr           *ie.Ambr
	ratType        *ie.RatType
	servingNetwork *ie.ServingNetwork
	pdnType        *ie.PdnType
}

// this function is for GoRoutine
func (session *GtpSession) gtpSessionRoutine() {
	myLog := log.WithField("SessionID", session.id)
	myLog.Debug("Start a GTP session goroutine")

	for msg := range session.cmdReqChan {
		myLog.Debugf("Received CMD : %v", msg)

		switch cmd := msg.(type) {
		case createSessionReq:
			err := session.procCreateSession(cmd, myLog)
			if err != nil {
				log.Error(err)
			}
		}
	}
	myLog.Debug("Stop a GTP session goroutine")
}

// this function is for GoRoutine
func (session *GtpSession) receivePacketRoutine() {
	myLog := log.WithFields(logrus.Fields{
		"SessionID": session.ID(),
		"routine":   "ReceivePacket",
	})
	myLog.Debug("Start a GTP session's receive packet goroutine")

	for recv := range session.fromCtrlReceiverChan {
		// Ensure received packet has sent from correct PGW address
		if !recv.raddr.IP.Equal(session.pgwCtrlAddr.IP) ||
			recv.raddr.Port != session.pgwCtrlAddr.Port {
			myLog.Debugf("Received invalid GTPv2-C packet from unkown address : %#v , expected : %v", recv.raddr, session.pgwCtrlAddr)
			continue
		}

		// Unmarchal received packet
		msg, _, err := gtpv2c.Unmarshal(recv.body)
		if err != nil {
			myLog.Debugf("Received invalid GTPv2-C packet : %#v", err)
			continue
		}
		myLog.Debugf("Received GTPv2-C packet : %#v", msg)

		switch typedMsg := msg.(type) {
		case *gtpv2c.CreateSessionResponse:
			session.receiveCSresChan <- typedMsg
		case *gtpv2c.DeleteBearerRequest:
			err := session.procDeleteBearer(recv.raddr, typedMsg, myLog)
			if err != nil {
				myLog.Error(err)
			}
		default:
			myLog.Error("Don't know how to precess the packet")
		}
	}
}

func (session *GtpSession) procCreateSession(cmd createSessionReq, myLog *logrus.Entry) error {
	session.status = GssSendingCSReq
	seqNum := session.sgwCtrl.nextSeqNum()

	recoveryIE, err := ie.NewRecovery(0, session.sgwCtrl.recovery)
	if err != nil {
		return err
	}
	bearerContextTBCIE, err := ie.NewBearerContextToBeCreatedWithinCSReq(
		ie.BearerContextToBeCreatedWithinCSReqArg{
			Ebi:          session.ebi,
			BearerQoS:    cmd.bearerQoS,
			SgwDataFteid: session.sgwDataFTEID,
		})
	if err != nil {
		return err
	}

	csReqArg := gtpv2c.CreateSessionRequestArg{
		Imsi:             session.imsi,
		Msisdn:           session.msisdn,
		Mei:              cmd.mei,
		Uli:              cmd.uli,
		ServingNetwork:   session.servingNetwork,
		RatType:          session.ratType,
		Indication:       cmd.indication,
		SgwCtrlFteid:     session.sgwCtrlFTEID,
		Apn:              session.apn,
		SelectionMode:    cmd.selectionMode,
		PdnType:          session.pdnType,
		Paa:              session.paa,
		ApnRestriction:   cmd.apnRestriction,
		ApnAmbr:          session.ambr,
		Pco:              cmd.pco,
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

	var res GscRes
	afterChan := time.After(config.Gtpv2cTimeoutDuration())

	select {
	case csRes := <-session.receiveCSresChan:
		if csReq.SeqNum() != seqNum {
			res = GscRes{Code: GscResNG, Msg: "The response's sequense number is invalid"}
			break
		}
		causeType, causeMsg := ie.CauseDetail(csRes.Cause().Value())
		switch causeType {
		case ie.CauseTypeAcceptance:
			// Set PGW's F-TEIDs into the session
			session.pgwCtrlFTEID = csRes.PgwCtrlFteid()
			session.pgwDataFTEID = csRes.BearerContextCeated().PgwDataFteid()
			// Set PDN Address Allocation into the session
			session.paa = csRes.Paa()

			res = GscRes{Code: GscResOK, Msg: causeMsg, Session: session}
		case ie.CauseTypeRetryableRejection:
			res = GscRes{Code: GscResRetryableNG, Msg: causeMsg}
		default:
			res = GscRes{Code: GscResNG, Msg: causeMsg}
		}

	case <-afterChan:
		myLog.Error("The Create Session Response timed out")
		res = GscRes{Code: GscResTimeout, Msg: "Create Session Request timed out"}
	}
	session.cmdResChan <- res

	return nil
}

func (session *GtpSession) procDeleteBearer(raddr net.UDPAddr, dbReq *gtpv2c.DeleteBearerRequest, myLog *logrus.Entry) error {
	myLog.Debugf("DeleteBearer process")
	_, pgwTeid := session.PgwCtrlFTEID()
	dbRes, err := gtpv2c.NewDeleteBearerResponse(pgwTeid, dbReq.SeqNum(), ie.CauseRequestAccepted)
	if err != nil {
		return err
	}
	fmt.Println("Sending packet ...")
	session.toCtrlSenderChan <- UDPpacket{raddr, dbRes.Marshal()}
	fmt.Println("DONE.")
	err = session.repo.deleteSession(session.ID())
	if err != nil {
		return err
	}
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

func (s *GtpSession) Mei() string {
	return s.mei.Value()
}

func (s *GtpSession) Ebi() byte {
	return s.ebi.Value()
}

func (s *GtpSession) Paa() net.IP {
	return s.paa.IPv4()
}

func (s *GtpSession) Apn() string {
	return s.apn.String()
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

func (s *GtpSession) Mcc() string {
	return s.servingNetwork.Mcc()
}

func (s *GtpSession) Mnc() string {
	return s.servingNetwork.Mnc()
}

func (s *GtpSession) PdnType() string {
	return s.pdnType.String()
}
