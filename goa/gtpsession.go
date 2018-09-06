package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/craftone/gojiko/domain"
	"github.com/craftone/gojiko/domain/apns"
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

	sgwCtrl, err := querySgw(ctx.SgwAddr)
	if err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}

	payload := ctx.Payload
	var sgwDataAddr *net.IP
	if payload.ExternalSgwDataAddr != nil {
		ip := net.ParseIP(*payload.ExternalSgwDataAddr)
		sgwDataAddr = &ip
	}
	csRes, sess, err := sgwCtrl.CreateSession(
		payload.Imsi, payload.Msisdn, payload.Mei, payload.Mcc, payload.Mnc,
		payload.Apn, byte(payload.Ebi), sgwDataAddr)
	if err != nil {
		switch err.(type) {
		case *domain.DuplicateSessionError:
			return ctx.Conflict(goa.NewErrorClass("conflict", 409)(err))
		case *apns.NoSuchAPNError:
			return ctx.BadRequest(goa.ErrBadRequest(err))
		default:
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	}

	switch csRes.Code {
	case domain.GsResOK:
		return ctx.OK(&app.Gtpv2cCsres{
			Cause:       newCauseMedia(csRes),
			SessionInfo: newGtpsessionMedia(sess),
		})
	case domain.GsResRetryableNG:
		return ctx.ServiceUnavailable(newCauseMedia(csRes))
	case domain.GsResTimeout:
		return ctx.GatewayTimeout(goa.NewErrorClass("gateway_timeout", 504)(errors.New(csRes.Msg)))
	}
	return ctx.InternalServerError(goa.ErrInternal(csRes.Msg))

	// GtpsessionController_Create: end_implement
}

// DeleteByIMSIandEBI runs the deleteByIMSIandEBI action.
func (c *GtpsessionController) DeleteByIMSIandEBI(ctx *app.DeleteByIMSIandEBIGtpsessionContext) error {
	// GtpsessionController_DeleteByIMSIandEBI: start_implement

	sgwCtrl, err := querySgw(ctx.SgwAddr)
	if err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}

	_, err = querySessionByIMSIandEBI(ctx.SgwAddr, ctx.Imsi, ctx.Ebi)
	if err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}

	gsRes, err := sgwCtrl.DeleteSession(ctx.Imsi, byte(ctx.Ebi))
	if err != nil {
		switch err.(type) {
		case *domain.InvalidGtpSessionStateError:
			return ctx.Conflict(goa.NewErrorClass("conflict", 409)(err))
		default:
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	}

	return ctx.OK(newCauseMedia(gsRes))

	// GtpsessionController_DeleteByIMSIandEBI: end_implement
}

// ShowByID runs the showByID action.
func (c *GtpsessionController) ShowByID(ctx *app.ShowByIDGtpsessionContext) error {
	// GtpsessionController_ShowByID: start_implement

	sgwCtrl, err := querySgw(ctx.SgwAddr)
	if err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}
	sess := sgwCtrl.FindBySessionID(domain.SessionID(ctx.Sid))
	if sess == nil {
		return ctx.NotFound(goa.ErrNotFound(fmt.Errorf("There is no session that's ID is %d", ctx.Sid)))
	}
	res := newGtpsessionMedia(sess)
	return ctx.OK(res)

	// GtpsessionController_ShowByID: end_implement
}

// ShowByIMSIandEBI runs the showByIMSIandEBI action.
func (c *GtpsessionController) ShowByIMSIandEBI(ctx *app.ShowByIMSIandEBIGtpsessionContext) error {
	// GtpsessionController_ShowByIMSIandEBI: start_implement

	sess, err := querySessionByIMSIandEBI(ctx.SgwAddr, ctx.Imsi, ctx.Ebi)
	if err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}
	return ctx.OK(newGtpsessionMedia(sess))

	// GtpsessionController_ShowByIMSIandEBI: end_implement
}
