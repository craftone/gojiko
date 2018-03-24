package domain

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/craftone/gojiko/config"
	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/craftone/gojiko/domain/stats"
	humanize "github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
)

type GtpSessionStatus byte

const (
	GssNewed GtpSessionStatus = iota
	GssCSReqSending
	GssCSReqSend
	GssCSResReceived
	GssConnected
	GssDSReqSending
	GssDSReqSend
	GssDSResReceived
)

var gssString = map[GtpSessionStatus]string{
	GssNewed:         "Newed",
	GssCSReqSending:  "CSReqSending",
	GssCSReqSend:     "CSReqSend",
	GssCSResReceived: "CSResReceived",
	GssConnected:     "Connected",
	GssDSReqSending:  "DSReqSending",
	GssDSReqSend:     "DSReqSend",
	GssDSResReceived: "DSResReceived",
}

type GtpSession struct {
	id         SessionID
	status     GtpSessionStatus
	mtx4status sync.RWMutex

	receiveCSresChan        chan *gtpv2c.CreateSessionResponse
	toSgwCtrlSenderChan     chan UDPpacket
	fromSgwCtrlReceiverChan chan UDPpacket
	toSgwDataSenderChan     chan UDPpacket
	fromSgwDataReceiverChan chan UDPpacket

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

	udpFlow     *UdpEchoFlow
	lastUDPFlow *UdpEchoFlow
	mtx4flow    sync.RWMutex
}

func (s *GtpSession) changeState(curState, nextState GtpSessionStatus) error {
	s.mtx4status.Lock()
	defer s.mtx4status.Unlock()
	if s.status != curState {
		return fmt.Errorf("Cannot change state from %s to %s : current State is %s",
			gssString[curState], gssString[nextState], gssString[s.status])
	}
	s.status = nextState
	log.WithField("SessionID", s.id).Debugf("Change GTP session state : %s -> %s", gssString[curState], gssString[nextState])
	return nil
}

// this function is for GoRoutine
func (s *GtpSession) receiveCtrlPacketRoutine() {
	myLog := log.WithFields(logrus.Fields{
		"SessionID": s.ID(),
		"routine":   "CtrlPacketReceiver",
	})
	myLog.Debug("Start a GTP session's ctrl packet receiver")

	for recv := range s.fromSgwCtrlReceiverChan {
		// Ensure received packet has sent from correct PGW address
		if !recv.raddr.IP.Equal(s.pgwCtrlAddr.IP) {
			myLog.Debugf("Received invalid GTPv2-C packet from unkown address : %s , expected : %s", recv.raddr.String(), s.pgwCtrlAddr.String())
			continue
		}

		// Unmarshal received packet
		msg, _, err := gtpv2c.Unmarshal(recv.body)
		if err != nil {
			myLog.Debugf("Received invalid GTPv2-C packet : %s", err)
			continue
		}
		myLog.Debugf("Received GTPv2-C packet : %#v", msg)

		switch typedMsg := msg.(type) {
		case *gtpv2c.CreateSessionResponse:
			err := s.changeState(GssCSReqSend, GssCSResReceived)
			if err != nil {
				myLog.Error("Received CreateSessionResponse in unexpected state")
			} else {
				s.receiveCSresChan <- typedMsg
			}
		case *gtpv2c.DeleteBearerRequest:
			err := s.procDeleteBearer(recv.raddr, typedMsg, myLog)
			if err != nil {
				myLog.Error(err)
			}
		default:
			myLog.Error("Don't know how to precess the packet")
		}
	}
}

// this function is for GoRoutine
func (s *GtpSession) receiveDataPacketRoutine() {
	log := log.WithFields(logrus.Fields{
		"SessionID": s.ID(),
		"routine":   "DataPacketReceiver",
	})
	log.Debug("Start a GTP session's data packet receiver")

	for recv := range s.fromSgwDataReceiverChan {
		// Ensure received packet has sent from correct PGW address
		if !recv.raddr.IP.Equal(s.pgwDataAddr.IP) {
			log.Debugf("Received invalid GTPv1-U packet from unkown address : %s , expected : %s", recv.raddr.String(), s.pgwDataAddr.String())
			continue
		}

		if s.udpFlow == nil {
			log.Debug("There is no flow process in this session")
			continue
		}

		s.udpFlow.fromSessDataReceiver <- recv
	}
	log.Debug("End a GTP session's data packet receiver")
}

