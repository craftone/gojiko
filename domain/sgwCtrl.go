package domain

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/craftone/gojiko/domain/apns"
	"github.com/craftone/gojiko/domain/gtp"
	"github.com/craftone/gojiko/domain/gtpv2c"
	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/sirupsen/logrus"
)

type SgwCtrl struct {
	*absSPgw
	*GtpSessionRepo
}

// newSgwCtrl creates a SgwCtrl and a paired SgwData that have same
// IP addr and recovery value.
func newSgwCtrl(addr net.UDPAddr, dataPort int, recovery byte) (*SgwCtrl, error) {
	myLog := log.WithFields(logrus.Fields{
		"addr":     addr,
		"recovery": recovery,
	})
	myLog.Info("A new SGW Ctrl has created")

	absSPgw, err := newAbsSPgw(addr, recovery, nil)
	if err != nil {
		return nil, err
	}
	gtpSessionRepo := newGtpSessionRepo()
	sgwCtrl := &SgwCtrl{absSPgw, gtpSessionRepo}

	sgwDataUDPAddr := net.UDPAddr{IP: addr.IP, Port: dataPort}
	sgwCtrl.pair, err = newSgwData(sgwDataUDPAddr, recovery, sgwCtrl)
	if err != nil {
		return nil, err
	}

	go sgwCtrl.sgwCtrlReceiverRoutine()
	go sgwCtrl.echoReceiver()

	return sgwCtrl, nil
}

/*CreateSession creates a new GTP session.
This method creates CreateSessionRequest Message,
sends the message to PGW , waits for CreateSessionResponse
and then returns the result (GsRes).

When pseudoSgwDataIP is nil, SGW-DATA F-TEID's IP Address in
a CreateSessionRequest Message this method will create
is same as SGW-CTRL's IP Address.

When pseudoSgwDataTEID is 0, SGW-DATA's TEID is generated
automatically.
*/
func (s *SgwCtrl) CreateSession(
	imsi, msisdn, mei, mcc, mnc, apnNI string,
	ebi byte, ratType byte,
	taiIE *ie.Tai, ecgiIE *ie.Ecgi,
	pseudoSgwDataIP *net.IP,
	pseudoSgwDataTEID gtp.Teid,
) (GsRes, *GtpSession, error) {
	// Query APN's IP address
	apn, err := apns.TheRepo().Find(apnNI, mcc, mnc)
	if err != nil {
		return GsRes{}, nil, err
	}
	pgwCtrlIPv4 := apn.GetIP()
	pgwCtrlAddr := net.UDPAddr{IP: pgwCtrlIPv4, Port: GtpControlPort}

	// Find or Create OpPgwCtrl
	_, err = s.findOrCreateOpSPgw(pgwCtrlAddr)
	if err != nil {
		return GsRes{}, nil, err
	}

	// Make SGW Ctrl F-TEID
	sgwCtrlFTEID, err := ie.NewFteid(0, s.laddr.IP, nil, ie.S5S8SgwGtpCIf, s.nextTeid())
	if err != nil {
		return GsRes{}, nil, err
	}

	// Make SGW Data F-TEID
	sgwData := s.pair
	sgwDataTeid := sgwData.nextTeid()
	var sgwDataIP net.IP
	if pseudoSgwDataIP == nil {
		sgwDataIP = sgwData.UDPAddr().IP
	} else {
		// when using external pseudo sgw-data
		sgwDataIP = *pseudoSgwDataIP
		if pseudoSgwDataTEID != 0 {
			sgwDataTeid = pseudoSgwDataTEID
		}
	}
	sgwDataFTEID, err := ie.NewFteid(0, sgwDataIP, nil, ie.S5S8SgwGtpUIf, sgwDataTeid)
	if err != nil {
		return GsRes{}, nil, err
	}

	// Make IMSI, MSISDN, etc
	imsiIE, err := ie.NewImsi(0, imsi)
	if err != nil {
		return GsRes{}, nil, err
	}

	msisdnIE, err := ie.NewMsisdn(0, msisdn)
	if err != nil {
		return GsRes{}, nil, err
	}

	meiIE, err := ie.NewMei(0, mei)
	if err != nil {
		return GsRes{}, nil, err
	}

	ebiIE, err := ie.NewEbi(0, ebi)
	if err != nil {
		return GsRes{}, nil, err
	}

	paaIE, err := ie.NewPaa(0, ie.PdnTypeIPv4, net.IPv4(0, 0, 0, 0), nil)
	if err != nil {
		return GsRes{}, nil, err
	}

	apnIE, err := ie.NewApn(0, apn.FullString())
	if err != nil {
		return GsRes{}, nil, err
	}

	ambrIE, err := ie.NewAmbr(0, 4294967, 4294967)
	if err != nil {
		return GsRes{}, nil, err
	}

	ratTypeIE, err := ie.NewRatType(0, ie.RatTypeValue(ratType))
	if err != nil {
		return GsRes{}, nil, err
	}

	servingNetworkID, err := ie.NewServingNetwork(0, mcc, mnc)
	if err != nil {
		return GsRes{}, nil, err
	}

	pdnType, err := ie.NewPdnType(0, 1)
	if err != nil {
		return GsRes{}, nil, err
	}

	// make a new session to the GTP Session Repo
	gsid, err := s.GtpSessionRepo.newSession(
		s, pgwCtrlIPv4,
		s.toSender,
		sgwCtrlFTEID,
		sgwDataFTEID,
		imsiIE,
		msisdnIE,
		meiIE,
		ebiIE,
		paaIE,
		apnIE,
		ambrIE,
		ratTypeIE,
		taiIE,
		ecgiIE,
		servingNetworkID,
		pdnType,
	)
	if err != nil {
		return GsRes{}, nil, err
	}

	// Make Create Session Request argument
	arg, err := NewCreateSessionReq(mcc, mnc, mei)
	if err != nil {
		return GsRes{}, nil, err
	}

	// Send CSreq and receive CSreq
	resChan := make(chan GsRes)
	session := s.GtpSessionRepo.FindBySessionID(gsid)
	go session.procCreateSession(arg, log, resChan)

	// Receive result of the process send CSreq and receive CSres
	res := <-resChan
	if res.err != nil {
		log.Error(res.err)
		s.GtpSessionRepo.deleteSession(session)
		return GsRes{}, nil, res.err
	}
	if res.Code != GsResOK {
		s.GtpSessionRepo.deleteSession(session)
	}

	return res, session, nil
}

