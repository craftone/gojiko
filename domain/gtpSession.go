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
	GssMBReqSending
	GssMBReqSend
	GssMBResReceived
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
	GssMBReqSending:  "MBReqSending",
	GssMBReqSend:     "MBReqSend",
	GssMBResReceived: "MBResReceived",
	GssDSReqSending:  "DSReqSending",
	GssDSReqSend:     "DSReqSend",
	GssDSResReceived: "DSResReceived",
}

func (s GtpSessionStatus) String() string {
	return gssString[s]
}

type GtpSession struct {
	id         SessionID
	status     GtpSessionStatus
	mtx4status sync.RWMutex

	receiveCSresChan        chan *gtpv2c.CreateSessionResponse
	receiveMBresChan        chan *gtpv2c.ModifyBearerResponse
	receiveDSresChan        chan *gtpv2c.DeleteSessionResponse
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
	tai            *ie.Tai
	ecgi           *ie.Ecgi

	udpFlow     *UdpEchoFlow
	lastUDPFlow *UdpEchoFlow
	mtx4flow    sync.RWMutex
}

func (s *GtpSession) changeState(curState, nextState GtpSessionStatus) error {
	s.mtx4status.Lock()
	defer s.mtx4status.Unlock()
	if s.status != curState {
		return fmt.Errorf("Cannot change state from %s to %s : current State is %s",
			curState.String(), nextState.String(), s.status.String())
	}
	s.status = nextState
	log.WithField("SessionID", s.id).Debugf("Change GTP session state : %s -> %s",
		curState.String(), nextState.String())
	return nil
}

// this function is for GoRoutine
func (s *GtpSession) receiveCtrlPacketRoutine() {
	log := log.WithFields(logrus.Fields{
		"SessionID": s.ID(),
		"routine":   "CtrlPacketReceiver",
	})
	log.Debug("Start a GTP session's ctrl packet receiver")

	for recv := range s.fromSgwCtrlReceiverChan {
		// Ensure received packet has sent from correct PGW address
		if !recv.raddr.IP.Equal(s.pgwCtrlAddr.IP) {
			log.Debugf("Received invalid GTPv2-C packet from unkown address : %s , expected : %s", recv.raddr.String(), s.pgwCtrlAddr.String())
			continue
		}

		// Unmarshal received packet
		msg, _, err := gtpv2c.Unmarshal(recv.body)
		if err != nil {
			log.Debugf("Received invalid GTPv2-C packet : %s", err)
			continue
		}
		log.Debugf("Received GTPv2-C packet : %#v", msg)

		switch typedMsg := msg.(type) {
		case *gtpv2c.CreateSessionResponse:
			err := s.changeState(GssCSReqSend, GssCSResReceived)
			if err != nil {
				log.Error("Received CreateSessionResponse in unexpected state")
			} else {
				s.receiveCSresChan <- typedMsg
			}
		case *gtpv2c.ModifyBearerResponse:
			err := s.changeState(GssMBReqSend, GssMBResReceived)
			if err != nil {
				log.Error("Received ModifyBearerResponse in unexpected state")
			} else {
				s.receiveMBresChan <- typedMsg
			}
		case *gtpv2c.DeleteSessionResponse:
			err := s.changeState(GssDSReqSend, GssDSResReceived)
			if err != nil {
				log.Error("Received DeleteSessionResponse in unexpected state")
			} else {
				s.receiveDSresChan <- typedMsg
			}
		case *gtpv2c.DeleteBearerRequest:
			err := s.procDeleteBearer(recv.raddr, typedMsg, log)
			if err != nil {
				log.Error(err)
			}
		default:
			log.Error("Don't know how to precess the packet")
		}
	}
	log.Debug("End a GTP session's control packet receiver")
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
	myLog = myLog.WithField("routine", "procCreateSession")

	// Change state from NEWED to SENDING
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
		causeValue := csRes.Cause().Value()
		gsRes := GsRes{Value: causeValue, Msg: causeValue.Detail()}
		switch gsRes.Value.Type() {
		case ie.CauseTypeAcceptance:
			// Set PGW's F-TEIDs into the session
			s.pgwCtrlFTEID = csRes.PgwCtrlFteid()
			s.pgwDataFTEID = csRes.BearerContextCeated().PgwDataFteid()
			// Set PGW's Data Addr into the session
			s.pgwDataAddr.IP = s.pgwDataFTEID.Ipv4()
			// Set PDN Address Allocation into the session
			s.paa = csRes.Paa()

			if err = s.changeState(GssCSResReceived, GssConnected); err != nil {
				gscResChan <- GsRes{err: err}
				return
			}
			gsRes.Code = GsResOK
		case ie.CauseTypeRetryableRejection:
			gsRes.Code = GsResRetryableNG
		default:
			gsRes.Code = GsResNG
		}
		gscResChan <- gsRes

	case <-timeoutChan:
		myLog.Info("Waiting for Create Session Response is timed out")
		err = s.changeState(GssCSReqSend, GssCSReqSending)
		if err != nil {
			myLog.Debugf("Current state is not CSReqSend ( is %s), so it seemed to received packet", s.Status().String())
			timeoutChan = time.After(time.Second)
			goto loop
		}
		retryCount++
		if retryCount <= config.Gtpv2cRetry() {
			log.Debugf("Waiting for Create Session Response timed out and retry : %s time", humanize.Ordinal(retryCount))
			goto retry
		}
		gscResChan <- GsRes{Code: GsResTimeout, Msg: "Waiting for Create Session Response timed out and retry out"}
	}
}

