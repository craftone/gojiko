package domain

import (
	"fmt"
	"net"

	"github.com/craftone/gojiko/domain/gtpv2c/ie"
	"github.com/craftone/gojiko/domain/gtpv2c/ie/pco"
)

// createSessionReq respresents the command message that
// asks the session routine to send a Create Session Request.
type createSessionReq struct {
	mei            *ie.Mei
	indication     *ie.Indication
	pco            *ie.PcoMsToNetwork
	uli            *ie.Uli
	bearerQoS      *ie.BearerQoS
	apnRestriction *ie.ApnRestriction
	selectionMode  *ie.SelectionMode
}

func (g createSessionReq) gscType() string {
	return "CreateSession"
}

func (g createSessionReq) String() string {
	return fmt.Sprintf("TYPE=%s MEI=%s", g.gscType(), g.mei.Value())
}

func NewCreateSessionReq(mcc, mnc, mei string) (createSessionReq, error) {
	dummy := createSessionReq{}

	meiIE, err := ie.NewMei(0, mei)
	if err != nil {
		return dummy, err
	}

	indicationIE, _ := ie.NewIndication(0, ie.IndicationArg{})

	msToNetwork := pco.NewMsToNetwork(
		pco.NewIpcp(pco.ConfigureRequest, 0, net.IPv4(0, 0, 0, 0), net.IPv4(0, 0, 0, 0)),
		true, false, true,
	)
	pcoMsToNetwork, _ := ie.NewPcoMsToNetwork(0, msToNetwork)

	ecgiIE, err := ie.NewEcgi(mcc, mnc, 0x22D6600)
	if err != nil {
		return dummy, err
	}

	taiIE, err := ie.NewTai(mcc, mnc, 0x1421)
	if err != nil {
		return dummy, err
	}

	uliArg := ie.UliArg{Ecgi: ecgiIE, Tai: taiIE}
	uliIE, _ := ie.NewUli(0, uliArg)

	bearerQoSArg := ie.BearerQoSArg{
		Pci:         true,
		Pl:          15,
		Pvi:         false,
		Label:       9,
		UplinkMBR:   0,
		DownlinkMBR: 0,
		UplinkGBR:   0,
		DownlinkGBR: 0,
	}
	bearerQoS, _ := ie.NewBearerQoS(0, bearerQoSArg)

	apnRestriction, _ := ie.NewApnRestriction(0, 0)
	selectionMode, _ := ie.NewSelectionMode(0, 1)

	return createSessionReq{
		mei:            meiIE,
		indication:     indicationIE,
		pco:            pcoMsToNetwork,
		uli:            uliIE,
		bearerQoS:      bearerQoS,
		apnRestriction: apnRestriction,
		selectionMode:  selectionMode,
	}, nil
}
