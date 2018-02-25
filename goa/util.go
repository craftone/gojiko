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