// procTAUwoSgwRelocation is for goroutine
func (s *GtpSession) procTAUwoSgwRelocation(taiIE *ie.Tai, ecgiIE *ie.Ecgi,
	log *logrus.Entry, gscResChan chan GsRes) {
	log = log.WithField("SessionID", s.id)
	log = log.WithField("routine", "procTAUwoSgwRelocation")

	// Change state from CONNECTED to SENDING
	err := s.changeState(GssConnected, GssMBReqSending)
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}

	// Prepare DeleteSessionRequest message
	seqNum := s.sgwCtrl.nextSeqNum()

	_, pgwTeid := s.PgwCtrlFTEID()
	uliArg := ie.UliArg{
		Tai:  taiIE,
		Ecgi: ecgiIE,
	}
	uliIE, err := ie.NewUli(0, uliArg)
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}

	indicationIE, err := ie.NewIndication(0, ie.IndicationArg{
		CLII: true,
	})
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}

	mbReqArg := gtpv2c.ModifyBearerRequestArg{
		Imsi:         s.imsi,
		PgwCtrlTeid:  pgwTeid,
		Uli:          uliIE,
		Indication:   indicationIE,
		SgwCtrlFteid: s.sgwCtrlFTEID,
	}
	mbReq, err := gtpv2c.NewModifyBearerRequest(seqNum, mbReqArg)
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}
	mbReqBin := mbReq.Marshal()
	raddr := s.pgwCtrlAddr

	// Send a Modify Bearer Request message to the PGW
	retryCount := 0
retry:
	s.toSgwCtrlSenderChan <- UDPpacket{raddr, mbReqBin}
	if err = s.changeState(GssMBReqSending, GssMBReqSend); err != nil {
		gscResChan <- GsRes{err: err}
		return
	}
	timeoutChan := time.After(config.Gtpv2cTimeoutDuration())

loop:
	select {
	case mbRes := <-s.receiveMBresChan:
		if mbRes.SeqNum() != seqNum {
			log.Debug("The response's sequense number is invalid")
			goto loop
		}
		causeValue := mbRes.Cause().Value()
		if err = s.changeState(GssMBResReceived, GssConnected); err != nil {
			gscResChan <- GsRes{err: err}
			return
		}
		gsRes := GsRes{Value: causeValue, Msg: causeValue.Detail()}
		switch gsRes.Value.Type() {
		case ie.CauseTypeAcceptance:
			// update session information
			s.tai = taiIE
			s.ecgi = ecgiIE

			gsRes.Code = GsResOK
		case ie.CauseTypeRetryableRejection:
			gsRes.Code = GsResRetryableNG
		default:
			gsRes.Code = GsResNG
		}
		gscResChan <- gsRes

	case <-timeoutChan:
		log.Info("Waiting for Modify Bearer Response is timed out")
		if err = s.changeState(GssMBReqSend, GssMBReqSending); err != nil {
			log.Debugf("Current state is not MBReqSend ( is %s), so it seemed to received packet", s.Status().String())
			timeoutChan = time.After(time.Second)
			goto loop
		}
		retryCount++
		if retryCount <= config.Gtpv2cRetry() {
			log.Debugf("Waiting for Modify Bearer Response timed out and retry : %s time", humanize.Ordinal(retryCount))
			goto retry
		}
		if err = s.changeState(GssMBReqSending, GssConnected); err != nil {
			gscResChan <- GsRes{err: err}
			return
		}
		gscResChan <- GsRes{Code: GsResTimeout, Msg: "Waiting for Modify Bearer Response timed out and retry out"}
	}
}

