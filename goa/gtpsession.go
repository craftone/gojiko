package main

import (
	"fmt"
	"net"

	"github.com/craftone/gojiko/gtp"

	"github.com/craftone/gojiko/domain"
	"github.com/craftone/gojiko/goa/app"
	"github.com/goadesign/goa"
)

// GtpsessionController implements the gtpsession resource.
type GtpsessionController struct {
	*goa.Controller
}

// NewGtpsessionController creates a gtpsession controller.
func NewGtpsessionController(service *goa.Service) *GtpsessionController {
	return &GtpsessionController{Controller: service.NewController("GtpsessionController")}
}

// Create runs the create action.
func (c *GtpsessionController) Create(ctx *app.CreateGtpsessionContext) error {
	// GtpsessionController_Create: start_implement

	sgwCtrlAddr := net.UDPAddr{IP: net.ParseIP(ctx.SgwAddr), Port: domain.GtpControlPort}
	theSgwCtrlRepo := domain.TheSgwCtrlRepo()
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(sgwCtrlAddr)
	if sgwCtrl == nil {
		return ctx.BadRequest(goa.ErrBadRequest(fmt.Errorf("There are no SGW that's IP address is %s", sgwCtrlAddr.String())))
	}

	payload := ctx.Payload
	csRes, err := sgwCtrl.CreateSession(
		payload.Imsi, payload.Msisdn, payload.Mei, payload.Mcc, payload.Mnc,
		payload.Apn, byte(payload.Ebi))
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	switch csRes.Code {
	case domain.GscResOK:
		sess := csRes.Session

		res := &app.Gtpsession{
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
		return ctx.OK(res)
	}
	return ctx.InternalServerError(goa.ErrInternal(csRes.Msg))

	// GtpsessionController_Create: end_implement
}

func newFteid(ip net.IP, teid gtp.Teid) *app.Fteid {
	return &app.Fteid{Ipv4: ip.String(), Teid: fmt.Sprintf("0x%08X", teid)}
}
