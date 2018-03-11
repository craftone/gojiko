package main

import (
	"errors"
	"net"

	"github.com/craftone/gojiko/domain"
	"github.com/craftone/gojiko/goa/app"
	"github.com/goadesign/goa"
)

// UDPEchoFlowByIMSIandEBIController implements the udpEchoFlowByIMSIandEBI resource.
type UDPEchoFlowByIMSIandEBIController struct {
	*goa.Controller
}

// NewUDPEchoFlowByIMSIandEBIController creates a udpEchoFlowByIMSIandEBI controller.
func NewUDPEchoFlowByIMSIandEBIController(service *goa.Service) *UDPEchoFlowByIMSIandEBIController {
	return &UDPEchoFlowByIMSIandEBIController{Controller: service.NewController("UDPEchoFlowByIMSIandEBIController")}
}

// Create runs the create action.
func (c *UDPEchoFlowByIMSIandEBIController) Create(ctx *app.CreateUDPEchoFlowByIMSIandEBIContext) error {
	// UDPEchoFlowByIMSIandEBIController_Create: start_implement

	sess, err := querySessionByIMSIandEBI(ctx.SgwAddr, ctx.Imsi, ctx.Ebi)
	if err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}

	pl := ctx.Payload
	destAddr := net.UDPAddr{IP: net.ParseIP(pl.DestAddr), Port: pl.DestPort}

	udpFlowArg := domain.UdpEchoFlowArg{
		DestAddr:       destAddr,
		SourcePort:     uint16(pl.SourcePort),
		SendPacketSize: uint16(pl.SendPacketSize),
		Tos:            byte(pl.Tos),
		Ttl:            byte(pl.TTL),
		TargetBps:      uint64(pl.TargetBps),
		NumOfSend:      pl.NumOfSend,
		RecvPacketSize: uint16(pl.RecvPacketSize),
	}

	err = sess.NewUdpFlow(udpFlowArg)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(&app.Udpechoflow{
		Param: pl,
	})

	// UDPEchoFlowByIMSIandEBIController_Create: end_implement
}

// Delete runs the delete action.
func (c *UDPEchoFlowByIMSIandEBIController) Delete(ctx *app.DeleteUDPEchoFlowByIMSIandEBIContext) error {
	// UDPEchoFlowByIMSIandEBIController_Delete: start_implement

	sess, err := querySessionByIMSIandEBI(ctx.SgwAddr, ctx.Imsi, ctx.Ebi)
	if err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}

	udpFlow, ok := sess.UdpFlow()
	if !ok {
		return ctx.NotFound(goa.ErrNotFound(errors.New("No UDP ECHO flow")))
	}

	udpFlow.Stop()

	udpFlowArg := &app.UDPEchoFlowPayload{
		DestAddr:       udpFlow.Arg.DestAddr.IP.String(),
		DestPort:       udpFlow.Arg.DestAddr.Port,
		NumOfSend:      udpFlow.Arg.NumOfSend,
		RecvPacketSize: int(udpFlow.Arg.RecvPacketSize),
		SendPacketSize: int(udpFlow.Arg.SendPacketSize),
		SourcePort:     int(udpFlow.Arg.SourcePort),
		TargetBps:      int(udpFlow.Arg.TargetBps),
		Tos:            int(udpFlow.Arg.Tos),
		TTL:            int(udpFlow.Arg.Ttl),
	}
	res := &app.Udpechoflow{
		Param: udpFlowArg,
	}
	return ctx.OK(res)
	// UDPEchoFlowByIMSIandEBIController_Delete: end_implement
}