func (s *SgwCtrl) TrackingAreaUpdateWithoutSgwRelocation(
	imsi string, ebi byte, taiIE *ie.Tai, ecgiIE *ie.Ecgi) (GsRes, error) {

	session := s.GtpSessionRepo.FindByImsiEbi(imsi, ebi)
	if session == nil {
		return GsRes{}, fmt.Errorf("There is no session whose imsi is %s and ebi is %d", imsi, ebi)
	}
	if session.Status() != GssConnected {
		return GsRes{}, NewInvalidGtpSessionStateError(GssConnected, session.Status())
	}

	// Send MBreq and receive MBres
	resChan := make(chan GsRes)
	go session.procTAUwoSgwRelocation(taiIE, ecgiIE, log, resChan)

	// Receive result of the process send MBreq and receive MBres
	res := <-resChan
	return res, nil
}

func (s *SgwCtrl) DeleteSession(imsi string, ebi byte) (GsRes, error) {
	session := s.GtpSessionRepo.FindByImsiEbi(imsi, ebi)
	if session == nil {
		return GsRes{}, fmt.Errorf("There is no session that's imsi is %s and ebi is %d", imsi, ebi)
	}
	if session.Status() != GssConnected {
		return GsRes{}, NewInvalidGtpSessionStateError(GssConnected, session.Status())
	}

	// Send DSreq and receive DSres
	resChan := make(chan GsRes)
	go session.procDeleteSession(resChan, log)

	// Receive result of the process send DSreq and receive DSres
	res := <-resChan
	s.GtpSessionRepo.deleteSession(session)
	if res.err != nil {
		log.Error(res.err)
		return GsRes{}, res.err
	}
	return res, nil
}

// sgwCtrlReceiverRoutine is for GoRoutine
func (s *SgwCtrl) sgwCtrlReceiverRoutine() {
	myLog := log.WithFields(logrus.Fields{
		"laddr":   s.laddr.String(),
		"routine": "SgwCtrlReceiver",
	})
	myLog.Info("Start a SGW Ctrl Receiver goroutine")

	buf := make([]byte, 2000)
	for {
		n, raddr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			myLog.Error(err)
			continue
		}
		myLog.Debugf("Received packet from %s : %v", raddr.String(), buf[:n])

		if n < 8 {
			myLog.Errorf("Too short packet : %v", buf[:n])
			continue
		}
		msgType := gtpv2c.MessageTypeNum(buf[1])

		switch msgType {
		case gtpv2c.EchoRequestNum:
			received := make([]byte, n)
			copy(received, buf[:n])
			s.toEchoReceiver <- UDPpacket{*raddr, received}
		case gtpv2c.EchoResponseNum:
			myLog.Error("Not yet implemented!")
		case gtpv2c.CreateSessionResponseNum,
			gtpv2c.DeleteSessionResponseNum,
			gtpv2c.DeleteBearerRequestNum:
			teid := gtp.Teid(binary.BigEndian.Uint32(buf[4:8]))
			sess := s.FindByCtrlTeid(teid)
			if sess == nil {
				myLog.Debugf("No session that have the teid : %04x", teid)
				continue
			}
			received := make([]byte, n)
			copy(received, buf[:n])
			sess.fromSgwCtrlReceiverChan <- UDPpacket{*raddr, received}
		default:
			myLog.Debugf("Unkown Message Type : %d", msgType)
		}
	}
}

// echoReceiver is for GoRoutine
func (s *SgwCtrl) echoReceiver() {
	myLog := log.WithFields(logrus.Fields{
		"laddr":   s.laddr.String(),
		"routine": "SPgwEchoReceiver",
	})
	myLog.Info("Start a SgwCtrl ECHO Receiver goroutine")

	for pkt := range s.toEchoReceiver {
		// ensure valid GTPv2-C ECHO Request
		req, _, err := gtpv2c.Unmarshal(pkt.body)
		if err != nil {
			myLog.Debugf("Received an invalid ECHO-C Request from %s", pkt.raddr.String())
			continue
		}

		myLog.Debugf("Received ECHO Request : %#v", req)

		// make ECHO Response
		echoRes, err := gtpv2c.NewEchoResponse(req.SeqNum(), s.recovery)
		if err != nil {
			myLog.Panicf("Making ECHO Response Failure : %v", err)
		}
		res := UDPpacket{
			raddr: pkt.raddr,
			body:  echoRes.Marshal(),
		}
		s.toSender <- res
	}
}
