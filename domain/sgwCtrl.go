package domain

import (
	"encoding/binary"
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

	return sgwCtrl, nil
}

func (s *SgwCtrl) CreateSession(
	imsi, msisdn, mei, mcc, mnc, apnNI string,
	ebi byte,
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

	// Make SGW Ctrl F-TEID and SGW Data F-TEID
	sgwCtrlFTEID, err := ie.NewFteid(0, s.addr.IP, nil, ie.S5S8SgwGtpCIf, s.nextTeid())
	if err != nil {
		return GsRes{}, nil, err
	}

	sgwData := s.pair
	sgwDataFTEID, err := ie.NewFteid(0, sgwData.UDPAddr().IP, nil, ie.S5S8SgwGtpUIf, sgwData.nextTeid())
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

	ratTypeIE, err := ie.NewRatType(0, 6)
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
		servingNetworkID,
		pdnType,
	)
	if err != nil {
		return GsRes{}, nil, err
	}

	// Make GTP Session CMD of Create Session Request
	cmd, err := NewCreateSessionReq(mcc, mnc, mei)
	if err != nil {
		return GsRes{}, nil, err
	}

	// Send CSreq and receive CSreq
	resChan := make(chan GsRes)
	session := s.GtpSessionRepo.FindBySessionID(gsid)
	go session.procCreateSession(cmd, log, resChan)

	// Receive result of the process send CSreq and receive CSres
	res := <-resChan
	if res.err != nil {
		log.Error(res.err)
		s.GtpSessionRepo.deleteSession(session.id)
		return GsRes{}, nil, res.err
	}
	if res.Code != GsResOK {
		s.GtpSessionRepo.deleteSession(session.id)
	}

	return res, session, nil
}

// sgwCtrlReceiverRoutine is for GoRoutine
func (sgwCtrl *SgwCtrl) sgwCtrlReceiverRoutine() {
	myLog := log.WithFields(logrus.Fields{
		"laddr":   sgwCtrl.addr.String(),
		"routine": "SgwCtrlReceiver",
	})
	myLog.Info("Start a SGW Ctrl Receiver goroutine")

	buf := make([]byte, 2000)
	for {
		n, raddr, err := sgwCtrl.conn.ReadFromUDP(buf)
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
			sgwCtrl.toEchoReceiver <- UDPpacket{*raddr, received}
		case gtpv2c.EchoResponseNum:
			myLog.Error("Not yet implemented!")
			// Not yet be implemented
		case gtpv2c.CreateSessionResponseNum, gtpv2c.DeleteBearerRequestNum:
			teid := gtp.Teid(binary.BigEndian.Uint32(buf[4:8]))
			sess := sgwCtrl.FindByCtrlTeid(teid)
			if sess == nil {
				myLog.Debug("No session that have the teid : %04x", teid)
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
