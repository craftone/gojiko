package main

import (
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
		return ctx.NotFound(err)
	}

	pl := ctx.Payload
	destAddr := net.UDPAddr{IP: net.ParseIP(pl.DestAddr), Port: pl.DestPort}

	udpFlow := domain.UdpEchoFlowArg{
		DestAddr:       destAddr,
		SourcePort:     uint16(pl.SourcePort),
		SendPacketSize: uint16(pl.SendPacketSize),
		Tos:            byte(pl.Tos),
		Ttl:            byte(pl.TTL),
		TargetBps:      uint64(pl.TargetBps),
		NumOfSend:      pl.NumOfSend,
		RecvPacketSize: uint16(pl.RecvPacketSize),
	}

	err = sess.NewUdpFlow(udpFlow)
	if err != nil {
		return ctx.InternalServerError(err)
	}

	return ctx.OK(&app.Udpechoflow{
		UDPEchoFlowArg: pl,
	})

	// UDPEchoFlowByIMSIandEBIController_Create: end_implement
}