// procCreateSession is for gorutine
func (s *GtpSession) procCreateSession(cmd createSessionReq, myLog *logrus.Entry, gscResChan chan GsRes) {
	myLog = myLog.WithField("SessionID", s.id)
	myLog = myLog.WithField("func", "procCreateSession")

	// Change state from IDLE to SENDING
	err := s.changeState(GssNewed, GssCSReqSending)
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}

	// Prepare CreateSessionRequest message
	seqNum := s.sgwCtrl.nextSeqNum()

	recoveryIE, err := ie.NewRecovery(0, s.sgwCtrl.recovery)
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}
	bearerContextTBCIE, err := ie.NewBearerContextToBeCreatedWithinCSReq(
		ie.BearerContextToBeCreatedWithinCSReqArg{
			Ebi:          s.ebi,
			BearerQoS:    cmd.bearerQoS,
			SgwDataFteid: s.sgwDataFTEID,
		})
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}

	csReqArg := gtpv2c.CreateSessionRequestArg{
		Imsi:             s.imsi,
		Msisdn:           s.msisdn,
		Mei:              cmd.mei,
		Uli:              cmd.uli,
		ServingNetwork:   s.servingNetwork,
		RatType:          s.ratType,
		Indication:       cmd.indication,
		SgwCtrlFteid:     s.sgwCtrlFTEID,
		Apn:              s.apn,
		SelectionMode:    cmd.selectionMode,
		PdnType:          s.pdnType,
		Paa:              s.paa,
		ApnRestriction:   cmd.apnRestriction,
		ApnAmbr:          s.ambr,
		Pco:              cmd.pco,
		BearerContextTBC: bearerContextTBCIE,
		Recovery:         recoveryIE,
	}

	csReq, err := gtpv2c.NewCreateSessionRequest(seqNum, csReqArg)
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}
	csReqBin := csReq.Marshal()

	raddr := s.pgwCtrlAddr

	// Send a Create Session Request message to the PGW
	retryCount := 0
retry:
	s.toSgwCtrlSenderChan <- UDPpacket{raddr, csReqBin}
	err = s.changeState(GssCSReqSending, GssCSReqSend)
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}
	timeoutChan := time.After(config.Gtpv2cTimeoutDuration())

loop:
	select {
	case csRes := <-s.receiveCSresChan:
		if csReq.SeqNum() != seqNum {
			myLog.Debug("The response's sequense number is invalid")
			goto loop
		}
		causeType, causeMsg := ie.CauseDetail(csRes.Cause().Value())
		switch causeType {
		case ie.CauseTypeAcceptance:
			// Set PGW's F-TEIDs into the session
			s.pgwCtrlFTEID = csRes.PgwCtrlFteid()
			s.pgwDataFTEID = csRes.BearerContextCeated().PgwDataFteid()
			// Set PDN Address Allocation into the session
			s.paa = csRes.Paa()

			if err = s.changeState(GssCSResReceived, GssConnected); err != nil {
				gscResChan <- GsRes{err: err}
				return
			}
			gscResChan <- GsRes{Code: GsResOK, Msg: causeMsg}
		case ie.CauseTypeRetryableRejection:
			gscResChan <- GsRes{Code: GsResRetryableNG, Msg: causeMsg}
		default:
			gscResChan <- GsRes{Code: GsResNG, Msg: causeMsg}
		}

	case <-timeoutChan:
		myLog.Error("The Create Session Response timed out")
		s.changeState(GssCSReqSend, GssCSReqSending)
		retryCount++
		if retryCount <= config.Gtpv2cRetry() {
			log.Debugf("Create Session Response timed out and retry : %s time", humanize.Ordinal(retryCount))
			goto retry
		}
		gscResChan <- GsRes{Code: GsResTimeout, Msg: "Create Session Response timed out and retry out"}
	}
}

