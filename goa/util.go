package main

import (
	"fmt"
	"net"

	"github.com/craftone/gojiko/domain"
	"github.com/craftone/gojiko/goa/app"
	"github.com/craftone/gojiko/gtp"
	"github.com/goadesign/goa"
)

func newFteid(ip net.IP, teid gtp.Teid) *app.Fteid {
	return &app.Fteid{Ipv4: ip.String(), Teid: fmt.Sprintf("0x%08X", teid)}
}

func querySgw(sgwAddr string) (*domain.SgwCtrl, error) {
	sgwCtrlAddr := net.UDPAddr{IP: net.ParseIP(sgwAddr), Port: domain.GtpControlPort}
	theSgwCtrlRepo := domain.TheSgwCtrlRepo()
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(sgwCtrlAddr)
	if sgwCtrl == nil {
		return nil, goa.ErrBadRequest(fmt.Errorf("There is no SGW that's IP address is %s", sgwCtrlAddr.String()))
	}
	return sgwCtrl, nil
}

func querySessionByIMSIandEBI(sgwAddr, imsi string, ebi int) (*domain.GtpSession, error) {
	sgwCtrl, err := querySgw(sgwAddr)
	if err != nil {
		return nil, err
	}
	sess := sgwCtrl.FindByImsiEbi(imsi, byte(ebi))
	if sess == nil {
		return nil, goa.ErrNotFound(fmt.Errorf("There is no session that's IMSI is %s and EBI is %d", imsi, ebi))
	}
	return sess, nil
}

func newGtpsessionMedia(sess *domain.GtpSession) *app.Gtpsession {
	return &app.Gtpsession{
		Apn: sess.Apn(),
		Ebi: int(sess.Ebi()),
		Fteid: &app.GtpSessionFTEIDs{
			PgwCtrlFTEID: newFteid(sess.PgwCtrlFTEID()),
			PgwDataFTEID: newFteid(sess.PgwDataFTEID()),
			SgwCtrlFTEID: newFteid(sess.SgwCtrlFTEID()),
			SgwDataFTEID: newFteid(sess.SgwDataFTEID()),
		},
		ID:     int(sess.ID()),
		Imsi:   sess.Imsi(),
		Mcc:    sess.Mcc(),
		Mei:    sess.Mei(),
		Mnc:    sess.Mnc(),
		Msisdn: sess.Msisdn(),
	}
}
