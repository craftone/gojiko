package domain

import (
	"errors"
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
	id         SessionID
	status     GtpSessionStatus
	mtx4status sync.RWMutex

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

	udpFlow  *UdpEchoFlow
	mtx4flow sync.RWMutex
}

func (sess *GtpSession) changeState(curState, nextState GtpSessionStatus) error {
	sess.mtx4status.Lock()
	defer sess.mtx4status.Unlock()
	if sess.status != curState {
		return fmt.Errorf("Current State is not %v", curState)
	}
	sess.status = nextState
	log.WithField("SessionID", sess.id).Debugf("Change GTP session state : %v -> %v", curState, nextState)
	return nil
}

// this function is for GoRoutine
func (session *GtpSession) gtpSessionRoutine() {
	myLog := log.WithField("SessionID", session.id)
	myLog.Debug("Start a GTP session CMD goroutine")

	for msg := range session.cmdReqChan {
		myLog.Debugf("Received CMD : %v", msg)

		switch cmd := msg.(type) {
		case createSessionReq:
			err := session.procCreateSession(cmd, myLog)
			if err != nil {
				log.Error(err)
				session.changeState(GssSendingCSReq, GssIdle)
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
		if !recv.raddr.IP.Equal(session.pgwCtrlAddr.IP) {
			myLog.Debugf("Received invalid GTPv2-C packet from unkown address : %v , expected : %v", recv.raddr, session.pgwCtrlAddr)
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
	// Change state from IDLE to SENDING
	err := session.changeState(GssIdle, GssSendingCSReq)
	if err != nil {
		myLog.Debug(err)
	}

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

	raddr := session.pgwCtrlAddr

	// Send a CSReq packet to the PGW
	session.toCtrlSenderChan <- UDPpacket{raddr, csReqBin}

	var res GscRes
	afterChan := time.After(config.Gtpv2cTimeoutDuration())

	select {
	case csRes := <-session.receiveCSresChan:
		if csReq.SeqNum() != seqNum {
			res = GscRes{Code: GscResNG, Msg: "The response's sequense number is invalid"}
			session.changeState(GssSendingCSReq, GssIdle)
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

			session.changeState(GssSendingCSReq, GssConnected)
			res = GscRes{Code: GscResOK, Msg: causeMsg, Session: session}
		case ie.CauseTypeRetryableRejection:
			session.changeState(GssSendingCSReq, GssIdle)
			res = GscRes{Code: GscResRetryableNG, Msg: causeMsg}
		default:
			session.changeState(GssSendingCSReq, GssIdle)
			res = GscRes{Code: GscResNG, Msg: causeMsg}
		}

	case <-afterChan:
		myLog.Error("The Create Session Response timed out")
		session.changeState(GssSendingCSReq, GssIdle)
		res = GscRes{Code: GscResTimeout, Msg: "Create Session Request timed out"}
	}
	session.cmdResChan <- res

	return nil
}

func (session *GtpSession) procDeleteBearer(raddr net.UDPAddr, dbReq *gtpv2c.DeleteBearerRequest, myLog *logrus.Entry) error {
	_, pgwTeid := session.PgwCtrlFTEID()
	dbRes, err := gtpv2c.NewDeleteBearerResponse(pgwTeid, dbReq.SeqNum(), ie.CauseRequestAccepted, session.Ebi(), session.sgwCtrl.recovery)
	if err != nil {
		return err
	}
	session.toCtrlSenderChan <- UDPpacket{raddr, dbRes.Marshal()}
	myLog.Debugf("Send Delete Bearer Response : %#v", dbRes)
	err = session.sgwCtrl.GtpSessionRepo.deleteSession(session.ID())
	if err != nil {
		return err
	}
	myLog.Debug("Delete the sessions records")
	return nil
}

func (sess *GtpSession) setUdpFlow(udpEchoFlow *UdpEchoFlow) error {
	sess.mtx4flow.Lock()
	defer sess.mtx4flow.Unlock()
	if sess.udpFlow != nil {
		return errors.New("This session already have a UdpFlow")
	}
	sess.udpFlow = udpEchoFlow
	return nil
}

func (sess *GtpSession) NewUdpFlow(udpEchoFlowArg UdpEchoFlowArg) error {
	if sess.status != GssConnected {
		return errors.New("This session is not connected")
	}
	if udpEchoFlowArg.SendPacketSize < MIN_UDP_ECHO_PACKET_SIZE {
		return fmt.Errorf("SendPacketSize must be bigger than %d", MIN_UDP_ECHO_PACKET_SIZE)
	}
	if udpEchoFlowArg.RecvPacketSize < MIN_UDP_ECHO_PACKET_SIZE {
		return fmt.Errorf("RecvPacketSize must be bigger than %d", MIN_UDP_ECHO_PACKET_SIZE)
	}
	udpEchoFlow := &UdpEchoFlow{udpEchoFlowArg}
	err := sess.setUdpFlow(udpEchoFlow)
	if err != nil {
		return err
	}

	go sess.udpFlow.sender(sess)

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