func (s *GtpSession) procDeleteBearer(raddr net.UDPAddr, dbReq *gtpv2c.DeleteBearerRequest, myLog *logrus.Entry) error {
	_, pgwTeid := s.PgwCtrlFTEID()
	dbRes, err := gtpv2c.NewDeleteBearerResponse(pgwTeid, dbReq.SeqNum(), ie.CauseRequestAccepted, s.Ebi(), s.sgwCtrl.recovery)
	if err != nil {
		return err
	}
	s.toSgwCtrlSenderChan <- UDPpacket{raddr, dbRes.Marshal()}
	myLog.Infof("Send Delete Bearer Response : %#v", dbRes)
	err = s.sgwCtrl.GtpSessionRepo.deleteSession(s.ID())
	if err != nil {
		return err
	}
	myLog.Info("Delete the sessions records")
	return nil
}

func (s *GtpSession) setUdpFlow(udpEchoFlow *UdpEchoFlow) error {
	s.mtx4flow.Lock()
	defer s.mtx4flow.Unlock()
	if s.udpFlow != nil {
		return errors.New("This session already have a UdpFlow")
	}
	s.udpFlow = udpEchoFlow
	return nil
}

func (s *GtpSession) NewUdpFlow(udpEchoFlowArg UdpEchoFlowArg) error {
	if s.status != GssConnected {
		return errors.New("This session is not connected")
	}
	if udpEchoFlowArg.SendPacketSize < MIN_UDP_ECHO_PACKET_SIZE {
		return fmt.Errorf("SendPacketSize must be bigger than %d", MIN_UDP_ECHO_PACKET_SIZE)
	}
	if udpEchoFlowArg.RecvPacketSize < MIN_UDP_ECHO_PACKET_SIZE {
		return fmt.Errorf("RecvPacketSize must be bigger than %d", MIN_UDP_ECHO_PACKET_SIZE)
	}
	ctx, cancel := context.WithCancel(context.Background())
	udpEchoFlow := &UdpEchoFlow{
		Arg:                  udpEchoFlowArg,
		session:              s,
		fromSessDataReceiver: make(chan UDPpacket, 100),
		ctxCencel:            cancel,
		stats:                stats.NewFlowStats(ctx),
	}
	err := s.setUdpFlow(udpEchoFlow)
	if err != nil {
		return err
	}

	go s.udpFlow.sender(ctx)
	go s.udpFlow.receiver(ctx)

	return nil
}

func (s *GtpSession) UDPFlow() (*UdpEchoFlow, bool) {
	s.mtx4flow.RLock()
	defer s.mtx4flow.RUnlock()
	if s.udpFlow == nil {
		return nil, false
	}
	return s.udpFlow, true
}

func (s *GtpSession) LastUDPFlow() (*UdpEchoFlow, bool) {
	s.mtx4flow.RLock()
	defer s.mtx4flow.RUnlock()
	if s.lastUDPFlow == nil {
		return nil, false
	}
	return s.lastUDPFlow, true
}

func (s *GtpSession) StopUDPFlow() error {
	s.mtx4flow.Lock()
	defer s.mtx4flow.Unlock()
	if s.udpFlow == nil {
		return errors.New("This session already stopped a UdpFlow")
	}
	s.udpFlow.ctxCencel()
	s.udpFlow.stats.SendTimeMsg(stats.EndTime, time.Now())
	s.lastUDPFlow = s.udpFlow
	s.udpFlow = nil
	return nil
}

// setter & getter

func (s *GtpSession) ID() SessionID {
	return s.id
}

func (s *GtpSession) Status() GtpSessionStatus {
	s.mtx4status.RLock()
	defer s.mtx4status.RUnlock()
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
