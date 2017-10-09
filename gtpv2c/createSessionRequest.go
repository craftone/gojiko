package gtpv2c

import (
	"errors"

	"github.com/craftone/gojiko/gtpv2c/ie"
)

type CreateSessionRequest struct {
	header
	imsi              *ie.Imsi
	msisdn            *ie.Msisdn
	mei               *ie.Mei
	uli               *ie.Uli
	servingNetwork    *ie.ServingNetwork
	ratType           *ie.RatType
	indication        *ie.Indication
	senderFteid       *ie.Fteid
	apn               *ie.Apn
	selectionMode     *ie.SelectionMode
	pdnType           *ie.PdnType
	paa               *ie.Paa
	maxApnRestriction *ie.ApnRestriction
	apnAmbr           *ie.Ambr
	ebi               *ie.Ebi
	pco               *ie.PcoMsToNetwork
	bearerContextTBC  *ie.BearerContextToBeCreatedWithinCSReq
	recovery          *ie.Recovery
}

type CreateSessionRequestArg struct {
	Imsi              *ie.Imsi
	Msisdn            *ie.Msisdn
	Mei               *ie.Mei
	Uli               *ie.Uli
	ServingNetwork    *ie.ServingNetwork
	RatType           *ie.RatType
	Indication        *ie.Indication
	SenderFteid       *ie.Fteid
	Apn               *ie.Apn
	SelectionMode     *ie.SelectionMode
	PdnType           *ie.PdnType
	Paa               *ie.Paa
	MaxApnRestriction *ie.ApnRestriction
	ApnAmbr           *ie.Ambr
	Ebi               *ie.Ebi
	Pco               *ie.PcoMsToNetwork
	BearerContextTBC  *ie.BearerContextToBeCreatedWithinCSReq
	Recovery          *ie.Recovery
}

func NewCreateSessionRequest(seqNum uint32, csReqArg CreateSessionRequestArg) (*CreateSessionRequest, error) {
	if err := checkCreateSessionRequestArg(csReqArg); err != nil {
		return nil, err
	}

	return &CreateSessionRequest{
		newHeader(createSessionRequestNum, false, true, 0, seqNum),
		csReqArg.Imsi,
		csReqArg.Msisdn,
		csReqArg.Mei,
		csReqArg.Uli,
		csReqArg.ServingNetwork,
		csReqArg.RatType,
		csReqArg.Indication,
		csReqArg.SenderFteid,
		csReqArg.Apn,
		csReqArg.SelectionMode,
		csReqArg.PdnType,
		csReqArg.Paa,
		csReqArg.MaxApnRestriction,
		csReqArg.ApnAmbr,
		csReqArg.Ebi,
		csReqArg.Pco,
		csReqArg.BearerContextTBC,
		csReqArg.Recovery,
	}, nil
}

func checkCreateSessionRequestArg(csReqArg CreateSessionRequestArg) error {
	errMsgs := make([]string, 0)
	if csReqArg.Imsi == nil {
		errMsgs = append(errMsgs, "IMSI")
	}
	if csReqArg.Msisdn == nil {
		errMsgs = append(errMsgs, "MSISDN")
	}
	if csReqArg.Mei == nil {
		errMsgs = append(errMsgs, "MEI")
	}
	if csReqArg.Uli == nil {
		errMsgs = append(errMsgs, "ULI")
	}
	if csReqArg.ServingNetwork == nil {
		errMsgs = append(errMsgs, "Serving Network")
	}
	if csReqArg.RatType == nil {
		errMsgs = append(errMsgs, "RatType")
	}
	if csReqArg.Indication == nil {
		errMsgs = append(errMsgs, "Indication")
	}
	if csReqArg.SenderFteid == nil {
		errMsgs = append(errMsgs, "Sender F-TEID")
	}
	if csReqArg.Apn == nil {
		errMsgs = append(errMsgs, "APN")
	}
	if csReqArg.SelectionMode == nil {
		errMsgs = append(errMsgs, "Selection Mode")
	}
	if csReqArg.PdnType == nil {
		errMsgs = append(errMsgs, "PDN Type")
	}
	if csReqArg.Paa == nil {
		errMsgs = append(errMsgs, "PAA")
	}
	if csReqArg.MaxApnRestriction == nil {
		errMsgs = append(errMsgs, "Max APN Restriction")
	}
	if csReqArg.ApnAmbr == nil {
		errMsgs = append(errMsgs, "APN-AMBR")
	}
	if csReqArg.Ebi == nil {
		errMsgs = append(errMsgs, "EBI")
	}
	if csReqArg.Pco == nil {
		errMsgs = append(errMsgs, "PCO")
	}
	if csReqArg.BearerContextTBC == nil {
		errMsgs = append(errMsgs, "Bearer Context to be created")
	}
	if csReqArg.Recovery == nil {
		errMsgs = append(errMsgs, "Recovery")
	}
	if len(errMsgs) == 0 {
		return nil
	}
	errMsg := ""
	for _, msg := range errMsgs {
		errMsg += msg + " must be specified. "
	}
	return errors.New(errMsg)
}

func MakeCSReqArg(imsi, msisdn, mei string) (CreateSessionRequestArg, error) {
	imsiIE, err := ie.NewImsi(0, imsi)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	msisdnIE, err := ie.NewMsisdn(0, msisdn)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	meiIE, err := ie.NewMei(0, mei)
	if err != nil {
		return CreateSessionRequestArg{}, err
	}

	return CreateSessionRequestArg{
		Imsi:   imsiIE,
		Msisdn: msisdnIE,
		Mei:    meiIE,
	}, nil
}
