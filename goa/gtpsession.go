package main

import (
	"net"

	"github.com/craftone/gojiko/domain"
	"github.com/craftone/gojiko/domain/gtpSessionCmd"
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

	sgwCtrlAddr := net.UDPAddr{IP: net.ParseIP(ctx.Payload.SgwAddr), Port: domain.GtpControlPort}
	theSgwCtrlRepo := domain.TheSgwCtrlRepo()
	sgwCtrl := theSgwCtrlRepo.GetSgwCtrl(sgwCtrlAddr)
	payload := ctx.Payload
	csRes, err := sgwCtrl.CreateSession(
		payload.Imsi, payload.Msisdn, payload.Mei, payload.Mcc, payload.Mnc,
		payload.Apn, byte(payload.Ebi))
	if err != nil {
		return goa.ErrInternal(err)
	}
	switch csRes.Code {
	case gtpSessionCmd.ResOK:
		res := &app.Gtpsession{}
		return ctx.OK(res)
	case gtpSessionCmd.ResTimeout:
		return goa.ErrInternal("Request timed out")
	}
	return goa.ErrInvalidRequest("Invalid request")

	// GtpsessionController_Create: end_implement
}