// procDeleteSession is for goroutine
func (s *GtpSession) procDeleteSession(gscResChan chan GsRes, log *logrus.Entry) {
	log = log.WithField("SessionID", s.id)
	log = log.WithField("routine", "procDeleteSession")

	// Change state from CONNECTED to SENDING
	err := s.changeState(GssConnected, GssDSReqSending)
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}

	// Prepare DeleteSessionRequest message
	seqNum := s.sgwCtrl.nextSeqNum()

	_, pgwTeid := s.PgwCtrlFTEID()
	dsReq, err := gtpv2c.NewDeleteSessionRequest(pgwTeid, seqNum, s.Ebi())
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}
	dsReqBin := dsReq.Marshal()
	raddr := s.pgwCtrlAddr

	// Send a Delete Session Request message to the PGW
	retryCount := 0
retry:
	s.toSgwCtrlSenderChan <- UDPpacket{raddr, dsReqBin}
	err = s.changeState(GssDSReqSending, GssDSReqSend)
	if err != nil {
		gscResChan <- GsRes{err: err}
		return
	}
	timeoutChan := time.After(config.Gtpv2cTimeoutDuration())

loop:
	select {
	case dsRes := <-s.receiveDSresChan:
		if dsReq.SeqNum() != seqNum {
			log.Debug("The response's sequense number is invalid")
			goto loop
		}
		causeValue := dsRes.Cause().Value()
		gsRes := GsRes{Value: causeValue, Msg: causeValue.Detail()}
		switch gsRes.Value.Type() {
		case ie.CauseTypeAcceptance:
			gsRes.Code = GsResOK
		case ie.CauseTypeRetryableRejection:
			gsRes.Code = GsResRetryableNG
		default:
			gsRes.Code = GsResNG
		}
		gscResChan <- gsRes

	case <-timeoutChan:
		log.Info("Waiting for Delete Session Response is timed out")
		err = s.changeState(GssDSReqSend, GssDSReqSending)
		if err != nil {
			log.Debugf("Current state is not DSReqSend ( is %s), so it seemed to received packet", s.Status().String())
			timeoutChan = time.After(time.Second)
			goto loop
		}
		retryCount++
		if retryCount <= config.Gtpv2cRetry() {
			log.Debugf("Waiting for Delete Session Response timed out and retry : %s time", humanize.Ordinal(retryCount))
			goto retry
		}
		gscResChan <- GsRes{Code: GsResTimeout, Msg: "Waiting for Delete Session Response timed out and retry out"}
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
	err = s.sgwCtrl.GtpSessionRepo.deleteSession(s)
	if err != nil {
		return err
	}
	myLog.Info("Delete the bearer(session)'s records")
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

func (s *GtpSession) RatTypeValue() ie.RatTypeValue {
	return s.ratType.Value()
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

func (s *GtpSession) TaiMcc() string {
	return s.tai.Mcc()
}

func (s *GtpSession) TaiMnc() string {
	return s.tai.Mnc()
}

func (s *GtpSession) TaiTac() uint16 {
	return s.tai.Tac()
}

func (s *GtpSession) EcgiMcc() string {
	return s.ecgi.Mcc()
}

func (s *GtpSession) EcgiMnc() string {
	return s.ecgi.Mnc()
}

func (s *GtpSession) EcgiEci() uint32 {
	return s.ecgi.Eci()
}
