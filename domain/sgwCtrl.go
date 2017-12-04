package domain

import (
	"fmt"
	"net"

	"github.com/craftone/gojiko/domain/gtpSessionCmd"

	"github.com/craftone/gojiko/domain/apns"
	"github.com/craftone/gojiko/gtpv2c/ie"
	"github.com/sirupsen/logrus"
)

type SgwCtrl struct {
	*absSPgw
}

// newSgwCtrl creates a SgwCtrl and a paired SgwData that have same
// IP addr and recovery value.
func newSgwCtrl(addr net.UDPAddr, dataPort int, recovery byte) (*SgwCtrl, error) {
	myLog := log.WithFields(logrus.Fields{
		"addr":     addr,
		"recovery": recovery,
	})
	myLog.Debug("A new SGW Ctrl has created")

	absSPgw, err := newAbsSPgw(addr, recovery, nil)
	if err != nil {
		return nil, err
	}
	sgwCtrl := &SgwCtrl{absSPgw}

	sgwDataUDPAddr := net.UDPAddr{IP: addr.IP, Port: dataPort}
	sgwCtrl.pair, err = newSgwData(sgwDataUDPAddr, recovery, sgwCtrl)
	if err != nil {
		return nil, err
	}

	return sgwCtrl, nil
}

func (s *SgwCtrl) CreateSession(
	imsi, msisdn, mei, mcc, mnc, apnNI string,
	ebi byte,
) error {
	// Query APN's IP address
	apn, err := apns.TheRepo().Find(apnNI, mcc, mnc)
	if err != nil {
		return err
	}
	pgwCtrlIPv4 := apn.GetIP()
	pgwCtrlAddr := net.UDPAddr{IP: pgwCtrlIPv4, Port: GtpControlPort}

	// Find or Create OpPgwCtrl
	_, err = s.findOrCreateOpSPgw(pgwCtrlAddr)
	if err != nil {
		return err
	}

	// Take SGW Ctrl F-TEID and SGW Data F-TEID
	sgwCtrlFTEID, err := ie.NewFteid(0, s.addr.IP, nil, ie.S5S8SgwGtpCIf, s.nextTeid())
	if err != nil {
		return err
	}

	sgwData := s.pair
	sgwDataFTEID, err := ie.NewFteid(0, sgwData.UDPAddr().IP, nil, ie.S5S8SgwGtpUIf, sgwData.nextTeid())
	if err != nil {
		return err
	}

	// make IMSI, MSISDN, etc
	imsiIE, err := ie.NewImsi(0, imsi)
	if err != nil {
		return err
	}

	msisdnIE, err := ie.NewMsisdn(0, msisdn)
	if err != nil {
		return err
	}

	ebiIE, err := ie.NewEbi(0, ebi)
	if err != nil {
		return err
	}

	paaIE, err := ie.NewPaa(0, ie.PdnTypeIPv4, net.IPv4(0, 0, 0, 0), nil)
	if err != nil {
		return err
	}

	apnIE, err := ie.NewApn(0, apn.FullString())
	if err != nil {
		return err
	}

	ambrIE, err := ie.NewAmbr(0, 4294967, 4294967)
	if err != nil {
		return err
	}

	ratTypeIE, err := ie.NewRatType(0, 6)
	if err != nil {
		return err
	}

	servingNetworkID, err := ie.NewServingNetwork(0, mcc, mnc)
	if err != nil {
		return err
	}

	pdnType, err := ie.NewPdnType(0, 1)
	if err != nil {
		return err
	}

	// make a new session to the GTP Session Repo
	gsid, err := theGtpSessionRepo.newSession(
		s, pgwCtrlIPv4,
		s.toSender,
		sgwCtrlFTEID,
		sgwDataFTEID,
		imsiIE,
		msisdnIE,
		ebiIE,
		paaIE,
		apnIE,
		ambrIE,
		ratTypeIE,
		servingNetworkID,
		pdnType,
	)
	if err != nil {
		return err
	}

	// Make GTP Session CMD
	cmd, err := gtpSessionCmd.NewCreateSessionReq(mcc, mnc, mei)
	if err != nil {
		return err
	}

	// Send the CMD to the session's CMD chan
	cmdChan := theGtpSessionRepo.getBySessionID(gsid).cmdChan
	cmdChan <- cmd

	res := <-cmdChan
	fmt.Print(res)

	return fmt.Errorf("Now be implementing")
}
